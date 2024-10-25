package depsdotdev

import (
	"context"
	"crypto/tls"
	"fmt"

	api "deps.dev/api/v3"
	aggregdepscore "github.com/DataDog/aggregated-dependency-score"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type Client struct {
	depsdotdev api.InsightsClient
	converter  aggregdepscore.ScoreTrustworthinessConverter
}

// compile-time interface checks
var _ aggregdepscore.IntrinsicTrustworthinessEvaluator = &Client{}
var _ aggregdepscore.DependencyResolver = &Client{}

func NewClient() (*Client, error) {
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

	return &Client{
		depsdotdev: api.NewInsightsClient(connection),
		converter:  &aggregdepscore.DefaultScoreTrustworthinessConverter{},
	}, nil
}

func (c *Client) EvaluateIntrinsicTrustworthiness(ctx context.Context, p aggregdepscore.Package) (float64, error) {
	ecosystem, err := depsdotdevEcosystem(p.Ecosystem)
	if err != nil {
		return 0, fmt.Errorf("converting ecosystem: %w", err)
	}

	version, err := c.depsdotdev.GetVersion(ctx, &api.GetVersionRequest{
		VersionKey: &api.VersionKey{
			System:  ecosystem,
			Name:    p.Name,
			Version: p.Version,
		},
	})
	if err != nil {
		return 0, fmt.Errorf("fetching package version: %w", err)
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

	if repository == "" {
		return 0, fmt.Errorf("no source repository found for package version")
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

	return c.converter.TrustworthinessFromScore(score), nil

}

func (c *Client) GetDirectDependencies(ctx context.Context, p aggregdepscore.Package) ([]aggregdepscore.Package, error) {
	ecosystem, err := depsdotdevEcosystem(p.Ecosystem)
	if err != nil {
		return nil, fmt.Errorf("converting ecosystem: %w", err)
	}

	dependencies, err := c.depsdotdev.GetDependencies(ctx, &api.GetDependenciesRequest{
		VersionKey: &api.VersionKey{
			System:  ecosystem,
			Name:    p.Name,
			Version: p.Version,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("fetching dependencies: %w", err)
	}

	var result []aggregdepscore.Package

	for _, dep := range dependencies.Nodes {
		if dep == nil || dep.Relation != api.DependencyRelation_DIRECT || dep.VersionKey == nil {
			continue
		}

		depEcosystem, err := depsdotdevEcosystemString(dep.VersionKey.System)
		if err != nil {
			return nil, fmt.Errorf("converting ecosystem of dependency %v to string: %w", dep.VersionKey, err)
		}

		result = append(result, aggregdepscore.Package{
			Ecosystem: depEcosystem,
			Name:      dep.VersionKey.Name,
			Version:   dep.VersionKey.Version,
		})
	}

	return result, nil
}
