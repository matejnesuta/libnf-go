package file

import (
	"C"
	"libnf/api/errors"
	"libnf/internal"
	"unsafe"
)
import (
	"bytes"
	"libnf/api/record"
	"time"
)

// File struct represents an Nfdump. It provides methods to open, read, write, and retrieve metadata about the file.
type File struct {
	ptr    uintptr // Pointer to the underlying file structure.
	opened bool    // Indicates whether the file is currently opened.
}

// GetPtr returns the pointer to the file.
func (f *File) GetPtr() uintptr {
	return f.ptr
}

// Opened returns whether the file is currently opened.
func (f *File) Opened() bool {
	return f.opened
}

// A list of possible compression types, which can be used with the OpenWrite method.
const (
	NoComp  int = 0                 // No compression.
	CompLZO int = internal.COMP_LZO // LZO compression.
	CompBZ2 int = internal.COMP_BZ2 // BZ2 compression.
)

func getInfo(info int, f *File) ([]byte, error) {
	if !f.opened {
		return nil, errors.ErrFileNotOpened
	}
	buf := make([]byte, internal.INFO_BUFSIZE)
	data := uintptr(unsafe.Pointer(&buf[0]))

	status := internal.Info(f.ptr, info, data, int64(len(buf)))

	switch status {
	case internal.ERR_NOMEM:
		return nil, errors.ErrNoMem
	case internal.ERR_OTHER:
		return nil, errors.ErrOther
	default:
		return buf, nil
	}
}

func getStringInfo(info int, f *File) (string, error) {
	buf, err := getInfo(info, f)
	if err != nil {
		return "", err
	}
	// Look for null terminator
	n := bytes.IndexByte(buf, 0)
	if n < 0 {
		n = len(buf)
	}
	return string(buf[:n]), nil
}

func getBoolInfo(info int, f *File) (bool, error) {
	buf, err := getInfo(info, f)
	if err != nil {
		return false, err
	}
	val := *(*int)(unsafe.Pointer(&buf[0]))
	return val != 0, nil
}

func getUint64Info(info int, f *File) (uint64, error) {
	buf, err := getInfo(info, f)
	if err != nil {
		return 0, err
	}
	return *(*uint64)(unsafe.Pointer(&buf[0])), nil
}

func getTimestamp(info int, f *File) (time.Time, error) {
	data, err := getUint64Info(info, f)
	if err != nil {
		return time.Time{}, err
	}
	return time.Unix(0, int64(time.Millisecond)*int64(data)), nil
}

// GetLibnfVersion returns the version of the C libnf library, which is used under the hood.
func (f *File) GetLibnfVersion() (string, error) {
	return getStringInfo(internal.INFO_VERSION, f)
}

// GetNfdumpVersion returns the version of Nfdump used in the file.
func (f *File) GetNfdumpVersion() (string, error) {
	return getStringInfo(internal.INFO_NFDUMP_VERSION, f)
}

// GetIdent returns the string identification of the file.
func (f *File) GetIdent() (string, error) {
	return getStringInfo(internal.INFO_IDENT, f)
}

// IsCompressed returns whether the file is compressed with LZO or BZ2.
func (f *File) IsCompressed() (bool, error) {
	return getBoolInfo(internal.INFO_COMPRESSED, f)
}

// IsAnonymized returns whether IP adresses in the file are anonymized or not..
func (f *File) IsAnonymized() (bool, error) {
	return getBoolInfo(internal.INFO_ANONYMIZED, f)
}

// HasCatalog return s whether the file has a catalog.
func (f *File) HasCatalog() (bool, error) {
	return getBoolInfo(internal.INFO_CATALOG, f)
}

// Various methods for retrieving uint64 metadata about the file.
func (f *File) GetFileVersion() (uint64, error) { return getUint64Info(internal.INFO_FILE_VERSION, f) }

// Get total number of blocks in the file.
func (f *File) GetBlocks() (uint64, error) { return getUint64Info(internal.INFO_BLOCKS, f) }

// Get timestamp of the first packet in the file.
func (f *File) GetFirst() (time.Time, error) { return getTimestamp(internal.INFO_FIRST, f) }

// Get timestamp of the last packet in the file.
func (f *File) GetLast() (time.Time, error) { return getTimestamp(internal.INFO_LAST, f) }

// Get total number of sequence failures in the file.
func (f *File) GetFailures() (uint64, error) { return getUint64Info(internal.INFO_FAILURES, f) }

// Get total number of flows in the file.
func (f *File) GetFlows() (uint64, error) { return getUint64Info(internal.INFO_FLOWS, f) }

// Get total number of bytes in the file.
func (f *File) GetBytes() (uint64, error) { return getUint64Info(internal.INFO_BYTES, f) }

// Get total number of packets in the file.
func (f *File) GetPackets() (uint64, error) { return getUint64Info(internal.INFO_PACKETS, f) }

// Get total number of processed blocks from the file.
func (f *File) GetProcBlocks() (uint64, error) { return getUint64Info(internal.INFO_PROC_BLOCKS, f) }

// Get total number of TCP flows from the file.
func (f *File) GetFlowsTcp() (uint64, error) { return getUint64Info(internal.INFO_FLOWS_TCP, f) }

