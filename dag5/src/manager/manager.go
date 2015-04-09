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
			// mottok systemdata fra com
			// gjør noe fornuftig
		case new_auction_data :=  <- System_data_sendToManagerCh:
			// mottok auksjonsdata fra com
			// vurder inkommet bud mot eget bud
		}
	}
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



