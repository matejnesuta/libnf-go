package libnf

import (
	libnf "libnf/api"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Info from the file was retrieved using the nfdump version 1.6.25 tool
func TestGetInfo(t *testing.T) {
	var file libnf.File
	err := file.OpenRead("testfiles/nfcapd.201705281555", false, false)
	assert.Equal(t, nil, err)

	compressed, err := file.IsCompressed()
	assert.Equal(t, nil, err)
	assert.Equal(t, false, compressed)

	anonymized, err := file.IsAnonymized()
	assert.Equal(t, nil, err)
	assert.Equal(t, false, anonymized)

	blocks, err := file.GetBlocks()
	assert.Equal(t, nil, err)
	assert.Equal(t, uint64(2), blocks)

	version, err := file.GetLibnfVersion()
	assert.Equal(t, nil, err)
	assert.Equal(t, "1.33", version)

	flows, err := file.GetFlows()
	assert.Equal(t, nil, err)
	assert.Equal(t, uint64(2035), flows)

	bytes, err := file.GetBytes()
	assert.Equal(t, nil, err)
	assert.Equal(t, uint64(148619), bytes)

	packets, err := file.GetPackets()
	assert.Equal(t, nil, err)
	assert.Equal(t, uint64(2161), packets)

	assert.Equal(t, nil, file.Close())
}

func TestOpenFileMultipleTimes(t *testing.T) {
	var file libnf.File
	err := file.OpenRead("testfiles/nfcapd.201705281555", false, false)
	assert.Equal(t, nil, err)
	err = file.OpenRead("testfiles/nfcapd.201705281555", false, false)
	assert.Equal(t, libnf.ErrFileAlreadyOpened, err)
	err = file.OpenAppend("testfiles/nfcapd.201705281555", false)
	assert.Equal(t, libnf.ErrFileAlreadyOpened, err)
	err = file.OpenWrite("testfiles/nfcapd.201705281555", "testfiles/nfcapd.201705281555", false, 0, false)
	assert.Equal(t, libnf.ErrFileAlreadyOpened, err)
	assert.Equal(t, nil, file.Close())
	err = file.OpenAppend("testfiles/nfcapd.201705281555", false)
	assert.Equal(t, nil, err)
	assert.Equal(t, nil, file.Close())
}

func TestCloseFileMultipleTimes(t *testing.T) {
	var file libnf.File
	err := file.OpenRead("testfiles/nfcapd.201705281555", false, false)
	assert.Equal(t, nil, err)
	assert.Equal(t, nil, file.Close())
	err = file.Close()
	assert.Equal(t, libnf.ErrFileNotOpened, err)
}

func TestGetInfoFromUnopenedFile(t *testing.T) {
	var file libnf.File
	_, err := file.GetFlows()
	assert.Equal(t, libnf.ErrFileNotOpened, err)
	_, err = file.GetBytes()
	assert.Equal(t, libnf.ErrFileNotOpened, err)
	_, err = file.GetPackets()
	assert.Equal(t, libnf.ErrFileNotOpened, err)
	assert.Equal(t, libnf.ErrFileNotOpened, file.Close())
}