// Get total number of UDP flows from the file.
func (f *File) GetFlowsUdp() (uint64, error) { return getUint64Info(internal.INFO_FLOWS_UDP, f) }

// Get total number of ICMP flows from the file.
func (f *File) GetFlowsIcmp() (uint64, error) { return getUint64Info(internal.INFO_FLOWS_ICMP, f) }

// Get total number of other flows from the file.
func (f *File) GetFlowsOther() (uint64, error) { return getUint64Info(internal.INFO_FLOWS_OTHER, f) }

// Get total number of TCP bytes from the file.
func (f *File) GetBytesTcp() (uint64, error) { return getUint64Info(internal.INFO_BYTES_TCP, f) }

// Get total number of UDP bytes from the file.
func (f *File) GetBytesUdp() (uint64, error) { return getUint64Info(internal.INFO_BYTES_UDP, f) }

// Get total number of ICMP bytes from the file.
func (f *File) GetBytesIcmp() (uint64, error) { return getUint64Info(internal.INFO_BYTES_ICMP, f) }

// Get total number of other bytes from the file.
func (f *File) GetBytesOther() (uint64, error) { return getUint64Info(internal.INFO_BYTES_OTHER, f) }

// Get total number of TCP packets from the file.
func (f *File) GetPacketsTcp() (uint64, error) { return getUint64Info(internal.INFO_PACKETS_TCP, f) }

// Get total number of UDP packets from the file.
func (f *File) GetPacketsUdp() (uint64, error) { return getUint64Info(internal.INFO_PACKETS_UDP, f) }

// Get total number of ICMP packets from the file.
func (f *File) GetPacketsIcmp() (uint64, error) { return getUint64Info(internal.INFO_PACKETS_ICMP, f) }

// Get total number of other packets from the file.
func (f *File) GetPacketsOther() (uint64, error) {
	return getUint64Info(internal.INFO_PACKETS_OTHER, f)
}

// GetCompressionType returns the compression type of the file.
func (f *File) GetCompressionType() (int, error) {
	data, err := getBoolInfo(internal.INFO_COMPRESSED, f)
	if err != nil {
		return NoComp, err
	} else if !data {
		return NoComp, nil
	}

	data, err = getBoolInfo(internal.INFO_LZO_COMPRESSED, f)
	if err != nil {
		return NoComp, err
	} else if data {
		return CompLZO, nil
	} else {
		return CompBZ2, nil
	}
}

// OpenRead opens the file in read mode.
func (f *File) OpenRead(inputFile string, readLoop bool, weakErr bool) error {
	if f.opened {
		return errors.ErrFileAlreadyOpened
	}

	flags := internal.READ
	if readLoop {
		flags |= internal.READ_LOOP
	}
	if weakErr {
		flags |= internal.WEAKERR
	}
	status := internal.Open(&f.ptr, inputFile, uint(flags), "")
	if status != internal.OK {
		return errors.ErrCannotOpenFile
	}
	f.opened = true
	return nil
}

// OpenAppend opens the file in append mode.
func (f *File) OpenAppend(inputFile string, weakErr bool) error {
	if f.opened {
		return errors.ErrFileAlreadyOpened
	}

	flags := internal.APPEND
	if weakErr {
		flags |= internal.WEAKERR
	}
	status := internal.Open(&f.ptr, inputFile, uint(flags), "")
	if status != internal.OK {
		return errors.ErrCannotOpenFile
	}
	f.opened = true
	return nil
}

// OpenWrite opens the file in write mode.
func (f *File) OpenWrite(outputFile string, ident string, anon bool, comp int, weakErr bool) error {
	if f.opened {
		return errors.ErrFileAlreadyOpened
	}

	flags := internal.WRITE
	if anon {
		flags |= internal.ANON
	}
	if comp == CompLZO {
		flags |= internal.COMP_LZO
	}
	if comp == CompBZ2 {
		flags |= internal.COMP_BZ2
	}
	if weakErr {
		flags |= internal.WEAKERR
	}
	status := internal.Open(&f.ptr, outputFile, uint(flags), ident)
	if status != internal.OK {
		return errors.ErrCannotOpenFile
	}
	f.opened = true
	return nil
}

// Close closes the file and frees resources.
func (file *File) Close() error {
	if file.opened {
		internal.Close(file.ptr)
		file.opened = false
		return nil
	}
	return errors.ErrFileNotOpened
}

// GetNextRecord reads the next record from the file.
func (file *File) GetNextRecord(r *record.Record) error {
	if !file.opened {
		return errors.ErrFileNotOpened
	} else if !r.Allocated() {
		return errors.ErrRecordNotAllocated
	}
	status := internal.Read(file.ptr, r.GetPtr())
	if status == internal.ERR_NOMEM {
		return errors.ErrNoMem
	} else if status == internal.EOF {
		return errors.ErrFileEof
	}
	return nil
}

// WriteRecord writes a record to the file.
func (file *File) WriteRecord(r *record.Record) error {
	if !file.opened {
		return errors.ErrFileNotOpened
	} else if !r.Allocated() {
		return errors.ErrRecordNotAllocated
	}
	status := internal.Write(file.ptr, r.GetPtr())
	if status == internal.ERR_NOMEM {
		return errors.ErrNoMem
	} else if status == internal.ERR_WRITE {
		return errors.ErrWrite
	}
	return nil
}
