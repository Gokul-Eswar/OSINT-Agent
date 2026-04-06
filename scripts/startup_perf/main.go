package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"
)

type startupStats struct {
	Command string  `json:"command"`
	Runs    int     `json:"runs"`
	MeanMS  float64 `json:"mean_ms"`
	MinMS   float64 `json:"min_ms"`
	MaxMS   float64 `json:"max_ms"`
}

type report struct {
	Binary  string       `json:"binary"`
	Version startupStats `json:"version"`
	Help    startupStats `json:"help"`
	TimeUTC time.Time    `json:"time_utc"`
}

func measure(binary string, args []string, label string, runs int) (startupStats, error) {
	if runs <= 0 {
		return startupStats{}, fmt.Errorf("runs must be > 0")
	}

	minVal := 1e12
	maxVal := 0.0
	total := 0.0

	for i := 0; i < runs; i++ {
		start := time.Now()
		cmd := exec.Command(binary, args...)
		cmd.Stdout = io.Discard
		cmd.Stderr = io.Discard
		if err := cmd.Run(); err != nil {
			return startupStats{}, fmt.Errorf("failed running %q: %w", label, err)
		}
		elapsed := float64(time.Since(start).Microseconds()) / 1000.0

		total += elapsed
		if elapsed < minVal {
			minVal = elapsed
		}
		if elapsed > maxVal {
			maxVal = elapsed
		}
	}

	return startupStats{
		Command: label,
		Runs:    runs,
		MeanMS:  total / float64(runs),
		MinMS:   minVal,
		MaxMS:   maxVal,
	}, nil
}

func requiredMean(baselineMS float64, minImprovementPct float64) float64 {
	return baselineMS * (1.0 - (minImprovementPct / 100.0))
}

func main() {
	binary := flag.String("binary", "./bin/spectre", "Path to spectre binary")
	runs := flag.Int("runs", 20, "Number of runs per command")
	versionBaselineMS := flag.Float64("version-baseline-ms", 0, "Baseline mean startup ms for `spectre version`")
	helpBaselineMS := flag.Float64("help-baseline-ms", 0, "Baseline mean startup ms for `spectre --help`")
	minImprovementPct := flag.Float64("min-improvement-pct", 0, "Required improvement percent vs baseline for gate checks")
	jsonOut := flag.String("json-out", "", "Optional path to write JSON report")
	flag.Parse()

	versionStats, err := measure(*binary, []string{"version"}, "version", *runs)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	helpStats, err := measure(*binary, []string{"--help"}, "--help", *runs)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	rep := report{
		Binary:  *binary,
		Version: versionStats,
		Help:    helpStats,
		TimeUTC: time.Now().UTC(),
	}

	fmt.Printf("Startup benchmark (%d runs)\n", *runs)
	fmt.Printf("  version  mean=%.2fms min=%.2fms max=%.2fms\n", rep.Version.MeanMS, rep.Version.MinMS, rep.Version.MaxMS)
	fmt.Printf("  --help   mean=%.2fms min=%.2fms max=%.2fms\n", rep.Help.MeanMS, rep.Help.MinMS, rep.Help.MaxMS)

	if *jsonOut != "" {
		data, err := json.MarshalIndent(rep, "", "  ")
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		if err := os.WriteFile(*jsonOut, data, 0644); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Printf("Wrote benchmark report: %s\n", *jsonOut)
	}

	if *minImprovementPct > 0 {
		failed := false
		if *versionBaselineMS > 0 {
			required := requiredMean(*versionBaselineMS, *minImprovementPct)
			if rep.Version.MeanMS > required {
				failed = true
				fmt.Fprintf(os.Stderr, "version gate failed: mean %.2fms > required %.2fms (baseline %.2fms, improvement %.2f%%)\n", rep.Version.MeanMS, required, *versionBaselineMS, *minImprovementPct)
			}
		}
		if *helpBaselineMS > 0 {
			required := requiredMean(*helpBaselineMS, *minImprovementPct)
			if rep.Help.MeanMS > required {
				failed = true
				fmt.Fprintf(os.Stderr, "help gate failed: mean %.2fms > required %.2fms (baseline %.2fms, improvement %.2f%%)\n", rep.Help.MeanMS, required, *helpBaselineMS, *minImprovementPct)
			}
		}

		if failed {
			os.Exit(1)
		}
		fmt.Println("startup performance gate passed")
	}
}
