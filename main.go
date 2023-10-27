package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"
)

type Ticket struct {
	BusTime   string
	Seat      int
	Available bool
}
type Passenger struct {
	ID      int
	BusTime string
	Seat    int
}
type Bus struct {
	BusTime        string
	AvailableSeats int
}
type TicketSystem struct {
	Tickets     []Ticket
	Passengers  []Passenger
	Buses       []Bus
	Mutex       sync.Mutex
	PassengerWG sync.WaitGroup
}

var totalPassageirosReservados int
var nextPassengerID int

func generatePassengerID() int {
	nextPassengerID++
	return nextPassengerID
}

func generateBuses() []Bus {
	buses := make([]Bus, 0)

	for hour := 7; hour <= 21; hour++ {
		for minute := 0; minute < 60; minute += 60 {
			busTime := fmt.Sprintf("%02d:%02d", hour, minute)
			bus := Bus{
				BusTime:        busTime,
				AvailableSeats: 40,
			}
			buses = append(buses, bus)
		}
	}
	return buses
}

func GeradorPassageiros(system *TicketSystem, numPassageiros int) {
	for i := 0; i < numPassageiros; i++ {
		if totalPassageirosReservados >= numPassageiros {
			break //foram todos
		}

		time.Sleep(time.Millisecond * time.Duration(rand.Intn(1000)))

		system.Mutex.Lock()
		horarios := make([]string, 0)
		for _, bus := range system.Buses {
			if bus.AvailableSeats > 0 {
				horarios = append(horarios, bus.BusTime)
			}
		}

		if len(horarios) == 0 {
			system.Mutex.Unlock()
			continue // NAO TEM HORARIO DISPONVEL PARA ESSE PASSAGEIRO
		}

		horarioEscolhido := horarios[rand.Intn(len(horarios))]

		// Lógica para gerar um passageiro com um horário de ônibus e uma poltrona aleatórios
		// Marque a poltrona como indisponível e adicione o passageiro ao sistema
		novoPassageiro := Passenger{
			BusTime: horarioEscolhido,
			Seat:    rand.Intn(40), // PRECISA CERTIFICAR QUE ESTA DISPONIVEL ESSA
			ID:      generatePassengerID(),
		}

		system.Passengers = append(system.Passengers, novoPassageiro)
		totalPassageirosReservados++
		system.Mutex.Unlock()
		system.PassengerWG.Add(1)
	}
}
func verificar_poltronas_disponiveis(system *TicketSystem, horario string) []int {
	poltronasDisponiveis := []int{}

	system.Mutex.Lock()
	defer system.Mutex.Unlock()

	for _, ticket := range system.Tickets {
		if ticket.BusTime == horario && ticket.Available {
			poltronasDisponiveis = append(poltronasDisponiveis, ticket.Seat)
		}
	}

	return poltronasDisponiveis
}
func reservar_passagem(system *TicketSystem, horario string, poltrona int, passengerID int) bool {
	system.Mutex.Lock()
	defer system.Mutex.Unlock()

	for i, ticket := range system.Tickets {
		if ticket.BusTime == horario && ticket.Seat == poltrona && ticket.Available {
			system.Tickets[i].Available = false

			fmt.Printf("PASSAGEIRO %d RESERVOU A POLTRONA %d DO ÔNIBUS PARTINDO AS %s HORAS\n", passengerID, poltrona, horario)

			return true
		}
	}
	// poltrona indisponível ou horário inexistente
	return false
}
func main() {
	// tratamentos para argumento invalido
	if len(os.Args) != 2 {
		fmt.Println("Uso: programa <quantidade de passageiros>")
		os.Exit(1)
	}

	numPassageiros, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println("O numero de passageiros deve ser um inteiro valido")
		os.Exit(1)
	}
	totalPassageirosReservados = numPassageiros

	system := &TicketSystem{}
	system.Buses = generateBuses()

	go GeradorPassageiros(system, numPassageiros)
	system.PassengerWG.Wait()

}
