package main

import (
	"Driver-go/elevio"
	"Network-go/network/peers"
	//"fmt"
	//"fmt"
	"os"
	"time"
)

func SendData(orderTx chan SendNewOrder, assignedOrderOtherElevator chan SendNewOrder,
	updatedState chan Elevator, elevatorTx chan Elevator) {
	//interval := 15 * time.Millisecond

	for {
		select {
		case sendOrder := <-assignedOrderOtherElevator:
			for i := 0; i < 10; i++ { //spammer pakken (antar at dette forhindrer pakketap)
				orderTx <- sendOrder
			}

		case sendState := <-updatedState:
			for i := 0; i < 10; i++ {
				elevatorTx <- sendState
			}
			//time.Sleep(time.Second)
		}
	}
}

func ReceiveData(elevatorRx chan Elevator, orderRx chan SendNewOrder, assignedOrder chan elevio.ButtonEvent, ConnectedElevators map[string]Elevator, connectedElevatorsCh chan map[string]Elevator, updatedState chan Elevator, peerUpdateCh chan peers.PeerUpdate, assignedOrderOtherElevator chan SendNewOrder) {
	elevator := <-updatedState
	//fmt.Println(elevator.Orders)
	Timers := make(map[string]*time.Timer)
	Timers[elevator.Id] = time.NewTimer(10 * time.Millisecond)

	for {
		select {
		case receiveOrder := <-orderRx:
			if receiveOrder.Id == elevator.Id {
				//fmt.Printf("test")
				assignedOrder <- receiveOrder.NewOrder
				ConnectedElevators[elevator.Id] = elevator
				connectedElevatorsCh <- ConnectedElevators // Kanskje noe feil med kanalen?
			}
			
		case receiveState := <-elevatorRx:
			if receiveState.Id != elevator.Id {

				//fmt.Println("Received state from:", receiveState.Id)
				//fmt.Println("Elevator", receiveState.Id, "has orders", receiveState.Orders)
				ConnectedElevators[receiveState.Id] = receiveState // sjekk om det faktisk legges til en ny verdi i mappet
				//fmt.Println(ConnectedElevators)
				//fmt.Println(ConnectedElevators[receiveState.Id])
				connectedElevatorsCh <- ConnectedElevators

				_, found := Timers[receiveState.Id]

				if !found {
					Timers[receiveState.Id] = time.NewTimer(10 * time.Millisecond)
				} else {
					Timers[receiveState.Id].Reset(10 * time.Millisecond)
				}
			}
		case elevator := <-updatedState: //!!!!!!! case 
			ConnectedElevators[elevator.Id] = elevator
			connectedElevatorsCh <- ConnectedElevators
		
		case peer := <- peerUpdateCh:
			if peer.New != "" {
				newElevator := Elevator{
					Id:        peer.New,
					Floor:     0,
					Direction: elevio.MD_Stop,
					State:  IDLE,
					Orders: [NumberFloors][NumberButtonTypes]int{}, 
					MechanicalError: false,} // Backup cab orders?
				
				ConnectedElevators[peer.New] = newElevator
				connectedElevatorsCh <- ConnectedElevators
				if _, err := os.Stat("backup"+peer.New+".txt"); err == nil{ // Endre path
					ReadBackup(newElevator, assignedOrderOtherElevator)
				}
			if len(peer.Lost) > 0 {
				for _, lostElevatorId := range peer.Lost {
					writeToBackup(ConnectedElevators[lostElevatorId])
					delete(ConnectedElevators, lostElevatorId)
					transferHallOrder(ConnectedElevators, ConnectedElevators[lostElevatorId], assignedOrderOtherElevator)
				}
				connectedElevatorsCh <- ConnectedElevators
			}	
			}

		//case <-Timers[elevator.Id].C:
		//	DisconnectElevators(elevator, ConnectedElevators)
			//connectedElevatorsCh <- ConnectedElevators
			// Omfordele bestillinger til heisen som ikke lenger er koblet til netverket
		
		}
	}
}
/*
func Timer(connectedElevatorsCh chan map[string]Elevator, updatedState chan Elevator) {
	connectedElevators :=<- connectedElevatorsCh
	elevator := <- updatedState
	Timers := make(map[string]*time.Timer)
	Timers[elevator.Id] = time.NewTimer(10 * time.Millisecond)
	for {
		for elevId, elev := range connectedElevators {
			select{
			case <-Timers[elevId].C:
				if elevId != elevator.Id {
					DisconnectElevators(elev, connectedElevators)
				}
			}
	}
}
}*/
