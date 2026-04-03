package report

import (
	"database/sql"
	"os"
	"path/filepath"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/spectre/spectre/internal/core"
	"github.com/spectre/spectre/internal/storage"
	"github.com/spf13/viper"
)

func setupReportTestDB(t *testing.T) {
	t.Helper()

	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = db.Close()
	})

	storage.DB = db
	if err := storage.Migrate(); err != nil {
		t.Fatalf("migrate failed: %v", err)
	}
}

func TestGeneratePDFReport_CaseNotFound(t *testing.T) {
	setupReportTestDB(t)

	err := GeneratePDFReport("missing-case", filepath.Join(t.TempDir(), "missing.pdf"))
	if err == nil {
		t.Fatal("expected error for missing case")
	}
}

func TestGeneratePDFReport_Success(t *testing.T) {
	setupReportTestDB(t)
	viper.Reset()
	viper.Set("report.branding.company", "Acme")
	viper.Set("report.branding.header", "Investigation Report")
	viper.Set("report.branding.footer", "Internal")

	caseID := "pdf-test-case"
	err := storage.CreateCase(&core.Case{
		ID:          caseID,
		Name:        "PDF Test Case",
		Description: "A case to test PDF report generation",
		Status:      "active",
	})
	if err != nil {
		t.Fatalf("failed to create case: %v", err)
	}

	err = storage.CreateEntity(&core.Entity{
		ID:     "pdf-e1",
		CaseID: caseID,
		Type:   "domain",
		Value:  "example.com",
		Source: "manual",
	})
	if err != nil {
		t.Fatalf("failed to create entity: %v", err)
	}

	outPath := filepath.Join(t.TempDir(), "report.pdf")
	err = GeneratePDFReport(caseID, outPath)
	if err != nil {
		t.Fatalf("GeneratePDFReport failed: %v", err)
	}

	info, err := os.Stat(outPath)
	if err != nil {
		t.Fatalf("expected PDF output file, got stat error: %v", err)
	}
	if info.Size() == 0 {
		t.Fatal("expected non-empty PDF output")
	}
}