package memheapv2

import (
	"bytes"
	"net"
	"sort"
	"time"

	"github.com/matejnesuta/libnf-go/api/fields"
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

func calculateSortValue(m *MemHeapV2, key string, deps map[int]int) {
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
		first := arr[deps[fields.First]].(time.Time)
		last := arr[deps[fields.Last]].(time.Time)
		arr[m.sortOffset] = last.Sub(first).Milliseconds()
	} else if m.sortField == fields.CalcBps || m.sortField == fields.CalcPps {
		first := arr[deps[fields.First]].(time.Time)
		last := arr[deps[fields.Last]].(time.Time)
		duration := last.Sub(first).Seconds()
		if duration == 0 {
			arr[m.sortOffset] = float64(0)
		} else if m.sortField == fields.CalcBps {
			arr[m.sortOffset] = float64(arr[deps[fields.Doctets]].(uint64)) * 8 / duration
		} else {
			arr[m.sortOffset] = float64(arr[deps[fields.Dpkts]].(uint64)) / duration
		}
	} else {
		first := arr[deps[fields.Dpkts]].(uint64)
		last := arr[deps[fields.Doctets]].(uint64)
		arr[m.sortOffset] = float64(last) / float64(first)
	}

	// TODO add more cases
}

func getDepIndexes(m *MemHeapV2, field int, depIndexes map[int]int) map[int]int {
	deps, ok := dependencies[field]
	if !ok {
		return nil
	}

	for _, dep := range deps {
		index := searchList(&m.keyTemplateList, dep)
		if index == -1 {
			index = searchList(&m.valueTemplateList, dep)
		}
		if index == -1 {
			continue
		}
		_, ok := dependencies[dep]
		if ok {
			depIndexes = getDepIndexes(m, dep, depIndexes)
		} else {
			depIndexes[dep] = int(index)
		}
	}
	return depIndexes
}

func sortRecords(m *MemHeapV2) {
	m.sortedKeys = make([]string, 0, m.table.itemCount())

	deps := make(map[int]int, 0)
	deps = getDepIndexes(m, m.sortField, deps)

	for _, shard := range m.table {
		for key := range shard.m {
			if len(deps) > 0 {
				calculateSortValue(m, key, deps)
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
