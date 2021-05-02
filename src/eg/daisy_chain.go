package example

import "fmt"

// DaisyChain
// https://youtu.be/f6kdp27TYZs?t=1623
func DaisyChain(left, right chan int) {
	left <- 1 + <-right
}

// DemoDaisyChain
func DemoDaisyChain() {
	const n = 10000
	leftmost := make(chan int)
	right := leftmost
	left := leftmost
	for i := 0; i < n; i++ {
		right = make(chan int)
		go DaisyChain(left, right)
		left = right
	}
	go func(c chan int) { c <- 1 }(right)
	fmt.Println(<-leftmost)
}
