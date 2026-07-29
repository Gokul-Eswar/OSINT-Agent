package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateCaseID(t *testing.T) {
	assert.NoError(t, ValidateCaseID("case-123_abc"))
	assert.Error(t, ValidateCaseID("../invalid"))
	assert.Error(t, ValidateCaseID("case/id"))
	assert.Error(t, ValidateCaseID(""))
}

func TestSanitizeFileName(t *testing.T) {
	assert.Equal(t, "user_repo", SanitizeFileName("user/repo"))
	assert.Equal(t, "file_name.json", SanitizeFileName("file name.json"))
	assert.Equal(t, "file", SanitizeFileName(""))
}
