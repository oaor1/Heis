package elevator

import (
	"../driver"
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
	//Setter opp kananler
	doorTimerStartCh = make(chan int) //as bool
	doorTimerStoppCh = make(chan int) //as bool
	idleCh = make(chan int) //as bool
	openCh = make(chan int) //as bool
	upCh = make(chan int) //as bool
	downCh = make(chan int) //as bool
	//osChan = chan os.Signal
)

func Run(){
	done = make(chan bool)
	//initialize
	driver.Elev_init()
}