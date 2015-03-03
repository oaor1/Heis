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
	//Setter opp kananler
	doorTimerStartCh = make(chan bool)
	doorTimerStoppCh = make(chan bool)
	idleCh = make(chan bool)
	openCh = make(chan bool)
	upCh = make(chan bool)
	downCh = make(chan bool)
	//osChan = chan os.Signal
)

func Run(){
	done = make(chan bool)
	//initialize
	driver.Elev_init()
}