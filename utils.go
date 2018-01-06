package rady

import (
	"reflect"
	"strings"
	"io/ioutil"
	"github.com/ghodss/yaml"
	"fmt"
	"unicode"
)

// ContainsField return true when Mother has a child as same type as filed
func ContainsField(Mother reflect.Type, field interface{}) bool {
	fieldType := reflect.TypeOf(field)
	for i := 0; i < Mother.NumField(); i++ {
		if Mother.Field(i).Type == fieldType {
			return true
		}
	}
	return false
}

// ContainsFields return true when Mother has a child with a type Set contains
func ContainsFields(Mother reflect.Type, Set map[reflect.Type]bool) bool {
	for i := 0; i < Mother.NumField(); i++ {
		if _, ok := Set[Mother.Field(i).Type]; ok {
			return true
		}
	}
	return false
}

// CheckFieldPtr return true when fieldType is kind of Ptr
func CheckFieldPtr(fieldType reflect.Type) bool {
	return fieldType.Kind() == reflect.Ptr
}

// CheckConfiguration return true when type in its tag is CONFIGURATION or ContainsField(field.Type.Elem(), Configuration{})
func CheckConfiguration(field reflect.StructField) bool {
	return CheckFieldPtr(field.Type) && CheckStruct(field.Type.Elem()) && (field.Tag.Get("type") == "" && ContainsField(field.Type.Elem(), Configuration{}) || field.Tag.Get("type") == CONFIGURATION)
}

// CheckController return true when type in its tag is CONTROLLER or ContainsField(field.Type.Elem(), Controller{})
func CheckController(field reflect.StructField) bool {
	return CheckFieldPtr(field.Type) && CheckStruct(field.Type.Elem()) && (field.Tag.Get("type") == "" && ContainsField(field.Type.Elem(), Controller{}) || field.Tag.Get("type") == CONTROLLER)
}

// CheckMiddleware return true when type in its tag is MIDDLEWARE or ContainsField(field.Type.Elem(), Middleware{})
func CheckMiddleware(field reflect.StructField) bool {
	return CheckFieldPtr(field.Type) && CheckStruct(field.Type.Elem()) && (field.Tag.Get("type") == "" && ContainsField(field.Type.Elem(), Middleware{}) || field.Tag.Get("type") == MIDDLEWARE)
}

// CheckRouter return true when type in its tag is ROUTER or ContainsField(field.Type.Elem(), Router{})
func CheckRouter(field reflect.StructField) bool {
	return CheckFieldPtr(field.Type) && CheckStruct(field.Type.Elem()) && (field.Tag.Get("type") == "" && ContainsField(field.Type.Elem(), Router{}) || field.Tag.Get("type") == ROUTER)
}

// CheckComponents return true when type in its tag is in COMPONENTS or ContainsFields(field.Type.Elem(), ComponentTypes)
func CheckComponents(field reflect.StructField) bool {
	_, ok := COMPONENTS[field.Tag.Get("type")]
	return CheckFieldPtr(field.Type) && CheckStruct(field.Type.Elem()) && (ok || field.Tag.Get("type") == "" && ContainsFields(field.Type.Elem(), ComponentTypes))
}

func CheckStruct(fieldType reflect.Type) bool {
	return fieldType.Kind() == reflect.Struct
}

func CheckValues(field reflect.StructField) bool {
	_, ok := GJsonTypesSet[field.Type]
	return ok
}

// GetBeanName get name from tag or from Type
func GetBeanName(Type reflect.Type, tag reflect.StructTag) string {
	if tag != *new(reflect.StructTag) {
		if aliasName := tag.Get("name"); strings.Trim(aliasName, " ") != "" {
			return aliasName
		}
	}
	return Type.String()
}

// GetTagFromName generate tag from name
func GetTagFromName(name string) reflect.StructTag {
	return (reflect.StructTag)(fmt.Sprintf(`name:"%s"`, name))
}

/*
ConfirmAddBeanMap return true when BeanMap[fieldType] == nil or BeanMap[fieldType][name] doesn't exist

and this function will make a map if BeanMap[fieldType] == nil
 */
