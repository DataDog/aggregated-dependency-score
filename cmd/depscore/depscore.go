package main

import (
	"context"
	"fmt"
	"os"

	aggregdepscore "github.com/DataDog/aggregated-dependency-score"
	"github.com/DataDog/aggregated-dependency-score/depsdotdev"
)

func main() {
	err := run()
	if err != nil {
		fmt.Fprintln(os.Stderr, "ERROR: "+err.Error())
		os.Exit(1)
	}
}

func run() error {
	depsdotdevClient, err := depsdotdev.NewClient()
	if err != nil {
		return fmt.Errorf("creating deps.dev client: %w", err)
	}

	evaluator := aggregdepscore.NewEvaluator(depsdotdevClient, depsdotdevClient)
	score, err := evaluator.EvaluateScore(context.Background(), aggregdepscore.Package{
		Ecosystem: "pypi",
		Name:      "requests",
		Version:   "2.28.1",
	})
	if err != nil {
		return fmt.Errorf("evaluating score: %w", err)
	}

	fmt.Println(score)

	return nil
}
