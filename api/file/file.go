package file

import (
	"C"
	LnfErr "libnf/api/errors"
	LnfRec "libnf/api/record"
	"libnf/internal"
	"unsafe"
)

// File represents a network flow file.
type File struct {
	ptr    uintptr
	opened bool
}

// GetPtr returns the pointer to the file.
func (f *File) GetPtr() uintptr {
	return f.ptr
}

// Opened returns whether the file is opened.
func (f *File) Opened() bool {
	return f.opened
}

// A list of possible compression methods, which can be used with OpenWrite method.
const (
	NoComp  int = 0
	CompLZO int = internal.COMP_LZO
	CompBZ2 int = internal.COMP_BZ2
)

func getInfo(info int, f *File) (uintptr, error) {
	if !(*f).opened {
		return 0, LnfErr.ErrFileNotOpened
	}
	buf := make([]byte, internal.INFO_BUFSIZE) // Allocate memory
	data := uintptr(unsafe.Pointer(&buf[0]))
	status := internal.Info((*f).ptr, info, data, int64(internal.INFO_BUFSIZE))
	if status == internal.ERR_NOMEM {
		return data, LnfErr.ErrNoMem
	}
	if status == internal.ERR_OTHER {
		return data, LnfErr.ErrOther
	}
	return data, nil
}

func getStringInfo(info int, f *File) (string, error) {
	data, err := getInfo(info, f)
	if err != nil {
		return "", err
	}
	return C.GoString((*C.char)(unsafe.Pointer(data))), nil
}

func getBoolInfo(info int, f *File) (bool, error) {
	data, err := getInfo(info, f)
	if err != nil {
		return false, err
	}
	value := *(*int)(unsafe.Pointer(data))
	return value != 0, nil
}

func getUint64Info(info int, f *File) (uint64, error) {
	data, err := getInfo(info, f)
	if err != nil {
		return 0, err
	}
	return *(*uint64)(unsafe.Pointer(data)), nil
}

// GetLibnfVersion return s the libnf version used underhood.
func (f *File) GetLibnfVersion() (string, error) {
	return getStringInfo(internal.INFO_VERSION, f)
}

// GetNfdumpVersion return s the version of Nfdump used in the file.
func (f *File) GetNfdumpVersion() (string, error) {
	return getStringInfo(internal.INFO_NFDUMP_VERSION, f)
}

// GetIdent return s the string identification of the file.
func (f *File) GetIdent() (string, error) {
	return getStringInfo(internal.INFO_IDENT, f)
}

// IsCompressed return s whether the file is compressed.
func (f *File) IsCompressed() (bool, error) {
	return getBoolInfo(internal.INFO_COMPRESSED, f)
}

// IsAnonymized return s whether the file is anonymized.
func (f *File) IsAnonymized() (bool, error) {
	return getBoolInfo(internal.INFO_ANONYMIZED, f)
}

// HasCatalog return s whether the file has a catalog.
func (f *File) HasCatalog() (bool, error) {
	return getBoolInfo(internal.INFO_CATALOG, f)
}

