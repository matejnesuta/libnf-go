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

func TestInitFilterTwice(t *testing.T) {
	var filter libnf.Filter
	err := filter.Init("src port 80")
	assert.Equal(t, nil, err)
	err = filter.Init("src port 80")
	assert.Equal(t, libnf.ErrFilterAlreadyInit, err)
	err = filter.Free()
	assert.Equal(t, nil, err)
}

func TestMatchUninitializedFilter(t *testing.T) {
	var filter libnf.Filter
	record, err := libnf.NewRecord()
	assert.Equal(t, nil, err)
	_, err = filter.Match(record)
	assert.Equal(t, libnf.ErrFilterNotInit, err)
	err = record.Free()
	assert.Equal(t, nil, err)
}

func TestMatchUninitializedRecord(t *testing.T) {
	var filter libnf.Filter
	err := filter.Init("src port 80")
	assert.Equal(t, nil, err)
	_, err = filter.Match(libnf.Record{})
	assert.Equal(t, libnf.ErrRecordNotAllocated, err)
	err = filter.Free()
	assert.Equal(t, nil, err)
}

func TestMatchClearFilter(t *testing.T) {
	var filter libnf.Filter
	record, err := libnf.NewRecord()
	assert.Equal(t, nil, err)
	err = filter.Init("src port 80")
	assert.Equal(t, nil, err)
	err = filter.Free()
	assert.Equal(t, nil, err)
	_, err = filter.Match(record)
	assert.Equal(t, libnf.ErrFilterNotInit, err)
	err = record.Free()
	assert.Equal(t, nil, err)
}

func TestDoubleClearFilter(t *testing.T) {
	var filter libnf.Filter
	err := filter.Init("src port 80")
	assert.Equal(t, nil, err)
	err = filter.Free()
	assert.Equal(t, nil, err)
	err = filter.Free()
	assert.Equal(t, libnf.ErrFilterNotInit, err)
}

func TestInitFilterWithWrongValue(t *testing.T) {
	var filter libnf.Filter
	err := filter.Init("uhhhhhhhhhhhhhhhhhhhhhhhhh")
	assert.Equal(t, libnf.ErrOtherMsg, err)
}

func TestMatchFilter(t *testing.T) {
	var file libnf.File
	err := file.OpenRead("testfiles/nfcapd.201705281555", false, false)
	assert.Equal(t, nil, err)

	var filter libnf.Filter
	err = filter.Init("src port 80")
	assert.Equal(t, nil, err)

	rec, err := libnf.NewRecord()
	assert.Equal(t, nil, err)

	num_of_matches := 0

	for {
		err = file.GetNextRecord(&rec)
		if err != nil {
			break
		}
		if match, _ := filter.Match(rec); match {
			num_of_matches++
		}
	}

	assert.Equal(t, int(4), num_of_matches)
}
