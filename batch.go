package batch

import (
	"context"
	"time"
)

// Batch represents the go-batch struct.
type Batch[T any] struct {
	opts *Options
	// Input represents the batch input channel.
	Input chan<- T
	// Output represents the batch output channel.
	Output <-chan []T
}

// New creates a new go-batch instance.
func New[T any](optionFuncs ...OptionFunc) *Batch[T] {
	opts := &Options{
		Size:    DefaultBatchSize,
		MaxWait: DefaultMaxWait,
		Context: context.Background(),
	}

	for _, optionFunc := range optionFuncs {
		optionFunc(opts)
	}

	input := make(chan T)
	output := make(chan []T)

	oneWayInput := make(chan<- T)
	oneWayOutput := make(<-chan []T)

	oneWayInput = input
	oneWayOutput = output

	b := &Batch[T]{
		opts:   opts,
		Input:  oneWayInput,
		Output: oneWayOutput,
	}

	go b.processor(input, output)

	return b
}

func (b *Batch[T]) processor(input chan T, output chan []T) {
	buffer := make([]T, 0, b.opts.Size)
	ticker := time.NewTicker(b.opts.MaxWait)

Loop:
	for {
		if b.opts.Context.Err() != nil {
			break
		}

		select {
		case event := <-input:
			buffer = append(buffer, event)

			if len(buffer) == b.opts.Size {
				output <- buffer

				buffer = make([]T, 0, b.opts.Size)

				ticker.Reset(b.opts.MaxWait)
			}
		case <-ticker.C:
			if len(buffer) > 0 {
				output <- buffer

				buffer = make([]T, 0, b.opts.Size)

				ticker.Reset(b.opts.MaxWait)
			}
		case <-b.opts.Context.Done():
			break Loop
		}
	}

	if len(buffer) > 0 {
		output <- buffer
	}

	close(output)
}
