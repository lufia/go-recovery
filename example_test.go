package recovery_test

import (
	"fmt"
	"log"

	"github.com/lufia/go-recovery"
)

type myPropagator struct {
	l *log.Logger
	i int
}

func (p *myPropagator) Propagate(v any) {
	if p.i >= 0 {
		p.l.Printf("[i=%d]: ", p.i)
	}
	p.l.Printf("recovered: %v", v)
}

var logger = &myPropagator{
	l: log.Default(),
	i: -1,
}

func parseOptions(i int) []recovery.Option {
	l := *logger
	l.i = i
	return []recovery.Option{recovery.WithPropagator(&l)}
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
		if i == 2 {
			panic("panic!")
		}
		return true
	}, recovery.WithRangeValueParser(parseOptions))
	// Output:
	// 0
	// 1
	// 2
}
