package defaults

import (
	"reflect"
)

// Defaults is a dummy struct to hold all "GetDefault" functions.
type Defaults struct{}

func (d *Defaults) GetDefaultFor(kind string) reflect.Value {
	fname := "Get" + kind + "Default"
	t := reflect.TypeOf(d)
	v := reflect.ValueOf(d)

	m, ok := t.MethodByName(fname)
	if !ok {
		panic("no default function defaults.Defaults." + fname + " found")
	}

	value := m.Func.Call([]reflect.Value{v})[0]

	return value
}
