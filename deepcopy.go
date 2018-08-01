package deepcopy

import (
	"fmt"
	"reflect"
)

// ErrMapDepth is a overflow maxdepth error
var ErrMapDepth = fmt.Errorf("overflow the maxdepth")

var errInvalidValue = fmt.Errorf("invalid value")

// Copy return a copy from value
func Copy(value interface{}) (interface{}, error) {
	v, err := deepCopy(reflect.ValueOf(value), 0, 1024)
	if err != nil {
		if err == errInvalidValue {
			return nil, nil
		}
		return nil, err
	}
	return v.Interface(), nil
}

func deepCopy(value reflect.Value, dep, maxDep int) (reflect.Value, error) {
	if !value.IsValid() {
		return reflect.ValueOf(nil), errInvalidValue
	}
	if dep > maxDep {
		return reflect.Zero(value.Type()), ErrMapDepth
	}

	typ := value.Type()
	kind := typ.Kind()
	if isBasicType(kind) {
		v := value.Interface()
		return reflect.ValueOf(v), nil
	}

	switch kind {
	case reflect.Interface:
		if value.IsNil() {
			return reflect.Zero(value.Type()), nil
		}
		return deepCopy(value.Elem(), dep+1, maxDep)
	case reflect.Ptr:
		if value.IsNil() {
			return reflect.Zero(value.Type()), nil
		}
		v := reflect.New(typ.Elem())
		value, err := deepCopy(value.Elem(), dep+1, maxDep)
		if err != nil {
			if err != errInvalidValue {
				return value, err
			}
			value = reflect.ValueOf(nil)
		}
		v.Elem().Set(value)
		return v, nil
	case reflect.Map:
		if value.IsNil() {
			return reflect.Zero(typ), nil
		}
		v := reflect.MakeMap(typ)
		for _, key := range value.MapKeys() {
			value := value.MapIndex(key)
			value, err := deepCopy(value, dep+1, maxDep)
			if err != nil {
				if err != errInvalidValue {
					return value, err
				}
				value = reflect.ValueOf(nil)
			}
			v.SetMapIndex(key, value)
		}
		return v, nil
	case reflect.Slice:
		if value.IsNil() {
			return reflect.Zero(typ), nil
		}
		v := reflect.MakeSlice(typ, value.Len(), value.Cap())
		for i := 0; i < v.Len(); i++ {
			value := value.Index(i)
			value, err := deepCopy(value, dep+1, maxDep)
			if err != nil {
				if err != errInvalidValue {
					return value, err
				}
				value = reflect.ValueOf(nil)
			}
			v.Index(i).Set(value)
		}
		return v, nil
	case reflect.Array:
		v := reflect.New(typ).Elem()
		for i := 0; i < v.Len(); i++ {
			value := value.Index(i)
			value, err := deepCopy(value, dep+1, maxDep)
			if err != nil {
				if err != errInvalidValue {
					return value, err
				}
				value = reflect.ValueOf(nil)
			}
			v.Index(i).Set(value)
		}
		return v, nil
	case reflect.Struct:
		v := reflect.New(typ).Elem()
		for i := 0; i < value.NumField(); i++ {
			field := typ.Field(i)
			if field.PkgPath != "" {
				continue
			}
			value := value.Field(i)
			value, err := deepCopy(value, dep+1, maxDep)
			if err != nil {
				if err != errInvalidValue {
					return value, err
				}
				value = reflect.ValueOf(nil)
			}
			v.FieldByName(field.Name).Set(value)
		}
		return v, nil
	}

	return reflect.Zero(value.Type()), fmt.Errorf("not support kind, %s", kind)
}

func isBasicType(kind reflect.Kind) bool {
	return kind >= reflect.Bool && kind <= reflect.Complex128 || reflect.String == kind
}
