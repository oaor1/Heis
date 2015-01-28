package main

import (
	"net"
	"fmt"
	"time"
	//"strconv"

)

func recive(port string){
	var buf []byte = make([]byte, 1024)

	

	time.Sleep(100 * time.Millisecond)
		
        udpAddress, err := net.ResolveUDPAddr("udp4",port)
	if err != nil{
		fmt.Println("ResolveUDPAdresse failed \n", err, "\n")
		return
	}

	sock, err := net.ListenUDP("udp", udpAddress)
	if err != nil{
		fmt.Println("ListenUDP failed \n", err, "\n")
		return		
	}
	

	for {
		rlen,adresse,err := sock.ReadFromUDP(buf[:])
		if err != nil{
			fmt.Println("ReadFromUDP failed, not able to recive from", adresse, "\n")
			return
		}
		fmt.Println("Recived ", rlen, " bytes from " ,adresse,"\n")
		fmt.Println(string(buf[:]))
				
	}	
}

func main() {

        port := "129.241.187.255:30000"

	recive(port)
}


