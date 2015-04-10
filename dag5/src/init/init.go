package init

import(
	"fmt"
	"../driver"
	"../com"
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

	UP = 1
	STOPP = 0
	DOWN = -1
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
		Add_order bool //true for add and flase for delete
		Matrix_type int //0-3, Upauction, downauction, handel, internal out
		Elevator_n int
		floor_n int
	}

	Auction_data struct{
		bid int
		Auction_object [] int // [up/down, etasje]
		elevator_number int

	}

func init(){
	go recive()//merge flere lister
	time.Sleep(20*time.millisecond)
	//Lytt etter struct, pr. tid, if not in list
		//append IP in IP_list, send struct
	//else
		//merge struct

	if driver.Elev_init() == 0{
		fmt.Println("Elev_init failed")
	}
	// identifisere heisnummer ved init kalt: Local_elevator_number
	//fra 0-N antall heiser
	//
//Spørre etter siste oppdaterte versjon av System data.
}