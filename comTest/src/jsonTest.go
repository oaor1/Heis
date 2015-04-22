package main

import (
	"./types"
	"fmt"
	"encoding/json"	
)

var (
	auction types.Auction_data
)

func main(){

/*
	posission := types.Elevator_state{
		Direction: types.RUNDOWN,
		Last_floor: 3}

*/

	auction.Bid = types.RUNUP
	auction.Direction = 3
	auction.Elevator_number = 143
	fmt.Printf("auction: %+v \n", auction)

//	Some_kind_of_struct := <- Some_kind_of_structCh
	resMarshal, _ := json.Marshal(auction)

	var buffer []byte = make([]byte, len(resMarshal)+1)
	buffer [0] = 1
	for i := 0; i < len(resMarshal); i++ {
		buffer [i+1] = resMarshal [i]
	}
	fmt.Println(buffer)
//	fmt.Println(len(buffer))


//	var buffer []byte = make([]byte, 256)

	var resUnmarshal types.Auction_data
	errunm := json.Unmarshal(buffer[1:len(resMarshal)+1], &resUnmarshal)
		if errunm != nil{
			fmt.Printf("resUnmarshal failed  %i \n", errunm)
			return
		}
	fmt.Printf("struct: %+v \n", resUnmarshal)
}
