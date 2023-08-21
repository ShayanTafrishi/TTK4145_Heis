package main

import (
	"Driver-go/elevio"
	"fmt"
	"time"
)

func ElevatorOperation2(assignedOrder chan elevio.ButtonEvent, drv_floors chan int, drv_obstr chan bool, elevator Elevator, updatedState chan Elevator) {
	doorTimer := time.NewTimer(3 * time.Second)
	doorTimer.Stop()
	//HardwareErrorTimer := time.NewTimer(30*time.Second)
	var currentDirection elevio.MotorDirection

	for {
		switch elevator.State {
		case IDLE:
			updatedState <- elevator
			currentDirection = elevator.Direction
			select {
			case newOrder := <-assignedOrder:
				fmt.Printf("%+v\n", newOrder)
				updatedState <- elevator
				elevio.SetButtonLamp(newOrder.Button, newOrder.Floor, true)
				if elevator.Floor == newOrder.Floor {

					elevio.SetButtonLamp(elevio.BT_HallUp, elevator.Floor, false)   //TEST
					elevio.SetButtonLamp(elevio.BT_HallDown, elevator.Floor, false) //TEST
					elevio.SetButtonLamp(elevio.BT_Cab, elevator.Floor, false)      //TEST

					elevator.State = DOOR_OPEN
					doorTimer.Reset(3 * time.Second)
				} else {
					elevator.Orders[newOrder.Floor][newOrder.Button] = 1
					elevator.State = MOVING
					//HardwareErrorTimer.Reset(10 * time.Second)
					//elevator.MechanicalError = false
					elevator.Direction = ChooseDirection(elevator, currentDirection)
					elevio.SetMotorDirection(elevator.Direction)
				}
				break
			case obstruction := <-drv_obstr:
				updatedState <- elevator
				elevator.Direction = ChooseDirection(elevator, currentDirection)
				elevio.SetMotorDirection(elevator.Direction)
				if obstruction {
					fmt.Println("Obstruction")
					doorTimer.Reset(3 * time.Second)
					elevator.Direction = elevio.MD_Stop
					elevio.SetMotorDirection(elevator.Direction)
					elevator.State = DOOR_OPEN
				}
			}
		case MOVING:
			updatedState <- elevator
			select {
			case newOrder := <-assignedOrder:
				updatedState <- elevator
				elevator.Orders[newOrder.Floor][newOrder.Button] = 1
				fmt.Printf("%+v\n", newOrder)
				elevio.SetButtonLamp(newOrder.Button, newOrder.Floor, true)
				break
			case newFloor := <-drv_floors:
				updatedState <- elevator
				//HardwareErrorTimer.Reset(10 * time.Second)
				//elevator.MechanicalError = false
				elevator.Floor = newFloor
				if StopElevator(elevator) {
					parameters := ClearOrders{}
					elevator = ClearOrderAtCurrentFloor(parameters, elevator)

					elevio.SetButtonLamp(elevio.BT_HallUp, newFloor, false)   //TEST
					elevio.SetButtonLamp(elevio.BT_HallDown, newFloor, false) //TEST
					elevio.SetButtonLamp(elevio.BT_Cab, newFloor, false)      //TEST

					// Kanal som sender bekreftelse på at denne bestillingen er gjort (unødvendig?)
					// sendAckCh <- true (ta inn kanal som input)

					currentDirection = elevator.Direction
					elevator.Direction = elevio.MD_Stop
					elevio.SetMotorDirection(elevator.Direction)
					elevator.State = DOOR_OPEN
					doorTimer.Reset(3 * time.Second)
				}
				break
			}
		case DOOR_OPEN:
			updatedState <- elevator
			select {
			case newOrder := <-assignedOrder: //må heller ha assigned order her (og alle steder i stedet for assignedOrder)
				updatedState <- elevator
				fmt.Printf("%+v\n", newOrder)
				elevio.SetButtonLamp(newOrder.Button, newOrder.Floor, true)
				if elevator.Floor == newOrder.Floor {

					elevio.SetButtonLamp(elevio.BT_HallUp, elevator.Floor, false)   //TEST
					elevio.SetButtonLamp(elevio.BT_HallDown, elevator.Floor, false) //TEST
					elevio.SetButtonLamp(elevio.BT_Cab, elevator.Floor, false)      //TEST

					elevator.State = DOOR_OPEN
					doorTimer.Reset(3 * time.Second)
				} else {
					elevator.Orders[newOrder.Floor][newOrder.Button] = 1
				}
				break
			case <-doorTimer.C:
				elevator.Direction = ChooseDirection(elevator, currentDirection)
				elevio.SetMotorDirection(elevator.Direction)

				if elevator.Direction == elevio.MD_Stop {
					elevator.State = IDLE
					// Hvis IDLE og ikke bestillinger, reset timer
					//HardwareErrorTimer.Reset(10 * time.Second)
					//elevator.MechanicalError = false
				} else {
					elevator.State = MOVING
					//HardwareErrorTimer.Reset(10 * time.Second)
					//elevator.MechanicalError = false
				}
				break
			case obstruction := <-drv_obstr:
				elevator.Direction = ChooseDirection(elevator, currentDirection)
				elevio.SetMotorDirection(elevator.Direction)
				if obstruction {
					fmt.Println("Obstruction")
					doorTimer.Reset(3 * time.Second)
					elevator.Direction = elevio.MD_Stop
					elevio.SetMotorDirection(elevator.Direction)
					elevator.State = DOOR_OPEN
				}
			}
		}
	}
}
