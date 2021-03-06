//run go [filename.go]
//fuser -k -n protocol portno
//sambaad-sindrevh/sanntid/heis/dag3
package main


import (
	"net"
	"fmt"
	"time"
	"runtime"
	//"strconv"

)

func recive(port string){
	var buf []byte = make([]byte, 1024)

	

	//time.Sleep(100 * time.Millisecond)
		
    udpAddress, err := net.ResolveUDPAddr("udp4",port)
	if err != nil{
		fmt.Println("ResolveUDPAdresse failed \n", err, "\n")
		return
	}

	sock, err := net.ListenUDP("udp", udpAddress)
	defer sock.Close()
	if err != nil{
		fmt.Println("ListenUDP failed \n", err, "\n")
		return		
	}
	

	for {
		rlen,adresse,err := sock.ReadFromUDP(buf[:])
		if err != nil{
			fmt.Println("ReadFromUDP failed, not able to recive from\n")
			return
		}
		fmt.Println("Recived ", rlen, " bytes from",adresse," \n")
		fmt.Println(string(buf[:]))
				
	}	
}

func send(port string){
	serverAddress, err := net.ResolveUDPAddr("udp4",port)
	if err != nil{
		fmt.Println("ResolveUDPAdresse failed \n", err, "\n")
		return
	}
	sendSock, err := net.DialUDP("udp", nil, serverAddress)
	defer sendSock.Close()
	if err != nil{
		fmt.Println("DialUDP failed \n", err, "\n")
		return
	}
	fmt.Println(serverAddress, "\n")

	for{
		n,err := sendSock.Write([]byte("jadda!\n"))
		if err!= nil{
			fmt.Println("WriteToUDP failed, ", err, "\n")
			return
		}
		time.Sleep(1*time.Second)
		if n==2{

		}
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	recivePort := "129.241.187.255:20015"
	sendPort := "129.241.187.136:20015"
	

	//recive(recivePort)
	go send(sendPort)
	go recive(recivePort)

	time.Sleep(3*time.Second)
}


