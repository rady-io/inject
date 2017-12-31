package summer

import "reflect"

func ContainsField(Mother reflect.Type, field interface{}) bool {
	fieldType := reflect.TypeOf(field)
	if innerField, ok := Mother.FieldByName(fieldType.Name()); ok {
		if innerField.Type == fieldType {
			return true
		}
	}
	return false
}

func ContainsFields(Mother reflect.Type, Set map[reflect.Type]bool) bool {
	for i := 0; i < Mother.NumField(); i++ {
		if _, ok := Set[Mother.Field(i).Type]; ok {
			return true
		}
	}
	return false
}
