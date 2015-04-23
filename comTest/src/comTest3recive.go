package main

import (
	"./types"
	"net"
	"fmt"
	"time"
	"encoding/json"	
)

const ( 
	CONN_PORT = "20408"
	CONN_REC = "30564"
	CONN_type = "udp"
	BROADCAST_IP = "129.241.187.255"
)

var (
//Oppretter globale chanels for å komunisere med manager
System_data_sendToManagerCh = make(chan types.System_data)
System_data_sendToComCh = make(chan types.System_data)

Auction_bid_sendToManagerCh = make(chan types.Auction_data)
Auction_bid_sendToComCh = make(chan types.Auction_data)

Update_system_data_sendToManagerCh = make (chan types.Update_system_data)
Update_system_data_sendToComCh = make (chan types.Update_system_data)
/*
Elevator_state_to_manager = make (chan types.Elevator_state)
Handle_confirmation_to_manager = make (chan types.Handle_confirmation)
Next_floor_to_elevator = make (chan int)
*/
)

func main(){
	var Struct1  types.Update_system_data
//	var Struct2  types.System_data
//	var Struct3  types.Auction_data

	go send()
	time.Sleep(300*time.Millisecond)
	go recive()
	go print()
//	Update_system_data_sendToComCh <- Struct1
	Auction_bid_sendToComCh <- Struct3

	time.Sleep(2*time.Second)

//	Auction_bid_sendToComCh <- Struct3
//	System_data_sendToComCh <- Struct2
	Update_system_data_sendToComCh <- Struct1


	time.Sleep(10*time.Second)


}

func print(){
	for{
		select{
		case UpdateO := <-Update_system_data_sendToManagerCh:
			fmt.Printf("struct: %+v\n",UpdateO)
		case SystemO := <-System_data_sendToManagerCh:
			fmt.Printf("struct: %+v\n",SystemO)
		case AuctionO := <-Auction_bid_sendToManagerCh:
			fmt.Printf("struct: %+v\n",AuctionO)
		default:
			time.Sleep(time.Millisecond)
		}
	}

}

func recive(){	
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
		var buffer []byte = make([]byte, 256)
		//er forrige pakke ulik den forrige?
		rlen,radr,err := socket.ReadFromUDP(buffer)
		if err != nil{
			fmt.Println("ReadFromUDP failed, not able to recive from\n")
			return
		}
		fmt.Println("Recived ", rlen, " bytes from",radr)
//		fmt.Println(buffer)
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

func send(){
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
			fmt.Println("System_data: ",serverAddress)
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

//Vi må sjekke timestamp på alt som gjøres, for å være sikker på at det er siste versjon
//for å unngå å overskrive esensiell data.
