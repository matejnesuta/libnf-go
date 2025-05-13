package filter_test

import (
	LnfErr "libnf-go/api/errors"
	LnfFile "libnf-go/api/file"
	LnfFilter "libnf-go/api/filter"
	LnfRec "libnf-go/api/record"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitFilter(t *testing.T) {
	var filter LnfFilter.Filter
	err := filter.Init("src port 80")
	assert.Equal(t, nil, err)
	err = filter.Free()
	assert.Equal(t, nil, err)
}

func TestInitFilterTwice(t *testing.T) {
	var filter LnfFilter.Filter
	err := filter.Init("src port 80")
	assert.Equal(t, nil, err)
	err = filter.Init("src port 80")
	assert.Equal(t, LnfErr.ErrFilterAlreadyInit, err)
	err = filter.Free()
	assert.Equal(t, nil, err)
}

func TestMatchUninitializedFilter(t *testing.T) {
	var filter LnfFilter.Filter
	record, err := LnfRec.NewRecord()
	assert.Equal(t, nil, err)
	_, err = filter.Match(record)
	assert.Equal(t, LnfErr.ErrFilterNotInit, err)
	err = record.Free()
	assert.Equal(t, nil, err)
}

func TestMatchUninitializedRecord(t *testing.T) {
	var filter LnfFilter.Filter
	err := filter.Init("src port 80")
	assert.Equal(t, nil, err)
	_, err = filter.Match(LnfRec.Record{})
	assert.Equal(t, LnfErr.ErrRecordNotAllocated, err)
	err = filter.Free()
	assert.Equal(t, nil, err)
}

func TestMatchClearFilter(t *testing.T) {
	var filter LnfFilter.Filter
	record, err := LnfRec.NewRecord()
	assert.Equal(t, nil, err)
	err = filter.Init("src port 80")
	assert.Equal(t, nil, err)
	err = filter.Free()
	assert.Equal(t, nil, err)
	assert.Equal(t, "", filter.String())
	_, err = filter.Match(record)
	assert.Equal(t, LnfErr.ErrFilterNotInit, err)
	err = record.Free()
	assert.Equal(t, nil, err)
}

func TestDoubleClearFilter(t *testing.T) {
	var filter LnfFilter.Filter
	err := filter.Init("src port 80")
	assert.Equal(t, nil, err)
	err = filter.Free()
	assert.Equal(t, nil, err)
	err = filter.Free()
	assert.Equal(t, LnfErr.ErrFilterNotInit, err)
}

func TestInitFilterWithWrongValue(t *testing.T) {
	var filter LnfFilter.Filter
	err := filter.Init("uhhhhhhhhhhhhhhhhhhhhhhhhh")
	assert.Equal(t, LnfErr.ErrOtherMsg, err)
}

func TestMatchFilter(t *testing.T) {
	var file LnfFile.File
	err := file.OpenRead("../testfiles/nfcapd.201705281555", false, false)
	assert.Equal(t, nil, err)

	var filter LnfFilter.Filter
	err = filter.Init("src port 80")
	assert.Equal(t, nil, err)

	rec, err := LnfRec.NewRecord()
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
