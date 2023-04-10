package recovery

type Logger interface {
	Error(msg string, args ...any)
}

type options struct {
	logger Logger
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

// WithLogger sets logger will be used for logging on panic.
func WithLogger(logger Logger) Option {
	return optionApplier(func(o *options) {
		o.logger = logger
	})
}

type discardLogger struct{}

func (*discardLogger) Error(msg string, args ...any) {}

var defaultOptions = options{
	logger: &discardLogger{},
}
