package main

import (
	"Driver-go/elevio"
	"fmt"
	//"strconv"
	//"time"

	//"main/distributor"
	"Network-go/network/bcast"
	"Network-go/network/peers"
	"flag"
)

func main() {
	numFloors := 4

	//elevio.Init("localhost:15657", numFloors)
	var elevio_port string
	var Id string
	flag.StringVar(&Id, "id", "", "id of this peer")

	flag.StringVar(&elevio_port, "elevio_port", "15657", "port to elevator server")
	flag.Parse()
	fmt.Println(elevio_port)

	elevio.Init("localhost:"+elevio_port, numFloors)

	drv_buttons := make(chan elevio.ButtonEvent)
	drv_floors := make(chan int)
	drv_obstr := make(chan bool)

	go elevio.PollButtons(drv_buttons)
	go elevio.PollFloorSensor(drv_floors)
	go elevio.PollObstructionSwitch(drv_obstr)

	InitElevator(drv_floors)

	
	floorNumber := <-drv_floors

	elevator := Elevator{
		Id:        Id,
		Floor:     floorNumber,
		Direction: elevio.MD_Stop,

		State:  IDLE,
		Orders: [NumberFloors][NumberButtonTypes]int{},
		MechanicalError : false}
	
		
	ConnectedElevators := make(map[string]Elevator)
	ConnectedElevators[elevator.Id] = elevator

	assignedOrder := make(chan elevio.ButtonEvent)
	connectedElevatorsCh := make(chan map[string]Elevator)
	
	
	assignedOrderOtherElevator := make(chan SendNewOrder)
	updatedState := make(chan Elevator)
	orderTx := make(chan SendNewOrder)
	orderRx := make(chan SendNewOrder)
	elevatorTx := make(chan Elevator)
	elevatorRx := make(chan Elevator)
	
	peerUpdateCh := make(chan peers.PeerUpdate)

	go bcast.Transmitter(16569, elevatorTx) 
	go bcast.Receiver(16569, elevatorRx)
	go bcast.Transmitter(16568, orderTx) //!!!!
	go bcast.Receiver(16568, orderRx)  //!!!!
	go ElevatorOperation2(assignedOrder, drv_floors, drv_obstr, elevator, updatedState) 
	go SendData(orderTx, assignedOrderOtherElevator, updatedState, elevatorTx) 
	go ReceiveData(elevatorRx, orderRx, assignedOrder,  ConnectedElevators, connectedElevatorsCh, updatedState, peerUpdateCh, assignedOrderOtherElevator) 

	/*
	elevator1 := Elevator{
		Id:        "1",
		Floor:     floorNumber,
		Direction: elevio.MD_Stop,
		State:  IDLE,
		Orders: [NumberFloors][NumberButtonTypes]int{}}
	elevator2 := Elevator{
		Id:        "2",
		Floor:     floorNumber,
		Direction: elevio.MD_Stop,
		State:  IDLE,
		Orders: [NumberFloors][NumberButtonTypes]int{}}
	connectedElevators := map[string]Elevator {
		elevator1.Id : elevator1, elevator2.Id : elevator2,  
	}*/
	go Distribute(drv_buttons, assignedOrder, assignedOrderOtherElevator, connectedElevatorsCh, updatedState) //!
	connectedElevatorsCh <- ConnectedElevators
	//go Timer(connectedElevatorsCh, updatedState)
	select {}

}
