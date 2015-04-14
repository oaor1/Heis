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
	"../init"
	"../com"
	"../elevator"
	"fmt"
	"time"
)

var(
	System_data init.System_data
	timeOutAuctionCh = make (chan bool)
	local_elevator_que [N_FLLORS] int
)

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

//go rutine for å sende sys_dat ved oppdatering og jevne mellomrom
func send_system_data_to_com(updated_system_data init.System_data){
	System_data_sendToComCh <- updated_system_data
}

//go rutine for å sende bud
func send_Auction_data_to_com(bid_offer init.Auction_data){
	Auction_bid_sendToComCh <- bid_offer
}

//go rutine for å sende updated system data
func send_System_data_update_to_com(system_data_update init.Auction_data){
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

//go rutine for å motta kvittering i fra heis på at etasje er besøkt /ordre er utført 



func manage_incomming_data(){
	for{
		select{
		case new_system_data := <- System_data_sendToManagerCh:
			fmt.printf ("Mottar system data fra com til Manager\n")
			// mottok systemdata fra com
			System_data = new_system_data
			// gjør noe fornuftig
		case New_auction_data :=  <- System_data_sendToManagerCh:
			fmt.printf ("Mottar auction data fra com til Manager\n")
			// mottok auksjonsdata fra com
			// vurderer inkommet bud mot eget bud, ekstern funksjon.

		case New_system_data_update := <- Update_System_data_sendToManagerCh:
			fmt.printf ("Mottar update til systemdata fra com til Manager\n")
			//mottok en oppdatering som skal legges til/slettes i system data
			
		}
	}
}

func Auction_round(auction_object init.Auction_data) bool{
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

//dette er overflødig:
func manage_outgoing_data(){
	for{
		select{
		case  <- System_data_sendToComCh:
			// registrerer ny versjon av system_data
			// send til com

		case <- Auction_data_sendToComCh:
			// registrerer utgående data på auction channel
			// send bud til com
		}
	}
}



