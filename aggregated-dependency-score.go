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

type EvaluatorInterface interface {
	EvaluateScore(ctx context.Context, p Package) (float64, error)
}

func (e *Evaluator) EvaluateScore(ctx context.Context, p Package) (float64, error) {
	aggregatedTrustworthiness, err := e.trustworthiness.evaluate(ctx, p)
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

// compile-time interface check
var _ EvaluatorInterface = &Evaluator{}

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

func (eval *trustwhorthinessEvaluator) evaluate(ctx context.Context, p Package) (float64, error) {
	intrinsic, err := eval.intrinsic.EvaluateIntrinsicTrustworthiness(ctx, p)
	if err != nil {
		return 0.0, fmt.Errorf("evaluating intrinsic trustworthiness of package: %w", err)
	}

	result := intrinsic

	deps, err := eval.deps.GetDirectDependencies(ctx, p)
	if err != nil {
		return 0.0, fmt.Errorf("getting direct dependencies of package: %w", err)
	}

	for _, dep := range deps {
		tPrimeQ, err := eval.evaluate(ctx, dep)
		if err != nil {
			return 0.0, fmt.Errorf("evaluating aggregated trustworthiness of %s: %w", dep.Name, err)
		}

		result *= math.Pow(tPrimeQ, transitiveTrustworthinessExponent)
	}

	return result, nil
}
