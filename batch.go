package batch

import (
	"time"
)

// Batch represents the go-batch struct.
type Batch struct {
	opts   *Options
	buffer []interface{}
	close  chan struct{}

	input  chan interface{}
	output chan []interface{}

	Input  chan<- interface{}
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
		opts:   opts,
		buffer: []interface{}{},
		close:  make(chan struct{}),

		input:  input,
		output: output,

		Input:  input,
		Output: output,
	}

	go b.processor()

	return b
}

func (b *Batch) processor() {
	ticker := time.NewTicker(b.opts.MaxWait)

	for {
		select {
		case event := <-b.input:
			b.buffer = append(b.buffer, event)

			if len(b.buffer) == b.opts.Size {
				b.output <- b.buffer

				b.buffer = []interface{}{}

				ticker.Reset(b.opts.MaxWait)
			}
		case <-ticker.C:
			if len(b.buffer) > 0 {
				b.output <- b.buffer

				b.buffer = []interface{}{}
			}
		case <-b.close:
			if len(b.buffer) > 0 {
				b.output <- b.buffer

				b.buffer = []interface{}{}
			}

			return
		}
	}
}

// Close drops all the batch events into the output channel and does not listen to the new events again.
func (b *Batch) Close() {
	b.close <- struct{}{}
}
