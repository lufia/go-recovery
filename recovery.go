// Package recovery provides some panic-safe functions.
package recovery

// Go starts a new goroutine then run f on it.
// Go is almost same as go statement except recover automatically if f panics.
func Go(f func()) {
	go Recover(f)
}

func Recover(f func()) (v any) {
	defer func() {
		if e := recover(); e != nil {
			v = e
		}
	}()
	f()
	return nil
}
