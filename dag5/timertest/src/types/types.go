package types

import(
	"time"
	"strings"
	"net"
	"strconv"
)

var(
	
	MY_IP = Get_last_byte_of_local_ip()
	MY_NUMBER = 0
)

const(
	
	MAX_N_ELEVATORS = 10
	N_FLOORS = 4
	N_BUTTONS = 3

	TIMEOUT = 3*time.Second

	LOOK_FOR_FRIENDS = 500

	HANDEL_Q = 0
	INTERNAL_ELEV_OUT = 1
	UPAUCTION_Q = 3
	DOWNAUCTION_Q = 4

	RUNUP = 1
	STOPP = 0
	RUNDOWN = -1

	UP = 0
	DOWN = 1
)

type(

	System_data struct{
		IP_list				[MAX_N_ELEVATORS] 		int		//
		M_handle_q			[N_FLOORS][2*MAX_N_ELEVATORS]	int
		M_internal_elev_out	[N_FLOORS][MAX_N_ELEVATORS]	int
	}

	Elevator_state struct{
		Direction int      //    RUNUP / RUNDOWN / STOP
		Last_floor int    //    0 - 3 
		Obstruction bool 
	}

	Handle_confirmation int  //Executed order

	Update_system_data struct{
		Elevator_IP int
		Elevator_number int 
		Add_order int //1 for add / 0 for delete
		Update_type int //0-3, IP_LIST - DELETE_TIMERS  -  UPDATE_HANDLE_q - Update-internal_elev_out
		Elevator_n int
		Floor_n int
		Direction int // 0 = up  / 1 DOWN
	}

	Auction_data struct{
		Bid int
		Floor int
		Direction int 
		Matrix_type int
		Elevator_IP int
		Elevator_number int 
		Add int
	}
)

func Get_last_byte_of_local_ip()int{
	conn, _ := net.Dial("udp", "google.com:80")  
    defer conn.Close()  
    var fullIP = strings.Split(conn.LocalAddr().String(), ":")[0]
    last_byte_of_local_IP := strings.Split(fullIP, ".")[3]  
   	i, _ := strconv.Atoi(last_byte_of_local_IP)
    return i
}
