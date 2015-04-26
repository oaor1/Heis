package cost

import(
	"../types"
	
)

func Calculate_cost(System_data types.System_data, auction_data types.Auction_data, elevator_state types.Elevator_state) int{
	cost := 0
	up_orders := 0
	dwn_orders := 0
	diff_cost := 0
	for i := 0; i < types.N_FLOORS; i++ {
		if System_data.M_handle_q[0][i] == 1{
			up_orders++
		}
		if System_data.M_handle_q[1][i] == 1{
			dwn_orders++

		}
	}
	if elevator_state.Direction != auction_data.Direction{
		diff := elevator_state.Last_floor - auction_data.Floor
		if diff < 0{
			diff_cost = diff*-1
		}else{
			diff_cost = diff
		}		
	}
	if elevator_state.Direction == 0{
		dwn_orders = dwn_orders*2
	}
	if elevator_state.Direction == 1{
		up_orders = up_orders*2
	}
	cost = up_orders + dwn_orders + diff_cost
	
	return cost
}