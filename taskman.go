package taskman

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"reflect"
	"runtime"
	"strconv"
	"strings"
)

func Run(taskPtrs ...interface{}) {
	log.SetFlags(0)
	log.SetOutput(os.Stdout)
	file := getTasksFile()
	tasks := extractTasks(file, taskPtrs)
	if len(os.Args) < 2 {
		usage(tasks)
		os.Exit(1)
	} else {
		name := os.Args[1]
		for _, t := range tasks {
			if t.name == name {
				err := processArgs(&t, os.Args[2:])
				if err != nil {
					log.Fatalln(err)
				}
				log.Printf("Task %s%s%s\n", Ok, t.name, Reset)
				t.call()
				os.Exit(0)
			}
		}
		log.Fatalf("task %s not found", name)
	}
}

func getTasksFile() string {
	pcs := make([]uintptr, 1)
	n := runtime.Callers(3, pcs)
	if n == 0 {
		log.Fatalln("No caller")
	}
	caller := runtime.FuncForPC(pcs[0] - 1)
	if caller == nil {
		log.Fatalln("Caller is empty")
	}
	file, _ := caller.FileLine(pcs[0] - 1)
	if file == "" {
		log.Fatalln("Can't resolve caller file")
	}
	return file
}

func extractTasks(file string, fptrs []interface{}) (tasks []task) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
	if err != nil {
		log.Fatalln(err)
	}
	for _, f := range node.Decls {
		fn, ok := f.(*ast.FuncDecl)
		if !ok || fn.Name.Name == "main" {
			continue
		}
		for _, fptr := range fptrs {
			fnc := runtime.FuncForPC(reflect.ValueOf(fptr).Pointer())
			name := strings.Split(fnc.Name(), ".")[1]
			if name == fn.Name.Name {
				t := task{
					name: fn.Name.Name,
					doc:  strings.Trim(strings.Replace(fn.Doc.Text(), fn.Name.Name, "", 1), "\n "),
					fn:   reflect.ValueOf(fptr),
				}
				for _, a := range fn.Type.Params.List {
					t.args = append(t.args, arg{
						name:   a.Names[0].String(),
						typeof: fmt.Sprint(a.Type),
					})
				}
				tasks = append(tasks, t)
				break
			}
		}
	}
	return
}

func usage(tasks []task) {
	log.Println()
	log.Printf("%sUsage:%s\n", Title, Reset)
	log.Printf("  %s [command] [arguments]\n", os.Args[0])
	log.Println()
	log.Printf("%sCommands:%s\n", Title, Reset)
	for _, t := range tasks {
		var flags []string
		for _, a := range t.args {
			if '*' == a.typeof[0] {
				flags = append(flags, fmt.Sprintf("-%s%s%s=%s", Info, a.name, Reset, a.typeof))
			} else {
				flags = append(flags, fmt.Sprintf("%s%s%s:%s", Info, a.name, Reset, a.typeof))
			}
		}
		log.Printf("  %s%s%s %s  - %s\n", Ok, t.name, Reset, strings.Join(flags, " "), t.doc)
	}
}

func processArgs(task *task, args []string) error {
	fs := flag.NewFlagSet(task.name, flag.ContinueOnError)
	var flags []interface{}
	for _, a := range task.args {
		switch a.typeof {
		case "*string":
			flags = append(flags, fs.String(a.name, "", ""))
		case "*int":
			flags = append(flags, fs.Int(a.name, 0, ""))
		case "*bool":
			flags = append(flags, fs.Bool(a.name, false, ""))
		}
	}

	err := fs.Parse(moveArgsFlagsFirst(args))
	if err != nil {
		return fmt.Errorf("bad command arguments: %s", err)
	}

	for parami, param := range fs.Args() {
		switch task.args[parami].typeof {
		case "string":
			task.args[parami].value = reflect.ValueOf(param)
		case "int":
			i, err := strconv.Atoi(param)
			if err != nil {
				return fmt.Errorf("can't convert param %s to int: `%s`", task.args[parami].name, err)
			}
			task.args[parami].value = reflect.ValueOf(i)
		default:
			return fmt.Errorf("unsupported param type `%s`", task.args[parami].typeof)
		}
	}

	for i, f := range flags {
		var val reflect.Value
		switch t := f.(type) {
		case *string:
			val = reflect.ValueOf(*t)
		case *int:
			val = reflect.ValueOf(*t)
		case *bool:
			val = reflect.ValueOf(*t)
		default:
			return fmt.Errorf("unsupported option type `%s`", t)
		}
		task.args[len(fs.Args()) + i].value = val
	}
	return nil
}

func moveArgsFlagsFirst(args []string) []string {
	var flags []string
	var params []string
	for _, a := range args {
		if '-' == a[0] {
			flags = append(flags, a)
		} else {
			params = append(params, a)
		}
	}
	return append(flags, params...)
}

type task struct {
	name string
	doc  string
	fn   reflect.Value
	args []arg
}

type arg struct {
	name   string
	typeof string
	value  reflect.Value
}

func (t *task) call() {
	in := make([]reflect.Value, len(t.args))
	for k, param := range t.args {
		in[k] = param.value
	}
	t.fn.Call(in)
}

type color string

const (
	Title color = "\u001b[35m"
	Info  color = "\u001b[36m"
	Ok    color = "\u001b[32m"
	Reset color = "\u001b[0m"
)
