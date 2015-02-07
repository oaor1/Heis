//run go [filename.go]
//fuser -k -n protocol portno
//sambaad-sindrevh/sanntid/heis/dag3
package main


import (
	"net"
	"fmt"
	"time"
	"runtime"


)

const (

	CONN_PORT = 34933
	CONN_TYPE = "tcp"
	)


func recive(port string){
	
		
    tcpAddress, err := net.ResolveTCPAddr("tcp4",port)
	if err != nil{
		fmt.Println("ResolvetcpAdresse failed \n", err, "\n")
		return
	}
	
	TCPConn, _ := net.DialTCP(CONN_TYPE, nil , tcpAddress)
	defer TCPConn.Close()

	for {
		var buf []byte = make([]byte, 1024)

		

		rlen,err := TCPConn.Read(buf[:])
		if err != nil{
			fmt.Println("ReadFromtcp failed, not able to recive from\n")
			return
		}
		
		

		fmt.Println("Recived ", rlen, " bytes from \n")
		fmt.Println(string(buf[:]))		
	}	
}

func send(port string){
	serverAddress, err := net.ResolveTCPAddr("tcp4",port)
	if err != nil{
		fmt.Println("ResolvetcpAdresse failed \n", err, "\n")
		return
	}
	sendtcpListener, err := net.DialTCP("tcp", nil, serverAddress)
	defer sendtcpListener.Close()
	if err != nil{
		fmt.Println("Dialtcp failed \n", err, "\n")
		return
	}
	fmt.Println(serverAddress)

	for{
		var buf []byte = make([]byte, 1024)
		message := "This is what i said\n" 
		copy(buf[:],message)
		_,err := sendtcpListener.Write(buf)
		if err!= nil{
			fmt.Println("WriteTotcp failed", err, "\n")
			return
		}
		sendtcpListener.Read(buf)
		fmt.Printf("%s",buf)
		time.Sleep(2*time.Second)
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	
	recivePort := "129.241.187.136:34933"
	sendPort := "129.241.187.136:34933"
	

	//recive(recivePort)
	go send(sendPort)
	go recive(recivePort)

	time.Sleep(3*time.Second)
}


