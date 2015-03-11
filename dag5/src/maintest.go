//export GOPATH=$(pwd)
//go install driver

package main

import(
	"time"
	."driver"
)

func main(){

	Elev_init()

	for{
		Elev_set_motor_direction(1)

		time.Sleep(2*time.Second)

		Elev_set_door_open_lamp(1)

		Elev_set_motor_direction(-1)

		time.Sleep(2*time.Second)
	}
}
