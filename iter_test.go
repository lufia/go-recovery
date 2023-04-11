package recovery

import (
	"reflect"
	"testing"
)

func TestChanReceiverRange(t *testing.T) {
	c := make(chan int)
	go func() {
		for i := 0; i < 3; i++ {
			c <- i
		}
		close(c)
	}()
	var a []int
	ChanIter(c).Range(func(i int) bool {
		a = append(a, i)
		return true
	})
	want := []int{0, 1, 2}
	if !reflect.DeepEqual(a, want) {
		t.Errorf("received %v; want %v", a, want)
	}
}

func TestChanReceiverRange_stop(t *testing.T) {
	c := make(chan int)
	go func() {
		for i := 0; i < 3; i++ {
			c <- i
		}
		close(c)
	}()
	var a []int
	ChanIter(c).Range(func(i int) bool {
		a = append(a, i)
		return i < 1
	})
	want := []int{0, 1}
	if !reflect.DeepEqual(a, want) {
		t.Errorf("received %v; want %v", a, want)
	}
}
