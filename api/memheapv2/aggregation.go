package memheapv2

import "time"

func getMin(a, b any) any {
	switch a.(type) {
	case uint8:
		if a.(uint8) < b.(uint8) {
			return a
		}
		return b
	case uint16:
		if a.(uint16) < b.(uint16) {
			return a
		}
		return b
	case uint32:
		if a.(uint32) < b.(uint32) {
			return a
		}
		return b
	case uint64:
		if a.(uint64) < b.(uint64) {
			return a
		}
		return b
	case float64:
		if a.(float64) < b.(float64) {
			return a
		}
		return b
	case time.Time:
		if a.(time.Time).Before(b.(time.Time)) {
			return a
		}
		return b
	default:
		return nil
	}
}

func getMax(a, b any) any {
	switch a.(type) {
	case uint8:
		if a.(uint8) > b.(uint8) {
			return a
		}
		return b
	case uint16:
		if a.(uint16) > b.(uint16) {
			return a
		}
		return b
	case uint32:
		if a.(uint32) > b.(uint32) {
			return a
		}
		return b
	case uint64:
		if a.(uint64) > b.(uint64) {
			return a
		}
		return b
	case float64:
		if a.(float64) > b.(float64) {
			return a
		}
		return b
	case time.Time:
		if a.(time.Time).After(b.(time.Time)) {
			return a
		}
		return b
	default:
		return nil
	}
}

func getSum(a, b any) any {
	switch a.(type) {
	case uint8:
		return a.(uint8) + b.(uint8)
	case uint16:
		return a.(uint16) + b.(uint16)
	case uint32:
		return a.(uint32) + b.(uint32)
	case uint64:
		return a.(uint64) + b.(uint64)
	case float64:
		return a.(float64) + b.(float64)
	default:
		return nil
	}
}

func getOr(a, b any) any {
	switch a.(type) {
	case uint8:
		return a.(uint8) | b.(uint8)
	case uint16:
		return a.(uint16) | b.(uint16)
	case uint32:
		return a.(uint32) | b.(uint32)
	case uint64:
		return a.(uint64) | b.(uint64)
	default:
		return nil
	}
}
