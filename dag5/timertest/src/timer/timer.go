package main

import (
	"../types"
	"fmt"
	"time"
)
const(
	AUCTIONTIME = 10  // the timer will count for something like AUCTIONTIME * DECREMENT_INTERVAL
	DECREMENT_INTERVAL = 40*time.Millisecond
)

var(
	NotifyWinningBidToManagerCh = make(chan types.Auction_data)
	NewAuctionInfoToTimerCh = make(chan types.Auction_data)
	NewPeripheralOrderCh = make(chan types.Auction_data)

	//Number_of_elevators int

	Up_order_timer [types.N_FLOORS-1] int
	Down_order_timer [types.N_FLOORS-1] int

	Standing_upwards_bid[ types.N_FLOORS-1] int
	Standing_downwards_bid[types.N_FLOORS-1] int
	
	Handle_q_timer [types.N_FLOORS-1][20]int
	
	
)

func main(){


	fmt.Print("hallo from main \n")
	go funk()
	go kanaliser()
	go decrement_and_check_handle_timers()
	go decrement_and_check_auction_timers()

	go remember_standing_bid()

	go add_handle_timer_for_new_system_data_update()

	var Tullebud types.Auction_data
	Tullebud.Bid = 70
	Tullebud.Direction = 1
	Tullebud.Add = 1
	Tullebud.Floor = 2
	fmt.Printf("%v  \n ",Tullebud)
	//fmt.Print("hallo etter tullebud \n")
	NewAuctionInfoToTimerCh <- Tullebud
	
	fmt.Print("hallo from main \n")
	
	NotifyWinningBidToManagerCh <- Tullebud
	

	time.Sleep(10*time.Second)

}

func kanaliser(){


	

}

func funk(){
	for{

		time.Sleep(1*time.Second)
		select{
		case ettelleranna := <- NotifyWinningBidToManagerCh:
			fmt.Print("Vi Vant Vi VaAnT !!! %v ",ettelleranna)

		}
	}
}

func remember_standing_bid(){
	for {
		time.Sleep(1*time.Second)
		select{
		

		
	
		case newAuctionObject := <- NewAuctionInfoToTimerCh:
			fmt.Print("mottok eit bud!-------\n")
			//Number_of_elevators = newAuctionObject.Number_of_elevators
			if newAuctionObject.Direction == 1 && Up_order_timer[newAuctionObject.Floor]== 0{
				Up_order_timer[newAuctionObject.Floor] = AUCTIONTIME*newAuctionObject.Add
				Standing_upwards_bid[newAuctionObject.Floor] = newAuctionObject.Bid*newAuctionObject.Add
			}else if newAuctionObject.Direction == 0 && Up_order_timer[newAuctionObject.Floor]== 0{
				Down_order_timer[newAuctionObject.Floor] = AUCTIONTIME*newAuctionObject.Add
				Standing_downwards_bid[newAuctionObject.Floor] = newAuctionObject.Bid*newAuctionObject.Add
			}
		}
	
	}
}


func add_handle_timer_for_new_system_data_update(){
	time.Sleep(1*time.Second)
	nPO := <- NewPeripheralOrderCh
	Handle_q_timer[nPO.Floor][(nPO.Elevator_number*2)+nPO.Direction] = nPO.Bid
}

func decrement_and_check_handle_timers(){
	time.Sleep(1*time.Second)
	fmt.Print("hallo from decrement_and_check_handle_timers \n")
	for i := 0; i < types.N_FLOORS-1; i++ {
		for j := 0; j < 20; j++ {
			if Handle_q_timer[i][j] > 0{
				Handle_q_timer[i][j] = Handle_q_timer[i][j] - 1
			}

		}
	}
}

func decrement_and_check_auction_timers(){
	

	for{
		time.Sleep(4*time.Second)

		for i := 0; i < 3; i++ {
			/*
			if Up_order_timer[i] > 0{
				Up_order_timer[i] = Up_order_timer[i] - 1
			}
			if Down_order_timer[i] > 0{
				Down_order_timer[i] = Down_order_timer[i] - 1
			}
			
			if Up_order_timer[i] == 1{
				Up_order_timer[i] = 0
				var WinningBid types.Auction_data
				WinningBid.Floor = i
				WinningBid.Matrix_type = 0 // maa agsaa legge til i handle timer
				NotifyWinningBidToManagerCh <- WinningBid

			}
			if Down_order_timer[i] ==1{
				Down_order_timer[i] = 0
				var WinningBid types.Auction_data
				WinningBid.Floor = i
				WinningBid.Matrix_type = 1
				NotifyWinningBidToManagerCh <- WinningBid
			}
			
			}*/
		}	
		//	time.Sleep(DECREMENT_INTERVAL) // maa sikkert endrest paa
		
	}
}

/*
Liker ikke at både main og go rutine sleeper
*/

