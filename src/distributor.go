package main

import (
	"Driver-go/elevio"
	"fmt"
	"time"
)

//var connectedElevators map[string]Elevator //map som inneholder alle de tilkoblede heisene
//TODO:må fylle denne på et vis
/*
func Distribute(drv_buttons chan elevio.ButtonEvent,
	assignedOrder chan elevio.ButtonEvent, assignedOrderOtherElevator chan SendNewOrder, connectedElevatorsCh chan map[string]Elevator, updatedState chan Elevator) {
	var connectedElevators map[string]Elevator
	var elevator Elevator
	//HardwareErrorTimer := time.NewTimer(30 * time.Second)

	//fmt.Println(elevator.Orders)
	for {
		select {
		case connectedElevators = <-connectedElevatorsCh:

		case elevator = <-updatedState:
			fmt.Println("te",connectedElevators)
			connectedElevators[elevator.Id] = elevator
		case newOrder := <-drv_buttons:
			fmt.Println(connectedElevators)
			if len(connectedElevators) > 1 { //den har andre heiser koblet til
				if newOrder.Button == elevio.BT_Cab {
					assignedOrder <- newOrder //del ut til seg selv
				} else {

					id := costfunc(connectedElevators, newOrder, updatedState) //bruk kostnadsfunksjon til å beregne hvilken heis som skal utføre bestillingen
					fmt.Println("costfunc ", id)
					//del ut bestillingen til denne heisen:
					if id == elevator.Id {
						assignedOrder <- newOrder
					} else {
						//TODO: bruk id til å sende til riktig heis
						sendNewOrder := SendNewOrder{NewOrder: newOrder, Id: id}
						assignedOrderOtherElevator <- sendNewOrder
						//TODO: hvordan vet man hvilken heis man skal sende til??
					}
				}
			} else {
				assignedOrder <- newOrder //del ut bestillingen til seg selv
			}


		case <-HardwareErrorTimer.C:
			fmt.Println("TIMER")
			HardwareErrorTimer.Reset(30*time.Second)
			elevator.MechanicalError = true
			transferHallOrder(connectedElevators, elevator, assignedOrderOtherElevator, updatedState) //!

		}
	}
}
*/

func costfunc(ConnectedElevators map[string]Elevator, order elevio.ButtonEvent) string {
	bestElevId := "none"
	lowestCost := 10000

	// Check for empty map?
	fmt.Println(ConnectedElevators)
	for elevId, elevator := range ConnectedElevators {
		floor := order.Floor
		for button := 0; button < NumberButtonTypes; button++ { //hvis en heis har bestillinger der fra før av, får denne bestillingen
			if elevator.Orders[floor][button] == 1 {
				fmt.Println("test1")
				return elevId
			}
		}
		if numberOfOrders(elevator) == 0 {
			fmt.Println("test2")
			return elevId
		}
		cost := numberOfOrders(elevator)
		if (elevator.Direction == elevio.MD_Up && order.Floor < elevator.Floor) ||
			(elevator.Direction == elevio.MD_Down && order.Floor > elevator.Floor) {
			cost++
			fmt.Println("test3")
		}
		if cost < lowestCost {
			lowestCost = cost
			bestElevId = elevId
		}
	}
	return bestElevId
}

func Distribute(drv_buttons chan elevio.ButtonEvent,
	assignedOrder chan elevio.ButtonEvent, assignedOrderOtherElevator chan SendNewOrder, connectedElevatorsCh chan map[string]Elevator, updatedState chan Elevator) {
	var connectedElevators map[string]Elevator
	var elevator Elevator
	HardwareErrorTimer := time.NewTimer(30 * time.Second)

	for {
		select {
		case connectedElevators = <-connectedElevatorsCh:
			/*for elevId, _ := range connectedElevators { // Overskriver i connectedElevators slik at det alltid er bare ett element??
				fmt.Println(connectedElevators[elevId])
			}*/
		case elevator = <-updatedState:
			//connectedElevators[elevator.Id] = elevator
		case newOrder := <-drv_buttons:
			//fmt.Println(connectedElevators)
			if len(connectedElevators) > 1 { //den har andre heiser koblet til
				if newOrder.Button == elevio.BT_Cab {
					assignedOrder <- newOrder //del ut til seg selv
				} else {

					id := costfunc(connectedElevators, newOrder) //bruk kostnadsfunksjon til å beregne hvilken heis som skal utføre bestillingen
					fmt.Println("costfunc ", id)
					//del ut bestillingen til denne heisen:
					if id == elevator.Id {
						assignedOrder <- newOrder
					} else {
						//TODO: bruk id til å sende til riktig heis
						sendNewOrder := SendNewOrder{NewOrder: newOrder, Id: id}
						assignedOrderOtherElevator <- sendNewOrder
						//TODO: hvordan vet man hvilken heis man skal sende til??
					}
				}
			} else {
				assignedOrder <- newOrder //del ut bestillingen til seg selv
			}

		case <-HardwareErrorTimer.C:
			fmt.Println("TIMER")
			elevator.MechanicalError = true
			fmt.Println(connectedElevators)
			transferHallOrder(connectedElevators, elevator, assignedOrderOtherElevator)
			HardwareErrorTimer.Reset(30*time.Second)
		}
	}
}

