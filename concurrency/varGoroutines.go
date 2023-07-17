package main

import (
	"fmt"
	"os"
	"strconv"
	"sync"
)

func main2() {
	count := 10
	arguments := os.Args
	if len(arguments) == 2 {
		t, err := strconv.Atoi(arguments[1])
		if err == nil {
			count = t
		} else {
			fmt.Printf(arguments[1], ": is not a number")
		}
	}

	fmt.Printf("Going to create %d goroutines.\n", count)

	var waitGroup sync.WaitGroup
	fmt.Printf("%#v\n", waitGroup)
	for i := 0; i < count; i++ {
		waitGroup.Add(1)
		go func(x int) {
			defer waitGroup.Done()
			fmt.Println(x)
		}(i)
	}

	fmt.Printf("%#v\n", waitGroup)
	waitGroup.Wait()
	fmt.Println("\nExiting...")
}
