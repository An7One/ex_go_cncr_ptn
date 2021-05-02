// Reference:
//	https://youtu.be/f6kdp27TYZs
//	https://blog.golang.org/pipelines
package example

import (
	"fmt"
	"math/rand"
	"time"
)

// FanIn
// Multiplexer
func FanIn(input1, input2 <-chan string) <-chan string {
	c := make(chan string)
	go func() {
		for {
			c <- <-input1
		}
	}()
	go func() {
		for {
			c <- <-input2
		}
	}()
	return c
}

// FanIn
// https://youtu.be/f6kdp27TYZs?t=1347
func FanInSingleGoroutine(input1, input2 <-chan string) <-chan string {
	c := make(chan string)
	go func() {
		for {
			select {
			case s := <-input1:
				c <- s
			case s := <-input2:
				c <- s
			// https://youtu.be/f6kdp27TYZs?t=1358
			case <-time.After(1 * time.Second):
				fmt.Println("You are too slow...")
				return
			}
		}
	}()
	return c
}

// boring returns receive-only channel of string
// Generator: the function that returns a channel
func boring(msg string) <-chan string {
	c := make(chan string)
	go func() {
		for i := 0; ; i++ {
			c <- fmt.Sprintf("%s %d", msg, i)
			time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
		}
	}()
	// to return the channel to the caller
	return c
}

func DemoFanIn() {
	c := FanIn(boring("Joe"), boring("Ann"))
	for i := 0; i < 10; i++ {
		fmt.Println(<-c)
	}
	fmt.Println("You are both boring; I am leaving")
}
