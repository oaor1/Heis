//insturere elevator om hvilken etasje som er den neste
//kalkulere og formidle kostanden og bud
//Avgjøre neste etasje
/*
Manager 
	
	kalkulerer og formidler kost til com.go
	mottar bud fra com og avgjør om budrunden ble vunnet
	avgjør neste etasje
	instruerer elevator
*/


package manager

import (
	"../types"
	"../cost"
	"../com"
	"../elevator"
	"fmt"
	"time"
)

var(
	System_data types.System_data
	timeOutAuctionCh = make (chan bool)
	local_elevator_que [types.N_FLOORS] int
	Elevator_state types.Elevator_state
)

func Run(){
	done := make (chan bool)
	go send_system_data_to_com()
	go send_Auction_data_to_com()
	go send_System_data_update_to_com()
	go manage_incomming_data()
	<- done

}

func timeOutAuction(){
	time.Sleep(40*time.Millisecond)
	timeOutAuctionCh <- false
}
/*
//go rutine for å ta i mot nye ordre i sys_dat
func recive_system_data_from_com(){
	new_system_data := <- System_data_sendToManagerCh
	return new_system_data
}
*/
func start_auction(){
	//trigger ved pakke fra auction channel
	new_external_bid <- Auction_bid_sendToManagerCh
	new_internal_bid = calculate_cost(System_data, new_external_bid)
}
//rutine for å sende sys_dat ved oppdatering og jevne mellomrom
func send_system_data_to_com(updated_system_data types.System_data){
	System_data_sendToComCh <- updated_system_data
}

//rutine for å sende bud
func send_Auction_data_to_com(bid_offer types.Auction_data){
	Auction_bid_sendToComCh <- bid_offer
}

//rutine for å sende updated system data
func send_System_data_update_to_com(system_data_update types.Auction_data){
	Update_system_data_sendToComCh <- system_data_update
}

/*
//go rutine for å motta bud
func recive_bid_from_com(){
	new_auction_data :=  <- System_data_sendToManagerCh
	return new_auction_data
}
*/

//go rutine for å instruere heis om neste etasje

/*
monitor the state of the elevator to determine optimal behavior

Heis bør sende oppdatering av tilstand kontinuerlig.

aldri endre retning før man har tatt øverste/nederste bestilling.


*/
func determine_next_floor() int{
// legge til metode for å endre direction i tilfelle liste i dir retning er tom
		if Elevator_state.Direction == RUNDOWN{
			for i := N_FLOORS; i > 0; i-- {
				if system_data.M_handel_q[1][i]==1{
					elevator.Next_floorCh <- i
				}
			}
		}
		else if Elevator_state.Direction == RUNUP{
			for i := 0; i > N_FLOORS; i++ {
				if system_data.M_handel_q[0][i]==1{
					elevator.Next_floorCh <- i
				}	
			}
		}
}

//go rutine for å motta kvittering i fra heis på at etasje er besøkt /ordre er utført 
func execute_order(){
	order_to_delete := <- Next_floor_doneCh 
//dette funker ikke helt som vi tenker at det skal, hais kan skifte retning
	System_data.M_handle_q[0][order_to_delete] = 0
	System_data.M_handle_q[1][order_to_delete] = 0
	System_data.M_internal_elev_out[/*myElevatorNumber*/][order_to_delete] = 0

}


func manage_incomming_data(){
	for{
		select{
		case new_system_data := <- System_data_sendToManagerCh:
			fmt.printf ("Mottar system data fra com til Manager\n")
			// mottok systemdata fra com
			System_data = new_system_data
			// gjør noe fornuftig
		case New_auction_data :=  <- Auction_data_sendToManagerCh:
			fmt.printf ("Mottar auction data fra com til Manager\n")
			// mottok auksjonsdata fra com

			// vurderer inkommet bud mot eget bud, ekstern funksjon.

		case New_system_data_update := <- Update_System_data_sendToManagerCh: //mulig man må pressisere type her 
			fmt.printf ("Mottar update til systemdata fra com til Manager\n")
			//mottok en oppdatering som skal legges til/slettes i system data
			select{
				case Matrix_type == 0: //  UpAuction_q
					System_data.M_UpAuction_q[New_system_data_update.floor_n] = New_system_data_update.Add_order
				case Matrix_type == 1: // DownAuction_q
					System_data.M_DownAuction_q[New_system_data_update.floor_n] = New_system_data_update.Add_order
				case Matrix_type == 2: // handel_q
					System_data.M_handel_q[New_system_data_update.floor_n] = New_system_data_update.Add_order
				case Matrix_type == 3: // internal out
					System_data.M_internal_elev_out[New_system_data_update.floor_n] = 1 // andre kan aldri slette interne ut-bestillinger
			}
			
		}
	}
}

func Auction_round(auction_object types.Auction_data){
	local_bid = calculate_cost(System_data, New_auction_data)
	local_best_bid bool = true
	if local_bid < New_auction_data.bid{
		New_auction_data.bid = local_bid
		send_Auction_data_to_com(New_auction_data)
	}
	else if local_bid == New_auction_data.bid && new_auction_data.elevator_number < Local_elevator_number
		local_best_bid = false
	}
	else{
		local_best_bid = false
	} 
	return local_best_bid
}

/*/dette er overflødig:
func manage_outgoing_data(){
	for{
		select{
		case  <- System_data_sendToComCh:
			// registrerer ny versjon av system_data
			// send til com

		case <- Auction_data_sendToComCh:
			// registrerer utgående data på auction channel
			// send bud til com

		case <- 
		}
	}
}
/*



