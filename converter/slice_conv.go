package converter

import "reflect"

func FieldToInt64s(in interface{}, field string) []int64 {
	val := reflect.ValueOf(in)

	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return []int64{}
		}
		val = val.Elem()
	}

	if val.Kind() != reflect.Slice {
		return []int64{}
	}

	ids := make([]int64, 0, val.Len())
	for i := 0; i < val.Len(); i++ {
		item := val.Index(i)
		if item.Kind() == reflect.Ptr {
			if item.IsNil() {
				continue
			}
			item = item.Elem()
		}

		id := item.FieldByName(field)
		if !id.IsValid() {
			return []int64{}
		}
		if id.Kind() < reflect.Int || id.Kind() > reflect.Int64 {
			return []int64{}
		}
		ids = append(ids, id.Int())
	}
	return ids
}
