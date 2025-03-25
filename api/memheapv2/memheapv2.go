package memheapv2

import (
	"fmt"
	"libnf/api/errors"
	"libnf/api/fields"
	"libnf/api/record"
	"net"
	"strconv"
	"sync"
	"time"
)

type aggrRecord struct {
	keys   []any
	values []any
}

type fieldOptions struct {
	field    int
	aggrType int
	sortType int
	numbits  uint
	numbits6 uint
}

type MemHeapV2 struct {
	keyTemplateList   []fieldOptions
	valueTemplateList []fieldOptions
	table             map[string]aggrRecord
	statsMode         bool
	sortOffset        int
	sortField         int
	sortByKey         bool
	sortType          int
	nfdumpComp        bool
	sortedKeys        []string
}

type MemHeapCursor struct {
	cursor uint64
}

var mapMux sync.Mutex

func NewMemHeapV2() *MemHeapV2 {
	return &MemHeapV2{
		table: make(map[string]aggrRecord),
	}
}

func searchList(list *[]fieldOptions, field int) int {
	i := 0
	for _, f := range *list {
		if f.field == field {
			return i
		}
		i++
	}
	return -1
}

func addOrUpdateList(list *[]fieldOptions, field fieldOptions) int {
	i := 0
	for _, f := range *list {
		if f.field == field.field {
			(*list)[i] = field
			return i
		}
		i++
	}
	*list = append(*list, field)
	return i
}

