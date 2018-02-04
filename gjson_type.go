package rady

import (
	"github.com/tidwall/gjson"
	"reflect"
	"time"
)

var (
	INT         int64
	UINT        uint64
	FLOAT       float64
	STRING      string
	BOOL        bool
	TIME        time.Time
	ARRAY       []gjson.Result
	MAP         map[string]gjson.Result
	ArrayInt    []int64
	ArrayUint   []uint64
	ArrayFloat  []float64
	ArrayString []string
	ArrayBool   []bool
	ArrayTime   []time.Time
)

var (
	IntType            = reflect.TypeOf(INT)
	UintType           = reflect.TypeOf(UINT)
	FloatType          = reflect.TypeOf(FLOAT)
	StringType         = reflect.TypeOf(STRING)
	BoolType           = reflect.TypeOf(BOOL)
	TimeType           = reflect.TypeOf(TIME)
	ArrayType          = reflect.TypeOf(ARRAY)
	ArrayIntType       = reflect.TypeOf(ArrayInt)
	ArrayUintType      = reflect.TypeOf(ArrayUint)
	ArrayFloatType     = reflect.TypeOf(ArrayFloat)
	ArrayStringType    = reflect.TypeOf(ArrayString)
	ArrayBoolType      = reflect.TypeOf(ArrayBool)
	ArrayTimeType      = reflect.TypeOf(ArrayTime)
	MapType            = reflect.TypeOf(MAP)
	IntPtrType         = reflect.TypeOf(&INT)
	UintPtrType        = reflect.TypeOf(&UINT)
	FloatPtrType       = reflect.TypeOf(&FLOAT)
	StringPtrType      = reflect.TypeOf(&STRING)
	BoolPtrType        = reflect.TypeOf(&BOOL)
	TimePtrType        = reflect.TypeOf(&TIME)
	ArrayPtrType       = reflect.TypeOf(&ARRAY)
	MapPtrType         = reflect.TypeOf(&MAP)
	ArrayIntPtrType    = reflect.TypeOf(&ArrayInt)
	ArrayUintPtrType   = reflect.TypeOf(&ArrayUint)
	ArrayFloatPtrType  = reflect.TypeOf(&ArrayFloat)
	ArrayStringPtrType = reflect.TypeOf(&ArrayString)
	ArrayBoolPtrType   = reflect.TypeOf(&ArrayBool)
	ArrayTimePtrType   = reflect.TypeOf(&ArrayTime)
)

var (
	GJsonPtrTypesSet = map[reflect.Type]bool{
		IntPtrType:         true,
		UintPtrType:        true,
		FloatPtrType:       true,
		StringPtrType:      true,
		BoolPtrType:        true,
		TimePtrType:        true,
		ArrayPtrType:       true,
		MapPtrType:         true,
		ArrayIntPtrType:    true,
		ArrayUintPtrType:   true,
		ArrayFloatPtrType:  true,
		ArrayStringPtrType: true,
		ArrayBoolPtrType:   true,
		ArrayTimePtrType:   true,
	}

	GJsonTypesSet = map[reflect.Type]bool{
		IntType:         true,
		UintType:        true,
		FloatType:       true,
		StringType:      true,
		BoolType:        true,
		TimeType:        true,
		ArrayType:       true,
		MapType:         true,
		ArrayIntType:    true,
		ArrayUintType:   true,
		ArrayFloatType:  true,
		ArrayStringType: true,
		ArrayBoolType:   true,
		ArrayTimeType:   true,
	}
)
