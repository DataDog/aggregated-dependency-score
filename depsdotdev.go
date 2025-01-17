package aggregdepscore

import (
	"context"
	"crypto/tls"
	"fmt"
	"strings"

	api "deps.dev/api/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type client struct {
	depsdotdev api.InsightsClient
	converter  ScoreTrustworthinessConverter
}

// compile-time interface checks
var _ IntrinsicTrustworthinessEvaluator = &client{}
var _ DependencyResolver = &client{}

// NewDepsDotDevClient creates an object that satisfies both the IntrinsicTrustworthinessEvaluator and DependencyResolver interfaces,
// using the deps.dev API as the source of data.
// The intrinsic trustworthiness is calculated based on the OSSF scorecard that is returned by the deps.dev API.
//
// Deprecated: in version 1 of package aggregdepscore,
// the deps.dev client will be moved to a new Go module, most likely in a new repository,
// and this function will be removed.
func NewDepsDotDevClient() (*client, error) {
	connection, err := grpc.NewClient(
		"api.deps.dev:443",
		grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{
			MinVersion: tls.VersionTLS12,
			MaxVersion: tls.VersionTLS13,
		})),
	)
	if err != nil {
		return nil, fmt.Errorf("creating grpc connection: %w", err)
	}

	// TODO let the user ask for caching
	// or provide its own cache

	return &client{
		depsdotdev: api.NewInsightsClient(connection),
		converter:  &DefaultScoreTrustworthinessConverter{},
	}, nil
}

func (c *client) getRespository(ctx context.Context, p Package) (string, error) {
	ecosystem, err := depsdotdevEcosystem(p.Ecosystem)
	if err != nil {
		return "", fmt.Errorf("converting ecosystem: %w", err)
	}

	version, err := c.depsdotdev.GetVersion(ctx, &api.GetVersionRequest{
		VersionKey: &api.VersionKey{
			System:  ecosystem,
			Name:    p.Name,
			Version: p.Version,
		},
	})
	if err != nil {
		return "", fmt.Errorf("fetching package version: %w", err)
	}

	var repository string

	for _, p := range version.RelatedProjects {
		if p == nil {
			continue
		}

		if p.RelationType != api.ProjectRelationType_SOURCE_REPO {
			continue
		}

		if p.ProjectKey == nil {
			continue
		}

		repository = p.ProjectKey.Id

		// XXX assumes there is only one source repository
		break
	}

	if repository != "" {
		return repository, nil
	}

	// deps.dev does not find the repository for gopkg.in packages
	// Cédric Van Rompay reported it to depsdev@google.com on 2025-01-03
	// in the meantime we use this workaround
	if p.Ecosystem == "go" && strings.HasPrefix(p.Name, "gopkg.in/") {
		repository, err := getGopkginRepository(p.Name)
		if err != nil {
			return "", fmt.Errorf("getting repository for gopkg.in package: %w", err)
		}

		return repository, nil
	}

	return "", fmt.Errorf("no source repository found for package version")
}

func (c *client) EvaluateIntrinsicTrustworthiness(ctx context.Context, p Package) (float64, error) {
	repository, err := c.getRespository(ctx, p)
	if err != nil {
		return 0, fmt.Errorf("getting repository: %w", err)
	}

	project, err := c.depsdotdev.GetProject(ctx, &api.GetProjectRequest{
		ProjectKey: &api.ProjectKey{
			Id: repository,
		},
	})
	if err != nil {
		return 0, fmt.Errorf("fetching project (%s): %w", repository, err)
	}

	if project.Scorecard == nil {
		return 0, fmt.Errorf("no scorecard found for project (%s)", repository)
	}

	score := float64(project.Scorecard.OverallScore) / 10.0

	// XXX OSSF scorecard tends to give pretty low scores
	// so we may want to adjust the trustworthiness
	// so that it better represents
	// "the probability that the package turns malicious one day"

	return c.converter.TrustworthinessFromScore(score), nil

}

