# Go Batch

[![Build](https://github.com/RashadAnsari/go-batch/actions/workflows/main.yml/badge.svg)](https://github.com/RashadAnsari/go-batch/actions/workflows/main.yml)

A simple batching library in Golang.

## Guid

### Installation

```bash
go get github.com/RashadAnsari/go-batch/v2
```

### Example

```go
package main

import (
	"context"
	"log"
	"math/rand"
	"reflect"
	"time"

	goBatch "github.com/RashadAnsari/go-batch/v2"
)

func main() {
	ctx, canl := context.WithCancel(context.Background())

	batch := goBatch.New[int](
		goBatch.WithSize(10),
		goBatch.WithMaxWait(1*time.Second),
		goBatch.WithContext(ctx),
	)

	go func() {
		for {
			output := <-batch.Output

			log.Printf("output: %v, size: %d\n",
				output, reflect.ValueOf(output).Len())
		}
	}()

	for i := 1; i <= 100; i++ {
		batch.Input <- i
	}

	time.Sleep(1 * time.Second)

	for i := 1; i <= 100; i++ {
		batch.Input <- i

		if rand.Intn(2) == 0 {
			time.Sleep(300 * time.Millisecond)
		}
	}

	canl()
}
```
