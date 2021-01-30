package container

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildPG(t *testing.T) {
	b, err := NewConBuilder()
	assert.NoError(t, err)
	_, err = b.RunPg("test", "test_db")
	assert.NoError(t, err)
	_, err = b.FindContainer("test")
	assert.NoError(t, err)
	err = b.PruneAll()
	assert.NoError(t, err)
}
