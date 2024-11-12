package aggregdepscore

import (
	"context"
	"fmt"
	"math"
	"testing"
)

type testIntrinsicTrustworthinessEvaluator struct {
	trustworthinessByName map[string]float64
}

func (eval *testIntrinsicTrustworthinessEvaluator) EvaluateIntrinsicTrustworthiness(ctx context.Context, p Package) (float64, error) {
	if s, ok := eval.trustworthinessByName[p.Name]; ok {
		return s, nil
	}

	return 0.0, fmt.Errorf("unknown package %q", p.Name)
}

type testDependencyResolver struct {
	directDependencyNamesByName map[string][]string
}

func (r *testDependencyResolver) GetDirectDependencies(ctx context.Context, p Package) ([]Package, error) {
	names, ok := r.directDependencyNamesByName[p.Name]
	if !ok {
		return nil, fmt.Errorf("unknown package %q", p.Name)
	}

	var deps []Package
	for _, name := range names {
		deps = append(deps, Package{Name: name})
	}

	return deps, nil
}

func TestAggregatedTrustworthinessEvaluation(t *testing.T) {
	t.Run("single package", func(t *testing.T) {
		eval := trustwhorthinessEvaluator{
			intrinsic: &testIntrinsicTrustworthinessEvaluator{
				trustworthinessByName: map[string]float64{
					"A": 0.92,
				},
			},
			deps: &testDependencyResolver{
				directDependencyNamesByName: map[string][]string{
					"A": {},
				},
			},
		}

		tPrimeA, err := eval.evaluate(context.Background(), Package{Name: "A"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		expected := 0.92

		if tPrimeA != expected {
			t.Fatalf("expected %g, got %g", expected, tPrimeA)
		}
	})

	t.Run("graph from blog post", func(t *testing.T) {
		// this represents the dependency graph
		// used as an example in https://cedricvanrompay.fr/blog/aggregated-dependency-score
		eval := trustwhorthinessEvaluator{
			intrinsic: &testIntrinsicTrustworthinessEvaluator{
				trustworthinessByName: map[string]float64{
					"A": 0.92,
					"B": 0.94,
					"C": 0.93,
					"D": 0.84,
					"E": 0.87,
					"F": 0.85,
					"G": 0.91,
					"H": 0.95,
				},
			},
			deps: &testDependencyResolver{
				directDependencyNamesByName: map[string][]string{
					"A": {"B", "C", "D", "E"},
					"B": {"G"},
					"C": {"F"},
					"D": {"F"},
					"E": {},
					"F": {},
					"G": {"H"},
					"H": {},
				},
			},
		}

		tPrimeA, err := eval.evaluate(context.Background(), Package{Name: "A"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		expected := 0.15379680196472223
		// due to floating point arithmetic, we need to allow for a small error
		// as the result may be a bit different from one CPU to another
		allowedError := 1e-10

		if math.Abs(tPrimeA-expected) > allowedError {
			// using %g for full precision
			t.Fatalf("expected %g, got %g (difference greater than %g)", expected, tPrimeA, allowedError)
		}
	})
}
