package main

import (
	"github.com/bitfield/script"
	"github.com/korchasa/taskman"
	"log"
)

// Hello says Hello
func Hello(who string, times *int, show *bool) {
	if !*show {
		return
	}
	for *times > 0 {
		log.Printf("Hello, %s!\n", who)
		*times--
	}
}

// Exec executes shell cmd
func Exec(cmd string) {
	p := script.Exec(cmd)
	output, err := p.String()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(output)
}

func main() {
	taskman.Run(Hello, Exec)
}
