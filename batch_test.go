package batch_test

import (
	"context"
	"reflect"
	"testing"
	"time"

	goBatch "github.com/RashadAnsari/go-batch/v2"
)

func TestBatchSize(t *testing.T) {
	batch := goBatch.New[int](
		goBatch.WithSize(10),
		goBatch.WithMaxWait(1*time.Minute),
	)

	go func() {
		for i := 1; i <= 100; i++ {
			batch.Input <- i
		}
	}()

	count := 0

	for i := 1; i <= 10; i++ {
		<-batch.Output

		count++
	}

	if count != 10 {
		t.Fatalf("invalid batch count: %d", count)
	}
}

func TestBatchMaxWait(t *testing.T) {
	batch := goBatch.New[int](
		goBatch.WithSize(100),
		goBatch.WithMaxWait(1*time.Second),
	)

	go func() {
		for i := 1; i <= 10; i++ {
			batch.Input <- i
		}
	}()

	output := <-batch.Output

	outputLen := reflect.ValueOf(output).Len()

	if outputLen != 10 {
		t.Fatalf("invalid batch size: %d", outputLen)
	}
}

func TestBatchClose(t *testing.T) {
	ctx, canl := context.WithCancel(context.Background())

	batch := goBatch.New[int](
		goBatch.WithSize(100),
		goBatch.WithMaxWait(100*time.Second),
		goBatch.WithContext(ctx),
	)

	for i := 1; i <= 10; i++ {
		batch.Input <- i
	}

	canl()

	output := <-batch.Output

	outputLen := reflect.ValueOf(output).Len()

	if outputLen != 10 {
		t.Fatalf("invalid batch size: %d", outputLen)
	}
}

func BenchmarkBatchSize(b *testing.B) {
	batch := goBatch.New[int](
		goBatch.WithSize(10),
		goBatch.WithMaxWait(1*time.Second),
	)

	go func() {
		for {
			<-batch.Output
		}
	}()

	for i := 0; i < b.N; i++ {
		batch.Input <- i
	}
}
