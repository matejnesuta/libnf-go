package ring_test

import (
	"os"
	"testing"

	"github.com/matejnesuta/libnf-go/api/errors"
	"github.com/matejnesuta/libnf-go/api/fields"
	"github.com/matejnesuta/libnf-go/api/record"
	"github.com/matejnesuta/libnf-go/api/ring"

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

func TestGetNextRecordOnUnallocatedRecord(t *testing.T) {
	ring, err := ring.NewRing("libnf-go", true, false, false)
	rec := record.Record{}
	assert.Nil(t, err)

	err = ring.GetNextRecord(&rec)
	assert.Equal(t, errors.ErrRecordNotAllocated, err)

	err = ring.Free()
	assert.Nil(t, err)
}

func TestWriteOnUnallocatedRecord(t *testing.T) {
	ring, err := ring.NewRing("libnf-go", true, false, false)
	rec := record.Record{}
	assert.Nil(t, err)

	err = ring.WriteRecord(&rec)
	assert.Equal(t, errors.ErrRecordNotAllocated, err)

	err = ring.Free()
	assert.Nil(t, err)
}

func TestHappyPath(t *testing.T) {
	r, err := ring.NewRing("libnf-go", true, false, false)
	assert.Nil(t, err)
	rec, err := record.NewRecord()
	assert.Nil(t, err)

	for i := 0; i < 10; i++ {
		record.SetField(&rec, fields.Dpkts, uint64(i))
		err = r.WriteRecord(&rec)
		assert.Nil(t, err)
	}

	_, err = r.Info(ring.RingTotal)
	assert.Nil(t, err)

	for i := 0; i < 10; i++ {
		err = r.GetNextRecord(&rec)
		assert.Nil(t, err)
		val, _ := rec.GetField(fields.Dpkts)
		assert.Equal(t, uint64(i), val)
	}
}
