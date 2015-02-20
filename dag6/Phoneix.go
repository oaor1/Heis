package UDPCOUNT

import (
	"net"
	"fmt"
	"time"
	"runtime"
	"encoding/json"	
)

const ( 
	CONN_PORT = "20021"
	CONN_REC = "30564"
	CONN_type = "udp"
	CONN_IP = "129.241.187.255"
	MY_IP = "129.241.187.153"
	LOCAL_IP = "127.0.0.1"
)

func initUDP(send_ch, recive_ch chan testStruct){

	//recive
    udpAddress, err := net.ResolveUDPAddr(CONN_type, CONN_IP+":"+CONN_REC)
	if err != nil{
		fmt.Println("ResolveUDPAdresse failed \n", err, "\n")
		return
	}

	rSocket, err := net.ListenUDP(CONN_type, udpAddress)
	defer rSocket.Close()
	if err != nil{
		fmt.Println("ListenUDP failed \n", err, "\n")
		return		
	}

	//Send
	serverAddress, err := net.ResolveUDPAddr(CONN_type, CONN_IP+":"+CONN_REC)
	if err != nil{
		fmt.Println("ResolveUDPAdresse failed \n", err, "\n")
		return
	}
	sSocket, err := net.DialUDP(CONN_type, nil, serverAddress)
	defer sSocket.Close()
	if err != nil{
		fmt.Println("DialUDP failed \n", err, "\n")
		return
	}
	fmt.Println(serverAddress, "\n")

	go recive(rSocket, recive_ch)
	go send (sSocket, send_ch)
}

func recive(rSocket *net.UDPConn, recive_ch chan testStruct){

	var buffer []byte = make([]byte, 32)
	var resUnmarshal testStruct
	for{
		rlen,radr,err := rSocket.ReadFromUDP(buffer)
		if err != nil{
			fmt.Println("ReadFromUDP failed, not able to recive from\n")
			return
		}
		fmt.Println("Recived ", rlen, " bytes from",radr," \n")
		
		errunm := json.Unmarshal(buffer[0:rlen], &resUnmarshal)
		if errunm != nil{
			fmt.Println("resUnmarshal failed  %i \n", errunm)
			return
		}
		recive_ch <- testStruct{tall: resUnmarshal.Tall}
	}
}

func send(sSocket *net.UDPConn, send_ch chan testStruct){
	
	midlertidigStruct := <- send_ch
	resMarshal, _ := json.Marshal(midlertidigStruct)

	_,err := sSocket.Write(resMarshal)
	if err!= nil{
		fmt.Println("WriteToUDP failed, ", err, "\n")
		return
	}
	time.Sleep(1*time.Second)
}


func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	resI := testStruct{
		Tall: 30}


	go send(resI)

	go recive()

	time.Sleep(5*time.Second)
}