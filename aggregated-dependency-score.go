package aggregdepscore

import (
	"context"
	"fmt"
	"math"
)

// transitiveTrustworthinessExponent is noted as "e" in the design paper
const transitiveTrustworthinessExponent = 1.5

type Package struct {
	ecosystem string
	name      string
	version   string
}

type IntrinsicTrustworthinessEvaluator interface {
	Evaluate(ctx context.Context, p Package) (float64, error)
}

type DependencyResolver interface {
	GetDirectDependencies(ctx context.Context, p Package) ([]Package, error)
}

type trustwhorthinessEvaluator struct {
	intrinsic IntrinsicTrustworthinessEvaluator
	deps      DependencyResolver
}

type Evaluator struct {
	trustworthiness trustwhorthinessEvaluator
}

func (eval *trustwhorthinessEvaluator) Evaluate(ctx context.Context, p Package) (float64, error) {
	intrinsic, err := eval.intrinsic.Evaluate(ctx, p)
	if err != nil {
		return 0.0, fmt.Errorf("evaluating intrinsic trustworthiness of package: %w", err)
	}

	result := intrinsic

	deps, err := eval.deps.GetDirectDependencies(ctx, p)
	if err != nil {
		return 0.0, fmt.Errorf("getting direct dependencies of package: %w", err)
	}

	for _, dep := range deps {
		s, err := eval.Evaluate(ctx, dep)
		if err != nil {
			return 0.0, fmt.Errorf("evaluating aggreated trustworthiness of %s: %w", dep.name, err)
		}

		result *= math.Pow(s, transitiveTrustworthinessExponent)
	}

	return result, nil
}
