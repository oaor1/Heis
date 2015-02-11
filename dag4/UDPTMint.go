package main

import (
	"net"
	"fmt"
	"time"
	"runtime"
	"encoding/json"
	
	
)

const ( 
	CONN_PORT = "20008"
	CONN_REC = "30564"
	CONN_type = "udp"
	CONN_IP = "129.241.187.255"
	MY_IP = "129.241.187.108"
	LOCAL_IP = "127.0.0.1"
)
/*
type orderList struct{
	IPAdresses [] string
	orderList [][] int
}
*/


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
		var resUnmarshal string
		errunm := json.Unmarshal(buffer[0:rlen], &resUnmarshal)
		if errunm != nil{
			fmt.Println("resUnmarshal failed  %i \n", errunm)
			return
		}
		mellomlagring := resUnmarshal
		fmt.Printf("Dette er konvertert : "+  mellomlagring+ "\n")

		time.Sleep(time.Second)	
	}	
}

func send(){

	testpakke := "dette er en test"
	resMarshal, _ := json.Marshal(testpakke)

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

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	

	go send()

	go recive()

	time.Sleep(5*time.Second)
}