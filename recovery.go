// Package recovery provides some panic-safe functions.
package recovery

// Go starts a new goroutine then run f on it.
// Go is almost same as go statement except recover automatically if f panics.
func Go(f func(), opts ...Option) {
	go Recover(f, opts...)
}

// Do runs f in the current goroutine.
// When f panicked, Do recover automatically.
func Do(f func(), opts ...Option) {
	Recover(f, opts...)
}

// Recover runs f and returns an error value if it panicked.
func Recover(f func(), opts ...Option) (v any) {
	var o = defaultOptions
	applyOptions(&o, opts...)
	defer func() {
		if e := recover(); e != nil {
			o.logger.Error("%v", e)
			v = e
		}
	}()
	f()
	return nil
}
