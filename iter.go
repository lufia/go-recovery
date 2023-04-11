package recovery

type receivable[T any] interface {
	~<-chan T | ~chan T
}

// ChanReceiver represents the receiver from a channel.
type ChanReceiver[C receivable[T], T any] struct {
	c       C
	options *options
}

// ChanIter returns a ChanReceiver with options.
func ChanIter[C receivable[T], T any](c C, opts ...Option) *ChanReceiver[C, T] {
	var o options
	applyOptions(&o, opts...)
	return &ChanReceiver[C, T]{
		c:       c,
		options: &o,
	}
}

type rangeOptions[T any] struct {
	valueParser func(v T) []Option
}

// RangeOption is a option for used on each loop.
type RangeOption[T any] interface {
	Apply(o *rangeOptions[T])
}

func applyRangeOptions[T any](o *rangeOptions[T], opts ...RangeOption[T]) {
	for _, opt := range opts {
		opt.Apply(o)
	}
}

type rangeOptionApplier[T any] func(o *rangeOptions[T])

func (opt rangeOptionApplier[T]) Apply(o *rangeOptions[T]) {
	opt(o)
}

// WithRangeValueParser sets f will be used for overwritting options on each loop.
func WithRangeValueParser[T any](f func(v T) []Option) RangeOption[T] {
	return rangeOptionApplier[T](func(o *rangeOptions[T]) {
		o.valueParser = f
	})
}

// Range calls f for each received items from c. If f returns false, range stops the iteration.
func (c *ChanReceiver[C, T]) Range(f func(v T) bool, opts ...RangeOption[T]) {
	o := *c.options
	for v := range c.c {
		var ro rangeOptions[T]
		applyRangeOptions(&ro, opts...)

		v, o := v, o
		if ro.valueParser != nil {
			applyOptions(&o, ro.valueParser(v)...)
		}
		var cont bool
		Do(func() {
			cont = f(v)
		})
		if !cont {
			break
		}
	}
}
