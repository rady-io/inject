package rady

import (
	"time"
	"github.com/tidwall/gjson"
	"reflect"
)

var (
	INT    int64
	UINT   uint64
	FLOAT  float64
	STRING string
	BOOL   bool
	TIME   time.Time
	ARRAY  []gjson.Result
	MAP    map[string]gjson.Result
)

var (
	IntPtrType    = reflect.TypeOf(&INT)
	UintPtrType   = reflect.TypeOf(&UINT)
	FloatPtrType  = reflect.TypeOf(&FLOAT)
	StringPtrType = reflect.TypeOf(&STRING)
	BoolPtrType   = reflect.TypeOf(&BOOL)
	TimePtrType   = reflect.TypeOf(&TIME)
	ArrayPtrType  = reflect.TypeOf(&ARRAY)
	MapPtrType    = reflect.TypeOf(&MAP)
)

var GJsonTypesSet = map[reflect.Type]bool{
	IntPtrType:    true,
	UintPtrType:   true,
	FloatPtrType:  true,
	StringPtrType: true,
	BoolPtrType:   true,
	TimePtrType:   true,
	ArrayPtrType:  true,
	MapPtrType:    true,
}
