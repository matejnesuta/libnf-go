package libnf

import (
	libnf "libnf/api"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitFilter(t *testing.T) {
	var filter libnf.Filter
	err := filter.Init("src port 80")
	assert.Equal(t, nil, err)
	err = filter.Free()
	assert.Equal(t, nil, err)
}
