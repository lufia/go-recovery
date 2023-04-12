package recovery_test

import (
	"fmt"

	"github.com/lufia/go-recovery"
)

type myPropagator struct{}

func (p *myPropagator) Propagate(v any) {
	fmt.Println("recovered:", v)
}

var logger = &myPropagator{}

func parseOptions(i int) []recovery.Option {
	return []recovery.Option{
		recovery.WithPropagator(&myPropagator{}),
	}
}

func Example_iter() {
	c := make(chan int)
	recovery.Go(func() {
		for i := 0; i < 3; i++ {
			c <- i
		}
		close(c)
	})
	recovery.ChanIter(c).Range(func(i int) bool {
		fmt.Println(i)
		if i == 1 {
			panic("panic!")
		}
		return true
	}, recovery.WithRangeValueParser(parseOptions))
	// Output:
	// 0
	// 1
	// recovered: panic!
	// 2
}
