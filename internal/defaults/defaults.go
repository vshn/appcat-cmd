package defaults

import (
	"reflect"

	"github.com/vshn/appcat-cli/internal/util"
)

// Defaults is a dummy struct to hold all "GetDefault" methods.
//
// This is required to call them dynamically by name.
type Defaults struct{}

func (d *Defaults) GetDefaultFor(kind string, input []util.Input) reflect.Value {
	fname := "Get" + kind + "Default"
	t := reflect.TypeOf(d)
	v := reflect.ValueOf(d)

	m, ok := t.MethodByName(fname)
	if !ok {
		panic("no default function defaults.Defaults." + fname + " found")
	}
	args := []reflect.Value{
		v,
		reflect.ValueOf(input),
	}
	value := m.Func.Call(args)[0]

	return value
}
