package recovery

import (
	"reflect"
	"sync"
	"testing"
)

type vPropagator struct {
	a  []any
	mu sync.Mutex
}

func (p *vPropagator) Propagate(v any) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.a = append(p.a, v)
}

func (p *vPropagator) Values() []any {
	p.mu.Lock()
	defer p.mu.Unlock()
	a := make([]any, len(p.a))
	copy(a, p.a)
	return a
}

func TestRecover(t *testing.T) {
	v := Recover(func() {
	})
	if v != nil {
		t.Errorf("Recover() = %v; want <nil>", v)
	}
}

func TestRecover_panic(t *testing.T) {
	msg := "should panic"
	v := Recover(func() {
		panic(msg)
	})
	if v == nil {
		t.Errorf("Recover() = %v; want panic(%q)", v, msg)
	} else if s, ok := v.(string); !ok || s != msg {
		t.Errorf("Recover() = %v; want panic(%q)", v, msg)
	}
}

func goSync(f func(), opts ...Option) {
	var wg sync.WaitGroup
	wg.Add(1)
	Go(func() {
		defer wg.Done()
		f()
	}, opts...)
	wg.Wait()
}

func testWithOptions(t *testing.T, check func(t testing.TB, do func(opts ...Option)), v any) {
	t.Helper()

	tests := map[string]func(opts ...Option){
		"Do": func(opts ...Option) {
			Do(func() { panic(v) }, opts...)
		},
		"Go": func(opts ...Option) {
			goSync(func() { panic(v) }, opts...)
		},
		"Recover": func(opts ...Option) {
			Recover(func() { panic(v) }, opts...)
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			check(t, tt)
		})
	}
}

func TestNoOptions(t *testing.T) {
	testWithOptions(t, func(t testing.TB, do func(opts ...Option)) {
		defer func() {
			if e := recover(); e != nil {
				t.Errorf("should recover in the wrapper function")
			}
		}()
		do()
	}, nil)
}

func TestWithPropagator(t *testing.T) {
	const msg = "intentionally panic"
	want := []any{msg}
	testWithOptions(t, func(t testing.TB, do func(opts ...Option)) {
		var p vPropagator
		do(WithPropagator(&p))
		if a := p.Values(); !reflect.DeepEqual(a, want) {
			t.Errorf("propagator receives %v; want %v", a, want)
		}
	}, msg)
}
