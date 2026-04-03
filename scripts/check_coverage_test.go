package main

import "testing"

func TestParseTotalCoverage_ValidOutput(t *testing.T) {
	input := "github.com/spectre/spectre/internal/server/server.go:23:\tStart\t33.3%\n" +
		"total:\t\t\t\t\t(statements)\t36.1%\n"

	got, err := parseTotalCoverage(input)
	if err != nil {
		t.Fatalf("parseTotalCoverage returned error: %v", err)
	}

	if got != 36.1 {
		t.Fatalf("expected 36.1, got %.1f", got)
	}
}

func TestParseTotalCoverage_InvalidOutput(t *testing.T) {
	_, err := parseTotalCoverage("no total line here")
	if err == nil {
		t.Fatal("expected parseTotalCoverage to fail for invalid output")
	}
}
