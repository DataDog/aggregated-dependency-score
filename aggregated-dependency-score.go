package aggregdepscore

import (
	"context"
	"fmt"
	"math"
)

// transitiveTrustworthinessExponent is noted as "e" in the design paper
const transitiveTrustworthinessExponent = 1.5

type Package struct {
	Ecosystem string
	Name      string
	Version   string
}

func (e *Evaluator) EvaluateScore(ctx context.Context, p Package) (float64, error) {
	aggregatedTrustworthiness, err := e.trustworthiness.evaluate(ctx, p, nil)
	if err != nil {
		return 0.0, err
	}

	return e.converter.ScoreFromTrustworthiness(aggregatedTrustworthiness), nil
}

func NewEvaluator(intrinsic IntrinsicTrustworthinessEvaluator, deps DependencyResolver) (*Evaluator, error) {
	if intrinsic == nil {
		return nil, fmt.Errorf("intrinsic trustworthiness evaluator is required")
	}

	if deps == nil {
		return nil, fmt.Errorf("dependency resolver is required")
	}

	return &Evaluator{
		trustworthiness: trustwhorthinessEvaluator{
			intrinsic: intrinsic,
			deps:      deps,
		},
		converter: &DefaultScoreTrustworthinessConverter{},
	}, nil
}

type Evaluator struct {
	trustworthiness trustwhorthinessEvaluator
	converter       ScoreTrustworthinessConverter
}

type trustwhorthinessEvaluator struct {
	intrinsic IntrinsicTrustworthinessEvaluator
	deps      DependencyResolver
}

type IntrinsicTrustworthinessEvaluator interface {
	EvaluateIntrinsicTrustworthiness(ctx context.Context, p Package) (float64, error)
}

type DependencyResolver interface {
	GetDirectDependencies(ctx context.Context, p Package) ([]Package, error)
}

func (evaluator *trustwhorthinessEvaluator) evaluate(ctx context.Context, p Package, ancestors map[string]struct{}) (float64, error) {
	intrinsic, err := evaluator.intrinsic.EvaluateIntrinsicTrustworthiness(ctx, p)
	if err != nil {
		return 0.0, fmt.Errorf("evaluating intrinsic trustworthiness of package: %w", err)
	}

	result := intrinsic

	deps, err := evaluator.deps.GetDirectDependencies(ctx, p)
	if err != nil {
		return 0.0, fmt.Errorf("getting direct dependencies of package: %w", err)
	}

	for _, dep := range deps {
		// XXX sometimes different names can refer to the same package,
		// for instance with gopkg.in URLs;
		// XXX should we consider the version as well?
		if _, ok := ancestors[dep.Name]; ok {
			// depedency cycle
			// TODO emit a log
			continue
		}

		// copy ancestors to avoid modifying the original map;
		// one day we may want to run the algorithm in parallel
		// so we will need to be careful with shared state
		childAncestors := make(map[string]struct{})
		for k, v := range ancestors {
			childAncestors[k] = v
		}
		childAncestors[p.Name] = struct{}{}

		tPrimeQ, err := evaluator.evaluate(ctx, dep, childAncestors)
		if err != nil {
			return 0.0, fmt.Errorf("evaluating aggregated trustworthiness of %s: %w", dep.Name, err)
		}

		result *= math.Pow(tPrimeQ, transitiveTrustworthinessExponent)
	}

	return result, nil
}
