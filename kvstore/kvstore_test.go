package kvstore

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStorage_Put(t *testing.T) {
	s := NewStorage()

	err := s.Put("key", "val")
	assert.NoError(t, err)

	v, err := s.Get("key")
	assert.NoError(t, err)
	assert.Equal(t, "val", v)

	assert.NoError(t, s.Delete("key"))
	v, err = s.Get("key")
	assert.ErrorIs(t, err, ErrNoSuchKey)
	assert.Equal(t, "", v)
}
