package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"sync"
)

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
	Passengers  []Passenger
	Buses       []Bus
	Mutex       sync.Mutex
	PassengerWG sync.WaitGroup
}

var totalPassageirosReservados int

// PASSO 2.7
func generateBuses() []Bus {
	buses := make([]Bus, 0)

	for hour := 7; hour <= 21; hour++ {
		busTime := fmt.Sprintf("%02d:00", hour)
		bus := Bus{
			BusTime:        busTime,
			AvailableSeats: 40,
		}
		buses = append(buses, bus)
	}
	return buses
}

// PASSO 2.3
func GeradorPassageiros(system *TicketSystem, numPassageiros int) {
	for i := 0; i < numPassageiros; i++ {
		// PASSO 2.6
		passageiro := Passenger{ID: i + 1}

		go func() {
			defer system.PassengerWG.Done()
			horarios := make([]string, 0)
			for _, bus := range system.Buses {
				if bus.AvailableSeats > 0 {
					horarios = append(horarios, bus.BusTime)
				}
			}
			if len(horarios) == 0 {
				return
			}
			horarioEscolhido := horarios[rand.Intn(len(horarios))]
			poltronasDisponiveis := verificar_poltronas_disponiveis(system, horarioEscolhido)
			if len(poltronasDisponiveis) == 0 {
				return
			}
			poltronaEscolhida := poltronasDisponiveis[rand.Intn(len(poltronasDisponiveis))] + 1
			passageiro.BusTime = horarioEscolhido
			passageiro.Seat = poltronaEscolhida
			reservar := reservar_passagem(system, passageiro.BusTime, passageiro.Seat, passageiro.ID)
			if reservar {
				fmt.Printf("PASSAGEIRO %d RESERVOU A POLTRONA %d DO ÔNIBUS PARTINDO ÀS %s HORAS\n", passageiro.ID, poltronaEscolhida, horarioEscolhido)
			}
			system.Passengers = append(system.Passengers, passageiro)
		}()
	}
}

// PASSO 2.4
func verificar_poltronas_disponiveis(system *TicketSystem, horarioEscolhido string) []int {
	poltronasDisponiveis := []int{}

	system.Mutex.Lock()
	defer system.Mutex.Unlock()

	for _, bus := range system.Buses {
		if bus.BusTime == horarioEscolhido && bus.AvailableSeats > 0 {
			for i := 0; i < bus.AvailableSeats; i++ {
				poltronasDisponiveis = append(poltronasDisponiveis, i)
			}
		}
	}

	return poltronasDisponiveis
}

// PASSO 2.5
func reservar_passagem(system *TicketSystem, horarioEscolhido string, poltronaEscolhida int, passengerID int) bool {
	system.Mutex.Lock()
	defer system.Mutex.Unlock()
	for i, bus := range system.Buses {
		if bus.BusTime == horarioEscolhido && bus.AvailableSeats > 0 {
			system.Buses[i].AvailableSeats--
			system.Passengers = append(system.Passengers, Passenger{ID: passengerID, BusTime: horarioEscolhido, Seat: poltronaEscolhida})
			totalPassageirosReservados++
			return true
		}
	}

	return false
	// Não encontrou o ônibus ou todos os assentos estão ocupados
}

func main() {
	// tratamentos para argumento inválido
	if len(os.Args) != 2 {
		fmt.Println("Uso: programa <quantidade de passageiros>")
		os.Exit(1)
	}

	numPassageiros, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println("O número de passageiros deve ser um inteiro válido")
		os.Exit(1)
	}

	system := &TicketSystem{}
	system.Buses = generateBuses()
	system.PassengerWG.Add(numPassageiros)
	go GeradorPassageiros(system, numPassageiros)
	system.PassengerWG.Wait()
}
