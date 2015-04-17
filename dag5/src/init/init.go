package init

import(
	"fmt"
	"../driver"
	"../com"
)

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