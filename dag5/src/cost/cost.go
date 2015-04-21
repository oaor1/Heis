







calculate_cost(System_data types.System_data, New_auction_data types.Auction_data) int{
	cost := 0
	up_orders := 0
	dwn_orders := 0
	overlap_orders := 0
	/*
	++++ obstruksjon
	++++ stoppknapp
	+ antall bestillinger i samme retning ^Y  * distandse * konst
	+ antall bestillinger i motsatt retning ^X * distanse * konst
	- overlappende bestillinger (bor splittes i to caser , en for hver retning)
	*/
	for i := 0; i < N_FLOORS; i++ {
		if System_data.M_handle_q[0][i] == 1{
			up_orders++
		}
		if System_data.M_handle_q[1][i] == 1{
			dwn_orders++
		}
		if System_data.internal_elev_out[/*my number*/][New_auction_data.]
		
	}
	return cost
}