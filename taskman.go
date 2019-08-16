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
	"strings"
)

func Run(taskPtrs ...interface{}) {
	log.SetFlags(0)
	log.SetOutput(os.Stdout)
	file := getTasksFile()
	taskCandidates, err := extractTasks(file)
	if err != nil {
		log.Fatalln(err)
	}
	tasks := attachFuncPointers(taskCandidates, taskPtrs)
	if len(os.Args) < 2 {
		usage(tasks)
		os.Exit(1)
	} else {
		task, err := resolve(tasks, os.Args)
		if err != nil {
			log.Fatalln(err)
		}
		log.Printf("Task %s%s%s\n", Ok, task.name, Reset)
		task.call()
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

func extractTasks(file string) (tasks []task, err error) {
	fset := token.NewFileSet()
	tree, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
	if nil != err {
		return
	}
	for _, f := range tree.Decls {
		fn, ok := f.(*ast.FuncDecl)
		if !ok || fn.Name.Name == "main" {
			continue
		}
		t := task{
			name: fn.Name.Name,
			doc:  strings.Trim(strings.Replace(fn.Doc.Text(), fn.Name.Name, "", 1), "\n "),
		}
		for _, a := range fn.Type.Params.List {
			n := arg{
				name: a.Names[0].String(),
			}
			switch tt := a.Type.(type) {
			case *ast.Ident:
				n.typeof = fmt.Sprint(a.Type)
			default:
				err = fmt.Errorf("unsupported argument type `%s`", tt)
			}
			t.args = append(t.args, n)
		}
		tasks = append(tasks, t)
	}
	return
}

func attachFuncPointers(taskCandidates []task, fptrs []interface{}) []task {
	for i, tc := range taskCandidates {
		for _, fptr := range fptrs {
			fnc := runtime.FuncForPC(reflect.ValueOf(fptr).Pointer())
			name := strings.Split(fnc.Name(), ".")[1]
			if name == tc.name {
				taskCandidates[i].fn = reflect.ValueOf(fptr)
			}
		}
	}
	return taskCandidates
}

func usage(tasks []task) {
	log.Println()
	log.Printf("%sUsage:%s\n", Title, Reset)
	log.Printf("  %s [command] [arguments]\n", os.Args[0])
	log.Println()

	log.Printf("%sCommands:%s\n", Title, Reset)
	for _, t := range tasks {
		var args []string
		for _, a := range t.args {
			args = append(args, fmt.Sprintf("-%s%s%s=%s", Info, a.name, Reset, a.typeof))
		}
		log.Printf("  %s%s%s %s  - %s\n", Ok, t.name, Reset, strings.Join(args, " "), t.doc)
	}
}

func resolve(tasks []task, args []string) (*task, error) {
	var ct *task

	for _, t := range tasks {
		if t.name == args[1] {
			ct = &t
			break
		}
	}
	if ct == nil {
		return nil, fmt.Errorf("task not found")
	}

	fs := flag.NewFlagSet(ct.name, flag.ExitOnError)
	var flags []interface{}
	for _, a := range ct.args {
		switch a.typeof {
		case "string":
			flags = append(flags, fs.String(a.name, "", ""))
		case "int":
			flags = append(flags, fs.Int(a.name, 0, ""))
		case "bool":
			flags = append(flags, fs.Bool(a.name, false, ""))
		default:
			return nil, fmt.Errorf("unsupported task argument type `%s`", a.typeof)
		}
	}

	err := fs.Parse(args[2:])
	if err != nil {
		return nil, fmt.Errorf("bad command arguments: %s", err)
	}

	for i, f := range flags {
		switch t := f.(type) {
		case *string:
			ct.args[i].value = reflect.ValueOf(*t)
		case *int:
			ct.args[i].value = reflect.ValueOf(*t)
		case *bool:
			ct.args[i].value = reflect.ValueOf(*t)
		default:
			return nil, fmt.Errorf("unsupported task argument type `%s`", t)
		}
	}

	return ct, nil
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
