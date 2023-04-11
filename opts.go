package recovery

// Propagator is the interface that wraps Propagate method for propagating panicking value as an error on panic.
type Propagator interface {
	Propagate(v any)
}

type options struct {
	propagator Propagator
}

// Option is a option for all provided functions.
type Option interface {
	Apply(o *options)
}

func applyOptions(o *options, opts ...Option) {
	for _, opt := range opts {
		opt.Apply(o)
	}
}

type optionApplier func(o *options)

func (opt optionApplier) Apply(o *options) {
	opt(o)
}

// WithPropagator sets p will be used for propagating on panic.
func WithPropagator(p Propagator) Option {
	return optionApplier(func(o *options) {
		o.propagator = p
	})
}
