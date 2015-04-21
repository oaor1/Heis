package timer

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

func timer(auction_data types.Auction_data){

	
	
}

func handle_timer_for_new_bids(){
	

}

func decrement_and_check_timers(){
	for i := 0; i < (types.N_FLOORS)-1; i++ {
		if up_order_timer[i] > 0{
			Up_order_timer[i] = up_order_timer[i] - 1
		}
		if Down_order_timer[i] > 0{
			Down_order_timer[i] = Down_order_timer[i] - 1
		}
		if Handle_q_timer[i] > 0{
			Handle_q_timer[i] = handle_q_timer[i] - 1
		}
		if up_order_timer[i] == 1{
			Up_order_timer[i] = 0
			WinningBid types.Auction_data
			WinningBid.Floor = i
			WinningBid.Matrix_type = 0
			NotifyWinningBidToManager <- WinningBid

		}
		if Down_order_timer[i] ==1{
			Down_order_timer[i] = 0
			WinningBid types.Auction_data
			WinningBid.Floor = i
			WinningBid.Matrix_type = 1
			NotifyWinningBidToManager <- WinningBid
		}
		if Handle_q_timer[i] ==1{
			Handle_q_timer[i] = 0
			WinningBid types.Auction_data
			WinningBid.Floor = i
			WinningBid.Matrix_type = 2
			NotifyWinningBidToManager <- WinningBid
		}
		
	}

	time.Sleep(40*time.Millisecond) // maa sikkert endrest paa
}