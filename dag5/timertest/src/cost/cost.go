package cost

import(
	"../types"
)

const(
	WRONG_DIR_COST = 2
)

func Calculate_cost(System_data types.System_data, auction_bid types.Auction_data) int{
	cost := 6
	up_orders := 0
	dwn_orders := 0
	//overlap_orders := 0
	//myNumber := New_auction_data.Elevator_number
	/*
	++++ obstruksjon
	++++ stoppknapp
	+ antall bestillinger i samme retning ^Y  * distandse * konst
	+ antall bestillinger i motsatt retning ^X * distanse * konst
	- overlappende bestillinger (bor splittes i to caser , en for hver retning)
	*/
	for i := 0; i < types.N_FLOORS; i++ {
		if System_data.M_handle_q[0][i] == 1{
			up_orders++
		}
		if System_data.M_handle_q[1][i] == 1{
			dwn_orders++

		}
		//if System_data.internal_elev_out[myNumber][New_auction_data.]
	/*	
	}
	if elevator_status.Direction == types.RUNUP{
		dwn_orders = dwn_orders * WRONG_DIR_COST
	}else if elevator_status:Direction == types.RUNUP{
		up_orders = up_orders * WRONG_DIR_COST
	}
	*/
	
	
	}
	return cost
}