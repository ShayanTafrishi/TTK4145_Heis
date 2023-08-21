package main

import (
	"Driver-go/elevio"
	"os"
	//"io/ioutil"
	"strconv"
	"strings"
)

//funksjon med switch for å ta inn hall og cab ordre
//en matrise med heisens bestillinger
//funksjon for å ta imot de andres struct, og lagrer bestillingene de har
//funksjon for å sende egen struct
//funksjon for å oppdatere lys dersom andre har ny hall-bestilling

func ReadBackup(elevator Elevator, assignedOrderOtherElevator chan SendNewOrder) { //kan hende vi må legge til delay med sleep
	// Mulig vi må bruke en annen kanal og sende på siden denne brukes i Distributer

	filename := "backup"+elevator.Id+".txt" //lager en fil
	content, _ := os.ReadFile(filename)
	contentStr := string(content)
	caborders := make([]int, NumberFloors)

	//må først gjøre om fra string og splitte opp til en array av ints igjen
	//må så sjekke om det er noen bestillinger, og vil isåfall legge dem til i orders

	s := strings.Split(contentStr, ",")

	for i := 0; i < len(s); i++ {
		caborders[i], _ = strconv.Atoi(s[i])
	}

	for i := 0; i < NumberFloors; i++ { //har ikke testet denne delen
		if caborders[i] == 1 {
			backupOrder := elevio.ButtonEvent{
				Floor:  i,
				Button: elevio.BT_Cab,
			}
			backup := SendNewOrder {
				NewOrder : backupOrder,
				Id : elevator.Id,
			}
			assignedOrderOtherElevator <- backup
		}
	}

}

func writeToBackup(elevator Elevator) {
	filename := "backup"+elevator.Id+".txt"
	file, _ := os.Create(filename)

	caborders := make([]int, NumberFloors) //lager en en liste med lengde 4, der hvert element skal representere hver etasje. 1 hvis bestilling i etasjen, ellers 0

	//flytter aktive caborders over i caborders
	for i := 0; i < 4; i++ {
		if elevator.Orders[2][i] == 1 {
			caborders[i] = 1
		} else {
			caborders[i] = 0 //usikker på om dette er nødvendig, mtp. om den er 0 fra før av?
		}
	}

	//gjør dette om til en tekststring for at det skal skrives til fil
	var content string
	for i := 0; i < NumberFloors; i++ {
		content += strconv.Itoa(caborders[i])
		if i != NumberFloors-1 {
			content += ","
		}
	}
	file.WriteString(content)
	defer file.Close()
}