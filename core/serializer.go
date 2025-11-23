package core

import (
	"encoding/json"
	"reflect"
	"strconv"

	"github.com/awesome-goose/contracts"
	"github.com/awesome-goose/platform/errors"
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

	if errVal, ok := value.(error); ok {
		return contracts.SerialTypeError, []byte(errVal.Error()), nil
	}

	rv := reflect.ValueOf(value)

	if rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return contracts.SerialTypeNil, nil, nil
		}
		return s.walk(rv.Elem().Interface())
	}

	if rv.Kind() == reflect.Slice && rv.Type().Elem().Kind() == reflect.Uint8 {
		return contracts.SerialTypeBinary, rv.Bytes(), nil
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
			return contracts.SerialTypeNil, nil, errors.JSONMARSHALFAILED.WithError(err)
		}
		return contracts.SerialTypeObject, b, nil

	default:
		return contracts.SerialTypeNil, nil, errors.UNSUPPORTEDTYPE.WithMeta(rv.Kind())
	}
}
