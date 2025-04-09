package memheapv2

import (
	"libnf/api/errors"
	"libnf/api/fields"
	"libnf/api/record"
	"net"
	"strconv"
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
	table             shardedMap[aggrRecord]
	statsMode         bool
	sortOffset        int
	sortField         int
	sortByKey         bool
	sortType          int
	nfdumpComp        bool
	sortedKeys        []string
	shards            uint
}

type MemHeapCursor struct {
	cursor uint64
}

func NewMemHeapV2(shards uint) *MemHeapV2 {
	if shards == 0 {
		shards = 1
	}

	return &MemHeapV2{
		table:  newShardedMap[aggrRecord](shards),
		shards: shards,
	}
}

func (m *MemHeapV2) SetShards(shards uint) {
	if shards == 0 {
		shards = 1
	}
	m.shards = shards
	m.table = newShardedMap[aggrRecord](shards)
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

		ip, ok := val.(net.IP)
		if ok {
			if ip.To4() != nil {
				val = ip.Mask(net.CIDRMask(int(x.numbits), 32))
			} else {
				val = ip.Mask(net.CIDRMask(int(x.numbits6), 128))
			}
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
			key += v.String()
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
			val = nil
		}
		valVals = append(valVals, val)
	}
	return valVals, nil
}

func insertOrUpdateRecord(table map[string]aggrRecord, key string, rec aggrRecord, values []fieldOptions) {
	// update record
	// fmt.Println("Values: ", rec.values)
	// fmt.Println("Key: ", key)
	// fmt.Println("Input: ", rec)
	if oldRec, ok := table[key]; ok {
		newValues := make([]any, len(oldRec.values))
		copy(newValues, oldRec.values)

		for i, val := range rec.values {
			switch values[i].aggrType {
			case AggrMin:
				newValues[i] = getMin(val, oldRec.values[i])
			case AggrMax:
				newValues[i] = getMax(val, oldRec.values[i])
			case AggrSum:
				newValues[i] = getSum(val, oldRec.values[i])
			case AggrOr:
				newValues[i] = getOr(val, oldRec.values[i])
			}
		}

		rec.values = newValues // Ensure `rec.values` is a fresh copy
		table[key] = rec
	} else {
		// fmt.Println("Inserting record")
		// fmt.Println("New record: ", rec)
		table[key] = rec
	}
	// fmt.Println("--------------------------------")

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

	shard := m.table.getShard(key)
	shard.Lock()
	insertOrUpdateRecord(shard.m, key, recs, m.valueTemplateList)
	shard.Unlock()
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
		shard2 := m.table.getShard(key2)
		shard2.Lock()
		insertOrUpdateRecord(shard2.m, key2, recs, m.valueTemplateList)
		shard2.Unlock()
	}

end:
	return nil
}

func (m *MemHeapV2) SetNfdumpComp(on bool) {
	m.nfdumpComp = on
}

func (m *MemHeapV2) FirstRecordPosition() (MemHeapCursor, error) {
	var cursor MemHeapCursor
	if m.table.itemCount() == 0 {
		return cursor, errors.ErrMemHeapEmpty
	}
	cursor.cursor = 0
	return cursor, nil
}

func (m *MemHeapV2) NextRecordPosition(cursor MemHeapCursor) (MemHeapCursor, error) {
	var newCursor MemHeapCursor
	if m.table.itemCount() == 0 {
		return newCursor, errors.ErrMemHeapEmpty
	}
	newCursor.cursor = cursor.cursor + 1
	if newCursor.cursor >= uint64(m.table.itemCount()) {
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

	if m.table.itemCount() == 0 {
		return errors.ErrMemHeapEmpty
	}

	// fmt.Println(m.valueTemplateList)
	// fmt.Println("sort offset: ", m.sortOffset)
	// fmt.Println("sort field: ", m.sortField)
	// fmt.Println("sort by key: ", m.sortByKey)

	if m.sortedKeys == nil {
		sortRecords(m)
	}

	if cursor.cursor >= uint64(m.table.itemCount()) {
		return errors.ErrMemHeapEnd
	}
	key := m.sortedKeys[cursor.cursor]
	recs := m.table.get(key)
	for i, val := range recs.keys {
		setFieldInRecord(rec, m.keyTemplateList[i].field, val)
	}

	for i, val := range recs.values {
		setFieldInRecord(rec, m.valueTemplateList[i].field, val)
	}

	return nil
}

func (m *MemHeapV2) Clear() {
	m.table = newShardedMap[aggrRecord](m.shards)
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
