package main

import (
	"github.com/bitfield/script"
	"log"
	"strconv"
	"github.com/korchasa/taskman"
)

// Hello says Hello
func Hello(who string, times string) {
	tt, err := strconv.Atoi(times)
	if err != nil {
		log.Fatal(err)
	}
	for tt > 0 {
		log.Printf("Hello, %s!\n", who)
		tt--
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