/*
func costfunc(ConnectedElevators map[string]Elevator, order elevio.ButtonEvent, updatedState chan Elevator) string {
	bestElevId := "none"
	lowestCost := 10000
	elev :=<- updatedState
	fmt.Println("test_cost")
	//Println(ConnectedElevators)
	//fmt.Println(len(ConnectedElevators))

	if len(ConnectedElevators) == 0{
		return elev.Id
	}

	for elevId, elevator := range ConnectedElevators {
		//fmt.Println(elevator)
		floor := order.Floor
		for button := 0; button < NumberButtonTypes; button++ { //hvis en heis har bestillinger der fra før av, får denne bestillingen
			if elevator.Orders[floor][button] == 1 {
				return elevId
			}
		}
		if numberOfOrders(elevator) == 0 {
			return elevId
		}
		cost := numberOfOrders(elevator)
		if (elevator.Direction == elevio.MD_Up && order.Floor < elevator.Floor) ||
			(elevator.Direction == elevio.MD_Down && order.Floor > elevator.Floor) {
			cost++
		}
		if cost < lowestCost {
			lowestCost = cost
			bestElevId = elevId
		}
	}
	return bestElevId
}*/

func transferHallOrder(ConnectedElevators map[string]Elevator, elevator Elevator, assignedOrderOtherElevator chan SendNewOrder) { // State = SYSTEM_ERROR, !! updatedState
	//var availableElevators map[string]Elevator
	availableElevators := make(map[string]Elevator)
	for elevId, elev := range ConnectedElevators {
		if elevId != elevator.Id && !elev.MechanicalError {
			availableElevators[elevId] = elev
		}
	}

	for floor := 0; floor < NumberFloors; floor++ {
		for btn := 0; btn < NumberButtonTypes-1; btn++ {
			order := elevio.ButtonEvent{Floor: floor, Button: elevio.ButtonType(btn)}
			//fmt.Println(availableElevators)
			id := costfunc(availableElevators, order)
			sendNewOrder := SendNewOrder{NewOrder: order, Id: id}
			assignedOrderOtherElevator <- sendNewOrder
			elevator.Orders[floor][btn] = 0
		}
	}
}

func hasFloorOrders(elevator Elevator) []int { //trenger kanskje ikke denne
	orders := make([]int, NumberFloors)
	//fmt.Println(elevator.Orders)

	for i := 0; i < NumberFloors; i++ {
		for j := 0; j < NumberButtonTypes; j++ {
			//fmt.Println(elevator.Orders[i][j])
			if elevator.Orders[i][j] == 1 {
				orders[i] = 1

			}
		}
	}
	//fmt.Println(orders)
	return orders
}

func numberOfOrders(elevator Elevator) int {
	orders := hasFloorOrders(elevator)
	numberOfOrders := 0
	for i := 0; i < len(orders); i++ {
		if orders[i] == 1 {
			numberOfOrders++
		}
	}
	return numberOfOrders
}

/*
func connectElevators(elevator Elevator, connectedElevators map[string]Elevator) {
	id := elevator.Id
	connectedElevators[id] = elevator

}*/

/*
func DisconnectElevators(elevator Elevator, ConnectedElevators map[string]Elevator) {
	delete(ConnectedElevators, elevator.Id)
}*/
