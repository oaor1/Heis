// go run makeThread.go

package main

import (
    . "fmt"     // Using '.' to avoid prefixing functions with their package names
                //   This is probably not a good idea for large projects...
    "runtime"
    "time"
)



func someGoroutine1(flag chan int) {
	i:= <- flag
	for k :=0; k < 1000000; k++{
		i++
	
	}
	flag <- i
}

func someGoroutine2(flag chan int) {
	i:= <- flag
	for j :=0; j < 1009990; j++{
		i--
	
	}
	flag <- i
}

func main() {
	var flag = make(chan int,  1)
	flag <- 0
	runtime.GOMAXPROCS(runtime.NumCPU())    // I guess this is a hint to what GOMAXPROCS does...
                                            // Try doing the exercise both with and without it!	
	go someGoroutine1(flag)                      // This spawns someGoroutine() as a goroutine
	go someGoroutine2(flag)

    // We have no way to wait for the completion of a goroutine (without additional syncronization of some sort)
    // We'll come back to using channels in Exercise 2. For now: Sleep.
    time.Sleep(100*time.Millisecond)
	i:= <- flag
    Println(i)
}
