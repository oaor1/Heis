package main

import (
	"fmt"
	"time"
	"runtime"
	"UDPCOUNT"
)

type testStruct struct {
	Tall int
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	send_ch := make(chan testStruct)
	recive_ch := make(chan testStruct)
	sisteTall := 0

	initUDP(send_ch, recive_ch)

	for{
		select{
			case update := <- recive_ch:
				Tall = slave(recive_ch)

			default:
				time.sleep(5*time.Millisecond)	
		}
	}

	time.Sleep(10*time.Second)
}

func master(send_ch chan testStruct) {
	var masterStruct testStruct
	var sisteTall int 
	fmt.Printf(i", ")
	for {
		i++
		masterStruct.Tall = i
		send_ch <- masterStruct
		fmt.Printf(masterStruct.Tall", ")
		time.Sleep(1*time.Second)
	}
}

func slave(recive_ch chan testStruct, send_ch chan testStrict, sisteTall) int {
	var reciveStruct testStruct
	for{
		reciveStruct  <- recive_ch
		sisteTall := reciveStruct.Tall
		if timer{
			master()
		}
		timer(3)
	}
	return sisteTall
}

func  timer(sek int) int{
	time.sleep(sek*time.Second)
	return 1
}