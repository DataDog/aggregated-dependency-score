package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	aggregdepscore "github.com/DataDog/aggregated-dependency-score"
)

var (
	ecosystem   = flag.String("ecosystem", "", "Ecosystem of the package")
	packageName = flag.String("package", "", "Name of the package")
	version     = flag.String("version", "", "Version of the package")
)

func main() {
	err := run()
	if err != nil {
		fmt.Fprintln(os.Stderr, "ERROR: "+err.Error())
		os.Exit(1)
	}
}

func run() error {
	flag.Parse()

	err := validateFlags()
	if err != nil {
		flag.Usage()
		return fmt.Errorf("validating flags: %w", err)
	}

	depsdotdev, err := aggregdepscore.NewDepsDotDevClient()
	if err != nil {
		return fmt.Errorf("creating deps.dev client: %w", err)
	}

	evaluator, err := aggregdepscore.NewEvaluator(depsdotdev, depsdotdev)
	if err != nil {
		return fmt.Errorf("creating evaluator: %w", err)
	}
	score, err := evaluator.EvaluateScore(context.Background(), aggregdepscore.Package{
		Ecosystem: *ecosystem,
		Name:      *packageName,
		Version:   *version,
	})
	if err != nil {
		return fmt.Errorf("evaluating score: %w", err)
	}

	fmt.Println(score)

	return nil
}

func validateFlags() error {
	if *ecosystem == "" {
		return fmt.Errorf("ecosystem is required")
	}
	if *packageName == "" {
		return fmt.Errorf("package is required")
	}
	if *version == "" {
		return fmt.Errorf("version is required")
	}

	return nil
}
