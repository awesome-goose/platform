package core

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"

	"github.com/awesome-goose/contracts"
)

type serializer struct{}

func NewSerializer() *serializer {
	return &serializer{}
}

func (s *serializer) Serialize(value any) (contracts.SerialType, []byte, error) {
	return s.walk(value)
}

func (s *serializer) walk(value any) (contracts.SerialType, []byte, error) {
	if value == nil {
		return contracts.SerialTypeNil, nil, nil
	}

	rv := reflect.ValueOf(value)

	if rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return contracts.SerialTypeNil, nil, nil
		}
		return s.walk(rv.Elem().Interface())
	}

	switch rv.Kind() {
	case reflect.String:
		return contracts.SerialTypeString, []byte(rv.String()), nil

	case reflect.Bool:
		return contracts.SerialTypeBool, []byte(strconv.FormatBool(rv.Bool())), nil

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return contracts.SerialTypeNumber, []byte(strconv.FormatInt(rv.Int(), 10)), nil

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return contracts.SerialTypeNumber, []byte(strconv.FormatUint(rv.Uint(), 10)), nil

	case reflect.Float32, reflect.Float64:
		return contracts.SerialTypeNumber, []byte(strconv.FormatFloat(rv.Float(), 'f', -1, 64)), nil

	case reflect.Struct, reflect.Map, reflect.Slice, reflect.Array:
		b, err := json.Marshal(value)
		if err != nil {
			return contracts.SerialTypeNil, nil, fmt.Errorf("serialize: json marshal failed: %w", err)
		}
		return contracts.SerialTypeObject, b, nil

	default:
		return contracts.SerialTypeNil, nil, fmt.Errorf("serialize: unsupported type: %s", rv.Kind())
	}
}
