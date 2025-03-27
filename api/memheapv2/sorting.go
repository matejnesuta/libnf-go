package memheapv2

import (
	"bytes"
	"net"
	"sort"
	"time"
)

func lessThan(a, b interface{}) bool {
	switch v1 := a.(type) {
	case uint8:
		v2, ok := b.(uint8)
		return ok && v1 < v2
	case uint16:
		v2, ok := b.(uint16)
		return ok && v1 < v2
	case uint32:
		v2, ok := b.(uint32)
		return ok && v1 < v2
	case uint64:
		v2, ok := b.(uint64)
		return ok && v1 < v2
	case float64:
		v2, ok := b.(float64)
		return ok && v1 < v2
	case string:
		v2, ok := b.(string)
		return ok && v1 < v2
	// this sorting is done based on how libnf sorts IP addresses, but that might not be correct
	case net.IP:
		v2, ok := b.(net.IP)
		if !ok {
			return false
		}
		if len(v1) > len(v2) {
			return true
		} else if len(v1) < len(v2) {
			return false
		}
		return bytes.Compare(v1, v2) > 0
	case time.Time:
		v2, ok := b.(time.Time)
		return ok && v1.Before(v2)
	case net.HardwareAddr:
		v2, ok := b.(net.HardwareAddr)
		return ok && bytes.Compare(v1, v2) < 0
	default:
		return false // Unsupported type
	}
}

func greaterThan(a, b interface{}) bool {
	switch v1 := a.(type) {
	case uint8:
		v2, ok := b.(uint8)
		return ok && v1 > v2
	case uint16:
		v2, ok := b.(uint16)
		return ok && v1 > v2
	case uint32:
		v2, ok := b.(uint32)
		return ok && v1 > v2
	case uint64:
		v2, ok := b.(uint64)
		return ok && v1 > v2
	case float64:
		v2, ok := b.(float64)
		return ok && v1 > v2
	case string:
		v2, ok := b.(string)
		return ok && v1 > v2
	// this sorting is done based on how libnf sorts IP addresses, but that might not be correct
	case net.IP:
		v2, ok := b.(net.IP)
		if len(v1) < len(v2) {
			return true
		} else if len(v1) > len(v2) {
			return false
		}
		return ok && bytes.Compare(v1, v2) < 0
	case time.Time:
		v2, ok := b.(time.Time)
		return ok && v1.After(v2)
	case net.HardwareAddr:
		v2, ok := b.(net.HardwareAddr)
		return ok && bytes.Compare(v1, v2) > 0
	default:
		return false // Unsupported type
	}
}

func sortRecords(m *MemHeapV2) {
	m.sortedKeys = make([]string, 0, len(m.table))
	for key := range m.table {
		m.sortedKeys = append(m.sortedKeys, key)
	}

	if m.sortType != SortNone {
		if m.sortByKey {
			if m.sortType == SortAsc {
				sort.Slice(m.sortedKeys, func(i, j int) bool {
					val1 := m.table[m.sortedKeys[i]].keys[m.sortOffset]
					val2 := m.table[m.sortedKeys[j]].keys[m.sortOffset]
					return lessThan(val1, val2)
				})
			} else {
				sort.Slice(m.sortedKeys, func(i, j int) bool {
					val1 := m.table[m.sortedKeys[i]].keys[m.sortOffset]
					val2 := m.table[m.sortedKeys[j]].keys[m.sortOffset]
					return greaterThan(val1, val2)
				})
			}
		} else {
			sort.Slice(m.sortedKeys, func(i, j int) bool {
				val1 := m.table[m.sortedKeys[i]].values[m.sortOffset]
				val2 := m.table[m.sortedKeys[j]].values[m.sortOffset]
				return lessThan(val1, val2)
			})
			if m.sortType == SortDesc {
				sort.Slice(m.sortedKeys, func(i, j int) bool {
					val1 := m.table[m.sortedKeys[i]].values[m.sortOffset]
					val2 := m.table[m.sortedKeys[j]].values[m.sortOffset]
					return greaterThan(val1, val2)
				})
			}
		}
	}
}
