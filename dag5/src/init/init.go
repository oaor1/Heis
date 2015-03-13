package init

import(
	"fmt"
	"../driver"

)

const(
	MAX_N_ELEVATORS = 10
	N_FLOORS = 4
	N_BUTTONS = 3

	TIMEOUT = 3*time.Second
)

type(

	System_data struct{
		IP_list				[] 		string
		M_internal_elev_out	[][]	int
		M_auction_q			[][]	int
		M_handle_q			[][]	int
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

	//
}