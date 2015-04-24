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
	"../elevator"
	"../timer"
	"fmt"
	"time"
	"runtime"
)

var(
	Bullshit_incrementor int
	System_data types.System_data
	timeOutAuctionCh = make (chan bool)
	local_elevator_que [types.N_FLOORS] int
	Elevator_state types.Elevator_state

)

func main(){
	runtime.GOMAXPROCS(runtime.NumCPU())
	fmt.Print("-------- hallo from main --------------\n")
	
	elevator.Run()
	go listen_for_timeout()
	go determine_next_floor()
	go manager_listen_for_elevator()
	
	go timer.Decrement_and_check_handle_timers()
	go timer.Decrement_and_check_auction_timers()
	go timer.Add_handle_timer_for_new_system_data_update()
	go timer.Listen_for_auctiondata_from_manager()

	
	go com.Com_listen_for_manager()
	
	go elevator.DoorTimer()
	go elevator.FloorLigths()
	go elevator.Update_channels()
	go elevator.Read_order_buttons()
	go elevator.Print()
	go elevator.Idle()
	
	

	//elevator.Run()
	
	time.Sleep(100*time.Second)

}





func manager_listen_for_elevator(){
	for{
		
		select{
			case order_to_delete := <- elevator.Next_floor_doneCh:
				fmt.Printf("dette skjedde")
				System_data.M_handle_q[order_to_delete][0] = 0  // dette slukker begge lys/køer ved besøk
				System_data.M_handle_q[order_to_delete][1] = 0
				//oppdater handle q i timer og send sys dat update til alle andre heisane
				System_data.M_internal_elev_out[order_to_delete][types.MY_NUMBER] = 0 
				timer.Executed_orderCh <- order_to_delete

			case new_external_auction_data := <- elevator.External_orderCh:
				fmt.Printf("fikk nokke paa external order ch \n")
				new_local_bid := cost.Calculate_cost(System_data, new_external_auction_data)
				new_external_auction_data.Elevator_IP = types.MY_IP            
				new_external_auction_data.Bid = new_local_bid
				new_external_auction_data.Add = 1
				com.Auction_bid_sendToComCh <- new_external_auction_data
				timer.NewAuctionInfoToTimerCh <- new_external_auction_data
			case new_internal_order := <- elevator.Internal_orderCh:
				fmt.Printf("denne også \n")
				System_data.M_internal_elev_out[types.MY_NUMBER][new_internal_order]=1
				System_data.M_handle_q[types.MY_NUMBER*2][new_internal_order]=1  // dette er litt juks
				System_data.M_handle_q[types.MY_NUMBER*2+1][new_internal_order]=1
			case Elevator_state := <- elevator.States_to_managerCh:
				Bullshit_incrementor += Elevator_state.Direction
				//fmt.Printf("mottok elevator state %v \n", Elevator_state)



		} 

	time.Sleep(100*time.Millisecond)
	}	
}

