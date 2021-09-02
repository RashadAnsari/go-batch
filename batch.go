package batch

import (
	"sync"
	"time"
)

// Batch represents the go-batch struct.
type Batch struct {
	mutex sync.Mutex

	opts       *Options
	buffer     []interface{}
	dropSignal chan struct{}

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

	dropSignal := make(chan struct{})
	input := make(chan interface{})
	output := make(chan []interface{})

	b := &Batch{
		opts:       opts,
		buffer:     []interface{}{},
		dropSignal: dropSignal,

		input:  input,
		output: output,

		Input:  input,
		Output: output,
	}

	go b.maxWaitHandler()
	go b.batchSizeHandler()

	return b
}

func (b *Batch) maxWaitHandler() {
	for {
		select {
		case <-time.NewTicker(b.opts.MaxWait).C:
		case <-b.dropSignal:
			continue
		}

		b.mutex.Lock()

		if len(b.buffer) > 0 {
			b.output <- b.buffer

			b.buffer = []interface{}{}
		}

		b.mutex.Unlock()
	}
}

func (b *Batch) batchSizeHandler() {
	for {
		event := <-b.input

		drop := false

		b.mutex.Lock()

		b.buffer = append(b.buffer, event)

		if len(b.buffer) == b.opts.Size {
			b.output <- b.buffer

			b.buffer = []interface{}{}

			drop = true
		}

		b.mutex.Unlock()

		if drop {
			b.dropSignal <- struct{}{}
		}
	}
}
