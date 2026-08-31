package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
)

func coverageOutput(profile string) (string, error) {
	cmd := exec.Command("go", "tool", "cover", "-func", profile)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to read coverage profile %q: %v\n%s", profile, err, string(output))
	}
	return string(output), nil
}

func parseTotalCoverage(output string) (float64, error) {
	re := regexp.MustCompile(`total:\s+\(statements\)\s+([0-9]+(?:\.[0-9]+)?)%`)
	matches := re.FindStringSubmatch(output)
	if len(matches) != 2 {
		return 0, errors.New("unable to parse total coverage from output")
	}

	actualCoverage, err := strconv.ParseFloat(matches[1], 64)
	if err != nil {
		return 0, fmt.Errorf("invalid coverage value %q: %w", matches[1], err)
	}

	return actualCoverage, nil
}

func main() {
	profile := flag.String("profile", "coverage.out", "Path to Go coverage profile")
	minCoverage := flag.Float64("min", 35.0, "Minimum total coverage percentage")
	flag.Parse()

	output, err := coverageOutput(*profile)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	actualCoverage, err := parseTotalCoverage(output)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n%s", err.Error(), output)
		os.Exit(1)
	}

	fmt.Printf("Total coverage: %.1f%% (required: %.1f%%)\n", actualCoverage, *minCoverage)
	if actualCoverage < *minCoverage {
		fmt.Fprintf(os.Stderr, "coverage gate failed: %.1f%% < %.1f%%\n", actualCoverage, *minCoverage)
		os.Exit(1)
	}

	fmt.Println("coverage gate passed")
}
