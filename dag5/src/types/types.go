package types

import(
	"time"
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
)

type(

	System_data struct{
		IP_list				[] 		string		//
		M_handle_q			[][]	int
		M_internal_elev_out	[][]	int
		M_UpAuction_q		[]		int
		M_DownAuction_q		[]		int
	}

	Update_system_data struct{
		Add_order bool //true for add and false for delete
		Matrix_type int //0-3, Upauction, downauction, handel, internal out
		Elevator_n int
		floor_n int
	}

	Auction_data struct{
		bid int
		Auction_object [] int // [up/down, etasje]
		elevator_number int

	}
)