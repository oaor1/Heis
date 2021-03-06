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
	initialized = false
	Bullshit_incrementor int
	System_data types.System_data
	timeOutAuctionCh = make (chan bool)
	no_one_is_aliveCH = make (chan bool)
	local_elevator_que [types.N_FLOORS] int
	Elevator_state types.Elevator_state

)

func main(){
	fmt.Print("-------- hallo from main --------------\n")
	done := make(chan bool)
	go timeout_check_for_life_on_network()
	go com.Listen_for_system_data()
	go init_manager()

	time.Sleep(types.LOOK_FOR_FRIENDS*2*time.Millisecond)
	
	runtime.GOMAXPROCS(runtime.NumCPU())

	fmt.Printf("this is my ip: %d\n", types.MY_IP)  
	fmt.Printf("this is my elevator number :%d\n", types.MY_NUMBER)
	
	go elevator.Run()
	
	go send_system_data_to_com()
	go manager_listen_for_elevator()
	go listen_for_timeout()
	go determine_next_floor()
	go manager_listen_for_com()
	go Update_button_lamps()
	
	go timer.Run()
	go com.Run()
	<-done
}

func timeout_check_for_life_on_network(){
	time.Sleep(types.LOOK_FOR_FRIENDS*time.Millisecond)
	no_one_is_aliveCH <- true
}

func send_system_data_to_com() {
	for{
		//fmt.Printf("this is my number %d \n", types.MY_NUMBER)
		//fmt.Printf("this is my ip %d \n", types.MY_IP)
		com.Dedicated_system_data_sendToComCh <- System_data
		time.Sleep(500*time.Millisecond)
	}
}

func init_manager(){
	fmt.Printf("initializing")
	for initialized==false{
		select{

		case load_system_data := <- com.Dedicated_system_data_sendToManagerCh:
			fmt.Printf("Init as one of several elevators \n")
			for i := 1; i < types.MAX_N_ELEVATORS; i++ {
				if load_system_data.IP_list[i]==types.MY_IP{
					types.MY_NUMBER=i
					System_data = load_system_data
					initialized = true
					com.Looking_for_other_elevators_on_network = false
					break
				}else if load_system_data.IP_list[i]==0{
					types.MY_NUMBER=i
					load_system_data.IP_list[i]=types.MY_IP
					var system_data_update types.Update_system_data
					system_data_update.Elevator_IP = types.MY_IP
					system_data_update.Elevator_number = types.MY_NUMBER
					com.Update_system_data_sendToComCh <- system_data_update
					System_data = load_system_data
					initialized = true
					com.Looking_for_other_elevators_on_network = false
					break
				}
			}

		case <- no_one_is_aliveCH:
			fmt.Printf(" Init as only elevator \n")
			types.MY_NUMBER = 0
			initialized = true
			com.Looking_for_other_elevators_on_network = false
		
		default:

		}
	time.Sleep(10*time.Millisecond)
	}
}

func manager_listen_for_elevator(){
	for{
		select{
			case order_to_delete := <- elevator.Next_floor_doneCh: 
				System_data.M_handle_q[order_to_delete][2*types.MY_NUMBER] = 0  // Deletes orders in both directions after visit
				System_data.M_handle_q[order_to_delete][2*types.MY_NUMBER+1] = 0
				System_data.M_internal_elev_out[order_to_delete][types.MY_NUMBER] = 0 
				timer.Executed_orderCh <- order_to_delete
				var system_data_update types.Update_system_data
				system_data_update.Add_order = 0
				system_data_update.Update_type = 1
				system_data_update.Floor_n = order_to_delete
				system_data_update.Elevator_number = types.MY_NUMBER
				com.Update_system_data_sendToComCh <- system_data_update



			case new_external_auction_data := <- elevator.External_orderCh:
				//fmt.Printf("fikk nokke paa external order ch \n")

				new_local_bid := cost.Calculate_cost(System_data, new_external_auction_data, Elevator_state)
				new_external_auction_data.Elevator_IP = types.MY_IP            
				new_external_auction_data.Bid = new_local_bid
				new_external_auction_data.Add = 1
				com.Auction_bid_sendToComCh <- new_external_auction_data
				timer.NewAuctionInfoToTimerCh <- new_external_auction_data

			case new_internal_order := <- elevator.Internal_orderCh:
				System_data.M_internal_elev_out[new_internal_order][types.MY_NUMBER]=1
//				System_data.M_handle_q[new_internal_order][types.MY_NUMBER*2]=1 
//				System_data.M_handle_q[new_internal_order][types.MY_NUMBER*2+1]=1
				var update_internal types.Update_system_data
				update_internal.Floor_n = new_internal_order
				update_internal.Elevator_number = types.MY_NUMBER
				update_internal.Update_type = 3
				com.Update_system_data_sendToComCh <- update_internal

			case Elevator_state := <- elevator.States_to_managerCh:
				Bullshit_incrementor += Elevator_state.Direction				
		} 
	time.Sleep(100*time.Millisecond)
	}	
}

