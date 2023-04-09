package recovery_test

import (
	"fmt"
	"log"

	"github.com/lufia/go-recovery"
)

type myLogger struct {
	l *log.Logger
	i int
}

func (l *myLogger) Error(msg string, args ...any) {
	if l.i >= 0 {
		l.l.Printf("[i=%d]: ", l.i)
	}
	l.l.Printf(msg, args...)
}

var logger = &myLogger{
	l: log.Default(),
	i: -1,
}

func parseOptions(i int) []recovery.Option {
	l := *logger
	l.i = i
	return []recovery.Option{recovery.WithLogger(&l)}
}

func Example_iter() {
	c := make(chan int)
	recovery.Go(func() {
		for i := 0; i < 3; i++ {
			c <- i
		}
		close(c)
	})
	recovery.ChanRange(c).Do(func(i int) {
		fmt.Println(i)
		if i == 2 {
			panic("panic!")
		}
	}, recovery.WithRangeValueParser(parseOptions))
	// Output:
	// 0
	// 1
	// 2
}