func (m *MemHeapV2) SortAggrOptions(field int, aggrType int, sortType int, numBits uint, numBits6 uint) error {
	m.sortedKeys = nil
	ret, ok := fields.FieldTypes[field]
	if !ok {
		return errors.ErrUnknownFld
	}
	if (field == fields.Username || field == fields.Brec1 || field == fields.MplsLabel || field == fields.EgressAcl || field == fields.IngressAcl) && aggrType != 0 {
		return errors.ErrUnknownFld
	}

	var fld fieldOptions
	fld.field = field
	fld.numbits = numBits
	fld.numbits6 = numBits6

	if aggrType == AggrAuto {
		_, ok := ret.(uint64)
		if ok && fld.numbits > 0 {
			fld.aggrType = AggrKey
		} else {
			fld.aggrType = defaults[field][0]
		}
	} else {
		fld.aggrType = aggrType
	}

	fld.sortType = sortType
	// here I would add an aggregation function to the field, but we have generics in Go
	offset := 0
	sortByKey := false
	if fld.aggrType == AggrKey {
		offset = addOrUpdateList(&m.keyTemplateList, fld)
		sortByKey = true
		_, ok := pairFields[field]
		if ok {
			m.statsMode = true
		}
	} else {
		offset = addOrUpdateList(&m.valueTemplateList, fld)
	}

	if fld.sortType != SortNone {
		m.sortField = field
		m.sortOffset = offset
		m.sortByKey = sortByKey
		m.sortType = fld.sortType
	}

	deps, ok := dependencies[field]
	if ok {
		for _, dep := range deps {
			if searchList(&m.keyTemplateList, dep) == -1 && searchList(&m.valueTemplateList, dep) == -1 {
				err := m.SortAggrOptions(dep, AggrAuto, SortNone, 0, 0)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func buildKey(record *record.Record, keyTemplateList []fieldOptions, pairset int) (string, []any, error) {
	var keyVals []any
	key := ""
	var field int
	for _, x := range keyTemplateList {
		f, ok := pairFields[x.field]
		if ok {
			field = f[pairset-1]
		} else {
			field = x.field
		}
		val, err := record.GetField(field)
		if err != nil {
			return "", nil, err
		}
		keyVals = append(keyVals, val)

		switch v := val.(type) {
		case uint8:
			key += strconv.FormatUint(uint64(v), 10)
		case uint16:
			key += strconv.FormatUint(uint64(v), 10)
		case uint32:
			key += strconv.FormatUint(uint64(v), 10)
		case uint64:
			key += strconv.FormatUint(v, 10)
		case float64:
			key += strconv.FormatFloat(v, 'f', -1, 64)
		case net.IP:
			if v.To4() != nil {
				ipv4Mask := net.CIDRMask(int(x.numbits), 32)
				key += v.Mask(ipv4Mask).String()
			} else {
				ipv6Mask := net.CIDRMask(int(x.numbits6), 128)
				key += v.Mask(ipv6Mask).String()
			}
		case time.Time:
			key += v.String()
		case net.HardwareAddr:
			key += v.String()
		default:
			return "", nil, errors.ErrUnknownFld
		}
		key += ";"
	}
	return key, keyVals, nil
}

func getValues(record *record.Record, valueTemplateList []fieldOptions, pairset int) ([]any, error) {
	var valVals []any
	for _, x := range valueTemplateList {
		val, err := record.GetField(x.field)
		if err != nil {
			return nil, err
		}
		valVals = append(valVals, val)
	}
	return valVals, nil
}

func insertOrUpdateRecord(table map[string]aggrRecord, key string, rec aggrRecord, values []fieldOptions) {
	// update record
	fmt.Println("Key: ", key)
	fmt.Println("Input: ", rec)
	if oldRec, ok := table[key]; ok {
		var tmpVal any
		fmt.Println("Updating record")
		fmt.Println("Old record: ", oldRec)
		for i, val := range rec.values {
			switch values[i].aggrType {
			case AggrMin:
				tmpVal = getMin(val, oldRec.values[i])
			case AggrMax:
				tmpVal = getMax(val, oldRec.values[i])
			case AggrSum:
				tmpVal = getSum(val, oldRec.values[i])
			case AggrOr:
				tmpVal = getOr(val, oldRec.values[i])
			}
			rec.values[i] = tmpVal
		}
		fmt.Println("New record: ", rec)
		table[key] = rec
		// insert new record
	} else {
		fmt.Println("Inserting record")
		fmt.Println("New record: ", rec)
		table[key] = rec
	}
	fmt.Println("--------------------------------")

}

func (m *MemHeapV2) WriteRecord(record *record.Record) error {
	if !record.Allocated() {
		return errors.ErrRecordNotAllocated
	}
	m.sortedKeys = nil
	pairset := 0
	if m.statsMode {
		pairset = 1
	}

	var recs aggrRecord
	key, keyVals, err := buildKey(record, m.keyTemplateList, pairset)
	if err != nil {
		return err
	}
	recs.keys = keyVals
	recs.values, err = getValues(record, m.valueTemplateList, pairset)
	if err != nil {
		return err
	}
	mapMux.Lock()
	insertOrUpdateRecord(m.table, key, recs, m.valueTemplateList)
	if pairset != 0 {
		key2, keyVals2, err := buildKey(record, m.keyTemplateList, 2)
		if err != nil {
			return err
		}
		recs.keys = keyVals2
		if m.nfdumpComp {
			if key == key2 {
				goto end
			}
		}
		insertOrUpdateRecord(m.table, key2, recs, m.valueTemplateList)
	}

end:
	mapMux.Unlock()
	return nil
}

func (m *MemHeapV2) SetNfdumpComp(on bool) {
	m.nfdumpComp = on
}

func (m *MemHeapV2) FirstRecordPosition() (MemHeapCursor, error) {
	var cursor MemHeapCursor
	if len(m.table) == 0 {
		return cursor, errors.ErrMemHeapEmpty
	}
	cursor.cursor = 0
	return cursor, nil
}

func (m *MemHeapV2) NextRecordPosition(cursor MemHeapCursor) (MemHeapCursor, error) {
	var newCursor MemHeapCursor
	if len(m.table) == 0 {
		return newCursor, errors.ErrMemHeapEmpty
	}
	newCursor.cursor = cursor.cursor + 1
	if newCursor.cursor >= uint64(len(m.table)) {
		return newCursor, errors.ErrMemHeapEnd
	}
	return newCursor, nil
}

func setFieldInRecord(rec *record.Record, field int, val any) error {
	var err error
	switch v := val.(type) {
	case uint8:
		err = record.SetField(rec, field, v)
	case uint16:
		err = record.SetField(rec, field, v)
	case uint32:
		err = record.SetField(rec, field, v)
	case uint64:
		err = record.SetField(rec, field, v)
	case time.Time:
		err = record.SetField(rec, field, v)
	case net.IP:
		err = record.SetField(rec, field, v)
	case net.HardwareAddr:
		err = record.SetField(rec, field, v)
	default:
		err = errors.ErrUnknownFld
	}
	if err != nil {
		return err
	}
	return nil
}

func (m *MemHeapV2) GetRecord(cursor *MemHeapCursor, rec *record.Record) error {
	if !rec.Allocated() {
		return errors.ErrRecordNotAllocated
	}

	rec.Clear()

	if len(m.table) == 0 {
		return errors.ErrMemHeapEmpty
	}

	if m.sortedKeys == nil {
		sortRecords(m)
	}

	if cursor.cursor >= uint64(len(m.table)) {
		return errors.ErrMemHeapEnd
	}
	key := m.sortedKeys[cursor.cursor]
	recs := m.table[key]
	for i, val := range recs.keys {
		setFieldInRecord(rec, m.keyTemplateList[i].field, val)
	}

	for i, val := range recs.values {
		setFieldInRecord(rec, m.valueTemplateList[i].field, val)
	}

	return nil
}

func (m *MemHeapV2) Clear() {
	m.table = make(map[string]aggrRecord)
	m.sortedKeys = nil
	m.keyTemplateList = nil
	m.valueTemplateList = nil
	m.sortOffset = 0
	m.sortField = 0
	m.sortByKey = false
	m.sortType = 0
	m.statsMode = false
	m.nfdumpComp = false
}
