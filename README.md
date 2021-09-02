# Go Batch

[![Build](https://github.com/RashadAnsari/go-batch/actions/workflows/main.yml/badge.svg)](https://github.com/RashadAnsari/go-batch/actions/workflows/main.yml)

A simple batching library in Golang.

## Guid

### Installation

```bash
go get github.com/RashadAnsari/go-batch
```

### Example

```go
package main

import (
	"log"
	"math/rand"
	"time"

	goBatch "github.com/RashadAnsari/go-batch"
)

func main() {
	batch := goBatch.New(
		goBatch.WithSize(10),
		goBatch.WithMaxWait(1*time.Second),
	)

	go func() {
		for {
			output := <-batch.Output

			log.Printf("output: %v, size: %d\n", output, len(output))
		}
	}()

	for i := 1; i <= 1000; i++ {
		batch.Input <- i
	}

	log.Println("====================")

	for i := 1; i <= 1000; i++ {
		batch.Input <- i

		if rand.Intn(2) == 0 {
			time.Sleep(300 * time.Millisecond)
		}
	}

	time.Sleep(10 * time.Second)
}
```
