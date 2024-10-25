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

type IntrinsicTrustworthinessEvaluator interface {
	EvaluateIntrinsicTrustworthiness(ctx context.Context, p Package) (float64, error)
}

type DependencyResolver interface {
	GetDirectDependencies(ctx context.Context, p Package) ([]Package, error)
}

type trustwhorthinessEvaluator struct {
	intrinsic IntrinsicTrustworthinessEvaluator
	deps      DependencyResolver
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
		s, err := eval.evaluate(ctx, dep)
		if err != nil {
			return 0.0, fmt.Errorf("evaluating aggreated trustworthiness of %s: %w", dep.Name, err)
		}

		result *= math.Pow(s, transitiveTrustworthinessExponent)
	}

	return result, nil
}

type Evaluator struct {
	trustworthiness trustwhorthinessEvaluator
	converter       ScoreTrustworthinessConverter
}

func NewEvaluator(intrinsic IntrinsicTrustworthinessEvaluator, deps DependencyResolver) *Evaluator {
	return &Evaluator{
		trustworthiness: trustwhorthinessEvaluator{
			intrinsic: intrinsic,
			deps:      deps,
		},
		converter: &DefaultScoreTrustworthinessConverter{},
	}
}

func (e *Evaluator) EvaluateScore(ctx context.Context, p Package) (float64, error) {
	aggregatedTrustworthiness, err := e.trustworthiness.evaluate(ctx, p)
	if err != nil {
		return 0.0, err
	}

	return e.converter.ScoreFromTrustworthiness(aggregatedTrustworthiness), nil
}
