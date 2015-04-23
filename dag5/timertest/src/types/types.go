package types

import(
	"time"
)

var(
	MY_IP = 123
	MY_NUMBER = 0

)
const(
	
	MAX_N_ELEVATORS = 10
	N_FLOORS = 4
	N_BUTTONS = 3

	TIMEOUT = 3*time.Second


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
		IP_list				[] 		string		//
		M_handle_q			[][]	int
		M_internal_elev_out	[][]	int
		M_UpAuction_q		[]		int
		M_DownAuction_q		[]		int
	}

	Elevator_state struct{
		Direction int     //    RUNUP / RUNDOWN / STOP
		Last_floor int   //    0 - 3 
	}

	Handle_confirmation int  //manager faar bekreftelse paa etasje besok

	Update_system_data struct{
		Add_order int //1 for add and 0 for delete
		Matrix_type int //0-3, Upauction, downauction, handel, internal out
		Elevator_n int
		Floor_n int
		Direction int // 0 = up  ----  1 DOWN
	}

	Auction_data struct{
		Bid int
		Floor int
		Direction int 
		Matrix_type int
		Elevator_IP int
		Elevator_number int 
		Add int
		//Number_of_elevators int

	}
)