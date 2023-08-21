package main

import (
	"Driver-go/elevio"
)

func InitElevator(drv_floors chan int) {
	ClearAllLights()
	elevio.SetMotorDirection(elevio.MD_Down) //må stå over løkka, pga får veldig mange iterasjoner ellers
	for elevio.GetFloor() == -1 {
	}
	elevio.SetMotorDirection(elevio.MD_Stop)
	elevio.SetFloorIndicator(elevio.GetFloor())

}

func ClearAllLights() {
	elevio.SetDoorOpenLamp(false)
	for floor := 0; floor < NumberButtonTypes; floor++ {
		for button := 0; button < NumberButtonTypes; button++ {
			elevio.SetButtonLamp(elevio.ButtonType(button), floor, false)
		}
	}
}


func StopElevator(elevator Elevator) bool {
	switch elevator.Direction {
	case elevio.MD_Down:
		if elevator.Orders[elevator.Floor][elevio.BT_HallDown] == 1 || elevator.Orders[elevator.Floor][elevio.BT_Cab] == 1 || !NewOrderBelowCurrentFloor(elevator) {
			return true
		} else {
			return false
		}
	case elevio.MD_Up:
		if elevator.Orders[elevator.Floor][elevio.BT_HallUp] == 1 || elevator.Orders[elevator.Floor][elevio.BT_Cab] == 1 || !NewOrderAboveCurrentFloor(elevator) {
			return true
		} else {
			return false
		}
	case elevio.MD_Stop: //?
	}

	return true
}

func NewOrderAboveCurrentFloor(elevator Elevator) bool {
	for floor := elevator.Floor + 1; floor < NumberFloors; floor++ { // All floor numbers above the current floor
		for button := 0; button < NumberButtonTypes; button++ { // All buttons
			if elevator.Orders[floor][button] == 1 {
				return true
			}
		}
	}
	return false // No orders above the current floor
}

func NewOrderBelowCurrentFloor(elevator Elevator) bool {
	for floor := 0; floor < elevator.Floor; floor++ {
		for button := 0; button < NumberButtonTypes; button++ {
			if elevator.Orders[floor][button] == 1 {
				return true
			}
		}
	}
	return false // No orders below the current floor
}

func ChooseDirection(elevator Elevator, currentDirection elevio.MotorDirection) elevio.MotorDirection {
	switch currentDirection { //?
	case elevio.MD_Up: // If current direction is up
		if NewOrderAboveCurrentFloor(elevator) {
			return elevio.MD_Up
		} else if NewOrderBelowCurrentFloor(elevator) {
			return elevio.MD_Down
		} else {
			return elevio.MD_Stop
		}
	case elevio.MD_Down:
		if NewOrderBelowCurrentFloor(elevator) {
			return elevio.MD_Down
		} else if NewOrderAboveCurrentFloor(elevator) {
			return elevio.MD_Up
		} else {
			return elevio.MD_Stop
		}
	case elevio.MD_Stop:
		if NewOrderAboveCurrentFloor(elevator) {
			return elevio.MD_Up
		} else if NewOrderBelowCurrentFloor(elevator) {
			return elevio.MD_Down
		} else {
			return elevio.MD_Stop
		}
	}
	return elevio.MD_Stop
}

type ClearOrders struct {
	ifEqual func(elevio.ButtonType, int) // ButtonType and floor number ?? ikke struct
}

func ClearOrderAtCurrentFloor(parameters ClearOrders, elevator Elevator) Elevator {
	elevator.Orders[elevator.Floor][elevio.BT_Cab] = 0
	haveFunction := !(parameters.ifEqual == nil)
	switch elevator.Direction {
	case elevio.MD_Up:
		if haveFunction {
			parameters.ifEqual(elevio.BT_HallUp, elevator.Floor)
		}
		elevator.Orders[elevator.Floor][elevio.BT_HallUp] = 0
		if !NewOrderAboveCurrentFloor(elevator) {
			if haveFunction {
				parameters.ifEqual(elevio.BT_HallDown, elevator.Floor)
			}
			elevator.Orders[elevator.Floor][elevio.BT_HallDown] = 0
		}
		break
	case elevio.MD_Down:
		if haveFunction {
			parameters.ifEqual(elevio.BT_HallUp, elevator.Floor)
		}
		elevator.Orders[elevator.Floor][elevio.BT_HallDown] = 0
		if !NewOrderBelowCurrentFloor(elevator) {
			if haveFunction {
				parameters.ifEqual(elevio.BT_HallUp, elevator.Floor)
			}
			elevator.Orders[elevator.Floor][elevio.BT_HallUp] = 0
		}
		break
	case elevio.MD_Stop:
		if haveFunction {
			parameters.ifEqual(elevio.BT_HallUp, elevator.Floor)
			parameters.ifEqual(elevio.BT_HallDown, elevator.Floor)
		}
		elevator.Orders[elevator.Floor][elevio.BT_HallUp] = 0
		elevator.Orders[elevator.Floor][elevio.BT_HallDown] = 0
		break
	default:
		if haveFunction {
			parameters.ifEqual(elevio.BT_HallUp, elevator.Floor)
			parameters.ifEqual(elevio.BT_HallDown, elevator.Floor)
		}
		elevator.Orders[elevator.Floor][elevio.BT_HallUp] = 0
		elevator.Orders[elevator.Floor][elevio.BT_HallDown] = 0
		break
	}
	return elevator
}