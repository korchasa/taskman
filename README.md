# taskman

## Write tasks with the power of go
```go
package main

import (
	"log"
	"strconv"
	"taskman"
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

func main() {
	taskman.Run(Hello)
}

```

## List tasks

```bash
$ go build && ./taskman

Usage:
  ./taskman [command] [arguments]

Commands:
  Hello  - says Hello. Arguments: who, times
```

## Run them

```bash
go build && ./taskman Hello -who=me -times=5
Task Hello
Hello, me!
Hello, me!
Hello, me!
Hello, me!
Hello, me!
```
