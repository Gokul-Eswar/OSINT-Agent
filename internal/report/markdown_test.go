package report

import (
	"database/sql"
	"strings"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/Gokul-Eswar/Spectre/internal/core"
	"github.com/Gokul-Eswar/Spectre/internal/storage"
)

func TestGenerateMarkdownReport(t *testing.T) {
	// Setup DB
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	storage.DB = db
	storage.Migrate()

	caseID := "report-test-case"
	storage.CreateCase(&core.Case{
		ID:          caseID,
		Name:        "Test Report Case",
		Description: "A case to test report generation",
		Status:      "active",
	})

	storage.CreateEntity(&core.Entity{
		ID:     "e1",
		CaseID: caseID,
		Type:   "ip",
		Value:  "1.2.3.4",
		Source: "manual",
	})

	storage.SaveAnalysis(&core.Analysis{
		CaseID:     caseID,
		Findings:   []string{"Finding X"},
		Confidence: 0.9,
	})

	report, err := GenerateMarkdownReport(caseID)
	if err != nil {
		t.Fatalf("GenerateMarkdownReport failed: %v", err)
	}

	if !strings.Contains(report, "Test Report Case") {
		t.Errorf("Report missing case name")
	}
	if !strings.Contains(report, "1.2.3.4") {
		t.Errorf("Report missing entity value")
	}
	if !strings.Contains(report, "Finding X") {
		t.Errorf("Report missing analysis finding")
	}
	if !strings.Contains(report, "0.90") {
		t.Errorf("Report missing confidence")
	}

}

func TestGenerateMarkdownReport_CaseNotFound(t *testing.T) {
	// Setup isolated DB
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	storage.DB = db
	if err := storage.Migrate(); err != nil {
		t.Fatalf("migrate failed: %v", err)
	}

	_, err = GenerateMarkdownReport("missing-case")
	if err == nil {
		t.Fatal("expected error for missing case")
	}
	if !strings.Contains(err.Error(), "case not found") {
		t.Fatalf("expected case not found error, got: %v", err)
	}
}

func TestEscapeMarkdown(t *testing.T) {
	input := "*bold* _italic_ `code` [link] <script> a|b"
	got := escapeMarkdown(input)

	checks := []string{"\\*bold\\*", "\\_italic\\_", "\\`code\\`", "\\[link\\]", "&lt;script&gt;", "a\\|b"}
	for _, want := range checks {
		if !strings.Contains(got, want) {
			t.Fatalf("expected escaped markdown to contain %q, got %q", want, got)
		}
	}
}
