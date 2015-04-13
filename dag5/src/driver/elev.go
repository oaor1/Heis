///run/user/11408/gvfs/smb-share:server=sambaad.stud.ntnu.no,share=sindrevh/sanntid/heis

package driver

import(
  "../init"
	"fmt"
)

const (
  /*
  MAX_N_ELEVATORS = 10
	N_FLOORS = 4
	N_BUTTONS = 3
  */

	DIRN_DOWN = -1
	DIRN_STOP = 0
	DIRN_UP = 1

	BUTTON_CALL_UP = 0
  BUTTON_CALL_DOWN = 1
  BUTTON_COMMAND = 2
)

var lamp_channel_matrix = [N_FLOORS][N_BUTTONS]int{
    {LIGHT_UP1, LIGHT_DOWN1, LIGHT_COMMAND1},
    {LIGHT_UP2, LIGHT_DOWN2, LIGHT_COMMAND2},
    {LIGHT_UP3, LIGHT_DOWN3, LIGHT_COMMAND3},
    {LIGHT_UP4, LIGHT_DOWN4, LIGHT_COMMAND4},
}

var button_channel_matrix = [N_FLOORS][N_BUTTONS]int{
    {BUTTON_UP1, BUTTON_DOWN1, BUTTON_COMMAND1},
    {BUTTON_UP2, BUTTON_DOWN2, BUTTON_COMMAND2},
    {BUTTON_UP3, BUTTON_DOWN3, BUTTON_COMMAND3},
    {BUTTON_UP4, BUTTON_DOWN4, BUTTON_COMMAND4},
}

/**
  Initialize elevator.
  @return Non-zero on success, 0 on failure.
*/
func Elev_init() int{
	if Io_init() != 0 {
		fmt.Printf("The elevator is initialized.\n")
	}

	//Slukker alle lamper under initialsiering
	for i := 0; i < N_FLOORS; i++ {
		if i != 0{
			Elev_set_button_lamp(BUTTON_CALL_DOWN, i , 0)
		}
		if i != N_FLOORS - 1{
			Elev_set_button_lamp(BUTTON_CALL_UP, i ,0)
		}
		Elev_set_button_lamp(BUTTON_COMMAND, i, 0)
	}
	Elev_set_stop_lamp(0)
	Elev_set_door_open_lamp(0)
	Elev_set_floor_indicator(0)

  return 1 //finished initialization
}

/**
  Sets the motor direction of the elevator.
  @param dirn New direction of the elevator.
*/
func Elev_set_motor_direction(dir int){
	if dir == 0{
    Io_write_analog(MOTOR, 0)
  } else if dir > 0 {
    Io_clear_bit(MOTORDIR)
    Io_write_analog(MOTOR, 2800)
  } else if dir < 0 {
    Io_set_bit(MOTORDIR)
    Io_write_analog(MOTOR, 2800)
  }
}

/**
  Turn door-open lamp on or off.
  @param value Non-zero value turns lamp on, 0 turns lamp off.
*/
func Elev_set_door_open_lamp(value int){
	if value != 0{
        Io_set_bit(LIGHT_DOOR_OPEN)
    } else {
        Io_clear_bit(LIGHT_DOOR_OPEN)
    }
}

/**
  Get signal from obstruction switch.
  @return 1 if obstruction is enabled. 0 if not.
*/
func Elev_get_obstruction_signal() int {
    return Io_read_bit(OBSTRUCTION)
}

/**
  Get signal from stop button.
  @return 1 if stop button is pushed, 0 if not.
*/
func Elev_get_stop_signal() int {
    return Io_read_bit(STOP)
}


/**
  Turn stop lamp on or off.
  @param value Non-zero value turns lamp on, 0 turns lamp off.
*/
func Elev_set_stop_lamp(value int) {
    if value != 0{
        Io_set_bit(LIGHT_STOP)
    } else {
        Io_clear_bit(LIGHT_STOP)
    }
}


/**
  Get floor sensor signal.
  @return -1 if elevator is not on a floor. 0-3 if elevator is on floor. 0 is
    ground floor, 3 is top floor.
*/
func Elev_get_floor_sensor_signal() int{
    if Io_read_bit(SENSOR_FLOOR1) != 0{
        return 0
    } else if Io_read_bit(SENSOR_FLOOR2) != 0{
        return 1
    } else if Io_read_bit(SENSOR_FLOOR3) != 0{
        return 2
    } else if Io_read_bit(SENSOR_FLOOR4) != 0{
        return 3
    } else {
        return -1
    }
}


// 2 bit representerer 4 indikatorer.
/**
  Set floor indicator lamp for a given floor.
  @param floor Which floor lamp to turn on. Other floor lamps are turned off.
*/
func Elev_set_floor_indicator(floor int) {

  if 0 <= floor && floor < N_FLOORS {
    if floor == 0{
      Io_clear_bit(LIGHT_FLOOR_IND1)
      Io_clear_bit(LIGHT_FLOOR_IND2)
    }

    if floor == 1{
      Io_set_bit(LIGHT_FLOOR_IND2)
      Io_clear_bit(LIGHT_FLOOR_IND1)
    }
    if floor == 2{
      Io_set_bit(LIGHT_FLOOR_IND1)
      Io_clear_bit(LIGHT_FLOOR_IND2)
    }
    if floor == 3{
      Io_set_bit(LIGHT_FLOOR_IND1)
      Io_set_bit(LIGHT_FLOOR_IND2)
    }
    if floor == -1{
      //elevator not in a floor
    }
  } else {
    fmt.Printf("input floor is out of range!\n")
  }
}


/**
  Gets a button signal.
  @param button Which button type to check. Can be BUTTON_CALL_UP,
    BUTTON_CALL_DOWN or BUTTON_COMMAND (button "inside the elevator).
  @param floor Which floor to check button. Must be 0-3.
  @return 0 if button is not pushed. 1 if button is pushed.
*/
func Elev_get_button_signal(button int, floor int) int {
    
	if 0 <= floor && floor< N_FLOORS {
		if Io_read_bit(button_channel_matrix[floor][button]) != 0{
      return 1
    } else {
      return 0
    }
	} else {
		fmt.Printf("input floor is out of range!\n")
    return 0
	}
}


/**
  Set a button lamp.
  @param lamp Which type of lamp to set. Can be BUTTON_CALL_UP,
    BUTTON_CALL_DOWN or BUTTON_COMMAND (button "inside" the elevator).
  @param floor Floor of lamp to set. Must be 0-3
  @param value Non-zero value turns lamp on, 0 turns lamp off.
*/
func Elev_set_button_lamp(button int, floor int, value int) {
	if 0 <= floor && floor< N_FLOORS {
		if value != 0{
    	Io_set_bit(lamp_channel_matrix[floor][button])
    } else {
      Io_clear_bit(lamp_channel_matrix[floor][button])
		}
  } else {
		fmt.Printf("input floor is out of range!\n")
	}
}