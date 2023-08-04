package main

import (
	"fmt"
	"sync"
	"time"
)

var password *secret
var waitGrp sync.WaitGroup

// until all of the readers of a sync.RWMutex mutex unlock
// that mutex, you cannot lock it for writing, same vice-versa

type secret struct {
	RWM  sync.RWMutex
	pswd string
}

func change(pass string) {
	fmt.Println("Change() function")
	password.RWM.Lock()
	fmt.Println("Change() Locked")
	time.Sleep(4 * time.Second)
	password.pswd = pass
	fmt.Println("Password Changed:", pass)
	password.RWM.Unlock()
	fmt.Println("Change() UnLocked")
}

func show() {
	defer waitGrp.Done()
	password.RWM.RLock()
	fmt.Println("Show function locked!")
	time.Sleep(2 * time.Second)
	fmt.Println("Pass value:", password.pswd)
	password.RWM.RUnlock()
}

func rwMutex() {
	password = &secret{
		pswd: "myPass",
	}
	for i := 0; i < 3; i++ {
		waitGrp.Add(1)
		go show()
	}

	waitGrp.Add(1)
	go func() {
		defer waitGrp.Done()
		change("123456")
	}()

	waitGrp.Add(1)
	go func() {
		defer waitGrp.Done()
		change("54321")
	}()

	waitGrp.Wait()

	fmt.Println("Current password value:", password.pswd)
}
