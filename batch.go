package batch

import (
	"context"
	"time"
)

// Batch represents the go-batch struct.
type Batch struct {
	opts *Options
	// Input represents the batch input channel.
	Input chan<- interface{}
	// Output represents the batch output channel.
	Output <-chan []interface{}
}

// New creates a new go-batch instance.
func New(optionFuncs ...OptionFunc) *Batch {
	opts := &Options{
		Size:    DefaultBatchSize,
		MaxWait: DefaultMaxWait,
		Context: context.Background(),
	}

	for _, optionFunc := range optionFuncs {
		optionFunc(opts)
	}

	input := make(chan interface{})
	output := make(chan []interface{})

	b := &Batch{
		opts:   opts,
		Input:  input,
		Output: output,
	}

	go b.processor(input, output)

	return b
}

func (b *Batch) processor(input chan interface{}, output chan []interface{}) {
	buffer := make([]interface{}, 0, b.opts.Size)
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

				buffer = make([]interface{}, 0, b.opts.Size)

				ticker.Reset(b.opts.MaxWait)
			}
		case <-ticker.C:
			if len(buffer) > 0 {
				output <- buffer

				buffer = make([]interface{}, 0, b.opts.Size)

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
