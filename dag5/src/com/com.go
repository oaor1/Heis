package com

import (
	"../init"
	"net"
	"fmt"
	"time"
	"encoding/json"	
)

const ( 
	CONN_PORT = "20008"
	CONN_REC = "30564"
	CONN_type = "udp"
	BROADCAST_IP = "129.241.187.255"
)

var (
//Oppretter globale chanels for å komunisere med manager
System_data_sendToManagerCh = make(chan init.System_data)
System_data_recivefromManagerCh = make(chan init.System_data)

Auction_bid_sendToManagerCh = make(chan init.Auction_data)
Auction_bid_reciveFromManagerCh = make(chan init.Auction_data)
)

func recive(){
	var buffer []byte = make([]byte, 256)
		
    udpAddress, err := net.ResolveUDPAddr(CONN_type, CONN_IP+":"+CONN_REC)
	if err != nil{
		fmt.Println("ResolveUDPAdresse failed \n", err, "\n")
		return
	}

	socket, err := net.ListenUDP(CONN_type, udpAddress)
	defer socket.Close()
	if err != nil{
		fmt.Println("ListenUDP failed \n", err, "\n")
		return		
	}
	
	for {
		rlen,radr,err := socket.ReadFromUDP(buffer)
		if err != nil{
			fmt.Println("ReadFromUDP failed, not able to recive from\n")
			return
		}
		fmt.Println("Recived ", rlen, " bytes from",radr," \n")
		fmt.Printf("%d  \n\n",buffer[:])
		var resUnmarshal testStruct
		errunm := json.Unmarshal(buffer[0:rlen], &resUnmarshal)
		if errunm != nil{
			fmt.Println("resUnmarshal failed  %i \n", errunm)
			return
		}
		mellomlagring := resUnmarshal.Tall
		fmt.Printf("Dette er konvertert %d  \n\n\n", mellomlagring)

		time.Sleep(time.Second)	
	}	
}

func send(inputStruct testStruct){

	
	resMarshal, _ := json.Marshal(inputStruct)

	serverAddress, err := net.ResolveUDPAddr(CONN_type, CONN_IP+":"+CONN_REC)
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

	for{

		_,err := socket.Write(resMarshal)
		if err!= nil{
			fmt.Println("WriteToUDP failed, ", err, "\n")
			return
		}
		time.Sleep(1*time.Second)
	}
}

//Vi må sjekke timestamp på alt som gjøres, for å være sikker på at det er siste versjon
//for å unngå å overskrive esensiell data.
