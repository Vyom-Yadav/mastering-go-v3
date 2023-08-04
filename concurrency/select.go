package main

import (
	"fmt"
	"math/rand"
	"time"
)

func gen(min, max int, createNumber chan int, end chan bool) {
	time.Sleep(time.Second)
	for {
		select {
		case createNumber <- rand.Intn(max-min) + min:
		case <-end:
			fmt.Println("Ended!")
			//return
		case <-time.After(4 * time.Second): // This is also a clever default branch only
			fmt.Println("time.After()!")
			return
		default:
			fmt.Println("none of them is ready")
			return
		}
	}
}
