package libnf

import (
	libnf "libnf/api"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRecordInit(t *testing.T) {
	rec, err := libnf.NewRecord()
	assert.Equal(t, nil, err)
	err = rec.Free()
	assert.Equal(t, nil, err)
}
