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


//go rutine for å ta i mot nye ordre i sys_dat
func recive_system_data(){
	new_system_data := <- System_data_sendToManagerCh
	//ta lure avgjørelser og vurder oppdatering av egen sys_dat

}

//go rutine for å sende sys_dat ved oppdatering og jevne mellomrom
func send_system_data(updated_system_data init.System_data){
	System_data_sendToComCh <- updated_system_data

}

//go rutine for å sende bud

//go rutine for å motta bud


//go rutine for å instruere heis om neste etasje

//go rutine for å motta kvittering i fra heis på at etasje er besøkt /ordre er utført 







