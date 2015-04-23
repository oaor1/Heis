package elevator

import (
	"../driver"
	"../types"
	"fmt"
	"time"
)

const(
	ON = 1
	OFF = 0
	SPEED = 300
)

var(
	//for comunication with manager
	next_floor int //husk kannal fremfor stor bokstav.
	Next_floorCh = make(chan int,1)
	Next_floor_doneCh = make (chan int)
	current_floor int
	Current_floorCh = make (chan int, 1)
	ElevDirectionCh = make (chan int)
	ElevDirection int = types.RUNDOWN
	External_orderCh = make (chan types.Auction_data)
	Internal_orderCh = make (chan int)
	States_to_managerCh = make (chan types.Elevator_state)

	//Setter opp interne kananler /"states"
	doorTimerStartCh = make(chan bool) //as bool
	doorTimerStoppCh = make(chan bool) //as bool
	//osChan = chan os.Signal
)

func Run(){
	done := make(chan bool)
	//initialise the elevator
	driver.Elev_init()
	go DoorTimer()
	go FloorLigths()
	driver.Elev_set_motor_direction(types.RUNDOWN)
	for driver.Elev_get_floor_sensor_signal() == types.RUNDOWN{
		time.Sleep(10*time.Millisecond)
	}
	driver.Elev_set_motor_direction(types.STOPP)
	//Elevator initialized and in a definite floor
	go Update_channels()
	go Read_order_buttons()
	go Print()
	Idle()
	<- done
}

func Idle(){
	for{
		fmt.Println("Idle entered")
		driver.Elev_set_motor_direction(types.STOPP)
		for{
			for driver.Elev_get_obstruction_signal() == 1{
				time.Sleep(100*time.Millisecond)
			}
			switch{
				case driver.Elev_get_floor_sensor_signal() == next_floor:
					Open()
					break
				case current_floor > next_floor: //neste etasje i ordre kø er under
					ElevDirection = types.RUNDOWN
					//	ElevDirectionCh <- types.RUNDOWN
					Down()
					break
				case current_floor < next_floor: //neste etasje i odrdre kø er over
					ElevDirection = types.RUNUP
					//	ElevDirectionCh <- types.RUNUP
					Up()
					break
				case driver.Elev_get_stop_signal() == 1:
					Stopp()
					break
				default:
				
			}	
		}
		time.Sleep(5*time.Millisecond)
	}
}

func Open(){
	fmt.Println("open enetered")
//		driver.Elev_set_motor_direction(ElevDirection *-1)
//		time.Sleep(50*time.Millisecond)
	driver.Elev_set_motor_direction(types.STOPP) 
	driver.Elev_set_door_open_lamp(ON)
	doorTimerStartCh <- true
	if driver.Elev_get_stop_signal() == 1{
	Stopp()
	}
	<-doorTimerStoppCh
	for driver.Elev_get_obstruction_signal() == 1{
		time.Sleep(100*time.Millisecond)
	}
	driver.Elev_set_door_open_lamp(OFF)
//	Next_floor_doneCh <- current_floor //Noen må lese denne, hvis ikke blir den stuck her
	Idle()
}

func Up(){
	fmt.Println("UP enetered")
	driver.Elev_set_door_open_lamp(OFF)
	driver.Elev_set_motor_direction(ElevDirection)
	for{
		for driver.Elev_get_obstruction_signal() == 1{
			driver.Elev_set_motor_direction(types.STOPP)
			fmt.Println("Obstruksjon RUNUP")
			time.Sleep(100*time.Millisecond)
		}
		switch{
			case current_floor == next_floor: //Her må det gjøres noe
				Open()
			case driver.Elev_get_floor_sensor_signal() == types.N_FLOORS-1:
				ElevDirection = types.RUNDOWN
				driver.Elev_set_motor_direction(types.STOPP)
				fmt.Println("saftey down")
				Idle()
			case driver.Elev_get_stop_signal() == 1:
				Stopp()
			default:
		}
		//Do we need any saftey feature to prevent the eleavtor crash into the roof
	}
}

func Down(){
	fmt.Println("UP enetered")
	driver.Elev_set_door_open_lamp(OFF)
	driver.Elev_set_motor_direction(ElevDirection)
	for{
		for driver.Elev_get_obstruction_signal() == 1{
			driver.Elev_set_motor_direction(types.STOPP)
			fmt.Println("Obstruksjon RUNUP")
			time.Sleep(100*time.Millisecond)
		}
		switch{
			case current_floor == next_floor:
				Open()
			case driver.Elev_get_floor_sensor_signal() == types.N_FLOORS-1:
				ElevDirection = types.RUNDOWN
				driver.Elev_set_motor_direction(types.STOPP)
				fmt.Println("saftey down")
				Idle()
			case driver.Elev_get_stop_signal() == 1:
				Stopp()
			default:
		}
		//Do we need any saftey feature to prevent the eleavtor crash into the roof
	}
}

func Stopp(){
	fmt.Println("Stopp pushed")
	driver.Elev_set_motor_direction(types.STOPP)
	for driver.Elev_get_stop_signal() == 1{
		driver.Elev_set_stop_lamp(1)
	}
	if driver.Elev_get_floor_sensor_signal() >= 0{
		driver.Elev_set_door_open_lamp (1)
	}
	time.Sleep(1*time.Second)
	for driver.Elev_get_stop_signal() == 0{
	}
	driver.Elev_set_door_open_lamp (0)
	driver.Elev_set_stop_lamp(0)
	time.Sleep(1*time.Second)
	Idle()
}

func Update_channels(){
	for{
		if driver.Elev_get_floor_sensor_signal() >= 0{
			current_floor = driver.Elev_get_floor_sensor_signal()
		}
		time.Sleep(10*time.Millisecond)
		select{
		case next_floor = <- Next_floorCh:
			fmt.Println("Mottok next_floor fra Next_floorCh")
		default:
		}
	}
}

func Door_safety(){
	for{
		if driver.Elev_get_floor_sensor_signal() == -1{
			driver.Elev_set_door_open_lamp(OFF)
			fmt.Println("Door safety")
		}
		time.Sleep(10*time.Second)
	}
}

func Read_order_buttons(){
	for{
		for i :=0; i < types.N_FLOORS; i++{
			if driver.Elev_get_button_signal(2,i) != 0{
				driver.Elev_set_button_lamp(2,i,1)
//				Internal_orderCh <- i
				next_floor = i
			}
			for j := 0; j < 3; j++{
				if driver.Elev_get_button_signal(j, i) != 0{
					driver.Elev_set_button_lamp(j,i,1)
//					var new_order types.Auction_data
//					new_order.Floor = i
//					new_order.Direction = j
//					External_orderCh <- new_order
				}
			}
		}
		time.Sleep(5*time.Millisecond)
	}
}

func DoorTimer(){
	for{
		<-doorTimerStartCh
		time.Sleep(3*time.Second)
		for driver.Elev_get_obstruction_signal() == 1{
			time.Sleep(1*time.Second)
			fmt.Printf("Obstruksjon av heis med dør åpen detektert.\n")
		}
		doorTimerStoppCh <- true
	}
}

func FloorLigths(){
	for{
		time.Sleep(50*time.Millisecond)
		driver.Elev_set_floor_indicator(driver.Elev_get_floor_sensor_signal())
	}
}

func Print(){
	for{
		fmt.Printf("current floor: %d, Direction: %d, next_floor: %d\n" ,current_floor, ElevDirection, next_floor)
/*
		var state types.Elevator_state
		state.Direction = ElevDirection
		state.Last_floor = current_floor
*/

		time.Sleep(1*time.Second)
	}
}