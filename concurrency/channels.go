package main

import (
	"fmt"
	"sync"
	"time"
)

func writeToChannel(c chan int, x int) {
	c <- x
	close(c)
}

func printer(ch chan bool, i int) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
		}
	}()
	fmt.Println("Sleep Start: ", i)
	time.Sleep(5000000000)
	fmt.Println("Sleep end: ", i)
	ch <- true
}

// BAD CODE, just for demonstration
func main() {
	c := make(chan int, 1)

	var waitGroup sync.WaitGroup
	waitGroup.Add(1)
	go func(c chan int) {
		defer waitGroup.Done()
		time.Sleep(5000000000)
		writeToChannel(c, 10)
		fmt.Println("Exit.")
	}(c)

	// close(c) // Writing to a closed channel would cause panic
	fmt.Println("Here, we are not blocked yet")
	fmt.Println("Read:", <-c)
	q, ok := <-c
	if ok {
		fmt.Println("Channel is open!", q)
	} else {
		// In this case it will get defaulted to zero value
		fmt.Println("Channel is closed!", q)
	}

	waitGroup.Wait()

	// Unbuffered channel
	// If the capacity is zero or absent, the channel is
	// unbuffered, communication succeeds only when
	// both a sender and receiver are ready.
	// If the channel is unbuffered, the sender blocks until
	// the receiver has received the value
	var ch chan bool = make(chan bool)
	for i := 0; i < 5; i++ {
		go printer(ch, i)
	}
	fmt.Println("Out") // This would be printed

	// If the channel is unbuffered, the sender blocks until the receiver has received the value.
	// If the channel is unbuffered, the receiver blocks unitl the sender has sent a value.
	// If the channel has a buffer, the sender blocks only until the value has been copied to the buffer;
	// if the buffer is full, this means waiting until some receiver has retrieved a value.

	// Synchronous: Unbuffered channels are synchronous, which means that a goroutine sending data to one will
	// block until another goroutine is prepared to accept it. A goroutine that accepts data from an unbuffered
	// channel will similarly block until new data becomes available

	// Range on channels
	// IMPORTANT: As the channel c is not closed,
	// the range loop does not exit by its own.

	close(ch) // If we close the channel, defauly values would be printed, if other go routines do not panic
	time.Sleep(10000000000)

	// Since we sleep and there is no data to accept, this blocks the main goroutine
	n := 0
	for i := range ch {
		fmt.Println("In Loop: ", i)
		if i {
			n++
		}
		if n > 2 {
			fmt.Println("n:", n)
			// close(ch)
			break
		}
	}

	for i := 0; i < 5; i++ {
		fmt.Println(<-ch)
	}
}
