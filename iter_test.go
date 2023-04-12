package recovery

import (
	"reflect"
	"testing"
)

func testChanReceiverWithRangeOptions[T any](t *testing.T, do func(v T) bool, items, want []T, opts ...RangeOption[T]) {
	t.Helper()
	c := make(chan T)
	go func() {
		for _, v := range items {
			c <- v
		}
		close(c)
	}()
	a := make([]T, 0, len(items))
	ChanIter(c).Range(func(v T) bool {
		rv := do(v)
		a = append(a, v)
		return rv
	}, opts...)
	if !reflect.DeepEqual(a, want) {
		t.Errorf("received items = %v; want %v", a, want)
	}
}

func TestChanReceiverRange_noOptions(t *testing.T) {
	data := []int{0, 1, 2}
	want := []int{0, 1, 2}
	testChanReceiverWithRangeOptions(t, func(i int) bool {
		return true
	}, data, want)
}

func TestChanReceiverRange_withRangeValueParser(t *testing.T) {
	var p vPropagator
	parserOpt := WithRangeValueParser(func(v int) []Option {
		return []Option{
			WithPropagator(&p),
		}
	})

	data := []int{0, 1, 2}
	want := []int{}
	testChanReceiverWithRangeOptions(t, func(i int) bool {
		panic(i)
	}, data, want, parserOpt)

	panicks := []any{0, 1, 2}
	if a := p.Values(); !reflect.DeepEqual(a, panicks) {
		t.Errorf("propagator receives %v; want %v", a, panicks)
	}
}

func TestChanReceiverRange_stop(t *testing.T) {
	data := []int{0, 1, 2}
	want := []int{0, 1}
	testChanReceiverWithRangeOptions(t, func(i int) bool {
		return i < 1
	}, data, want)
}
