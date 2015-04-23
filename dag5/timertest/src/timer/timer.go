package timer

import (
	"../types"
	"fmt"
	"time"
)
const(
	AUCTIONTIME = 1000  // the timer will count for something like AUCTIONTIME * DECREMENT_INTERVAL
	DECREMENT_INTERVAL = 4*time.Millisecond
)

var(
	NotifyWinningBidToManagerCh = make(chan types.Auction_data)
	NewAuctionInfoToTimerCh = make(chan types.Auction_data)
	NewPeripheralOrderCh = make(chan types.Auction_data)
	Handle_q_timeoutCh = make(chan types.Auction_data)


	//Number_of_elevators int

	Up_order_timer [types.N_FLOORS-1] int
	Down_order_timer [types.N_FLOORS-1] int

	Standing_upwards_bid[ types.N_FLOORS-1] int
	Standing_downwards_bid[types.N_FLOORS-1] int
	
	Handle_q_timer [types.N_FLOORS-1][20]int
	
	
)



func Listen_for_auctiondata_from_manager(){
	for{
		time.Sleep(time.Millisecond*10)
		select{
		case new_auction_data := <- NewAuctionInfoToTimerCh:
			//mottar bud fra manager
			fmt.Printf("---mottar ny budinfo fra manager \n %v",new_auction_data)
			switch{
			case new_auction_data.Direction == types.UP:
				Up_order_timer[new_auction_data.Floor]=AUCTIONTIME*new_auction_data.Add
				Standing_upwards_bid[new_auction_data.Floor]=new_auction_data.Bid*new_auction_data.Add
			case new_auction_data.Direction == types.DOWN:
				Down_order_timer[new_auction_data.Floor]=AUCTIONTIME*new_auction_data.Add
				Standing_downwards_bid[new_auction_data.Floor]=new_auction_data.Bid*new_auction_data.Add
			default:
				fmt.Printf("Dette skal ikkje skje .. timer ln 35 mottok eit ugyldig bud")
			}
		case new_peripheral_order := <- NewPeripheralOrderCh:
			Handle_q_timer[new_peripheral_order.Floor][new_peripheral_order.Elevator_number*2+new_peripheral_order.Direction]=new_peripheral_order.Bid
			//ny perifer ordre maa legges til i timout queue
			fmt.Printf("---nokken har tatt på seg eit oppdrag, vi maa ta tida\n %v",new_peripheral_order)
		}
	}
}



/*
func Remember_standing_bid(){
	fmt.Print("-------- hallo from remember_standing_bid() \n")
	

	for {
		time.Sleep(10*time.Millisecond)
		select{
		

			case newAuctionObject := <- NewAuctionInfoToTimerCh:
			fmt.Print("----------mottok eit bud!-------\n")
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
*/
/*
func Funk(){
	fmt.Print("-------- hallo from funk() \n")
	for{

		time.Sleep(1*time.Millisecond)

		select{
		case vinnarbud := <- NotifyWinningBidToManagerCh:

			fmt.Printf("Vi Vant Vi VaAnT !!! \n %v \n",vinnarbud)
		case nokken_somla := <- Handle_q_timeoutCh:
			fmt.Printf("-----Nokken har somla, \n %v \n",nokken_somla)


		}
	}
}
*/
func Add_handle_timer_for_new_system_data_update(){
	fmt.Print("-------- hallo from add_handle_timer_for_new_system_data_update()-\n")
	for{
		time.Sleep(1*time.Second)
		nPO := <- NewPeripheralOrderCh
		Handle_q_timer[nPO.Floor][(nPO.Elevator_number*2)+nPO.Direction] = nPO.Bid
	}
}

func Decrement_and_check_handle_timers(){
	fmt.Print("-------- hallo from decrement_and_check_handle_timers \n")
	for{
	time.Sleep(1*time.Second)
	for i := 0; i < types.N_FLOORS-1; i++ {
		for j := 0; j < 20; j++ {
			if Handle_q_timer[i][j] == 1{
				Handle_q_timer[i][j] = 0
				var NewOrder types.Auction_data
				NewOrder.Direction = j%2  //dette betyr at up er 0 og down er 1
				NewOrder.Floor = i
				Handle_q_timeoutCh <- NewOrder
			}
			if Handle_q_timer[i][j] > 0{
				Handle_q_timer[i][j] = Handle_q_timer[i][j] - 1
			}

		}
	}
	}
}

func Decrement_and_check_auction_timers(){
	fmt.Print("-------- hallo from decrement_and_check_auction_timers()\n")
	

	for{
		

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
				Handle_q_timer[WinningBid.Floor][types.MY_NUMBER*2]=Standing_upwards_bid[WinningBid.Floor]
				NotifyWinningBidToManagerCh <- WinningBid

			}
			if Down_order_timer[i] ==1{
				Down_order_timer[i] = 0
				var WinningBid types.Auction_data
				WinningBid.Floor = i
				WinningBid.Matrix_type = 1
				Handle_q_timer[WinningBid.Floor][(types.MY_NUMBER*2)+1]=Standing_downwards_bid[WinningBid.Floor]
				NotifyWinningBidToManagerCh <- WinningBid
			}
			
			}
			time.Sleep(DECREMENT_INTERVAL) 
		}	
		
		
	}



