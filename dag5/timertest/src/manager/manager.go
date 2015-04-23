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


package main

import (
	"../types"
	"../cost"
	"../com"
	//"../elevator"
	"../timer"
	"fmt"
	"time"
)

func main(){


	fmt.Print("-------- hallo from main --------------\n")
	go timer.Funk()
	go timer.Kanaliser()
	go timer.Decrement_and_check_handle_timers()
	go timer.Decrement_and_check_auction_timers()
	go timer.Remember_standing_bid()
	go timer.Add_handle_timer_for_new_system_data_update()
	fmt.Print("-------- oppretter tullebud \n")
	var Tullebud types.Auction_data
	Tullebud.Bid = 3
	Tullebud.Direction = 0
	Tullebud.Add = 1
	Tullebud.Floor = 2
	fmt.Print("-------- dette er tullebud: \n")
	fmt.Printf("%v  \n ",Tullebud)
	fmt.Print("-------- sender tullebud på new auction channel \n")
	fmt.Print("-------- håper vi vinner :) :) :) \n")
	timer.NewAuctionInfoToTimerCh <- Tullebud
	fmt.Print("-------- sender inn ekstern ordre \n")
	timer.NewPeripheralOrderCh <- Tullebud
	
	
	//fmt.Print("-------- sender tullebud på notify winner channel \n")
	//NotifyWinningBidToManagerCh <- Tullebud
	

	time.Sleep(10*time.Second)

}


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

/*
//go rutine for å ta i mot nye ordre i sys_dat
func recive_system_data_from_com(){
	new_system_data := <- System_data_sendToManagerCh
	return new_system_data
}
*/



func start_auction(external_bid types.Auction_data){ //lage to funskjoner for forskjellige triggere
	new_internal_bid := cost.Calculate_cost(System_data, external_bid)
	if new_internal_bid < external_bid.Bid{

	}
}
//rutine for å sende sys_dat ved oppdatering og jevne mellomrom
func send_system_data_to_com(){
	//com.System_data_sendToComCh <- updated_system_data
}

//rutine for å sende bud
func send_Auction_data_to_com(){
	//Auction_bid_sendToComCh <- bid_offer
}

//rutine for å sende updated system data
func send_System_data_update_to_com(){
	//Update_system_data_sendToComCh <- system_data_update
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
		if Elevator_state.Direction == types.RUNDOWN{
			for i := types.N_FLOORS; i > 0; i-- {
				if System_data.M_handle_q[1][i]==1{
					elevator.Next_floorCh <- i
				}
			}
		}else if Elevator_state.Direction == types.RUNUP{
			for i := 0; i > types.N_FLOORS; i++ {
				if System_data.M_handle_q[0][i]==1{
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
	System_data.M_internal_elev_out[1][order_to_delete] = 0   /*myElevatorNumber must fiks ------*/

}
/* kan slettes
func external_auction_monitor(){

	new_external_bid <- Auction_bid_sendToManagerCh
	start_auction(new_external_bid)
}

func local_auction_monitor(){
	new_external_bid <- Auction_bid_sendToManagerCh
	start_auction(new_external_bid)
}
*/


func manage_incomming_data(){
	for{
		select{

		
		case new_system_data := <- com.System_data_sendToManagerCh:
			fmt.Printf("Mottar system data fra com til Manager\n")
			// mottok systemdata fra com
			System_data = new_system_data
			// gjør noe fornuftig

		case new_Local_order := <- elevator.Local_orderCh:
			new_internal_bid = calculate_cost(System_data, new_Local_order)
			new_Local_order.bid = new_internal_bid
			new_Local_order.elevator_number = //_____________
			Auction_bid_sendToComCh <- new_Local_order


		case new_external_auction_data :=  <- com.Auction_bid_sendToManagerCh:
			fmt.Printf("Mottar auction data fra com til Manager\n")
			// mottok auksjonsdata fra com
			new_internal_bid = calculate_cost(System_data, new_external_auction_data)
			if new_internal_bid < new_external_auction_data.bid{
				new_external_auction_data.elevator_number = 3                   //___________    OBS!
				new_external_auction_data.bid = new_internal_bid
				Auction_bid_sendToComCh <- new_external_auction_data //er egentlig den lokale veriden
				//Jeg leder, men må vente på timeout
			}else if new_internal_bid == new_external_auction_data.bid{ //
				if new_external_auction_data.elevator_number < local_elevator_number{
					new_external_auction_data.elevator_number = 2                   //___________    OBS!
					new_external_auction_data.bid = new_internal_bid
					Auction_bid_sendToComCh <- new_external_auction_data //er egentlig den lokale veriden
					//Jeg leder, men må vente på timeout
				}
			}
			// vurderer inkommet bud mot eget bud, ekstern funksjon.

		case new_system_data_update := <- Update_System_data_sendToManagerCh: //mulig man må pressisere type her 
			fmt.Printf("Mottar update til systemdata fra com til Manager\n")
			//mottok en oppdatering som skal legges til/slettes i system data
			select{
				case Matrix_type == 0: //  UpAuction_q
					System_data.M_UpAuction_q[new_system_data_update.floor_n] = new_system_data_update.Add_order
				case Matrix_type == 1: // DownAuction_q
					System_data.M_DownAuction_q[new_system_data_update.floor_n] = new_system_data_update.Add_order
				case Matrix_type == 2: // handel_q
					System_data.M_handel_q[new_system_data_update.floor_n] = new_system_data_update.Add_order
				case Matrix_type == 3: // internal out
					System_data.M_internal_elev_out[new_system_data_update.floor_n] = 1 // andre kan aldri slette interne ut-bestillinger
		//default???	
			}
			
		}
	}
}

func Auction_round(auction_object types.Auction_data)bool{
	local_bid = calculate_cost(System_data, New_auction_data)
	var local_best_bid bool = true
	if local_bid < New_auction_data.bid{
		New_auction_data.bid = local_bid
		send_Auction_data_to_com(New_auction_data)
	}else if local_bid == New_auction_data.bid && new_auction_data.elevator_number < Local_elevator_number{
		local_best_bid = false
	}else{
		local_best_bid = false
	}
	return local_best_bid
}

/*/dette er overflodig:
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
*/



