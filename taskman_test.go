package taskman

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestExtractTasks(t *testing.T) {
	code := []byte(`
package main
// Hello says Hello
func Hello(who string, times *int, show *bool) {}
`)
	dir, err := ioutil.TempDir("", "example")
	assert.NoError(t, err)
	file := filepath.Join(dir, "tmpfile")
	err = ioutil.WriteFile(file, code, 0666)
	assert.NoError(t, err)
	tasks, err := extractTasks(file)
	assert.NoError(t, err)
	assert.Equal(t, task{
		name: "Hello",
		doc:  "says Hello",
		args: []arg{
			{name: "who", typeof: "string"},
			{name: "times", typeof: "int", optional: true},
			{name: "show", typeof: "bool", optional: true},
		},
	}, tasks[0])
}

func TestProcessArgsOnly(t *testing.T) {
	cases := []struct {
		name             string
		givenTaskArgs    string
		givenOsArgs      string
		wantedArgsValues string
		wantedErr        error
	}{
		{"one string nothing", "s1", "", "", nil},
		{"one string param", "s1", "v1", "v1", nil},
		{"one string option", "s1", "-s1=v1", "v1", nil},
		{"one string option alt", "s1", "-s1 v1", "v1", nil},
		{"one param and one option, option first", "s1 | s2", "-s2=v2 | s1", "s1 | s2", nil},
		{"one param and one option, param first", "s1 | s2", "s1 | -s2=v2", "s1 | s2", nil},
	}
	for _, tt := range cases {
		t.Run(fmt.Sprintf("`%s`", tt.name), func(t *testing.T) {
			task := task{name: tt.name}
			for _, n := range strings.Split(tt.givenTaskArgs, " | ") {
				task.args = append(task.args, arg{name: n, typeof: "string"})
			}
			err := processArgs(&task, strings.Split(tt.givenOsArgs, " | "))
			for i := range strings.Split(tt.wantedArgsValues, " | ") {
				v := reflect.ValueOf(tt.wantedArgsValues[i])
				assert.Equal(t, v.Interface(), task.args[i].value.Interface())
			}
			assert.Equal(t, tt.wantedErr, err)
		})
	}
}

func TestProcessArgsCombinations(t *testing.T) {
	cases := []struct {
		name             string
		givenTaskArgs    string
		givenOsArgs      string
		wantedArgsValues string
		wantedErr        error
	}{
		{"one string nothing", "s1", "", "", nil},
		{"one string param", "s1", "v1", "v1", nil},
		{"one string option", "s1", "-s1=v1", "v1", nil},
		{"one string option alt", "s1", "-s1 v1", "v1", nil},
		{"one param and one option, option first", "s1 | s2", "-s2=v2 | s1", "s1 | s2", nil},
		{"one param and one option, param first", "s1 | s2", "s1 | -s2=v2", "s1 | s2", nil},
	}
	for _, tt := range cases {
		t.Run(fmt.Sprintf("`%s`", tt.name), func(t *testing.T) {
			task := task{name: tt.name}
			for _, n := range strings.Split(tt.givenTaskArgs, " | ") {
				task.args = append(task.args, arg{name: n, typeof: "string"})
			}
			err := processArgs(&task, strings.Split(tt.givenOsArgs, " | "))
			for i := range strings.Split(tt.wantedArgsValues, " | ") {
				v := reflect.ValueOf(tt.wantedArgsValues[i])
				assert.Equal(t, v.Interface(), task.args[i].value.Interface())
			}
			assert.Equal(t, tt.wantedErr, err)
		})
	}
}

//func TestProcessArgsParamsTypeCasting(t *testing.T) {
//	var cases = []struct {
//		name string
//		givenTaskArgs []arg
//		givenOsArgs string
//		wantedArgsValues []interface{}
//		wantedErr error
//	}{
//		{
//			"all types params",
//			[]arg{
//				{name: "str", typeof: "string"},
//				{name: "int", typeof: "int"},
//			},
//			"bar 5",
//			[]interface{}{"bar", 5},
//			nil,
//		},
//	}
//	for _, tt := range cases {
//		t.Run(fmt.Sprintf("`%s`", tt.name), func(t *testing.T) {
//			task := task{
//				name: tt.name,
//				args: tt.givenTaskArgs,
//			}
//			err := processArgs(&task, strings.Split(tt.givenOsArgs, " "))
//			for i := range tt.wantedArgsValues {
//				v := reflect.ValueOf(tt.wantedArgsValues[i])
//				assert.Equal(t, v.Interface(), task.args[i].value.Interface())
//			}
//			assert.Equal(t, tt.wantedErr, err)
//		})
//	}
//}
//
//func TestProcessArgsOptionsTypeCasting(t *testing.T) {
//	var cases = []struct {
//		name string
//		givenTaskArgs []arg
//		givenOsArgs string
//		wantedArgsValues []interface{}
//		wantedErr error
//	}{
//		{
//			"all types options",
//			[]arg{
//				{name: "str", typeof: "string"},
//				{name: "int", typeof: "int"},
//				{name: "bool", typeof: "bool"},
//			},
//			"-str=bar -int=5 -bool",
//			[]interface{}{"bar", 5, true},
//			nil,
//		},
//	}
//	for _, tt := range cases {
//		t.Run(fmt.Sprintf("`%s`", tt.name), func(t *testing.T) {
//			task := task{
//				name: tt.name,
//				args: tt.givenTaskArgs,
//			}
//			err := processArgs(&task, strings.Split(tt.givenOsArgs, " "))
//			for i := range tt.wantedArgsValues {
//				v := reflect.ValueOf(tt.wantedArgsValues[i])
//				assert.Equal(t, v.Interface(), task.args[i].value.Interface())
//			}
//			assert.Equal(t, tt.wantedErr, err)
//		})
//	}
//}