func manager_listen_for_com(){
	for{
		time.Sleep(time.Millisecond*100)
		select{

		case new_system_data := <- com.System_data_sendToManagerCh:
			fmt.Printf("---Ny systemdata fra com til manager\n %v",new_system_data)


		case new_external_auction_data := <- com.Auction_bid_sendToManagerCh:
			fmt.Printf("---nytt bud fra com til manager \n sender videre til timer\n %v",new_external_auction_data)
			new_internal_bid := cost.Calculate_cost(System_data, new_external_auction_data)
			
			if new_internal_bid < new_external_auction_data.Bid{
				new_external_auction_data.Elevator_IP = types.MY_IP            
				new_external_auction_data.Bid = new_internal_bid
				new_external_auction_data.Add = 1
				com.Auction_bid_sendToComCh <- new_external_auction_data
				timer.NewAuctionInfoToTimerCh <- new_external_auction_data
				
			}else if new_internal_bid > new_external_auction_data.Bid{
				new_external_auction_data.Elevator_IP = types.MY_IP            
				new_external_auction_data.Bid = 0
				new_external_auction_data.Add = 0
				com.Auction_bid_sendToComCh <- new_external_auction_data
				timer.NewAuctionInfoToTimerCh <- new_external_auction_data

			}else if new_internal_bid == new_external_auction_data.Bid{ 

				if new_external_auction_data.Elevator_IP < types.MY_IP{
					new_external_auction_data.Elevator_number = types.MY_IP                   
					new_external_auction_data.Bid = new_internal_bid
					new_external_auction_data.Add = 1
					com.Auction_bid_sendToComCh <- new_external_auction_data
					timer.NewAuctionInfoToTimerCh <- new_external_auction_data 
				}
			}
			
		case new_system_data_update := <- com.Update_system_data_sendToManagerCh:
			fmt.Printf("---ny system data update fra com til manager\n %v",new_system_data_update)
			switch{
				case new_system_data_update.Matrix_type == 0: //  UpAuction_q
					System_data.M_UpAuction_q[new_system_data_update.Floor_n] = new_system_data_update.Add_order
				case new_system_data_update.Matrix_type == 1: // DownAuction_q
					System_data.M_DownAuction_q[new_system_data_update.Floor_n] = new_system_data_update.Add_order
				//case new_system_data_update.Matrix_type == 2: // handel_q    //usikker paa om denne funksjonaliteten er oensket 
				//	System_data.M_handle_q[new_system_data_update.Floor_n] = new_system_data_update.Add_order
				case new_system_data_update.Matrix_type == 3: // internal out
					System_data.M_internal_elev_out[new_system_data_update.Floor_n][new_system_data_update.Direction+types.MY_NUMBER*2] = 1 // andre kan aldri slette interne ut-bestillinger
			}
		}
	}
}

func listen_for_timeout(){
	for{
		time.Sleep(time.Millisecond*100)
		select{
		case peripheral_timout := <- timer.Handle_q_timeoutCh:
			//trigge ny budrunde
			fmt.Printf("---Nokken har somla vi maa trigge ny budrunde\n %v",peripheral_timout)
			new_local_bid := cost.Calculate_cost(System_data, peripheral_timout)
			peripheral_timout.Elevator_IP = types.MY_IP            
			peripheral_timout.Bid = new_local_bid
			peripheral_timout.Add = 1
			com.Auction_bid_sendToComCh <- peripheral_timout
			timer.NewAuctionInfoToTimerCh <- peripheral_timout

		case won_assignment := <- timer.NotifyWinningBidToManagerCh:
			fmt.Printf("---Vi vant budrunda, det maa vi fikse \n %v " , won_assignment)
			if won_assignment.Direction==0{
				System_data.M_handle_q[won_assignment.Floor][types.MY_NUMBER*2+won_assignment.Direction]=1

			}else{
				System_data.M_handle_q[won_assignment.Floor][types.MY_NUMBER*2+won_assignment.Direction]=1

			}
			
		default:
		}
	}
}
//timer.Standing_upwards_bid[won_assignment.Floor]


func determine_next_floor(){
// legge til metode for å endre direction i tilfelle liste i dir retning er tom
	for{
		if Elevator_state.Direction == types.RUNDOWN{				
			for i := types.N_FLOORS; i > 0; i-- {
				if System_data.M_handle_q[i][0]>0||System_data.M_internal_elev_out[types.MY_NUMBER][i]==1{						
					elevator.Next_floorCh <- i
					}
						}
				}else if Elevator_state.Direction == types.RUNUP{
				
					for j := 0; j < types.N_FLOORS; j++ {
						if System_data.M_handle_q[j][0]>0{
							elevator.Next_floorCh <- j
						}	
				}
				}else if Elevator_state.Direction == types.STOPP{
				
					for k := 0; k < types.N_FLOORS; k++ {
						if System_data.M_handle_q[k][0]>0 || System_data.M_handle_q[k][1]>0{	
						elevator.Next_floorCh <- k
						}	
					}
		}

			
		
			time.Sleep(1000*time.Millisecond)
		
		}
		
}



