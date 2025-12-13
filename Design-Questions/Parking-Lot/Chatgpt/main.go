package main

import (
	"errors"
	"fmt"
	"time"
)

//////////////////// ENUMS ////////////////////

type VehicleType string
type SpotStatus string

const (
	Bicycle VehicleType = "BICYCLE"
	Bike    VehicleType = "BIKE"
	Car     VehicleType = "CAR"
	SUV     VehicleType = "SUV"
)

const (
	Free     SpotStatus = "FREE"
	Occupied SpotStatus = "OCCUPIED"
)

//////////////////// ENTITIES ////////////////////

type Vehicle struct {
	Number string
	Type   VehicleType
}

type ParkingSpot struct {
	ID     int
	Type   VehicleType
	Status SpotStatus
	Vehicle *Vehicle
}

func (p *ParkingSpot) CanFit(v *Vehicle) bool {
	return p.Status == Free && p.Type == v.Type
}

func (p *ParkingSpot) Park(v *Vehicle) error {
	if !p.CanFit(v) {
		return errors.New("spot not available")
	}
	p.Vehicle = v
	p.Status = Occupied
	return nil
}

func (p *ParkingSpot) Vacate() {
	p.Vehicle = nil
	p.Status = Free
}

type ParkingTicket struct {
	ID        string
	Vehicle   *Vehicle
	SpotID    int
	EntryTime time.Time
	ExitTime  time.Time
}

//////////////////// STRATEGIES ////////////////////

// --- Spot Allocation Strategy ---
type SpotAllocator interface {
	FindSpot(spots []*ParkingSpot, v *Vehicle) (*ParkingSpot, error)
}

type FirstFreeSpotAllocator struct{}

func (f *FirstFreeSpotAllocator) FindSpot(spots []*ParkingSpot, v *Vehicle) (*ParkingSpot, error) {
	for _, s := range spots {
		if s.CanFit(v) {
			return s, nil
		}
	}
	return nil, errors.New("no spot available")
}

// --- Pricing Strategy ---
type PricingStrategy interface {
	Calculate(entry, exit time.Time) int
}

type HourlyPricing struct{}

func (h *HourlyPricing) Calculate(entry, exit time.Time) int {
	hours := int(exit.Sub(entry).Hours())
	if hours == 0 {
		hours = 1
	}
	return hours * 50
}

//////////////////// PARKING LOT (SERVICE) ////////////////////

type ParkingLot struct {
	Spots       []*ParkingSpot
	Tickets     map[string]*ParkingTicket
	Allocator   SpotAllocator
	Pricing     PricingStrategy
}

func NewParkingLot(
	spots []*ParkingSpot,
	allocator SpotAllocator,
	pricing PricingStrategy,
) *ParkingLot {
	return &ParkingLot{
		Spots:     spots,
		Tickets:   make(map[string]*ParkingTicket),
		Allocator: allocator,
		Pricing:   pricing,
	}
}

func (p *ParkingLot) ParkVehicle(v *Vehicle) (*ParkingTicket, error) {
	spot, err := p.Allocator.FindSpot(p.Spots, v)
	if err != nil {
		return nil, err
	}

	err = spot.Park(v)
	if err != nil {
		return nil, err
	}

	ticket := &ParkingTicket{
		ID:        fmt.Sprintf("%s-%d", v.Number, time.Now().UnixNano()),
		Vehicle:   v,
		SpotID:    spot.ID,
		EntryTime: time.Now(),
	}

	p.Tickets[ticket.ID] = ticket
	return ticket, nil
}

func (p *ParkingLot) ExitVehicle(ticketID string) (int, error) {
	ticket, ok := p.Tickets[ticketID]
	if !ok {
		return 0, errors.New("invalid ticket")
	}

	ticket.ExitTime = time.Now()

	var spot *ParkingSpot
	for _, s := range p.Spots {
		if s.ID == ticket.SpotID {
			spot = s
			break
		}
	}

	if spot == nil {
		return 0, errors.New("spot not found")
	}

	spot.Vacate()
	amount := p.Pricing.Calculate(ticket.EntryTime, ticket.ExitTime)
	delete(p.Tickets, ticketID)

	return amount, nil
}

//////////////////// MAIN ////////////////////

func main() {
	spots := []*ParkingSpot{
		{ID: 1, Type: Car, Status: Free},
		{ID: 2, Type: Bike, Status: Free},
	}

	parkingLot := NewParkingLot(
		spots,
		&FirstFreeSpotAllocator{},
		&HourlyPricing{},
	)

	vehicle := &Vehicle{Number: "UP32AA1234", Type: Car}

	ticket, _ := parkingLot.ParkVehicle(vehicle)
	fmt.Println("Ticket Issued:", ticket.ID)

	time.Sleep(2 * time.Second)

	amount, _ := parkingLot.ExitVehicle(ticket.ID)
	fmt.Println("Total Amount:", amount)
}
