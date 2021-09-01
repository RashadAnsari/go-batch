package batch

import "time"

const (
	// DefaultBatchSize represents the default value for the size of each batch.
	DefaultBatchSize = 10
	// DefaultMaxWait represents the default value of the maximum waiting time for filling each batch.
	DefaultMaxWait = 1 * time.Second
)

// Option defines type for Options setter function.
type Option func(*Options)

// Options sets configuration options for the go-batch instance.
type Options struct {
	// Size specifies the maximum number of items that can be in a batch.
	//
	// Default: 10.
	Size int
	// MaxWait specifies the maximum waiting time for filling a batch.
	//
	// Default: 1 * time.Second.
	MaxWait time.Duration
}

// WithSize sets the maximum number of items that can be in each batch.
func WithSize(size int) Option {
	return func(opts *Options) {
		opts.Size = size
	}
}

// WithMaxWait sets the maximum waiting time for filling each batch.
func WithMaxWait(maxWait time.Duration) Option {
	return func(opts *Options) {
		opts.MaxWait = maxWait
	}
}
