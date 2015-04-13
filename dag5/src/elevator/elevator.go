package elevator

import (
	"../driver"
	"../fmt"
)

const(
	UP = 0
	DOWN = 1
	IDLE = 2
	OPEN = 3

	ON = 1
	OFF = 0
	SPEED = 300
)

var(
	//for comunication with manager
	Next_floorCh = make(chan int)
	Current_floorCh = make (chan int)
	current_floor int
	ElevDirectionCh = make (chan int)
	ElevDirection int = DOWN
	//Setter opp kananler
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
	driver.ElevSetSpeed(-SPEED)
	for driver.Elev_get_floor_sensor_signal() == -1{
		time.Sleep(10*time.Millisecond)
	}
	driver.ElevSetSpeed(0)
	//Elevator initialized and in a definite floor
	go Idle()
	go Open()
	go Down()
	go Up()
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