package main

import (
	"./types"
	"net"
	"fmt"
	"time"
	"encoding/json"	
)

const ( 
	CONN_PORT = "20008"
	CONN_REC = "30564"
	CONN_type = "udp"
	BROADCAST_IP = "78.91.51.255"
	RECIVE_IP = "78.91.51.196"
)

type(
	TestStruct struct{
		Value int
		Tekst string
		Liste []int
	}
)

var (
//	posission types.Elevator_state
)

func main(){
	/*
	try1 := TestStruct{
		Value: 7,
		Tekst: "haleluja, det virker",
		Liste: []int {1,2,3,4,5,7}}
*/
	posission := types.Elevator_state{
		Direction: types.RUNDOWN,
		Last_floor: 3}

//	posission.Direction = types.RUNUP
//	posission.Last_floor = 3

//	go send(try1)
	go recive()
	go send(posission)
	go recive()
	time.Sleep(40*time.Millisecond)
}

/*
var (
//Oppretter globale chanels for å komunisere med manager
System_data_sendToManagerCh = make(chan init.System_data)
System_data_sendToComCh = make(chan init.System_data)

Auction_bid_sendToManagerCh = make(chan init.Auction_data)
Auction_bid_sendToComCh = make(chan init.Auction_data)

Update_system_data_sendToManagerCh = make (chan init.Update_system_data)
Update_system_data_sendToComCh = make (chan init.Update_system_data)

)
*/

func recive(){
	var buffer []byte = make([]byte, 256)
		
    udpAddress, err := net.ResolveUDPAddr(CONN_type, ":"+CONN_REC)
	if err != nil{
		fmt.Println("ResolveUDPAdresse failed \n", err, "\n")
		return
	}

	socket, err := net.ListenUDP(CONN_type, udpAddress)
//	defer socket.Close()
	if err != nil{
		fmt.Println("ListenUDP failed \n", err, "\n")
		return		
	}

	for i := 0; i<4; i++{
		rlen,radr,err := socket.ReadFromUDP(buffer)
		if err != nil{
			fmt.Println("ReadFromUDP failed, not able to recive from\n")
			return
		}
		fmt.Println("Recived ", rlen, " bytes from",radr," \n")
//		fmt.Printf("%d  \n\n",buffer[:])
		var resUnmarshal types.Elevator_state
		errunm := json.Unmarshal(buffer[0:rlen], &resUnmarshal)
		if errunm != nil{
			fmt.Println("resUnmarshal failed  %i \n", errunm)
			return
		}
		fmt.Printf("Dette er konvertert %v  \n\n\n",resUnmarshal)

		time.Sleep(10*time.Millisecond)	
	}	
}

func send(inputStruct types.Elevator_state){

	
	resMarshal, _ := json.Marshal(inputStruct)

	serverAddress, err := net.ResolveUDPAddr(CONN_type, BROADCAST_IP+":"+CONN_REC)
	if err != nil{
		fmt.Println("ResolveUDPAdresse failed \n", err, "\n")
		return
	}

	socket, err := net.DialUDP(CONN_type, nil, serverAddress)
	defer socket.Close()
	if err != nil{
		fmt.Println("DialUDP failed \n", err, "\n")
		return
	}
	fmt.Println(serverAddress, "\n")

	for i := 0; i<4; i++{

		_,err := socket.Write(resMarshal)
		if err!= nil{
			fmt.Println("WriteToUDP failed, ", err, "\n")
			return
		}
		time.Sleep(10*time.Millisecond)
	}
}

//Vi må sjekke timestamp på alt som gjøres, for å være sikker på at det er siste versjon
//for å unngå å overskrive esensiell data.
