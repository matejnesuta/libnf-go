package memheap_test

import (
	"libnf/api/memheap"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMemHeap(t *testing.T) {
	memHeap, err := memheap.NewMemHeap()
	assert.Equal(t, nil, err)
	assert.Equal(t, true, memHeap.Allocated())

	// if err != nil {
	// 	t.Errorf("NewMemHeap() failed: %v", err)
	// }
	err = memHeap.Free()
	assert.Equal(t, nil, err)
	assert.Equal(t, false, memHeap.Allocated())
	// if err != nil {
	// 	t.Errorf("MemHeap.Free() failed: %v", err)
	// }
}
