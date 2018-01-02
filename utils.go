package rhapsody

import (
	"reflect"
	"strings"
	"io/ioutil"
	"github.com/ghodss/yaml"
)

func ContainsField(Mother reflect.Type, field interface{}) bool {
	fieldType := reflect.TypeOf(field)
	for i := 0; i < Mother.NumField(); i++ {
		if Mother.Field(i).Type == fieldType {
			return true
		}
	}
	return false
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
	return CheckFieldPtr(field.Type) && (field.Tag.Get("type") == "" && ContainsField(field.Type.Elem(), Configuration{}) || field.Tag.Get("type") == CONFIGURATION)
}

func CheckComponents(field reflect.StructField) bool {
	_, ok := COMPONENTS[field.Tag.Get("type")]
	return CheckFieldPtr(field.Type) && (ok || field.Tag.Get("type") == "" && ContainsFields(field.Type.Elem(), COMPONENT_TYPES))
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
func ConfirmAddBeanMap(BeanMap map[reflect.Type]map[string]*Bean, fieldType reflect.Type, name string) bool {
	if BeanMap[fieldType] == nil {
		BeanMap[fieldType] = make(map[string]*Bean)
	} else if _, ok := BeanMap[fieldType][name]; ok {
		return false
	}
	return true
}

func GetJSONFromAnyFile(path string, fileType string) (string, error) {
	fileBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	if fileType != JSON {
		if fileType == YAML || strings.HasSuffix(path, ".yml") || strings.HasSuffix(path, ".yaml") {
			fileBytes, err = yaml.YAMLToJSON(fileBytes)
		}
	}
	return string(fileBytes), err
}
