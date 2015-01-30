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

const (
	//CONN_HOST = addr()
	CONN_PORT = 34933
	CONN_TYPE = "tcp"
	)


func recive(port string){
	//time.Sleep(100 * time.Millisecond)
		
    tcpAddress, err := net.ResolveTCPAddr("tcp4",port)
	if err != nil{
		fmt.Println("ResolvetcpAdresse failed \n", err, "\n")
		return
	}
	
	TCPConn, _ := net.DialTCP(CONN_TYPE, nil , tcpAddress)
	/*
	tcpListener, err := net.ListenTCP("tcp", tcpAddress)
	defer tcpListener.Close()
	if err != nil{
		fmt.Println("ListenTCP failed \n", err, "\n")
		return		
	}
	*/

	for {
		var buf []byte = make([]byte, 1024)

		/*
		conn, err := tcpListener.Accept()
		if err != nil{
			fmt.Println("tcpListener.accept failed \n", err, "\n")
			return
		}
		*/

		rlen,err := TCPConn.Read(buf[:])
		if err != nil{
			fmt.Println("ReadFromtcp failed, not able to recive from\n")
			return
		}
		
		//conn.Write([]byte("message recived."))

		//conn.Close()

		fmt.Println("Recived ", rlen, " bytes from \n")
		fmt.Println(string(buf[:]))		
	}	
}
/*
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
	fmt.Println(serverAddress, "\n")

	for{
		_,err := sendtcpListener.Write([]byte("jadda!\n"))
		if err!= nil{
			fmt.Println("WriteTotcp failed", err, "\n")
			return
		}
		time.Sleep(1*time.Second)
	}
}
*/
func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	
	recivePort := "129.241.187.136:34933"
	//sendPort := "129.241.187.136:34933"
	

	//recive(recivePort)
	//go send(sendPort)
	go recive(recivePort)

	time.Sleep(3*time.Second)
}


