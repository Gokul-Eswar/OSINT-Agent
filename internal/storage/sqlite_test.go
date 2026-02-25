package storage

import (
	"os"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitDB_Integration(t *testing.T) {
	// Setup temporary config
	viper.Set("database.path", "test_spectre.db")
	defer os.Remove("test_spectre.db")
	
	// Ensure we close and reset
	defer CloseDB()

	// Execute
	err := InitDB()
	require.NoError(t, err)
	require.NotNil(t, DB)

	// Verify we can ping
	err = DB.Ping()
	assert.NoError(t, err)

	// Verify tables were created (Migrate was called)
	var name string
	err = DB.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='cases'").Scan(&name)
	assert.NoError(t, err)
	assert.Equal(t, "cases", name)
}

func TestCloseDB(t *testing.T) {
	viper.Set("database.path", ":memory:")
	InitDB()
	
	err := CloseDB()
	assert.NoError(t, err)
}
