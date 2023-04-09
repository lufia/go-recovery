package recovery

import (
	"sync"
	"testing"
)

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

func TestGo_panic(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	Go(func() {
		defer wg.Done()
		panic("intentionally panic")
	})
	wg.Wait()
}
