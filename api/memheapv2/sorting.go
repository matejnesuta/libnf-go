package memheapv2

import (
	"bytes"
	"fmt"
	"libnf/api/fields"
	"net"
	"sort"
	"time"
)

func lessThan(a, b interface{}) bool {
	switch v1 := a.(type) {
	case int64:
		v2, ok := b.(int64)
		return ok && v1 < v2
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
	case int64:
		v2, ok := b.(int64)
		return ok && v1 > v2
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

func calculateSortValue(m *MemHeapV2, key string) {
	if m.sortType == SortNone {
		return
	}

	var arr []any

	if m.sortByKey {
		arr = m.table.get(key).keys
	} else {
		arr = m.table.get(key).values
	}

	if m.sortField == fields.CalcDuration {
		first := arr[m.sortOffset+1].(time.Time)
		last := arr[m.sortOffset+2].(time.Time)
		arr[m.sortOffset] = last.Sub(first).Milliseconds()
	} else if m.sortField == fields.CalcBps || m.sortField == fields.CalcPps {
		item := arr[m.sortOffset+1].(uint64)
		first := arr[m.sortOffset+3].(time.Time)
		last := arr[m.sortOffset+4].(time.Time)
		duration := last.Sub(first).Seconds()
		if duration == 0 {
			arr[m.sortOffset] = float64(0)
		} else if m.sortField == fields.CalcBps {
			arr[m.sortOffset] = float64(item) * 8 / duration
		} else {
			arr[m.sortOffset] = float64(item) / duration
		}
	} else {
		first := arr[m.sortOffset+1].(uint64)
		last := arr[m.sortOffset+2].(uint64)
		arr[m.sortOffset] = float64(last) / float64(first)
	}

	// TODO add more cases
}

func sortRecords(m *MemHeapV2) {
	fmt.Println(m.table.itemCount())
	m.sortedKeys = make([]string, 0, m.table.itemCount())

	_, ok := dependencies[m.sortField]

	for _, shard := range m.table {
		for key := range shard.m {
			if ok {
				calculateSortValue(m, key)
			}
			m.sortedKeys = append(m.sortedKeys, key)
		}
	}

	if m.sortType != SortNone {
		if m.sortByKey {
			if m.sortType == SortAsc {
				sort.Slice(m.sortedKeys, func(i, j int) bool {
					val1 := m.table.get(m.sortedKeys[i]).keys[m.sortOffset]
					val2 := m.table.get(m.sortedKeys[j]).keys[m.sortOffset]
					return lessThan(val1, val2)
				})
			} else {
				sort.Slice(m.sortedKeys, func(i, j int) bool {
					val1 := m.table.get(m.sortedKeys[i]).keys[m.sortOffset]
					val2 := m.table.get(m.sortedKeys[j]).keys[m.sortOffset]
					return greaterThan(val1, val2)
				})
			}
		} else {
			if m.sortType == SortAsc {
				sort.Slice(m.sortedKeys, func(i, j int) bool {
					val1 := m.table.get(m.sortedKeys[i]).values[m.sortOffset]
					val2 := m.table.get(m.sortedKeys[j]).values[m.sortOffset]
					return lessThan(val1, val2)
				})
			}
			if m.sortType == SortDesc {

				sort.Slice(m.sortedKeys, func(i, j int) bool {
					key1 := m.table.get(m.sortedKeys[i])
					key2 := m.table.get(m.sortedKeys[j])
					val1 := key1.values[m.sortOffset]
					val2 := key2.values[m.sortOffset]
					return greaterThan(val1, val2)
				})
			}
		}
	}
}
