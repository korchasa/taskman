package taskman

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"path/filepath"
	"testing"
)

func TestExtractTasks(t *testing.T) {
	code := []byte(`
package main
// Hello says Hello
func Hello(who string, times int, show bool) {}
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
			{name: "times", typeof: "int"},
			{name: "show", typeof: "bool"},
		},
	}, tasks[0])
}