func (c *client) GetDirectDependencies(ctx context.Context, p Package) ([]Package, error) {
	ecosystem, err := depsdotdevEcosystem(p.Ecosystem)
	if err != nil {
		return nil, fmt.Errorf("converting ecosystem: %w", err)
	}

	versionKey := &api.VersionKey{
		System:  ecosystem,
		Name:    p.Name,
		Version: p.Version,
	}

	dependencies, err := c.depsdotdev.GetDependencies(ctx, &api.GetDependenciesRequest{VersionKey: versionKey})
	if err != nil {
		return nil, fmt.Errorf("fetching dependencies: %w", err)
	}

	var result []Package

	hasBundledDependencies := false

	for _, dep := range dependencies.Nodes {
		if dep == nil || dep.Relation != api.DependencyRelation_DIRECT || dep.VersionKey == nil {
			continue
		}

		if dep.Bundled {
			hasBundledDependencies = true
			continue
		}

		depEcosystem, err := depsdotdevEcosystemString(dep.VersionKey.System)
		if err != nil {
			return nil, fmt.Errorf("converting ecosystem of dependency %v to string: %w", dep.VersionKey, err)
		}

		result = append(result, Package{
			Ecosystem: depEcosystem,
			Name:      dep.VersionKey.Name,
			Version:   dep.VersionKey.Version,
		})
	}

	if hasBundledDependencies {
		// TODO use string constants for the ecosystem names
		if p.Ecosystem == "npm" {
			bundledDependencies, err := c.getNPMBundledDependencies(ctx, versionKey)
			if err != nil {
				return nil, fmt.Errorf("getting bundled dependencies: %w", err)
			}

			result = append(result, bundledDependencies...)
		}

		// TODO (https://github.com/DataDog/aggregated-dependency-score/issues/19) log a warning if hasBundledDependencies is true and the ecosystem is NOT npm
	}

	return result, nil
}

func (c *client) getNPMBundledDependencies(ctx context.Context, versionKey *api.VersionKey) ([]Package, error) {
	ecosystem := "npm"

	if versionKey.System != api.System_NPM {
		return nil, fmt.Errorf("bundled dependencies are only supported for npm")
	}

	requirements, err := c.depsdotdev.GetRequirements(ctx, &api.GetRequirementsRequest{VersionKey: versionKey})
	if err != nil {
		return nil, fmt.Errorf("fetching requirements: %w", err)
	}

	if requirements.Npm == nil {
		return nil, fmt.Errorf("requirements.Npm is nil")
	}

	if requirements.Npm.Bundled == nil {
		return nil, fmt.Errorf("requirements.Npm.Bundled is nil")
	}

	var result []Package

	for _, dep := range requirements.Npm.Bundled {
		if dep == nil {
			// TODO (https://github.com/DataDog/aggregated-dependency-score/issues/19) log a warning
			continue
		}

		// variable dep *api.Requirements_NPM_Bundle has a "Name" field
		// but it does not include the NPM namespace,
		// for instance for package https://www.npmjs.com/package/@balena/dockerignore
		// the Name field is "dockerignore";
		// and we do need the namespace in Package.Name.
		// fortunately, the Path field seems to be "node_modules/" plus the name with namespace
		// so we use it instead.

		if !strings.HasPrefix(dep.Path, "node_modules/") {
			// TODO (https://github.com/DataDog/aggregated-dependency-score/issues/19) log a warning
			continue
		}

		result = append(result, Package{
			Ecosystem: ecosystem,
			Name:      dep.Path[len("node_modules/"):],
			Version:   dep.Version,
		})
	}

	return result, nil
}

func depsdotdevEcosystem(x string) (api.System, error) {
	switch x {
	case "pypi":
		return api.System_PYPI, nil
	case "npm":
		return api.System_NPM, nil
	case "maven":
		return api.System_MAVEN, nil
	case "cargo":
		return api.System_SYSTEM_UNSPECIFIED, fmt.Errorf(`please use "crates.io" instead of "cargo"`)
	case "crates.io":
		return api.System_CARGO, nil
	case "nuget":
		return api.System_NUGET, nil
	case "go":
		return api.System_GO, nil
	default:
		return api.System_SYSTEM_UNSPECIFIED, fmt.Errorf("unknown ecosystem: %q", x)
	}
}

func depsdotdevEcosystemString(x api.System) (string, error) {
	switch x {
	case api.System_PYPI:
		return "pypi", nil
	case api.System_NPM:
		return "npm", nil
	case api.System_MAVEN:
		return "maven", nil
	case api.System_CARGO:
		return "crates.io", nil
	case api.System_GO:
		return "go", nil
	case api.System_NUGET:
		return "nuget", nil
	default:
		return "", fmt.Errorf("unknown ecosystem: %v", x)
	}
}
