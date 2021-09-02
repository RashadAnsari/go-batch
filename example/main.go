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

	batch.Close()
}
