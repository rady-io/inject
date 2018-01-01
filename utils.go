package rhapsody

import (
	"reflect"
	"rhapsody/types"
	"strings"
	"rhapsody/bean"
)

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

func CheckFieldPtr(fieldType reflect.Type) bool {
	if fieldType.Kind() != reflect.Ptr {
		return false
	}
	return true
}

func CheckConfiguration(field reflect.StructField) bool {
	return CheckFieldPtr(field.Type) && (field.Tag.Get("type") == "" && ContainsField(field.Type.Elem(), types.Configuration{}) || field.Tag.Get("type") == types.CONFIGURATION)
}

func CheckComponents(field reflect.StructField) bool {
	_, ok := types.COMPONENTS[field.Tag.Get("type")]
	return CheckFieldPtr(field.Type) && (ok || field.Tag.Get("type") == "" && ContainsFields(field.Type.Elem(), types.COMPONENT_TYPES))
}

func GetBeanName(Type reflect.Type, tag reflect.StructTag) string {
	if tag != *new(reflect.StructTag) {
		if aliasName := tag.Get("name"); strings.Trim(aliasName, " ") != "" {
			return aliasName
		}
	}
	return Type.String()
}

// if return true, just go!
func ConfirmAddBeanMap(BeanMap map[reflect.Type]map[string]*bean.Bean, fieldType reflect.Type, name string) bool {
	if BeanMap[fieldType] == nil {
		BeanMap[fieldType] = make(map[string]*bean.Bean)
	} else if _, ok := BeanMap[fieldType][name]; ok {
		return false
	}
	return true
}