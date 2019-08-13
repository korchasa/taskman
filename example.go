package main

import (
	"fmt"
	"log"
	"strconv"
)

// TaskHello says Hello
func TaskHello(who string, times string) {
	tt, err := strconv.Atoi(times)
	if err != nil {
		log.Fatal(err)
	}
	for tt > 0 {
		fmt.Printf("Hello, %s!", who)
		tt--
	}
}

// NotATask says "Hello, %who!", %times times
func NotATask() {
}

func main() {
	Run(TaskHello)
}