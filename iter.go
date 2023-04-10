package recovery

type receivable[T any] interface {
	~<-chan T | ~chan T
}

type ChanIter[C receivable[T], T any] struct {
	c       C
	options *options
}

func ChanRange[C receivable[T], T any](c C, opts ...Option) *ChanIter[C, T] {
	o := defaultOptions
	applyOptions(&o, opts...)
	return &ChanIter[C, T]{
		c:       c,
		options: &o,
	}
}

type rangeOptions[T any] struct {
	valueParser func(v T) []Option
}

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

func WithRangeValueParser[T any](f func(v T) []Option) RangeOption[T] {
	return rangeOptionApplier[T](func(o *rangeOptions[T]) {
		o.valueParser = f
	})
}

func (i *ChanIter[C, T]) Do(f func(v T), opts ...RangeOption[T]) {
	o := *i.options
	for v := range i.c {
		var ro rangeOptions[T]
		applyRangeOptions(&ro, opts...)

		v, o := v, o
		if ro.valueParser != nil {
			applyOptions(&o, ro.valueParser(v)...)
		}
		Do(func() {
			f(v)
		})
	}
}