func ConfirmAddBeanMap(BeanMap map[reflect.Type]map[string]*Bean, fieldType reflect.Type, name string) bool {
	if BeanMap[fieldType] == nil {
		BeanMap[fieldType] = make(map[string]*Bean)
	} else if _, ok := BeanMap[fieldType][name]; ok {
		return false
	}
	return true
}

/*
ConfirmSameTypeInMap return true when len(BeanMap[fieldType]) > 0

and this function will also make a map if BeanMap[fieldType] == nil, but return false
 */
func ConfirmSameTypeInMap(BeanMap map[reflect.Type]map[string]*Bean, fieldType reflect.Type) bool {
	if BeanMap[fieldType] == nil {
		BeanMap[fieldType] = make(map[string]*Bean)
	} else if len(BeanMap[fieldType]) > 0 {
		return true
	}
	return false
}

// ConfirmBeanInMap return true when BeanMap[fieldType][name] exist
func ConfirmBeanInMap(BeanMap map[reflect.Type]map[string]*Bean, fieldType reflect.Type, name string) bool {
	if BeanMap[fieldType] != nil {
		if _, ok := BeanMap[fieldType][name]; ok {
			return true
		}
	}
	return false
}

/*
GetJSONFromAnyFile can get json string from file

This function work well only when:

	1. fileType == "yaml" or (fileType != "json" and path end with ".yml" or ".yaml"), and content in file is yaml
	2. content in file is json
 */
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

func GetNewPrefix(prefix string, path string) string {
	prefix = strings.TrimRight(prefix, "/")
	path = strings.Trim(path, "/")
	return strings.TrimRight(fmt.Sprintf("%s/%s", prefix, path), "/")
}

func GetPathFromType(field reflect.StructField, Type interface{}) string {
	prefix := field.Tag.Get("prefix")
	if prefix != "" {
		return prefix
	}

	for i := 0; i < field.Type.Elem().NumField(); i ++ {
		child := field.Type.Elem().Field(i)
		if child.Type == reflect.TypeOf(Type) {
			return child.Tag.Get("prefix")
		}
	}
	return ""
}

func ParseHandlerName(Name string) (ok bool, method interface{}, path string) {
	if method, ok = StrToMethod[Name]; ok {
		return
	}

	method, path = GetHttpMethodAndPath(Name)
	if method != nil {
		if path == "" {
			return true, method, path
		}
		pathSlice := SplitByUpper(path)
		for i, slot := range pathSlice {
			if IsStringAllUpper(slot) {
				pathSlice[i] = GetDynamicPath(slot)
			} else {
				pathSlice[i] = strings.ToLower(slot)
			}
		}
		return true, method, strings.Join(pathSlice, "/")
	}
	return
}

func SplitByUpper(raw string) []string {
	var start int
	result := make([]string, 0)
	runes := []rune(raw)
	for i, r := range runes {
		if unicode.IsUpper(r) && i != 0 && (i != len(runes) - 1 && unicode.IsLower(runes[i+1]) || unicode.IsLower(runes[i-1]) ) {
			result = append(result, string(runes[start: i]))
			start = i
		}

		if i == len(runes) - 1 && unicode.IsUpper(r) {
			result = append(result, string(runes[start: i+1]))
		}
	}
	return result
}

func GetHttpMethodAndPath(Name string) (interface{}, string) {
	splitMethod := func(r rune) bool { return r == '0' }
	result := strings.FieldsFunc(Name, splitMethod)
	if len(result) == 1 { // no "_"
		value := result[0]
		newResult := SplitByUpper(value)
		if len(newResult) > 0 {
			methodStr := newResult[0]
			if method, ok := StrToMethod[methodStr]; ok {
				return method, ""
			}
		}
	} else if len(result) == 2 {
		methodStr := result[0]
		if method, ok := StrToMethod[methodStr]; ok {
			return method, result[1]
		}
	}
	return nil, ""
}

func IsStringAllUpper(str string) bool {
	for _, u := range []rune(str) {
		if unicode.IsLower(u) {
			return false
		}
	}
	return true
}

func GetDynamicPath(upper string) string {
	return fmt.Sprintf(":%s", strings.ToLower(upper))
}

func CheckFilenameValid(Name string) bool {
	return Name != ""
}
