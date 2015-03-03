package main

import (
	"../elevator"
	"runtime"
	"fmt"
)

func main (){
	runtime.GOMAXPROCS(runtime.NumCPU())

	done := make(chan bool)

	go elevator.Run()

	<-done
	fmt.PrintLN("Ended")
}