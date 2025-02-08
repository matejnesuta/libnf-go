package libnf

import (
	"C"
	"libnf/internal"
	"unsafe"
)

type LnfFile uintptr
type LnfRecord uintptr
type LnfFileFlag uint

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

func Error() string {
	return internal.Error()
}

func getInfo(info int, f uintptr) (uintptr, error) {
	if f == 0 {
		return 0, ErrFileNotOpened
	}
	buf := make([]byte, internal.INFO_BUFSIZE) // Allocate memory
	data := uintptr(unsafe.Pointer(&buf[0]))
	status := internal.Info(f, info, data, int64(internal.INFO_BUFSIZE))
	if status == internal.ERR_NOMEM {
		return data, ErrNoMem
	}
	if status == internal.ERR_OTHER {
		return data, ErrOther
	}
	return data, nil
}

func getStringInfo(info int, f uintptr) (string, error) {
	data, err := getInfo(info, f)
	if err != nil {
		return "", err
	}
	return C.GoString((*C.char)(unsafe.Pointer(data))), nil
}

func getBoolInfo(info int, f uintptr) (bool, error) {
	data, err := getInfo(info, f)
	if err != nil {
		return false, err
	}
	value := *(*int)(unsafe.Pointer(data))
	return value == 1, nil
}

func getUint64Info(info int, f uintptr) (uint64, error) {
	data, err := getInfo(info, f)
	if err != nil {
		return 0, err
	}
	return *(*uint64)(unsafe.Pointer(data)), nil
}

func (f LnfFile) GetLibnfVersion() (string, error) {
	return getStringInfo(internal.INFO_VERSION, uintptr(f))
}

func (f LnfFile) GetNfdumpVersion() (string, error) {
	return getStringInfo(internal.INFO_NFDUMP_VERSION, uintptr(f))
}

func (f LnfFile) GetIdent() (string, error) {
	return getStringInfo(internal.INFO_IDENT, uintptr(f))
}

func (f LnfFile) IsCompressed() (bool, error) {
	return getBoolInfo(internal.INFO_COMPRESSED, uintptr(f))
}

func (f LnfFile) IsAnonimized() (bool, error) {
	return getBoolInfo(internal.INFO_ANONYMIZED, uintptr(f))
}

func (f LnfFile) HasCatalog() (bool, error) {
	return getBoolInfo(internal.INFO_CATALOG, uintptr(f))
}

func (f LnfFile) GetFileVersion() (uint64, error) {
	return getUint64Info(internal.INFO_FILE_VERSION, uintptr(f))
}

func (f LnfFile) GetBlocks() (uint64, error) {
	return getUint64Info(internal.INFO_BLOCKS, uintptr(f))
}

func (f LnfFile) GetFirst() (uint64, error) {
	return getUint64Info(internal.INFO_FIRST, uintptr(f))
}

func (f LnfFile) GetLast() (uint64, error) {
	return getUint64Info(internal.INFO_LAST, uintptr(f))
}

func (f LnfFile) GetFailures() (uint64, error) {
	return getUint64Info(internal.INFO_FAILURES, uintptr(f))
}

func (f LnfFile) GetFlows() (uint64, error) {
	return getUint64Info(internal.INFO_FLOWS, uintptr(f))
}

func (f LnfFile) GetBytes() (uint64, error) {
	return getUint64Info(internal.INFO_BYTES, uintptr(f))
}

func (f LnfFile) GetPackets() (uint64, error) {
	return getUint64Info(internal.INFO_PACKETS, uintptr(f))
}

func (f LnfFile) GetProcBlocks() (uint64, error) {
	return getUint64Info(internal.INFO_PROC_BLOCKS, uintptr(f))
}

func (f LnfFile) GetFlowsTcp() (uint64, error) {
	return getUint64Info(internal.INFO_FLOWS_TCP, uintptr(f))
}

func (f LnfFile) GetFlowsUdp() (uint64, error) {
	return getUint64Info(internal.INFO_FLOWS_UDP, uintptr(f))
}

func (f LnfFile) GetFlowsIcmp() (uint64, error) {
	return getUint64Info(internal.INFO_FLOWS_ICMP, uintptr(f))
}

func (f LnfFile) GetFlowsOther() (uint64, error) {
	return getUint64Info(internal.INFO_FLOWS_OTHER, uintptr(f))
}

func (f LnfFile) GetBytesTcp() (uint64, error) {
	return getUint64Info(internal.INFO_BYTES_TCP, uintptr(f))
}

func (f LnfFile) GetBytesUdp() (uint64, error) {
	return getUint64Info(internal.INFO_BYTES_UDP, uintptr(f))
}

func (f LnfFile) GetBytesIcmp() (uint64, error) {
	return getUint64Info(internal.INFO_BYTES_ICMP, uintptr(f))
}

func (f LnfFile) GetBytesOther() (uint64, error) {
	return getUint64Info(internal.INFO_BYTES_OTHER, uintptr(f))
}

func (f LnfFile) GetPacketsTcp() (uint64, error) {
	return getUint64Info(internal.INFO_PACKETS_TCP, uintptr(f))
}

func (f LnfFile) GetPacketsUdp() (uint64, error) {
	return getUint64Info(internal.INFO_PACKETS_UDP, uintptr(f))
}

func (f LnfFile) GetPacketsIcmp() (uint64, error) {
	return getUint64Info(internal.INFO_PACKETS_ICMP, uintptr(f))
}

func (f LnfFile) GetPacketsOther() (uint64, error) {
	return getUint64Info(internal.INFO_PACKETS_OTHER, uintptr(f))
}

func OpenRead(filename string, readLoop bool, weakErr bool) (LnfFile, error) {
	var ptr uintptr
	flags := internal.READ
	if readLoop {
		flags |= internal.READ_LOOP
	}
	if weakErr {
		flags |= internal.WEAKERR
	}
	status := internal.Open(&ptr, filename, uint(flags), "")
	if status != internal.OK {
		return LnfFile(ptr), ErrCannotOpenFile
	}
	return LnfFile(ptr), nil
}

func Close(file LnfFile) {
	internal.Close(uintptr(file))
}

func GetVersion() string {
	return "v1.33"
}
