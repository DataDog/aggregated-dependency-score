package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	aggregdepscore "github.com/DataDog/aggregated-dependency-score"
)

func main() {
	err := run()
	if err != nil {
		fmt.Fprintln(os.Stderr, "ERROR: "+err.Error())
		os.Exit(1)
	}
}

func run() (err error) {
	testData := struct {
		RealCases []testCase `json:"real_cases"`
	}{}

	testDataBytes, err := os.ReadFile("test-data.json")
	if err != nil {
		return fmt.Errorf("reading test data: %w", err)
	}

	err = json.Unmarshal(testDataBytes, &testData)
	if err != nil {
		return fmt.Errorf("unmarshalling test data: %w", err)
	}

	depsdotdev, err := aggregdepscore.NewDepsDotDevClient()
	if err != nil {
		return fmt.Errorf("creating deps.dev client: %w", err)
	}

	evaluator, err := aggregdepscore.NewEvaluator(depsdotdev, depsdotdev)
	if err != nil {
		return fmt.Errorf("creating evaluator: %w", err)
	}

	var caseName string
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic in test case %q: %v", caseName, r)
		}
	}()

	ctx := context.Background()
	nbCasesRun := 0
	for _, testCase := range testData.RealCases {
		caseName = testCase.Name

		_, err = evaluator.EvaluateScore(ctx, testCase.Package)
		if err != nil {
			return fmt.Errorf("running test case %q: %w", caseName, err)
		}

		nbCasesRun++
	}

	fmt.Printf("All %d test cases passed\n", nbCasesRun)
	return nil
}

type testCase struct {
	Name    string
	Package aggregdepscore.Package
}
