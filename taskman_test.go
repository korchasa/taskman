package taskman

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestResolveSuccess(t *testing.T) {
	var cases = []struct {
		name string
		givenTaskArgs []arg
		givenOsArgs []string
		wantedArgsValues []reflect.Value
		wantedErr error
	}{
		{
			"one string param",
			[]arg{{"str", "string", reflect.Value{}}},
			[]string{"bar"},
			[]reflect.Value{reflect.ValueOf("bar")},
			nil,
		},
		{
			"all types params",
			[]arg{
				{"str", "string", reflect.Value{}},
				{"int", "int", reflect.Value{}},
			},
			[]string{
				"bar", "5",
			},
			[]reflect.Value{
				reflect.ValueOf("bar"),
				reflect.ValueOf(5),
			},
			nil,
		},
		{
			"one string option",
			[]arg{{"str", "*string", reflect.Value{}}},
			[]string{"-str=bar"},
			[]reflect.Value{reflect.ValueOf("bar")},
			nil,
		},
		{
			"all types options",
			[]arg{
				{"str", "*string", reflect.Value{}},
				{"int", "*int", reflect.Value{}},
				{"bool", "*bool", reflect.Value{}},
			},
			[]string{
				"-str=bar",
				"-int=5",
				"-bool",
			},
			[]reflect.Value{
				reflect.ValueOf("bar"),
				reflect.ValueOf(5),
				reflect.ValueOf(true),
			},
			nil,
		},
		{
			"one param and one option, option first",
			[]arg{{"param", "string", reflect.Value{}}, {"opt", "*string", reflect.Value{}}},
			[]string{"-opt=baz", "bar"},
			[]reflect.Value{reflect.ValueOf("bar"), reflect.ValueOf("baz")},
			nil,
		},
		{
			"one param and one option, param first",
			[]arg{{"param", "string", reflect.Value{}}, {"opt", "*string", reflect.Value{}}},
			[]string{"bar", "-opt=baz"},
			[]reflect.Value{reflect.ValueOf("bar"), reflect.ValueOf("baz")},
			nil,
		},
		{
			"params and options together",
			[]arg{
				{"command", "string", reflect.Value{}},
				{"times", "int", reflect.Value{}},
				{"format", "*string", reflect.Value{}},
				{"limit", "*int", reflect.Value{}},
				{"verbose", "*bool", reflect.Value{}},
			},
			[]string{
				"replicate",
				"5",
				"-format=json",
				"-limit=100",
				"-verbose",
			},
			[]reflect.Value{
				reflect.ValueOf("replicate"),
				reflect.ValueOf(5),
				reflect.ValueOf("json"),
				reflect.ValueOf(100),
				reflect.ValueOf(true),
			},
			nil,
		},
	}
	for _, tt := range cases {
		t.Run(fmt.Sprintf("`%s`", tt.name), func(t *testing.T) {
			task := task{
				name: tt.name,
				args: tt.givenTaskArgs,
			}
			err := processArgs(&task, tt.givenOsArgs)
			for i := range tt.wantedArgsValues {
				assert.Equal(t, tt.wantedArgsValues[i].Interface(), task.args[i].value.Interface())
			}
			assert.Equal(t, tt.wantedErr, err)
		})
	}
}