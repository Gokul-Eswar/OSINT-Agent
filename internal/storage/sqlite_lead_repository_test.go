package storage

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/Gokul-Eswar/Spectre/internal/core"
)

func TestCreateAndGetLead(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	oldDB := DB
	DB = db
	defer func() { DB = oldDB }()

	if err := Migrate(); err != nil {
		t.Fatal(err)
	}

	// Create a case first for foreign key constraint
	c := &core.Case{ID: "case-1", Name: "Case 1"}
	if err := CreateCase(c); err != nil {
		t.Fatal(err)
	}

	lead := &core.IntelligenceLead{
		ID:          "lead-1",
		CaseID:      "case-1",
		Hypothesis:  "Target is hosting on a known vulnerable server",
		Confidence:  0.8,
		EvidenceIDs: []string{"evidence-1", "evidence-2"},
		Status:      "active",
	}

	// Trigger hook to verify it's called
	hookCalled := false
	OnLeadCreated = func(l *core.IntelligenceLead) {
		if l.ID == "lead-1" {
			hookCalled = true
		}
	}
	defer func() { OnLeadCreated = nil }()

	if err := CreateLead(lead); err != nil {
		t.Fatalf("CreateLead failed: %v", err)
	}

	if !hookCalled {
		t.Error("expected OnLeadCreated hook to be called")
	}

	retrieved, err := GetLead("lead-1")
	if err != nil {
		t.Fatalf("GetLead failed: %v", err)
	}

	if retrieved == nil {
		t.Fatal("expected retrieved lead to not be nil")
	}

	if retrieved.Hypothesis != "Target is hosting on a known vulnerable server" {
		t.Errorf("expected hypothesis to match, got: %s", retrieved.Hypothesis)
	}

	if retrieved.Confidence != 0.8 {
		t.Errorf("expected confidence 0.8, got %f", retrieved.Confidence)
	}

	if len(retrieved.EvidenceIDs) != 2 || retrieved.EvidenceIDs[0] != "evidence-1" {
		t.Errorf("expected evidence IDs ['evidence-1', 'evidence-2'], got: %v", retrieved.EvidenceIDs)
	}
}

func TestUpdateLead(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	oldDB := DB
	DB = db
	defer func() { DB = oldDB }()

	if err := Migrate(); err != nil {
		t.Fatal(err)
	}

	c := &core.Case{ID: "case-1", Name: "Case 1"}
	if err := CreateCase(c); err != nil {
		t.Fatal(err)
	}

	lead := &core.IntelligenceLead{
		ID:          "lead-1",
		CaseID:      "case-1",
		Hypothesis:  "Vulnerability hypothesis",
		Confidence:  0.5,
		EvidenceIDs: []string{"evidence-1"},
		Status:      "active",
	}

	if err := CreateLead(lead); err != nil {
		t.Fatal(err)
	}

	// Update lead fields
	lead.Hypothesis = "Confirmed vulnerability"
	lead.Confidence = 0.95
	lead.Status = "verified"
	lead.EvidenceIDs = append(lead.EvidenceIDs, "evidence-3")

	if err := UpdateLead(lead); err != nil {
		t.Fatalf("UpdateLead failed: %v", err)
	}

	retrieved, err := GetLead("lead-1")
	if err != nil {
		t.Fatal(err)
	}

	if retrieved.Hypothesis != "Confirmed vulnerability" {
		t.Errorf("expected updated hypothesis, got: %s", retrieved.Hypothesis)
	}

	if retrieved.Confidence != 0.95 {
		t.Errorf("expected updated confidence, got %f", retrieved.Confidence)
	}

	if retrieved.Status != "verified" {
		t.Errorf("expected status 'verified', got: %s", retrieved.Status)
	}

	if len(retrieved.EvidenceIDs) != 2 || retrieved.EvidenceIDs[1] != "evidence-3" {
		t.Errorf("expected updated evidence IDs, got: %v", retrieved.EvidenceIDs)
	}
}

func TestListLeadsByCase(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	oldDB := DB
	DB = db
	defer func() { DB = oldDB }()

	if err := Migrate(); err != nil {
		t.Fatal(err)
	}

	c := &core.Case{ID: "case-1", Name: "Case 1"}
	if err := CreateCase(c); err != nil {
		t.Fatal(err)
	}

	lead1 := &core.IntelligenceLead{
		ID:         "lead-1",
		CaseID:     "case-1",
		Hypothesis: "Hypothesis 1",
	}
	lead2 := &core.IntelligenceLead{
		ID:         "lead-2",
		CaseID:     "case-1",
		Hypothesis: "Hypothesis 2",
	}

	if err := CreateLead(lead1); err != nil {
		t.Fatal(err)
	}
	if err := CreateLead(lead2); err != nil {
		t.Fatal(err)
	}

	leads, err := ListLeadsByCase("case-1")
	if err != nil {
		t.Fatalf("ListLeadsByCase failed: %v", err)
	}

	if len(leads) != 2 {
		t.Errorf("expected 2 leads, got %d", len(leads))
	}
}