func start_auction(external_bid types.Auction_data){ //lage to funskjoner for forskjellige triggere
	new_internal_bid := cost.Calculate_cost(System_data, external_bid)
	if new_internal_bid < external_bid.Bid{

	}
}
//go rutine for å motta kvittering i fra heis på at etasje er besøkt /ordre er utført 
//func delete_executed_order(){
//	order_to_delete := <- elevator.Next_floor_doneCh 
//dette funker ikke helt som vi tenker at det skal, mulig å bare slukke en lampe
	  /*myElevatorNumber must fiks ------*/

//}

/*
//go rutine for å ta i mot nye ordre i sys_dat
func recive_system_data_from_com(){
	new_system_data := <- System_data_sendToManagerCh
	return new_system_data
}
*/

/*
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

/*
func manage_incomming_data(){
	for{
		select{

		case new_Local_order := <- elevator.Local_orderCh:
			new_internal_bid := cost.Calculate_cost(System_data, new_Local_order)
			new_Local_order.Bid = new_internal_bid
			new_Local_order.Elevator_number = types.MY_IP 
			com.Auction_bid_sendToComCh <- new_Local_order


		case new_external_auction_data :=  <- com.Auction_bid_sendToManagerCh:
			fmt.Printf("Mottar auction data fra com til Manager\n")
			// mottok auksjonsdata fra com
			new_internal_bid := cost.Calculate_cost(System_data, new_external_auction_data)
			if new_internal_bid < new_external_auction_data.Bid{
				new_external_auction_data.Elevator_number = 3                   //___________    OBS!
				new_external_auction_data.Bid = new_internal_bid
				com.Auction_bid_sendToComCh <- new_external_auction_data //er egentlig den lokale veriden
				//Jeg leder, men må vente på timeout
			}else if new_internal_bid == new_external_auction_data.Bid{ //
				if new_external_auction_data.Elevator_number < types.MY_IP{
					new_external_auction_data.Elevator_number = 2                   //___________    OBS!
					new_external_auction_data.Bid = new_internal_bid
					com.Auction_bid_sendToComCh <- new_external_auction_data //er egentlig den lokale veriden
					//Jeg leder, men må vente på timeout
				}
			}
			// vurderer inkommet bud mot eget bud, ekstern funksjon.

		case new_system_data_update := <- com.Update_system_data_sendToManagerCh: //mulig man må pressisere type her 
			fmt.Printf("Mottar update til systemdata fra com til Manager\n")
			//mottok en oppdatering som skal legges til/slettes i system data
			switch{
				case new_system_data_update.Matrix_type == 0: //  UpAuction_q
					System_data.M_UpAuction_q[new_system_data_update.Floor_n] = new_system_data_update.Add_order
				case new_system_data_update.Matrix_type == 1: // DownAuction_q
					System_data.M_DownAuction_q[new_system_data_update.Floor_n] = new_system_data_update.Add_order
				//case new_system_data_update.Matrix_type == 2: // handel_q    //usikker paa om denne funksjonaliteten er oensket 
				//	System_data.M_handle_q[new_system_data_update.Floor_n] = new_system_data_update.Add_order
				case new_system_data_update.Matrix_type == 3: // internal out
					System_data.M_internal_elev_out[new_system_data_update.Floor_n][new_system_data_update.Direction] = 1 // andre kan aldri slette interne ut-bestillinger
		//default???	
			}
			
		}
	}
}
*/
/*
func Auction_round(New_auction_data types.Auction_data)bool{
	local_bid := cost.Calculate_cost(System_data, New_auction_data)
	var local_best_bid bool = true
	if local_bid < New_auction_data.Bid{
		New_auction_data.Bid = local_bid
		com.Auction_bid_sendToComCh <- New_auction_data
		//send_Auction_data_to_com(New_auction_data)  //droppa funksjon gjorde direkte
	}else if local_bid == New_auction_data.Bid && New_auction_data.Elevator_IP < types.MY_IP{
		local_best_bid = false
	}else{
		local_best_bid = false
	}
	return local_best_bid
}
*/
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



	/*
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
	*/
	
	//fmt.Print("-------- sender tullebud på notify winner channel \n")
	//NotifyWinningBidToManagerCh <- Tullebud
	
