package cli

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/Gokul-Eswar/Spectre/internal/version"
)

func TestContextAndTargetRoundTrip(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("SPECTRE_HOME", tmp)

	if err := SaveContext("case-123"); err != nil {
		t.Fatalf("SaveContext failed: %v", err)
	}
	if err := SaveTarget("example.com"); err != nil {
		t.Fatalf("SaveTarget failed: %v", err)
	}

	gotCase, err := LoadContext()
	if err != nil {
		t.Fatalf("LoadContext failed: %v", err)
	}
	if gotCase != "case-123" {
		t.Fatalf("expected case-123, got %q", gotCase)
	}

	gotTarget, err := LoadTarget()
	if err != nil {
		t.Fatalf("LoadTarget failed: %v", err)
	}
	if gotTarget != "example.com" {
		t.Fatalf("expected example.com, got %q", gotTarget)
	}

	if _, err := os.Stat(filepath.Join(tmp, contextFile)); err != nil {
		t.Fatalf("context file not found: %v", err)
	}
	if _, err := os.Stat(filepath.Join(tmp, targetFile)); err != nil {
		t.Fatalf("target file not found: %v", err)
	}
}

func TestLoadContextMissingFile(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("SPECTRE_HOME", tmp)

	if _, err := LoadContext(); err == nil {
		t.Fatal("expected LoadContext to fail when file does not exist")
	}
}

func TestVersionCommandOutput(t *testing.T) {
	oldOut := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	rootCmd.SetArgs([]string{"version"})
	err := rootCmd.Execute()

	_ = w.Close()
	os.Stdout = oldOut

	if err != nil {
		t.Fatalf("version command failed: %v", err)
	}

	var b bytes.Buffer
	_, _ = io.Copy(&b, r)
	output := b.String()

	want := "SPECTRE " + version.Version
	if !bytes.Contains([]byte(output), []byte(want)) {
		t.Fatalf("expected output to contain %q, got %q", want, output)
	}
}
