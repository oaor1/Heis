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
	CONN_REC = "35654"
	CONN_type = "tcp4"
	CONN_IP = "129.241.187.136"
	MY_IP = "129.241.187.161"
)

var laddr *net.UDPAddr



func send(){
	TCPConn, _ := net.Dial(CONN_type, CONN_IP+":"+CONN_PORT)
	

	var buffer []byte = make([]byte, 1024)
	message := "Connect to:"+MY_IP+":"+CONN_REC
	copy(buffer[:], message)
	ln ,_ := net.Listen(CONN_type,":"+CONN_REC)


	_,err := TCPConn.Write(buffer)
	//_,err := TCPConn.Write([]byte("Connect to:"+MY_IP+":"+CONN_REC))
		if err!= nil{
			fmt.Println("WriteTotcp failed", err, "\n")
			return
		}
		fmt.Printf("measaage sendt \n"	)
		time.Sleep(1*time.Second)
	

	
	CONNRec ,_:= ln.Accept()


	for{	
		CONNRec.Read(buffer[:])
		fmt.Printf("dette er den fyste meldinga : %s\n", buffer)
		message = "ka som helst \n \x00"
		copy(buffer[:], message)
		CONNRec.Write(buffer)
		time.Sleep(3*time.Second)
		CONNRec.Read(buffer[:])
		fmt.Printf("Dette er den andre medlinga: %s\n", buffer)

	}
	defer TCPConn.Close()
	defer CONNRec.Close()
}


func main() {

	fmt.Printf(CONN_IP+":"+CONN_PORT+"\n")
	runtime.GOMAXPROCS(runtime.NumCPU())
	fmt.Printf(MY_IP+"\n")

	go send()

	time.Sleep(10*time.Second)
}