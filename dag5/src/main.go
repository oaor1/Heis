package main

import (
	"./manager"
	"./elevator"
	"runtime"
//	"time"
	"fmt"
)

func main (){
	runtime.GOMAXPROCS(runtime.NumCPU())

	done := make(chan bool)

/*
	admin_snd := make (chan struct)//må endres til system_status
	admin_rec := make (chan struct)

	elev_snd := make (chan int)
	elev_snd := make (chan int)
*/
	go elevator.Run()
	go manager.Run()

	<-done
	fmt.Println("Ended")
}