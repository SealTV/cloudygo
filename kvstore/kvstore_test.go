package kvstore

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStorage_Put(t *testing.T) {
	err := Put("key", "val")
	assert.NoError(t, err)

	v, err := Get("key")
	assert.NoError(t, err)
	assert.Equal(t, "val", v)

	assert.NoError(t, Delete("key"))
	v, err = Get("key")
	assert.ErrorIs(t, err, ErrNoSuchKey)
	assert.Equal(t, "", v)
}
