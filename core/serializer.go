package platform

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
)

type serializer struct{}

func NewSerializer() *serializer {
	return &serializer{}
}

func (s *serializer) Serialize(value any) ([]byte, error) {
	return s.walk(value)
}

func (s *serializer) walk(value any) ([]byte, error) {
	if value == nil {
		return nil, nil
	}

	rv := reflect.ValueOf(value)

	if rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return nil, nil
		}
		return s.walk(rv.Elem().Interface())
	}

	switch rv.Kind() {
	case reflect.String:
		return []byte(rv.String()), nil

	case reflect.Bool:
		return []byte(strconv.FormatBool(rv.Bool())), nil

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return []byte(strconv.FormatInt(rv.Int(), 10)), nil

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return []byte(strconv.FormatUint(rv.Uint(), 10)), nil

	case reflect.Float32, reflect.Float64:
		return []byte(strconv.FormatFloat(rv.Float(), 'f', -1, 64)), nil

	case reflect.Struct, reflect.Map, reflect.Slice, reflect.Array:
		b, err := json.Marshal(value)
		if err != nil {
			return nil, fmt.Errorf("serialize: json marshal failed: %w", err)
		}
		return b, nil

	default:
		return nil, fmt.Errorf("serialize: unsupported type: %s", rv.Kind())
	}
}
