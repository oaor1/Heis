package elevator

import (
	"../driver"
	"../types"
	"fmt"
)

const(
	RUNUP = 1
	STOPP = 0
	RUNDOWN = -1

	ON = 1
	OFF = 0
	SPEED = 300
)

var(
	//for comunication with manager
	Next_floor int
	Next_floorCh = make(chan int)
	Next_floor_doneCh = make (chan bool)
	Current_floor int
	Current_floorCh = make (chan int)
	ElevDirectionCh = make (chan int)
	ElevDirection int = RUNDOWN

	//Setter opp interne kananler /"states"
	doorTimerStartCh = make(chan bool) //as bool
	doorTimerStoppCh = make(chan bool) //as bool
	idleCh = make(chan bool) //as bool
	openDoorCh = make(chan bool) //as bool
	elevUpCh = make(chan bool) //as bool
	elevDownCh = make(chan bool) //as bool
	//osChan = chan os.Signal
)

func Run(){
	done = make(chan bool)
	//initialise the elevator
	driver.Elev_init()
	go DoorTimer()
	go FloorLigths()
	driver.Elev_set_motor_direction(RUNDOWN)
	for driver.Elev_get_floor_sensor_signal() == RUNDOWN{
		time.Sleep(10*time.Millisecond)
	}
	driver.Elev_set_motor_direction(STOPP)
	//Elevator initialized and in a definite floor
	go Idle()
	go Open()
	go Down()
	go Up()
	go Range_safety()
	go Update_channels()
	idleCh <- true
	<- done

}

func DoorTimer(){
	for{
		<-doorTimerStartCh
		time.Sleep(3*time.Second)
		for driver.Elev_get_obstruction_signal() == 1{
			time.Sleep(1*second)¨
			fmt.Printf("Obstruksjon av heis med dør åpen detektert.\n")
		}
		doorTimerStoppCh <- true
	}
}

func FloorLigths(){
	for{
		time.Sleep(10*millisecond)
		driver.Elev_set_floor_indicator(Elev_get_floor_sensor_signal())
	}
}

func Idle(){
	for{
		<-idleCh
		driver.Elev_set_motor_direction(STOPP)
		for{
			time.Sleep(10*time.Millisecond)
			if Elev_get_floor_sensor_signal() == Next_floor{
				Open <- true
				break
			}
			else if current_floor > Next_floor {//neste etasje i ordre kø er under
				elevDownCh <- true
				break
			}
			else if current_floor < Next_floor {//neste etasje i odrdre kø er over
				elevUpCh <- true
				break
			}
			for Elev_get_obstruction_signal == 1{
				time.Sleep(100*time.Millisecond)
			}	
		}
	}
}

func Open(){
	for{
		<-openDoorCh
		Elev_set_motor_direction(ElevDirection *-1)
		time.Sleep(50*time.Millisecond)
		Elev_set_motor_direction(STOPP) 
		//stop func evt
		driver.Elev_set_door_open_lamp(ON)
		doorTimerStartCh <- true
		<-doorTimerStoppCh
		for Elev_get_obstruction_signal() = 1{
			time.Sleep(100*time.Millisecond)
		}
		driver.Elev_set_door_open_lamp(OFF)
		Next_floor_doneCh <- true
		idleCh <- true
	}
}

func Up(){
	<-elevUpCh
	ElevDirection = RUNUP
	driver.Elev_set_door_open_lamp(OFF)
	Elev_set_motor_direction(ElevDirection)
	for{
		for Elev_get_obstruction_signal == 1{
			Elev_set_motor_direction(STOPP)
			time.Sleep(100*time.Millisecond)
		}
		if Elev_get_floor_sensor_signal() == Next_floor{//Har nådd next_floor
			OpenDoorCh <- true
			break
		}
		//Do we need any saftey feature to prevent the eleavtor crash into the roof
	}
}

func Down(){
	<-elevDownCh
	ElevDirection = RUNDOWN
	ElevDirectionCh <- ElevDirection
	Elev_set_motor_direction(ElevDirection)
	driver.Elev_set_door_open_lamp(OFF)
	for{
		for Elev_get_obstruction_signal == 1{
			Elev_set_motor_direction(STOPP)
			time.Sleep(100*time.Millisecond)
		}
		if Elev_get_floor_sensor_signal() == Next_floor{//Har nådd next_floor
			OpenDoorCh <- true
			break
		}
		//Do we need any saftey feature to prevent the eleavtor crash into the roof
	}
}

func Update_channels(){
	for{
		Next_floor <- Next_floorCh
		if Elev_get_floor_sensor_signal() >= 0{
			Current_floorCh <- Elev_get_floor_sensor_signal()
			Current_floor = Elev_get_floor_sensor_signal()
		}
		time.Sleep(100*time.Millisecond)
	}
}

func Range_safety(){
	for{
		if Elev_get_floor_sensor_signal() == N_FLOORS{
			ElevDirection = RUNDOWN
		}
		else if Elev_get_floor_sensor_signal() = 0{
			ElevDirection = RUNUP
		}
		time.Sleep(10*time.Millisecond)
	}
}

func Door_safety(){
	for{
		if Elev_get_floor_sensor_signal == -1{
			Elev_set_door_open_lamp(OFF)
		}
	}
}