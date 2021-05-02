// Reference:
//	https://chai2010.gitbooks.io/advanced-go-programming-book/content/ch1-basic/ch1-06-goroutine.html
package prodr_et_cons

import (
	"fmt"
	"time"
)

// Producer produces multiples of `factor`
func Producer(factor int, out chan<- int) {
	for i := 0; ; i++ {
		out <- i * factor
	}
}

// Consumer consumes the input
func Consumer(in <-chan int) {
	for v := range in {
		fmt.Println(v)
	}
}

// DemoProdrNCons demos producers and consumers concurrency pattern
func DemoProdrNCons() {
	ch := make(chan int, 64) // buffered channel

	go Producer(3, ch)
	go Producer(4, ch)
	go Consumer(ch)

	time.Sleep(5 * time.Second)
}
