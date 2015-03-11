package driver  // where "driver" is the folder that contains io.go, io.c, io.h, channels.go, channels.c and driver.go
/*
#cgo LDFLAGS: -lcomedi -lm
#include "io.h"
*/
import "C"

/**
  Initialize libComedi in "Sanntidssalen"
  @return Non-zero on success and 0 on failure
*/
func Io_init() int{
	return int (int(C.io_init()))
}

/**
  Sets a digital channel bit.
  @param channel Channel bit to set.
*/
func Io_set_bit(channel int){
	C.io_set_bit(C.int(channel))
}

/**
  Clears a digital channel bit.
  @param channel Channel bit to set.
*/
func Io_clear_bit(channel int){
	C.io_clear_bit(C.int(channel))
}

/**
  Writes a value to an analog channel.
  @param channel Channel to write to.
  @param value Value to write.
*/
func Io_write_analog(channel int, value int){
	C.io_write_analog(C.int(channel), C.int(value))
}

/**
  Reads a bit value from a digital channel.
  @param channel Channel to read from.
  @return Value read.
*/
func Io_read_bit(channel int) int { //return type var int, men comiler klaget
	return int(C.io_read_bit(C.int(channel)))
}

/**
  Reads a bit value from an analog channel.
  @param channel Channel to read from.
  @return Value read.
*/
func Io_read_analog(channel int) int { //return type var int, men comiler klaget
	return int(C.io_read_analog(C.int(channel)))
}