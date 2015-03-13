//export GOPATH=$(pwd)
//go install driver

package main

import(
	"fmt"
	"time"
	"driver"
	"runtime"
)
func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	driver.Elev_init()

	go test_buttons()
	go set_level_ligths()
	go order_button_test()
	go run_elev()

	time.Sleep(60*time.Second)
	driver.Elev_set_stop_lamp(0)
	driver.Elev_set_motor_direction(0)
}

func run_elev(){
	for{
		driver.Elev_set_motor_direction(1)

		time.Sleep(1*time.Second)

		driver.Elev_set_door_open_lamp(1)

		driver.Elev_set_motor_direction(-1)

		time.Sleep(1*time.Second)
	}
}


func test_buttons(){
	for{
		if driver.Elev_get_obstruction_signal() != 0{
			fmt.Println ("YOLO Obstruction")
			driver.Elev_set_stop_lamp(1)
			time.Sleep(1*time.Second)
			driver.Elev_set_stop_lamp(1)
		}
	}
}

func set_level_ligths(){
	for{
		level := driver.Elev_get_floor_sensor_signal()
		if level >= 0{
			driver.Elev_set_floor_indicator(level)
		}
	}
}

func order_button_test(){
	for{
		for i :=0; i < 4; i++{
			for j := 0; j < 3; j++{
				if driver.Elev_get_button_signal(j, i) != 0{
					driver.Elev_set_button_lamp(j,i,1)
				}
			}
		}
	}
}