func manager_listen_for_com(){
	for{
		time.Sleep(time.Millisecond*100)
		select{

		case new_external_auction_data := <- com.Auction_bid_sendToManagerCh:
			new_internal_bid := cost.Calculate_cost(System_data, new_external_auction_data, Elevator_state)

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
				timer.NewAuctionInfoToTimerCh <- new_external_auction_data

			}else if new_internal_bid == new_external_auction_data.Bid{ 
				if new_external_auction_data.Elevator_IP > types.MY_IP{
					new_external_auction_data.Elevator_number = types.MY_IP                   
					new_external_auction_data.Bid = new_internal_bid
					new_external_auction_data.Add = 1
					com.Auction_bid_sendToComCh <- new_external_auction_data
					timer.NewAuctionInfoToTimerCh <- new_external_auction_data 
				}
			}
			
		case new_system_data_update := <- com.Update_system_data_sendToManagerCh:
			switch{

				case new_system_data_update.Update_type == 0: //  add elevator in ip list
					System_data.IP_list[new_system_data_update.Elevator_number]=new_system_data_update.Elevator_IP

				case new_system_data_update.Update_type == 1: // delete timers
					var set_peripheral_order types.Auction_data
					set_peripheral_order.Elevator_number = new_system_data_update.Elevator_number
					set_peripheral_order.Floor = new_system_data_update.Floor_n
					set_peripheral_order.Add = new_system_data_update.Add_order
					timer.NewPeripheralOrderCh <- set_peripheral_order
					
				case new_system_data_update.Update_type == 2: // handel_q   
					System_data.M_handle_q[new_system_data_update.Floor_n][2*new_system_data_update.Elevator_number+new_system_data_update.Direction] = new_system_data_update.Add_order

				case new_system_data_update.Update_type == 3:
					System_data.M_internal_elev_out[new_system_data_update.Floor_n][new_system_data_update.Elevator_number]=1
				default:
				
			}
		}
	}
}

func listen_for_timeout(){
	for{
		time.Sleep(time.Millisecond*100)
		select{

		case peripheral_timout := <- timer.Handle_q_timeoutCh:
			new_local_bid := cost.Calculate_cost(System_data, peripheral_timout, Elevator_state)
			peripheral_timout.Elevator_IP = types.MY_IP            
			peripheral_timout.Bid = new_local_bid
			peripheral_timout.Add = 1
			com.Auction_bid_sendToComCh <- peripheral_timout
			timer.NewAuctionInfoToTimerCh <- peripheral_timout

		case won_assignment := <- timer.NotifyWinningBidToManagerCh:
			if won_assignment.Direction==0{
				System_data.M_handle_q[won_assignment.Floor][types.MY_NUMBER*2+won_assignment.Direction]=1
			}else{
				System_data.M_handle_q[won_assignment.Floor][types.MY_NUMBER*2+won_assignment.Direction]=1
			}	
		default:
		}
	}
}

func determine_next_floor(){
	for{
		if Elevator_state.Direction == types.RUNDOWN{
			for i := Elevator_state.Last_floor; i >= 0; i-- {
				if System_data.M_handle_q[i][2*types.MY_NUMBER+types.DOWN]==1||System_data.M_internal_elev_out[i][types.MY_NUMBER]==1{
					elevator.Next_floorCh <- i
				break
				}
			}
		}else if Elevator_state.Direction == types.RUNUP{
			for i := Elevator_state.Last_floor; i < types.N_FLOORS; i++ {
				if System_data.M_handle_q[i][2*types.MY_NUMBER+types.UP]==1||System_data.M_internal_elev_out[i][types.MY_NUMBER]==1{
					elevator.Next_floorCh <- i
				break
				}
			}
		}else if Elevator_state.Direction == types.STOPP{
			for i := types.N_FLOORS-1; i >=0 ; i-- {
				if System_data.M_handle_q[i][2*types.MY_NUMBER+types.UP]==1||System_data.M_internal_elev_out[i][types.MY_NUMBER]==1||System_data.M_handle_q[i][2*types.MY_NUMBER+types.DOWN]==1{
					elevator.Next_floorCh <- i
				break
				}
			}
		}		
		time.Sleep(500*time.Millisecond)
	}
}

func start_auction(external_bid types.Auction_data){ //lage to funskjoner for forskjellige triggere
	new_internal_bid := cost.Calculate_cost(System_data, external_bid, Elevator_state)
	if new_internal_bid < external_bid.Bid{
	}
}

func Update_button_lamps(){
	for{		
		for i := 0; i < types.N_FLOORS; i++ {
			var zero_counter_up int = 0
			var zero_counter_down int = 0
			for j := 0; j < types.MAX_N_ELEVATORS*2; j++ {
				if System_data.M_internal_elev_out[i][types.MY_NUMBER] == 0{
					elevator.Set_button_lamps(2,i,0)
				}
				if System_data.M_internal_elev_out[i][types.MY_NUMBER] == 1{
					elevator.Set_button_lamps(2,i,1)
				}				
				if System_data.M_handle_q[i][j] == 0{
					if (j%2)==1{
						zero_counter_down = zero_counter_down + 1
					}else{
						zero_counter_up = zero_counter_up + 1
					}
				}				
				if System_data.M_handle_q[i][j] == 1{
					if (j%2)==0{
						elevator.Set_button_lamps(types.UP,i,1)
					}else{
						elevator.Set_button_lamps(types.DOWN,i,1)
					}
				}
			}
			if zero_counter_up == types.MAX_N_ELEVATORS{
				elevator.Set_button_lamps(0,i,0)
			}
			if zero_counter_down == types.MAX_N_ELEVATORS{
				elevator.Set_button_lamps(1,i,0)
			}			
		}
	time.Sleep(50*time.Millisecond)	
	}
}
