package web

import (
	"errors"
	"reflect"
	"strconv"
)

func assignFromString(s string, v reflect.Value) error {
	t := v.Type()
	switch t.Kind() {
	case reflect.Bool:
		val, err := strconv.ParseBool(s)
		v.SetBool(val)
		return err
	case reflect.Int:
		fallthrough
	case reflect.Int8:
		fallthrough
	case reflect.Int32:
		fallthrough
	case reflect.Int64:
		val, err := strconv.ParseInt(s, 10, t.Bits())
		v.SetInt(val)
		return err
	case reflect.Uint:
		fallthrough
	case reflect.Uint8:
		fallthrough
	case reflect.Uint32:
		fallthrough
	case reflect.Uint64:
		val, err := strconv.ParseUint(s, 10, t.Bits())
		v.SetUint(val)
		return err
	case reflect.Float32:
		fallthrough
	case reflect.Float64:
		val, err := strconv.ParseFloat(s, t.Bits())
		v.SetFloat(val)
		return err
	case reflect.String:
		v.SetString(s)
		return nil
	}
	return errors.New("invalid type")
}

func translateString(s string, t reflect.Type) (reflect.Value, error) {
	v := reflect.New(t).Elem()
	err := assignFromString(s, v)
	return v, err
}

func translateStringSlice(s []string, t reflect.Type) (reflect.Value, error) {
	// make a slice
	if t.Kind() != reflect.Slice {
		return reflect.Value{}, errors.New("not a slice")
	}
	slice := reflect.MakeSlice(t, 0, len(s))
	for _, strval := range s {
		val, err := translateString(strval, t.Elem())
		if err != nil {
			return slice, err
		}
		slice = reflect.Append(slice, val)
	}
	return slice, nil
}

func assignTo(to interface{}, from string) error {
	// "to" must be a pointer
	ptrval := reflect.ValueOf(to)
	if ptrval.Kind() != reflect.Ptr {
		return errors.New("type error: not a pointer")
	}
	return assignFromString(from, ptrval.Elem())
}

func assignToStruct(to interface{}, from map[string]string) error {
	// "to" must be a pointer to a struct
	ptrval := reflect.ValueOf(to)
	if ptrval.Kind() != reflect.Ptr {
		return errors.New("type error: not a pointer")
	}
	structval := ptrval.Elem()
	if structval.Kind() != reflect.Struct {
		return errors.New("value error: not a struct")
	}
	nfields := structval.NumField()
	for i := 0; i < nfields; i++ {
		fv := structval.Type().Field(i)
		if valstring, ok := from[fv.Name]; ok {
			if err := assignFromString(valstring, structval.Field(i)); err != nil {
				return err
			}
		}
	}
	return nil
}

