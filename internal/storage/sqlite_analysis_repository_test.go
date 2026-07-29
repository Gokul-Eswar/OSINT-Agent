package storage

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/Gokul-Eswar/Spectre/internal/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAnalysisOperations(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)
	defer db.Close()

	oldDB := DB
	DB = db
	defer func() { DB = oldDB }()

	err = Migrate()
	require.NoError(t, err)

	caseID := "test-case-analysis"
	analysis := &core.Analysis{
		ID:          "analysis-1",
		CaseID:      caseID,
		ContextHash: "hash-123",
		Findings:    []string{"Finding 1"},
		Confidence:  0.85,
	}

	t.Run("SaveAnalysis", func(t *testing.T) {
		err := SaveAnalysis(analysis)
		assert.NoError(t, err)
	})

	t.Run("GetLatestAnalysis", func(t *testing.T) {
		a, err := GetLatestAnalysis(caseID)
		assert.NoError(t, err)
		require.NotNil(t, a)
		assert.Equal(t, analysis.ID, a.ID)
		assert.Equal(t, "Finding 1", a.Findings[0])
	})

	t.Run("GetAnalysisByHash", func(t *testing.T) {
		a, err := GetAnalysisByHash(caseID, "hash-123")
		assert.NoError(t, err)
		require.NotNil(t, a)
		assert.Equal(t, analysis.ID, a.ID)
	})
}
