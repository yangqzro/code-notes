package utils

import (
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"strings"
)

var (
	ErrFuncNil         = errors.New("callable is nil")
	ErrNotFunc         = errors.New("callable is not a func")
	ErrFuncNotCallable = errors.New("callable is nil or is not a func")
)

type Func struct {
	raw  any
	val  reflect.Value
	meta *runtime.Func
}

func (f *Func) canCall() {
	if f.raw == nil || f.meta == nil {
		panic(ErrFuncNil)
	}
	if f.val.Kind() != reflect.Func {
		panic(ErrNotFunc)
	}
}

func (f *Func) Func() any {
	return f.raw
}

func (f *Func) CanCall() bool {
	return f.raw != nil && f.val.Kind() == reflect.Func
}

func (f *Func) Name() string {
	f.canCall()
	ret, _ := At(strings.Split(f.meta.Name(), "/"), -1)
	return ret
}

func (f *Func) FileLine() (string, int) {
	f.canCall()
	return f.meta.FileLine(f.meta.Entry())
}

func (f *Func) GetType() reflect.Type {
	f.canCall()
	return f.val.Type()
}

func (f *Func) InTypes() []reflect.Type {
	f.canCall()
	l := f.val.Type().NumIn()
	params := make([]reflect.Type, l)
	for i := 0; i < l; i++ {
		params[i] = f.val.Type().In(i)
	}
	return params
}

func (f *Func) OutTypes() []reflect.Type {
	f.canCall()
	l := f.val.Type().NumOut()
	results := make([]reflect.Type, l)
	for i := 0; i < l; i++ {
		results[i] = f.val.Type().Out(i)
	}
	return results
}

func (f *Func) Call(args ...any) []any {
	f.canCall()
	in := make([]reflect.Value, len(args))
	for i, arg := range args {
		in[i] = reflect.ValueOf(arg)
	}
	out := f.val.Call(in)
	if len(out) == 0 {
		return nil
	}
	rets := make([]any, len(out))
	for i, v := range out {
		rets[i] = v.Interface()
	}
	return rets
}

func (f *Func) Equals(other *Func) bool {
	if other == nil {
		return false
	}
	if f.raw == nil || other.raw == nil {
		return f.raw == other.raw
	}
	ap := f.val.Pointer()
	bp := other.val.Pointer()
	return ap == bp
}

func (f *Func) EqualsFunc(fn any) bool {
	if fn == nil {
		return f.raw == nil
	}
	if f.raw == nil {
		return false
	}
	v := reflect.ValueOf(fn)
	if v.Kind() != reflect.Func {
		return false
	}
	ap := f.val.Pointer()
	bp := v.Pointer()
	return ap == bp
}

func (f *Func) String() string {
	if !f.CanCall() {
		return "nil"
	}
	name := f.Name()
	file, line := f.FileLine()

	in := strings.Join(Map(f.InTypes(), func(t reflect.Type, _ int) string { return t.String() }), ", ")
	out := strings.Join(Map(f.OutTypes(), func(t reflect.Type, _ int) string { return t.String() }), ", ")
	if out != "" {
		out = " (" + out + ")"
	}
	return fmt.Sprintf("func %s(%s)%s at %s:%d", name, in, out, file, line)
}

func NewFunc(callable any) (*Func, error) {
	v := reflect.ValueOf(callable)
	var meta *runtime.Func
	if callable != nil {
		meta = runtime.FuncForPC(v.Pointer())
		if v.Kind() != reflect.Func {
			return nil, ErrNotFunc
		}
	}
	f := &Func{raw: callable, val: v, meta: meta}
	return f, nil
}

func FuncEquals(a, b any) bool {
	if a == nil || b == nil {
		return a == b
	}
	ap := reflect.ValueOf(a).Pointer()
	bp := reflect.ValueOf(b).Pointer()

	an := runtime.FuncForPC(ap).Name()
	bn := runtime.FuncForPC(bp).Name()
	return ap == bp && an == bn
}
