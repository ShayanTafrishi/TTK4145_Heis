package main

import (
	"Driver-go/elevio"
	"Network-go/network/localip"
	"time"
)
var ConnectedElevators map[string]Elevator //map som inneholder alle de tilkoblede heisene

const (
	NumberFloors      int = 4
	NumberButtonTypes     = 3
	NumberElevators       = 3
)

type ElevatorState int

const (
	IDLE         ElevatorState = 0
	MOVING                     = 1
	DOOR_OPEN                  = 2
	//SYSTEM_ERROR               = 3
)
type SendNewOrder struct {
	NewOrder elevio.ButtonEvent
	Id string
}

type Elevator struct {
	Id        string
	Floor     int
	Direction elevio.MotorDirection
	State     ElevatorState
	Orders    [NumberFloors][NumberButtonTypes]int // Number floors er radene, number buttons er kolonnene
	/*
			OPP	  NED	CAB
		1
		2
		3
		4
	*/
	//Disconnected bool
	MechanicalError bool
}

func elevatorIp() string {
	ip, _ := localip.LocalIP()
	return ip
}

var Timers map[string]*time.Timer

var HardwareErrorTimer *time.Timer
