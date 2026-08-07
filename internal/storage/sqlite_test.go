package storage

import (
	"os"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitDB_Integration(t *testing.T) {
	viper.Set("database.type", "sqlite")
	viper.Set("database.path", "test_spectre.db")

	t.Cleanup(func() {
		CloseDB()
		os.Remove("test_spectre.db")
	})

	err := InitDB()
	require.NoError(t, err)
	require.NotNil(t, DB)

	err = DB.Ping()
	assert.NoError(t, err)

	var name string
	err = DB.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='cases'").Scan(&name)
	assert.NoError(t, err)
	assert.Equal(t, "cases", name)
}

func TestCloseDB(t *testing.T) {
	viper.Set("database.type", "sqlite")
	viper.Set("database.path", ":memory:")
	InitDB()

	t.Cleanup(func() {
		CloseDB()
	})

	err := CloseDB()
	assert.NoError(t, err)
}