// Various methods for retrieving uint64 metadata about the file.
func (f *File) GetFileVersion() (uint64, error) { return getUint64Info(internal.INFO_FILE_VERSION, f) }
func (f *File) GetBlocks() (uint64, error)      { return getUint64Info(internal.INFO_BLOCKS, f) }
func (f *File) GetFirst() (uint64, error)       { return getUint64Info(internal.INFO_FIRST, f) }
func (f *File) GetLast() (uint64, error)        { return getUint64Info(internal.INFO_LAST, f) }
func (f *File) GetFailures() (uint64, error)    { return getUint64Info(internal.INFO_FAILURES, f) }
func (f *File) GetFlows() (uint64, error)       { return getUint64Info(internal.INFO_FLOWS, f) }
func (f *File) GetBytes() (uint64, error)       { return getUint64Info(internal.INFO_BYTES, f) }
func (f *File) GetPackets() (uint64, error)     { return getUint64Info(internal.INFO_PACKETS, f) }
func (f *File) GetProcBlocks() (uint64, error)  { return getUint64Info(internal.INFO_PROC_BLOCKS, f) }
func (f *File) GetFlowsTcp() (uint64, error)    { return getUint64Info(internal.INFO_FLOWS_TCP, f) }
func (f *File) GetFlowsUdp() (uint64, error)    { return getUint64Info(internal.INFO_FLOWS_UDP, f) }
func (f *File) GetFlowsIcmp() (uint64, error)   { return getUint64Info(internal.INFO_FLOWS_ICMP, f) }
func (f *File) GetFlowsOther() (uint64, error)  { return getUint64Info(internal.INFO_FLOWS_OTHER, f) }
func (f *File) GetBytesTcp() (uint64, error)    { return getUint64Info(internal.INFO_BYTES_TCP, f) }
func (f *File) GetBytesUdp() (uint64, error)    { return getUint64Info(internal.INFO_BYTES_UDP, f) }
func (f *File) GetBytesIcmp() (uint64, error)   { return getUint64Info(internal.INFO_BYTES_ICMP, f) }
func (f *File) GetBytesOther() (uint64, error)  { return getUint64Info(internal.INFO_BYTES_OTHER, f) }
func (f *File) GetPacketsTcp() (uint64, error)  { return getUint64Info(internal.INFO_PACKETS_TCP, f) }
func (f *File) GetPacketsUdp() (uint64, error)  { return getUint64Info(internal.INFO_PACKETS_UDP, f) }
func (f *File) GetPacketsIcmp() (uint64, error) { return getUint64Info(internal.INFO_PACKETS_ICMP, f) }
func (f *File) GetPacketsOther() (uint64, error) {
	return getUint64Info(internal.INFO_PACKETS_OTHER, f)
}

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
		return LnfErr.ErrFileAlreadyOpened
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
		return LnfErr.ErrCannotOpenFile
	}
	f.opened = true
	return nil
}

// OpenAppend opens the file in append mode.
func (f *File) OpenAppend(inputFile string, weakErr bool) error {
	if f.opened {
		return LnfErr.ErrFileAlreadyOpened
	}

	flags := internal.APPEND
	if weakErr {
		flags |= internal.WEAKERR
	}
	status := internal.Open(&f.ptr, inputFile, uint(flags), "")
	if status != internal.OK {
		return LnfErr.ErrCannotOpenFile
	}
	f.opened = true
	return nil
}

// OpenWrite opens the file in write mode.
func (f *File) OpenWrite(outputFile string, ident string, anon bool, comp int, weakErr bool) error {
	if f.opened {
		return LnfErr.ErrFileAlreadyOpened
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
		return LnfErr.ErrCannotOpenFile
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
	return LnfErr.ErrFileNotOpened
}

// GetNextRecord reads the next record from the file.
func (file *File) GetNextRecord(r *LnfRec.Record) error {
	if !file.opened {
		return LnfErr.ErrFileNotOpened
	} else if !r.Allocated() {
		return LnfErr.ErrRecordNotAllocated
	}
	status := internal.Read(file.ptr, r.GetPtr())
	if status == internal.ERR_NOMEM {
		return LnfErr.ErrNoMem
	} else if status == internal.EOF {
		return LnfErr.ErrFileEof
	}
	return nil
}

// WriteRecord writes a record to the file.
func (file *File) WriteRecord(r *LnfRec.Record) error {
	if !file.opened {
		return LnfErr.ErrFileNotOpened
	} else if !r.Allocated() {
		return LnfErr.ErrRecordNotAllocated
	}
	status := internal.Write(file.ptr, r.GetPtr())
	if status == internal.ERR_NOMEM {
		return LnfErr.ErrNoMem
	} else if status == internal.ERR_WRITE {
		return LnfErr.ErrWrite
	}
	return nil
}
