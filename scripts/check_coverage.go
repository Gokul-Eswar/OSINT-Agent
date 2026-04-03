package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
)

func main() {
	profile := flag.String("profile", "coverage.out", "Path to Go coverage profile")
	minCoverage := flag.Float64("min", 80.0, "Minimum total coverage percentage")
	flag.Parse()

	cmd := exec.Command("go", "tool", "cover", "-func", *profile)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read coverage profile %q: %v\n%s", *profile, err, string(output))
		os.Exit(1)
	}

	re := regexp.MustCompile(`total:\s+\(statements\)\s+([0-9]+(?:\.[0-9]+)?)%`)
	matches := re.FindStringSubmatch(string(output))
	if len(matches) != 2 {
		fmt.Fprintf(os.Stderr, "unable to parse total coverage from output:\n%s", string(output))
		os.Exit(1)
	}

	actualCoverage, err := strconv.ParseFloat(matches[1], 64)
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid coverage value %q: %v\n", matches[1], err)
		os.Exit(1)
	}

	fmt.Printf("Total coverage: %.1f%% (required: %.1f%%)\n", actualCoverage, *minCoverage)
	if actualCoverage < *minCoverage {
		fmt.Fprintf(os.Stderr, "coverage gate failed: %.1f%% < %.1f%%\n", actualCoverage, *minCoverage)
		os.Exit(1)
	}

	fmt.Println("coverage gate passed")
}
