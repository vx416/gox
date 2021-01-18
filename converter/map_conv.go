package converter

import "reflect"

// SliceToMap slice to map
func SliceToMap(from interface{}, keyField string) interface{} {
	slice := reflect.ValueOf(from)
	if slice.Len() == 0 {
		return nil
	}

	element := slice.Index(0)
	if element.Kind() == reflect.Ptr {
		element = element.Elem()
	}

	keyType := element.FieldByName(keyField).Type()
	valueType := slice.Index(0).Type()
	to := reflect.MakeMap(reflect.MapOf(keyType, valueType))

	for i := 0; i < slice.Len(); i++ {
		element := slice.Index(i)
		var field reflect.Value
		if element.Kind() == reflect.Ptr {
			field = element.Elem().FieldByName(keyField)
		} else {
			field = element.FieldByName(keyField)
		}

		to.SetMapIndex(field, element)
	}

	return to.Interface()
}
