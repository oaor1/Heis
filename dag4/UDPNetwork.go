package main

import (
	"net"
	"fmt"
	"time"
	"runtime"
	
)

const ( 
	CONN_PORT = "20008"
	CONN_REC = "20010"
	CONN_type = "udp"
	CONN_IP = "129.241.187.136"
	MY_IP = "129.241.187.161"
)

func send(){
	Socket , _ := net.Dial(CONN_type, CONN_IP+":"+CONN_PORT)
	var buffer []byte = make([]byte, 256)
	copy(buffer[:], "test")
//	ln , _ := net.Listen(CONN_type,":"+CONN_REC)

	for{
		Socket.Write(buffer)
/*
	if err!=nil{
		fmt.Println("WriteToUDP failed, ", err, "\n")
	}
*/
		time.Sleep(time.Second)
		fmt.Printf("sender\n")
	}

}

func recive(){
	Socket , _ := net.Dial(CONN_type, CONN_IP+":"+CONN_PORT)
	for{
		fmt.Printf("tar i mot\n")
		var buffer []byte = make([]byte, 256)
						fmt.Printf("tar i mot\n")
		Socket.Read(buffer[:])
				fmt.Printf("tar i mot\n")
		fmt.Printf("%s\n", buffer)
		time.Sleep(1*time.Second)
	}
}
//var laddr *net.UDPAddr

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	go send()
	go recive()

	time.Sleep(10*time.Second)
}
