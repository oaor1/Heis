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
	CONN_PORT = "34933"
	CONN_type = "tcp4"
	CONN_IP = "129.241.187.136"
	MY_IP = "129.241.187.148"
)

var laddr *net.UDPAddr

func recive(){
	TCPConn, _ := net.Dial(CONN_type, CONN_IP+":"+CONN_PORT)
	defer TCPConn.Close()

	for {
		var buffer []byte = make([]byte, 1024)

		rlen,err := TCPConn.Read(buffer[:])
		if err != nil{
			fmt.Println("ReadFromtcp failed, not able to recive from\n")
			return
		}
		fmt.Println("Recived ", rlen, " bytes from :")
		fmt.Println(string(buffer[:]),"\n")
	}
}

func send(){
	TCPConn, _ := net.Dial(CONN_type, CONN_IP+":"+CONN_PORT)
	

	var buffer []byte = make([]byte, 1024)
	message := "This is still what i said, Gonorea\x00"
	copy(buffer[:], message)



	_,err := TCPConn.Write(buffer)
		if err!= nil{
			fmt.Println("WriteTotcp failed", err, "\n")
			return
		}
		fmt.Printf("measaage sendt"	)
		time.Sleep(1*time.Second)
		
	TCPConn.Read(buffer[:])
	fmt.Printf("%s", buffer)
	time.Sleep(1*time.Second)
	TCPConn.Read(buffer[:])
	fmt.Printf("%s", buffer)
	defer TCPConn.Close()

}


func main() {

	fmt.Printf(CONN_IP+":"+CONN_PORT+"\n")
	runtime.GOMAXPROCS(runtime.NumCPU())
	fmt.Printf(MY_IP+"\n")

	go send()

	time.Sleep(10*time.Second)
}