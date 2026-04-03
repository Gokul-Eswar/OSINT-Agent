package core

import (
	"encoding/json"
	"testing"
	"time"
)

type mockCollector struct{}

func (m mockCollector) Name() string        { return "mock" }
func (m mockCollector) Description() string { return "mock collector" }
func (m mockCollector) Collect(caseID string, target string, options map[string]interface{}) ([]Evidence, error) {
	return []Evidence{{CaseID: caseID, Collector: "mock"}}, nil
}
func (m mockCollector) IsActive() bool { return false }

func TestCollectorInterface_Implemented(t *testing.T) {
	var c Collector = mockCollector{}
	if c.Name() != "mock" {
		t.Fatalf("expected collector name mock, got %q", c.Name())
	}
}

func TestEvidenceJSON_OmitsRawData(t *testing.T) {
	e := Evidence{
		ID:          "ev-1",
		CaseID:      "case-1",
		Collector:   "dns",
		CollectedAt: time.Now(),
		RawData:     map[string]string{"secret": "hidden"},
	}

	b, err := json.Marshal(e)
	if err != nil {
		t.Fatalf("failed to marshal evidence: %v", err)
	}

	if string(b) == "" {
		t.Fatal("expected non-empty JSON")
	}

	var payload map[string]interface{}
	if err := json.Unmarshal(b, &payload); err != nil {
		t.Fatalf("failed to unmarshal evidence JSON: %v", err)
	}

	if _, ok := payload["raw_data"]; ok {
		t.Fatal("raw_data should not be serialized")
	}
}

func TestCoreModelsJSONTags(t *testing.T) {
	entity := Entity{ID: "ent-1", Type: "domain", Value: "example.com"}
	rel := Relationship{ID: "rel-1", FromEntityID: "ent-1", ToEntityID: "ent-2", Type: "resolves_to"}
	timeline := TimelineEvent{Type: "entity_discovered", Description: "domain observed"}
	analysis := Analysis{ID: "an-1", CaseID: "case-1", Findings: []string{"f1"}}
	caseModel := Case{ID: "case-1", Name: "Test"}

	fixtures := []interface{}{entity, rel, timeline, analysis, caseModel}
	for i, fixture := range fixtures {
		b, err := json.Marshal(fixture)
		if err != nil {
			t.Fatalf("fixture %d failed to marshal: %v", i, err)
		}
		if len(b) == 0 {
			t.Fatalf("fixture %d marshaled to empty payload", i)
		}
	}
}
