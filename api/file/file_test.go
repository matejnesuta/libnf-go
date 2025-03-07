package file_test

import (
	LnfErr "libnf/api/errors"
	LnfFile "libnf/api/file"
	LnfRec "libnf/api/record"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Info from the file was retrieved using the nfdump version 1.6.25 tool
func TestGetInfo(t *testing.T) {
	var file LnfFile.File
	err := file.OpenRead("../testfiles/nfcapd.201705281555", false, false)
	assert.Equal(t, nil, err)
	assert.Equal(t, true, file.Opened())

	nfdumpVersion, err := file.GetNfdumpVersion()
	assert.Equal(t, nil, err)
	assert.Equal(t, "nfdump, 1.6.25, peter@people.ops-trust.net", nfdumpVersion)

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

	flowsUdp, err := file.GetFlowsUdp()
	assert.Equal(t, nil, err)
	assert.Equal(t, uint64(1), flowsUdp)

	flowsTcp, err := file.GetFlowsTcp()
	assert.Equal(t, nil, err)
	assert.Equal(t, uint64(2034), flowsTcp)

	assert.Equal(t, nil, file.Close())
	assert.Equal(t, false, file.Opened())
}

func TestOpenFileMultipleTimes(t *testing.T) {
	var file LnfFile.File
	err := file.OpenRead("../testfiles/nfcapd.201705281555", false, false)
	assert.Equal(t, nil, err)
	assert.Equal(t, true, file.Opened())
	err = file.OpenRead("../testfiles/nfcapd.201705281555", false, false)
	assert.Equal(t, LnfErr.ErrFileAlreadyOpened, err)
	assert.Equal(t, true, file.Opened())
	err = file.OpenAppend("../testfiles/nfcapd.201705281555", false)
	assert.Equal(t, LnfErr.ErrFileAlreadyOpened, err)
	assert.Equal(t, true, file.Opened())
	err = file.OpenWrite("../testfiles/nfcapd.201705281555", "testfiles/nfcapd.201705281555", false, 0, false)
	assert.Equal(t, LnfErr.ErrFileAlreadyOpened, err)
	assert.Equal(t, true, file.Opened())
	assert.Equal(t, nil, file.Close())
	assert.Equal(t, false, file.Opened())
	err = file.OpenAppend("../testfiles/nfcapd.201705281555", false)
	assert.Equal(t, true, file.Opened())
	assert.Equal(t, nil, err)
	assert.Equal(t, nil, file.Close())
	assert.Equal(t, false, file.Opened())
}

func TestCloseFileMultipleTimes(t *testing.T) {
	var file LnfFile.File
	assert.Equal(t, false, file.Opened())
	err := file.OpenRead("../testfiles/nfcapd.201705281555", false, false)
	assert.Equal(t, true, file.Opened())
	assert.Equal(t, nil, err)
	assert.Equal(t, nil, file.Close())
	assert.Equal(t, false, file.Opened())
	err = file.Close()
	assert.Equal(t, LnfErr.ErrFileNotOpened, err)
	assert.Equal(t, false, file.Opened())
}

func TestGetInfoFromUnopenedFile(t *testing.T) {
	var file LnfFile.File
	_, err := file.GetFlows()
	assert.Equal(t, LnfErr.ErrFileNotOpened, err)
	_, err = file.GetBytes()
	assert.Equal(t, LnfErr.ErrFileNotOpened, err)
	_, err = file.GetPackets()
	assert.Equal(t, LnfErr.ErrFileNotOpened, err)
	assert.Equal(t, LnfErr.ErrFileNotOpened, file.Close())
}

func TestReadRecordFromOpenedFile(t *testing.T) {
	var file LnfFile.File
	err := file.OpenRead("../testfiles/test-file.tmp", false, false)
	assert.Equal(t, nil, err)
	rec, _ := LnfRec.NewRecord()

	err = file.GetNextRecord(&rec)
	assert.Equal(t, nil, err)
}

func TestReadRecordFromUnopenedFile(t *testing.T) {
	var file LnfFile.File
	rec, _ := LnfRec.NewRecord()
	err := file.GetNextRecord(&rec)
	assert.Equal(t, LnfErr.ErrFileNotOpened, err)
}

func TestReadRecordFromClosedFile(t *testing.T) {
	var file LnfFile.File
	err := file.OpenRead("../testfiles/test-file.tmp", false, false)
	assert.Equal(t, nil, err)
	rec, _ := LnfRec.NewRecord()
	assert.Equal(t, nil, file.Close())
	err = file.GetNextRecord(&rec)
	assert.Equal(t, LnfErr.ErrFileNotOpened, err)
}

func TestReadRecordFromUnallocatedRecord(t *testing.T) {
	var file LnfFile.File
	err := file.OpenRead("../testfiles/test-file.tmp", false, false)
	assert.Equal(t, nil, err)
	var rec LnfRec.Record
	err = file.GetNextRecord(&rec)
	assert.Equal(t, LnfErr.ErrRecordNotAllocated, err)
}

func TestReadUntilEOF(t *testing.T) {
	var file LnfFile.File
	err := file.OpenRead("../testfiles/test-file.tmp", false, false)
	assert.Equal(t, nil, err)
	rec, _ := LnfRec.NewRecord()
	for {
		err = file.GetNextRecord(&rec)
		if err != nil {
			break
		}
	}
	assert.Equal(t, err, LnfErr.ErrFileEof)
}

func TestOpenNonexistentFile(t *testing.T) {
	var file LnfFile.File
	err := file.OpenRead("nonexistent-file.tmp", false, false)
	assert.Equal(t, LnfErr.ErrCannotOpenFile, err)
	assert.Equal(t, false, file.Opened())
}

func TestCreateFile(t *testing.T) {
	var file LnfFile.File
	err := file.OpenWrite("../tmp/test-file.tmp", "", false, 0, false)
	assert.Equal(t, nil, err)
	assert.Equal(t, true, file.Opened())
	assert.Equal(t, nil, file.Close())
	assert.Equal(t, false, file.Opened())
	err = file.OpenRead("../tmp/test-file.tmp", false, false)
	assert.Equal(t, nil, err)
	assert.Equal(t, true, file.Opened())
	compressed, err := file.IsCompressed()
	assert.Equal(t, nil, err)
	assert.Equal(t, false, compressed)
	anon, err := file.IsAnonymized()
	assert.Equal(t, nil, err)
	assert.Equal(t, false, anon)
	assert.Equal(t, nil, file.Close())
	assert.Equal(t, false, file.Opened())
}

func TestCreateFileWithIdent(t *testing.T) {
	var file LnfFile.File
	err := file.OpenWrite("../tmp/test-file.tmp", "test-file.tmp", false, 0, false)
	assert.Equal(t, nil, err)
	assert.Equal(t, true, file.Opened())
	assert.Equal(t, nil, file.Close())
	assert.Equal(t, false, file.Opened())
	err = file.OpenRead("../tmp/test-file.tmp", false, false)
	assert.Equal(t, nil, err)
	assert.Equal(t, true, file.Opened())
	ident, err := file.GetIdent()
	assert.Equal(t, nil, err)
	assert.Equal(t, "test-file.tmp", ident)
	assert.Equal(t, nil, file.Close())
	assert.Equal(t, false, file.Opened())
}

func TestCreateFileWithLZOCompression(t *testing.T) {
	var file LnfFile.File
	err := file.OpenWrite("../tmp/lzo-file.tmp", "", false, LnfFile.CompLZO, false)
	assert.Equal(t, nil, err)
	assert.Equal(t, true, file.Opened())
	assert.Equal(t, nil, file.Close())
	assert.Equal(t, false, file.Opened())
	err = file.OpenRead("../tmp/lzo-file.tmp", false, false)
	assert.Equal(t, nil, err)
	assert.Equal(t, true, file.Opened())
	compressed, err := file.IsCompressed()
	assert.Equal(t, nil, err)
	assert.Equal(t, true, compressed)
	assert.Equal(t, nil, file.Close())
	assert.Equal(t, false, file.Opened())
}

func TestCreateFileWithBZ2Compression(t *testing.T) {
	var file LnfFile.File
	err := file.OpenWrite("../tmp/bz2-file.tmp", "", false, LnfFile.CompBZ2, false)
	assert.Equal(t, nil, err)
	assert.Equal(t, true, file.Opened())
	assert.Equal(t, nil, file.Close())
	assert.Equal(t, false, file.Opened())
	err = file.OpenRead("../tmp/bz2-file.tmp", false, false)
	assert.Equal(t, nil, err)
	assert.Equal(t, true, file.Opened())
	compressed, err := file.IsCompressed()
	assert.Equal(t, nil, err)
	assert.Equal(t, true, compressed)
	assert.Equal(t, nil, file.Close())
	assert.Equal(t, false, file.Opened())
}

func TestCreateFileWithWeakErr(t *testing.T) {
	var file LnfFile.File
	err := file.OpenWrite("../tmp/weakerr-file.tmp", "", false, 0, true)
	assert.Equal(t, nil, err)
	assert.Equal(t, true, file.Opened())
	assert.Equal(t, nil, file.Close())
	assert.Equal(t, false, file.Opened())
	err = file.OpenRead("../tmp/weakerr-file.tmp", false, true)
	assert.Equal(t, nil, err)
	assert.Equal(t, true, file.Opened())
	assert.Equal(t, nil, file.Close())
	assert.Equal(t, false, file.Opened())
}

// func TestCreateAnonymizedFile(t *testing.T) {
// 	var file LnfFile.File
// 	err := file.OpenWrite("../tmp/anon-file.tmp", "", true, 0, false)
// 	assert.Equal(t, nil, err)
// 	assert.Equal(t, true, file.Opened())
// 	assert.Equal(t, nil, file.Close())
// 	assert.Equal(t, false, file.Opened())
// 	err = file.OpenRead("../tmp/anon-file.tmp", true, false)
// 	assert.Equal(t, nil, err)
// 	assert.Equal(t, true, file.Opened())
// 	anon, err := file.IsAnonymized()
// 	assert.Equal(t, nil, err)
// 	assert.Equal(t, true, anon)
// 	assert.Equal(t, nil, file.Close())
// 	assert.Equal(t, false, file.Opened())
// }

func TestWriteRecordToOpenedFile(t *testing.T) {
	var outputFile LnfFile.File
	err := outputFile.OpenWrite("../tmp/file-with-record.tmp", "", false, 0, false)
	assert.Equal(t, nil, err)
	assert.Equal(t, true, outputFile.Opened())

	var inputFile LnfFile.File
	err = inputFile.OpenRead("../testfiles/nfcapd.201705281555", false, false)
	assert.Equal(t, nil, err)
	assert.Equal(t, true, inputFile.Opened())
	rec, _ := LnfRec.NewRecord()
	err = inputFile.GetNextRecord(&rec)
	assert.Equal(t, nil, err)
	err = outputFile.WriteRecord(&rec)
	assert.Equal(t, nil, err)
	assert.Equal(t, nil, inputFile.Close())
	assert.Equal(t, false, inputFile.Opened())
	assert.Equal(t, nil, outputFile.Close())
	assert.Equal(t, false, outputFile.Opened())

	outputFile.OpenRead("../tmp/file-with-record.tmp", false, false)
	assert.Equal(t, true, outputFile.Opened())
	flows, err := outputFile.GetFlows()
	assert.Equal(t, uint64(1), flows)
	assert.Equal(t, nil, err)
	rec, _ = LnfRec.NewRecord()
	err = outputFile.GetNextRecord(&rec)

	assert.Equal(t, nil, err)
	assert.Equal(t, nil, outputFile.Close())
	assert.Equal(t, false, outputFile.Opened())

}

func TestOpenAppendNonexistentFile(t *testing.T) {
	var file LnfFile.File
	err := file.OpenAppend("nonexistent-file.tmp", false)
	assert.Equal(t, LnfErr.ErrCannotOpenFile, err)
	assert.Equal(t, false, file.Opened())
}

// func TestOpenAppendFile(t *testing.T) {
// 	var outputFile LnfFile.File
// 	err := outputFile.OpenWrite("../tmp/appended-file.tmp", "", false, 0, false)
// 	assert.Equal(t, nil, err)
// 	assert.Equal(t, true, outputFile.Opened())

// 	var inputFile LnfFile.File
// 	err = inputFile.OpenRead("../testfiles/nfcapd.201705281555", false, false)
// 	assert.Equal(t, nil, err)
// 	assert.Equal(t, true, inputFile.Opened())
// 	rec1, _ := LnfRec.NewRecord()
// 	rec2, _ := LnfRec.NewRecord()
// 	err = inputFile.GetNextRecord(&rec1)
// 	assert.Equal(t, nil, err)
// 	err = outputFile.WriteRecord(&rec1)
// 	assert.Equal(t, nil, err)
// 	assert.Equal(t, nil, inputFile.Close())
// 	assert.Equal(t, false, inputFile.Opened())
// 	assert.Equal(t, nil, outputFile.Close())
// 	assert.Equal(t, false, outputFile.Opened())

// 	outputFile.OpenAppend("../tmp/appended-file.tmp", false)
// 	assert.Equal(t, true, outputFile.Opened())
// 	err = outputFile.WriteRecord(&rec2)
// 	// flows, err := outputFile.GetFlows()
// 	// assert.Equal(t, uint64(1), flows)
// 	assert.Equal(t, nil, err)
// 	flows, err := outputFile.GetFlows()
// 	assert.Equal(t, uint64(2), flows)
// 	assert.Equal(t, nil, outputFile.Close())
// 	// assert.Equal(t, false, outputFile.Opened())
// rec, _ = LnfRec.NewRecord()

// outputFile.OpenRead("../tmp/appended-file.tmp", false, false)
// assert.Equal(t, true, outputFile.Opened())
// flows, err := outputFile.GetFlows()
// assert.Equal(t, nil, err)
// rec, _ = LnfRec.NewRecord()
// err = outputFile.GetNextRecord(&rec)

// }
