package ring_test

import (
	"libnf/api/ring"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateSharedFile(t *testing.T) {
	filename := "libnf-go"
	path := "/dev/shm/" + filename
	ring, err := ring.NewRing(filename, true, false, false)
	assert.Nil(t, err)

	_, err = os.Stat(path)
	assert.Nil(t, err)

	err = ring.Free()
	assert.Nil(t, err)

	_, err = os.Stat(path)
	assert.Equal(t, true, os.IsNotExist(err))
}
