package com

import (

	"../types"
	"net"
	"fmt"
	"time"
	"encoding/json"	
)

const (

	CONN_PORT = "20008"
	CONN_REC = "30564"
	CONN_REC_SYSTEM_DATA = "33445"
	CONN_type = "udp"
	BROADCAST_IP = "129.241.187.255"
)

var (

	Looking_for_other_elevators_on_network = true

	System_data_sendToManagerCh = make(chan types.System_data, 10)
	System_data_sendToComCh = make(chan types.System_data)
	
	Dedicated_system_data_sendToComCh = make(chan types.System_data)
	Dedicated_system_data_sendToManagerCh = make(chan types.System_data, 10)

	Auction_bid_sendToManagerCh = make(chan types.Auction_data, 10)
	Auction_bid_sendToComCh = make(chan types.Auction_data)

	Update_system_data_sendToManagerCh = make (chan types.Update_system_data, 10)
	Update_system_data_sendToComCh = make (chan types.Update_system_data)
)

func Run(){
	go Send_system_data()
	go Send()
	go Recive()

}

func Recive(){
    udpAddress, err := net.ResolveUDPAddr(CONN_type, BROADCAST_IP+":"+CONN_REC)
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
		var buffer []byte = make([]byte, 512)
		//er forrige pakke ulik den forrige?
		rlen,_,err := socket.ReadFromUDP(buffer)
		if err != nil{
			fmt.Println("ReadFromUDP failed, not able to recive from\n")
			return
		}
		switch {
    	case buffer[0] == 0:
    		var resUnmarshal types.Auction_data
			errunm := json.Unmarshal(buffer[1:rlen], &resUnmarshal)
			if errunm != nil{
				fmt.Printf("resUnmarshal failed  %i \n", errunm)
				return
			}
			Auction_bid_sendToManagerCh <- resUnmarshal
    	case buffer[0] == 1:
    		var resUnmarshal types.Update_system_data
			errunm := json.Unmarshal(buffer[1:rlen], &resUnmarshal)
			if errunm != nil{
				fmt.Printf("resUnmarshal failed  %i \n", errunm)
				return
			}
			Update_system_data_sendToManagerCh <- resUnmarshal
    	case buffer[0] == 2:
    		var resUnmarshal types.System_data
			errunm := json.Unmarshal(buffer[1:rlen], &resUnmarshal)
			if errunm != nil{
				fmt.Printf("resUnmarshal failed  %i \n", errunm)
				return
			}
			System_data_sendToManagerCh <- resUnmarshal
		default:
			fmt.Println("unknown type of struct")
		}
		time.Sleep(time.Millisecond)	
	}	
}

func Send(){
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
	for {
		select{
		case toMarshalAuction := <- Auction_bid_sendToComCh:
			resMarshal, _ := json.Marshal(toMarshalAuction)
			var buffer []byte = make([]byte, len(resMarshal)+1)
			buffer [0] = 0
			for i := 0; i < len(resMarshal); i++ {
				buffer [i+1] = resMarshal [i]
			}
			fmt.Println("Auction: ",serverAddress)
			for i := 0; i<1; i++{
				_,err := socket.Write(buffer)
				if err!= nil{
					fmt.Println("WriteToUDP failed, ", err, "\n")
					return
				}
			time.Sleep(5*time.Millisecond)
			}
		case toMarshalUpdate_system_data := <- Update_system_data_sendToComCh:
			resMarshal, _ := json.Marshal(toMarshalUpdate_system_data)
			var buffer []byte = make([]byte, len(resMarshal)+1)
			buffer [0] = 1
			for i := 0; i < len(resMarshal); i++ {
				buffer [i+1] = resMarshal [i]
			}
			fmt.Println("update: ",serverAddress)
			for i := 0; i<1; i++{
				_,err := socket.Write(buffer)
				if err!= nil{
					fmt.Println("WriteToUDP failed, ", err, "\n")
					return
				}
			time.Sleep(5*time.Millisecond)
			}
		case toMarshalSystem_data := <- System_data_sendToComCh:
			resMarshal, _ := json.Marshal(toMarshalSystem_data)
			var buffer []byte = make([]byte, len(resMarshal)+1)
			buffer [0] = 2
			for i := 0; i < len(resMarshal); i++ {
				buffer [i+1] = resMarshal [i]
			}
			//fmt.Println("System_data: ",serverAddress)
			for i := 0; i<1; i++{
				_,err := socket.Write(buffer)
				if err!= nil{
					fmt.Println("WriteToUDP failed, ", err, "\n")
					return
				}
			time.Sleep(5*time.Millisecond)
			}
		}
	}
}
func Send_system_data(){
	serverAddress, err := net.ResolveUDPAddr(CONN_type, BROADCAST_IP+":"+CONN_REC_SYSTEM_DATA)
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
	for {
		select{
		case toMarshalSystem_data := <- Dedicated_system_data_sendToComCh:
			resMarshal, _ := json.Marshal(toMarshalSystem_data)
			var buffer []byte = make([]byte, len(resMarshal)+1)
			buffer [0] = 2
			for i := 0; i < len(resMarshal); i++ {
				buffer [i+1] = resMarshal [i]
			}
			for i := 0; i<1; i++{
				_,err := socket.Write(buffer)
				if err!= nil{
					fmt.Println("WriteToUDP failed, ", err, "\n")
					return
				}
			time.Sleep(50*time.Millisecond)
			}
		}
	}
}
func Listen_for_system_data(){
	udpAddress, err := net.ResolveUDPAddr(CONN_type, BROADCAST_IP+":"+CONN_REC_SYSTEM_DATA)
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
	for Looking_for_other_elevators_on_network == true{
		
		var buffer []byte = make([]byte, 512)
		rlen,_,err := socket.ReadFromUDP(buffer)
		if err != nil{
			fmt.Println("ReadFromUDP failed, not able to recive from\n")
			return
		}
		switch {
    	case buffer[0] == 2:
    		var resUnmarshal types.System_data
			errunm := json.Unmarshal(buffer[1:rlen], &resUnmarshal)
			if errunm != nil{
				fmt.Printf("resUnmarshal failed  %i \n", errunm)
				return
			}
			Dedicated_system_data_sendToManagerCh <- resUnmarshal
		
		default:
			fmt.Println("unknown type of struct")
		}
		time.Sleep(time.Millisecond)	
	}	
}
