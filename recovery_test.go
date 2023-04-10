package recovery

import (
	"bytes"
	"fmt"
	"sync"
	"testing"
)

type memLogger struct {
	buf bytes.Buffer
	mu  sync.Mutex
}

func (m *memLogger) Error(msg string, args ...any) {
	m.mu.Lock()
	defer m.mu.Unlock()
	fmt.Fprintf(&m.buf, msg, args...)
}

func (m *memLogger) String() string {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.buf.String()
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

func TestRecover_logging(t *testing.T) {
	const msg = "should panic"
	var m memLogger
	Recover(func() {
		panic(msg)
	}, WithLogger(&m))
	if s := m.String(); s != msg {
		t.Errorf("Recover outputs %q; want %q", s, msg)
	}
}

func TestGo_panic(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	Go(func() {
		defer wg.Done()
		panic("intentionally panic")
	})
	wg.Wait()
}

func TestGo_logging(t *testing.T) {
	const msg = "intentionally panic"
	var (
		wg sync.WaitGroup
		m  memLogger
	)
	wg.Add(1)
	Go(func() {
		defer wg.Done()
		panic(msg)
	}, WithLogger(&m))
	wg.Wait()
	if s := m.String(); s != msg {
		t.Errorf("Go outputs %q; want %q", s, msg)
	}
}
