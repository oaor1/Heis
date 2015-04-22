package timer

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
			
			}
		}	
		//	time.Sleep(DECREMENT_INTERVAL) // maa sikkert endrest paa
		
	}
}

/*
Liker ikke at både main og go rutine sleeper
*/



//--------------------------------------------OLD CODE -----------------------------------------


/*package timer

import (
	"../types"
	"fmt"
	"time"
)

var(
	NotifyWinningBidToManager = make(chan types.Auction_data)
	NewAuctionInfoToTimer = make(chan types.Auction_data)
	Up_order_timer [types.N_FLOORS-1] int
	Down_order_timer [types.N_FLOORS-1] int
	Handle_q_timer [types.N_FLOORS-1] int
)
// vi maa utvide
func timer(auction_data types.Auction_data){

	
	
}

func handle_timer_for_new_bids(){
	newAuctionObject := <- NewAuctionInfoToTimer
	if newAuctionObject.



}
func decrement_and_check_handle_timers(){
	for i := 0; i < ; i++ {
		
	}

}

func decrement_and_check_auction_timers(){
	for i := 0; i < (types.N_FLOORS)-1; i++ {
		if up_order_timer[i] > 0{
			Up_order_timer[i] = up_order_timer[i] - 1
		}
		if Down_order_timer[i] > 0{
			Down_order_timer[i] = Down_order_timer[i] - 1
		}
		
		if up_order_timer[i] == 1{
			Up_order_timer[i] = 0
			WinningBid types.Auction_data
			WinningBid.Floor = i
			WinningBid.Matrix_type = 0 // maa agsaa legge til i handle timer
			NotifyWinningBidToManager <- WinningBid

		}
		if Down_order_timer[i] ==1{
			Down_order_timer[i] = 0
			WinningBid types.Auction_data
			WinningBid.Floor = i
			WinningBid.Matrix_type = 1
			NotifyWinningBidToManager <- WinningBid
		}
	}

	time.Sleep(40*time.Millisecond) // maa sikkert endrest paa
}*/