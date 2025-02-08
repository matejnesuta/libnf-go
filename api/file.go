package libnf

import (
	"C"
	"libnf/internal"
	"unsafe"
)

type File struct {
	ptr    uintptr
	opened bool
}
type LnfRecord uintptr

// read flags:
// #define LNF_READ			0x00
// #define LNF_READ_LOOP	0x40

// write flags:
// #define LNF_WRITE		0x01
// #define LNF_ANON			0x02
// #define LNF_COMP			0x04
// #define LNF_COMP_LZO		0x04
// #define LNF_COMP_BZ2		0x20

// flags for append:
// #define LNF_APPEND			0x10

// this flag is stored into the lnf_file_t structure, but it is not clear whether it is used anywhere
// #define LNF_WEAKERR			0x08

const (
	No  = 0
	LZO = internal.COMP_LZO
	BZ2 = internal.COMP_BZ2
)

func getInfo(info int, f *File) (uintptr, error) {
	if !(*f).opened {
		return 0, ErrFileNotOpened
	}
	buf := make([]byte, internal.INFO_BUFSIZE) // Allocate memory
	data := uintptr(unsafe.Pointer(&buf[0]))
	status := internal.Info((*f).ptr, info, data, int64(internal.INFO_BUFSIZE))
	if status == internal.ERR_NOMEM {
		return data, ErrNoMem
	}
	if status == internal.ERR_OTHER {
		return data, ErrOther
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
	return value == 1, nil
}

func getUint64Info(info int, f *File) (uint64, error) {
	data, err := getInfo(info, f)
	if err != nil {
		return 0, err
	}
	return *(*uint64)(unsafe.Pointer(data)), nil
}

func (f *File) GetLibnfVersion() (string, error) {
	return getStringInfo(internal.INFO_VERSION, f)
}

func (f *File) GetNfdumpVersion() (string, error) {
	return getStringInfo(internal.INFO_NFDUMP_VERSION, f)
}

func (f *File) GetIdent() (string, error) {
	return getStringInfo(internal.INFO_IDENT, f)
}

func (f *File) IsCompressed() (bool, error) {
	return getBoolInfo(internal.INFO_COMPRESSED, f)
}

func (f *File) IsAnonimized() (bool, error) {
	return getBoolInfo(internal.INFO_ANONYMIZED, f)
}

func (f *File) HasCatalog() (bool, error) {
	return getBoolInfo(internal.INFO_CATALOG, f)
}

func (f *File) GetFileVersion() (uint64, error) {
	return getUint64Info(internal.INFO_FILE_VERSION, f)
}

func (f *File) GetBlocks() (uint64, error) {
	return getUint64Info(internal.INFO_BLOCKS, f)
}

func (f *File) GetFirst() (uint64, error) {
	return getUint64Info(internal.INFO_FIRST, f)
}

func (f *File) GetLast() (uint64, error) {
	return getUint64Info(internal.INFO_LAST, f)
}

func (f *File) GetFailures() (uint64, error) {
	return getUint64Info(internal.INFO_FAILURES, f)
}

func (f *File) GetFlows() (uint64, error) {
	return getUint64Info(internal.INFO_FLOWS, f)
}

func (f *File) GetBytes() (uint64, error) {
	return getUint64Info(internal.INFO_BYTES, f)
}

func (f *File) GetPackets() (uint64, error) {
	return getUint64Info(internal.INFO_PACKETS, f)
}

func (f *File) GetProcBlocks() (uint64, error) {
	return getUint64Info(internal.INFO_PROC_BLOCKS, f)
}

func (f *File) GetFlowsTcp() (uint64, error) {
	return getUint64Info(internal.INFO_FLOWS_TCP, f)
}

func (f *File) GetFlowsUdp() (uint64, error) {
	return getUint64Info(internal.INFO_FLOWS_UDP, f)
}

func (f *File) GetFlowsIcmp() (uint64, error) {
	return getUint64Info(internal.INFO_FLOWS_ICMP, f)
}

func (f *File) GetFlowsOther() (uint64, error) {
	return getUint64Info(internal.INFO_FLOWS_OTHER, f)
}

func (f *File) GetBytesTcp() (uint64, error) {
	return getUint64Info(internal.INFO_BYTES_TCP, f)
}

func (f *File) GetBytesUdp() (uint64, error) {
	return getUint64Info(internal.INFO_BYTES_UDP, f)
}

func (f *File) GetBytesIcmp() (uint64, error) {
	return getUint64Info(internal.INFO_BYTES_ICMP, f)
}

func (f *File) GetBytesOther() (uint64, error) {
	return getUint64Info(internal.INFO_BYTES_OTHER, f)
}

func (f *File) GetPacketsTcp() (uint64, error) {
	return getUint64Info(internal.INFO_PACKETS_TCP, f)
}

func (f *File) GetPacketsUdp() (uint64, error) {
	return getUint64Info(internal.INFO_PACKETS_UDP, f)
}

func (f *File) GetPacketsIcmp() (uint64, error) {
	return getUint64Info(internal.INFO_PACKETS_ICMP, f)
}

func (f *File) GetPacketsOther() (uint64, error) {
	return getUint64Info(internal.INFO_PACKETS_OTHER, f)
}

func (f *File) OpenRead(inputFile string, readLoop bool, weakErr bool) error {
	if f.opened {
		return ErrFileAlreadyOpened
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
		return ErrCannotOpenFile
	}
	f.opened = true
	return nil
}

func (f *File) OpenAppend(inputFile string, weakErr bool) error {
	if f.opened {
		return ErrFileAlreadyOpened
	}

	flags := internal.APPEND
	if weakErr {
		flags |= internal.WEAKERR
	}
	status := internal.Open(&f.ptr, inputFile, uint(flags), "")
	if status != internal.OK {
		return ErrCannotOpenFile
	}
	f.opened = true
	return nil
}

func (f *File) OpenWrite(inputFile string, outputFile string, anon bool, comp int, weakErr bool) error {
	if f.opened {
		return ErrFileAlreadyOpened
	}

	flags := internal.WRITE
	if anon {
		flags |= internal.ANON
	}
	if comp == LZO {
		flags |= internal.COMP_LZO
	}
	if comp == BZ2 {
		flags |= internal.COMP_BZ2
	}
	if weakErr {
		flags |= internal.WEAKERR
	}
	status := internal.Open(&f.ptr, inputFile, uint(flags), outputFile)
	if status != internal.OK {
		return ErrCannotOpenFile
	}
	f.opened = true
	return nil
}

func (file *File) Close() error {
	if file.opened {
		internal.Close(file.ptr)
		file.opened = false
		return nil
	}
	return ErrFileNotOpened
}
