package storage

import (
	"database/sql"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/spectre/spectre/internal/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetCaseTimeline(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)
	defer db.Close()

	oldDB := DB
	DB = db
	defer func() { DB = oldDB }()

	err = Migrate()
	require.NoError(t, err)

	caseID := "test-case-timeline"
	
	// Create an entity
	err = CreateEntity(&core.Entity{
		CaseID:       caseID,
		Type:         "domain",
		Value:        "example.com",
		DiscoveredAt: time.Now().Add(-1 * time.Hour),
		Source:       "manual",
	})
	require.NoError(t, err)

	// Create evidence
	err = CreateEvidence(&core.Evidence{
		CaseID:      caseID,
		Collector:   "dns",
		FilePath:    "/tmp/dns.json",
		CollectedAt: time.Now(),
	})
	require.NoError(t, err)

	timeline, err := GetCaseTimeline(caseID)
	assert.NoError(t, err)
	assert.Len(t, timeline, 2)
	assert.Equal(t, "entity_discovered", timeline[0].Type)
	assert.Equal(t, "evidence_collected", timeline[1].Type)
}
