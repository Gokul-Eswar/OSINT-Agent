package storage

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/spectre/spectre/internal/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEvidenceOperations(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)
	defer db.Close()

	oldDB := DB
	DB = db
	defer func() { DB = oldDB }()

	err = Migrate()
	require.NoError(t, err)

	caseID := "test-case-evidence"
	evidence := &core.Evidence{
		ID:        "ev-1",
		CaseID:    caseID,
		Collector: "test-collector",
		FilePath:  "/tmp/test.json",
		Metadata:  map[string]interface{}{"foo": "bar"},
	}

	t.Run("CreateEvidence", func(t *testing.T) {
		err := CreateEvidence(evidence)
		assert.NoError(t, err)
	})

	t.Run("ListEvidenceByCase", func(t *testing.T) {
		list, err := ListEvidenceByCase(caseID)
		assert.NoError(t, err)
		assert.Len(t, list, 1)
		assert.Equal(t, evidence.ID, list[0].ID)
		assert.Equal(t, "bar", list[0].Metadata["foo"])
	})
}
