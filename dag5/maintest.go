//export GOPATH=$(pwd)
//go install driver

package main

import(
	"time"
	."../driver"
)

Elev_init()



for{
	driver.Elev_set_motor_direction(1)

	time.Sleep(2*time.Second)

	Set_door_open_lamp(1)

	Elev_set_motor_direction(-1)

	time.Sleep(2*time.Second)
}

