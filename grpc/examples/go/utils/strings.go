package utils

import (
	"fmt"
	"math/rand/v2"
	"reflect"
	"strings"
)

func RandString(n int) string {
	charset := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	sb := strings.Builder{}
	sb.Grow(n)
	for i := 0; i < n; i++ {
		sb.WriteByte(charset[rand.IntN(len(charset))])
	}
	return sb.String()
}

func GetPointElem(v reflect.Value) reflect.Value {
	if v.Kind() == reflect.Pointer {
		return GetPointElem(v.Elem())
	}
	return v
}

func String(obj any) string {
	switch obj := obj.(type) {
	case byte, rune:
		return fmt.Sprintf(`'%c'`, obj)
	}

	val := reflect.ValueOf(obj)
	ts := val.Type().String()
	if val.Kind() == reflect.Pointer {
		if val.IsNil() {
			return "nil"
		}
		val = GetPointElem(val)
	}

	switch val.Kind() {
	case reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64,
		reflect.Complex64, reflect.Complex128:
		return fmt.Sprintf("%v", val.Interface())
	case reflect.String:
		return fmt.Sprintf(`"%s"`, val.Interface())
	case reflect.Chan:
		return ts
	case reflect.Func:
		fn, _ := NewFunc(val.Interface())
		return fn.String()
	case reflect.Array, reflect.Slice:
		l := val.Len()
		pairs := make([]string, l)
		for i := 0; i < l; i++ {
			vv := val.Index(i)
			pairs[i] = strings.TrimPrefix(String(vv.Interface()), vv.Type().String())
		}
		return fmt.Sprintf("%s{%s}", ts, strings.Join(pairs, ", "))
	case reflect.Map:
		pairs := make([]string, len(val.MapKeys()))
		for i, k := range val.MapKeys() {
			kv := val.MapIndex(k)
			key, value := String(k.Interface()), strings.TrimPrefix(String(kv.Interface()), kv.Type().String())
			pairs[i] = fmt.Sprintf("%v: %v", key, value)
		}
		return fmt.Sprintf("%s{%s}", ts, strings.Join(pairs, ", "))
	case reflect.Struct:
		l := val.NumField()
		pairs := make([]string, 0, l)
		typ := val.Type()
		for i := 0; i < l; i++ {
			t, v := typ.Field(i), val.Field(i)
			if v.Kind() == reflect.Pointer && v.IsNil() {
				pairs = append(pairs, fmt.Sprintf("%v: %v", t.Name, "nil"))
				continue
			}
			if v.CanInterface() {
				pairs = append(pairs, fmt.Sprintf("%v: %v", t.Name, String(v.Interface())))
			}
		}
		return fmt.Sprintf("%s{%s}", ts, strings.Join(pairs, ", "))
	default:
		return fmt.Sprintf("%s{%v}", ts, val.Interface())
	}
}
