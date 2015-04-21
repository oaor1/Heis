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
	next_floor int = 1
	Next_floorCh = make(chan int)
	Next_floor_doneCh = make (chan int)
	current_floor int
	Current_floorCh = make (chan int)
	ElevDirectionCh = make (chan int)
	ElevDirection int = types.RUNDOWN
	External_orderCh = make (chan types.Auction_data)
	Internal_orderCh = make (chan int)

	//Setter opp interne kananler /"states"
	doorTimerStartCh = make(chan bool) //as bool
	doorTimerStoppCh = make(chan bool) //as bool
	idleCh = make(chan bool) //as bool
	openDoorCh = make(chan bool) //as bool
	elevUpCh = make(chan bool) //as bool
	elevDownCh = make(chan bool) //as bool
	stoppCh = make(chan bool)
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
	go Idle()
	go Open()
	go Down()
	go Up()
	go Stopp()
	go Update_channels()
	go Read_order_buttons()
	idleCh <- true
	<- done

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

func Idle(){
	for{
		<-idleCh
		fmt.Println("Idle entered")
		driver.Elev_set_motor_direction(types.STOPP)
		for{
			time.Sleep(10*time.Millisecond)
			if driver.Elev_get_floor_sensor_signal() == next_floor{
				openDoorCh <- true
				break
			}else if current_floor > next_floor {//neste etasje i ordre kø er under
				elevDownCh <- true
				break
			}else if current_floor < next_floor {//neste etasje i odrdre kø er over
				elevUpCh <- true
				break
			}else if driver.Elev_get_stop_signal() == 1{
				stoppCh <- true
				break
			}
			for driver.Elev_get_obstruction_signal() == 1{
				time.Sleep(100*time.Millisecond)
			fmt.Println("blalblalb")
			}	
		}
	}
}

func Open(){
	for{
		<-openDoorCh
		fmt.Println("open enetered")
//		driver.Elev_set_motor_direction(ElevDirection *-1)
//		time.Sleep(50*time.Millisecond)
		driver.Elev_set_motor_direction(types.STOPP) 
		driver.Elev_set_door_open_lamp(ON)
		doorTimerStartCh <- true
		<-doorTimerStoppCh
		for driver.Elev_get_obstruction_signal() == 1{
			time.Sleep(100*time.Millisecond)
		}
		driver.Elev_set_door_open_lamp(OFF)
//		Next_floor_doneCh <- current_floor //Noen må lese denne, hvis ikke blir den stuck her
		idleCh <- true
	}
}

func Up(){
	<-elevUpCh
	fmt.Println("UP enetered")
	ElevDirection = types.RUNUP
//	ElevDirectionCh <- types.RUNUP
	driver.Elev_set_door_open_lamp(OFF)
	driver.Elev_set_motor_direction(ElevDirection)
	for{
		for driver.Elev_get_obstruction_signal() == 1{
			driver.Elev_set_motor_direction(types.STOPP)
			fmt.Println("Obstruksjon RUNUP")
			time.Sleep(100*time.Millisecond)
		}
		if driver.Elev_get_floor_sensor_signal() == next_floor{//Har nådd next_floor
			openDoorCh <- true
			break
		}else if driver.Elev_get_floor_sensor_signal() == types.N_FLOORS-1{
			ElevDirection = types.RUNDOWN
			driver.Elev_set_motor_direction(types.STOPP)
			fmt.Println("saftey down")
		}else if driver.Elev_get_stop_signal() == 1{
			stoppCh <- true
		}
		//Do we need any saftey feature to prevent the eleavtor crash into the roof
	}
}

func Down(){
	<-elevDownCh
	fmt.Println("down enetered")
	ElevDirection = types.RUNDOWN
//	ElevDirectionCh <- ElevDirection
	driver.Elev_set_motor_direction(ElevDirection)
	driver.Elev_set_door_open_lamp(OFF)
	for{
		for driver.Elev_get_obstruction_signal() == 1{
			driver.Elev_set_motor_direction(types.STOPP)
			fmt.Println("Obstruksjon RUNUP")
			time.Sleep(100*time.Millisecond)
		}
		if driver.Elev_get_floor_sensor_signal() == next_floor{//Har nådd next_floor
			openDoorCh <- true
			break
		}else if driver.Elev_get_floor_sensor_signal() == 0{
			ElevDirection = types.RUNUP
			driver.Elev_set_motor_direction(types.STOPP)
			fmt.Println("saftey up")
			break
		}else if driver.Elev_get_stop_signal() == 1{
			stoppCh <- true
			break
		}
		//Do we need any saftey feature to prevent the eleavtor crash into the roof
	}
}

func Stopp(){
	<-stoppCh
	fmt.Println("Stopp pushed")
	driver.Elev_set_motor_direction(types.STOPP)
	for driver.Elev_get_stop_signal() == 1{
		driver.Elev_set_stop_lamp(1)
	}
	if driver.Elev_get_floor_sensor_signal() >= 0{
		driver.Elev_set_door_open_lamp (1)
	}
	time.Sleep(1*time.Second)
	for driver.Elev_get_stop_signal() == 1{
		driver.Elev_set_door_open_lamp (0)
		driver.Elev_set_stop_lamp(0)
	}
	idleCh <- true
}

func Update_channels(){
	for{
		selcet{
		case next_floor = <- Next_floorCh:
			fmt.Println("Mottok next_floor fra Next_floorCh")
		default:
		}

		if driver.Elev_get_floor_sensor_signal() >= 0{
//			Current_floorCh <- driver.Elev_get_floor_sensor_signal()
			current_floor = driver.Elev_get_floor_sensor_signal()
//			fmt.Println("Current floor: ", current_floor)
		}
		time.Sleep(100*time.Millisecond)
	}
}

func Door_safety(){
	for{
		if driver.Elev_get_floor_sensor_signal() == -1{
			driver.Elev_set_door_open_lamp(OFF)
			fmt.Println("Door safety")
		}
	}
}

func Read_order_buttons(){
	for{
		for i :=0; i < types.N_FLOORS-1; i++{
			if driver.Elev_get_button_signal(0,i) != 0{
				driver.Elev_get_button_signal(0,i,1)
				Internal_orderCh <- i
			}
			for j := 1; j < 3; j++{
				if driver.Elev_get_button_signal(j, i) != 0{
					driver.Elev_set_button_lamp(j,i,1)
					new_order types.Auction_data
					new_order.floor
					External_orderCh <- 
					next_floor = i
					fmt.Println("next floor: ", next_floor)
				}
			}
		}
		time.Sleep(40*time.Millisecond)
	}
}