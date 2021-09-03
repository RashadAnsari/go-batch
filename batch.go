package batch

import (
	"time"
)

// Batch represents the go-batch struct.
type Batch struct {
	opts  *Options
	close chan struct{}

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
	}

	for _, optionFunc := range optionFuncs {
		optionFunc(opts)
	}

	input := make(chan interface{})
	output := make(chan []interface{})

	b := &Batch{
		opts:  opts,
		close: make(chan struct{}),

		Input:  input,
		Output: output,
	}

	go b.processor(input, output)

	return b
}

func (b *Batch) processor(input chan interface{}, output chan []interface{}) {
	buffer := make([]interface{}, 0, b.opts.Size)
	ticker := time.NewTicker(b.opts.MaxWait)

	for {
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
		case <-b.close:
			if len(buffer) > 0 {
				output <- buffer
			}

			return
		}
	}
}

// Close drops all the batch events into the output channel and does not listen to the new events again.
func (b *Batch) Close() {
	b.close <- struct{}{}
}
