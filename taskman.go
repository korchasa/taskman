package main

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
	file := getTasksFile()
	tasks := extractTasks(file, taskPtrs)
	for _, t := range tasks {
		log.Println("Name:", t.name)
		log.Printf("Type:%#v\n", t.fn)
	}
	if len(os.Args) < 2 {
		usage(tasks)
		os.Exit(1)
	} else {
		task, args, err := resolve(tasks, os.Args)
		if err != nil {
			log.Fatalln(err)
		}
		log.Printf("Task %s%s%s: %s%s%s\n", Ok, task.name, Reset, Info, strings.Join(args, ", "), Reset)
		err = call(task.fn, args...)
		if err != nil {
			log.Fatalln(err)
		}
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

type task struct {
	name string
	doc  string
	fn   reflect.Value
	args []arg
}

type arg struct {
	name   string
	typeof string
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
			typ := reflect.TypeOf(fptr)
			//val := reflect.ValueOf(fptr)
			log.Printf("%#v %s", typ, fn.Name.Name)
			//if typ.Name() == fn.Name.Name {
			//	t := task{
			//		name: fn.Name.Name,
			//		doc:  strings.Trim(strings.Replace(fn.Doc.Text(), fn.Name.Name, "", 1), "\n "),
			//		fn: reflect.ValueOf(fptr),
			//	}
			//	for _, a := range fn.Type.Params.List {
			//		t.args = append(t.args, arg{
			//			name:   a.Names[0].String(),
			//			typeof: fmt.Sprint(a.Type),
			//		})
			//	}
			//	tasks = append(tasks, t)
			//	break
			//}
		}
	}
	return
}

func usage(tasks []task) {
	log.Println()
	log.Printf("%sUsage:%s\n", Info, Reset)
	log.Printf("  %s [command] [arguments]\n", os.Args[0])
	log.Println()

	//maxLen := 0
	//for _, t := range tasks {
	//	if len(t.name) > maxLen {
	//		maxLen = len(t.name)
	//	}
	//}

	log.Printf("%sCommands:%s\n", Info, Reset)
	for _, t := range tasks {
		var args []string
		for _, a := range t.args {
			args = append(args, fmt.Sprintf("%s%s%s:%s", Info, a.name, Reset, a.typeof))
		}
		log.Printf("  %s%s%s  - %s [ %s ]\n", Ok, t.name, Reset, t.doc, strings.Join(args, " "))
	}
}

func resolve(tasks []task, args []string) (*task, []string, error) {
	var ct *task

	for _, t := range tasks {
		if t.name == args[1] {
			ct = &t
			break
		}
	}
	if ct == nil {
		return nil, nil, fmt.Errorf("task not found")
	}

	fs := flag.NewFlagSet(ct.name, flag.ExitOnError)
	var flags []*string
	for _, a := range ct.args {
		flags = append(flags, fs.String(a.name, "", ""))
	}

	err := fs.Parse(args[2:])
	if err != nil {
		return nil, nil, fmt.Errorf("bad command arguments: %s", err)
	}

	var callArgs []string
	for _, f := range flags {
		callArgs = append(callArgs, *f)
	}

	return ct, callArgs, nil
}

func call(f reflect.Value, params ...string) error {
	in := make([]reflect.Value, len(params))
	for k, param := range params {
		in[k] = reflect.ValueOf(param)
	}
	f.Call(in)
	return nil
}

type color string

const (
	Err   color = "\u001b[31m"
	Info  color = "\u001b[36m"
	Ok    color = "\u001b[32m"
	Reset color = "\u001b[0m"
)
