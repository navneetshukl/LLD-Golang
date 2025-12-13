package main

import (
	"log"
	"time"
)

type VehicleType string
type ParkingStatus string

const (
	Bicycle   VehicleType = "Bicycle"
	MotoCycle VehicleType = "MotorCycle"
	Car       VehicleType = "Car"
	SUV       VehicleType = "SUV"
)

const (
	Free     ParkingStatus = "Free"
	Occupied ParkingStatus = "Occupied"
)

type Vehicle struct {
	VehicleNumber string
	VehicleType   VehicleType
}

func NewVehicle(num string, ty VehicleType) *Vehicle {
	return &Vehicle{
		VehicleNumber: num,
		VehicleType:   ty,
	}
}

func (v *Vehicle) GetVehicleNumber() string {
	return v.VehicleNumber
}

func (v *Vehicle) GetVehicleType() string {
	return string(v.VehicleType)
}

type ParkingSpot struct {
	Status     ParkingStatus
	SpotType   VehicleType
	SpotNumber int
	Vehicle    *Vehicle
}

func NewParkingSpot(v VehicleType, sp int) *ParkingSpot {
	return &ParkingSpot{
		Status:     Free,
		SpotType:   v,
		SpotNumber: sp,
	}
}

func (p *ParkingSpot) ChangeStatus(st ParkingStatus) {
	p.Status = st
}

func (p *ParkingSpot) GetStatus() string {
	return string(p.Status)
}

func (p *ParkingSpot) GetParkingType() string {
	return string(p.SpotType)
}

func (p *ParkingSpot) ParkVehicle(v *Vehicle) {
	p.Vehicle = v
}

func (p *ParkingSpot) RemoveVehicle() {
	p.Vehicle = nil
}

type VehicleTicket struct {
	VehicleId     string
	VehicleNumber string
	SpotNumber    int
	EntryTime     time.Time
}

func NewVehicleTicket(id, num string, spot int, ti time.Time) *VehicleTicket {
	return &VehicleTicket{
		VehicleId:     id,
		VehicleNumber: num,
		SpotNumber:    spot,
		EntryTime:     time.Now(),
	}
}

type ParkingManager struct {
	Spots         []*ParkingSpot
	VehicleTicket map[string]*VehicleTicket
}

func NewParkingManager() *ParkingManager {
	return &ParkingManager{
		Spots:         make([]*ParkingSpot, 0),
		VehicleTicket: make(map[string]*VehicleTicket),
	}
}

func (p *ParkingManager) AddSpots(spot ...*ParkingSpot) {
	for _, v := range spot {
		p.Spots = append(p.Spots, v)
	}
}

func (p *ParkingManager) EntryVehicle(ve Vehicle) bool {
	// check if there is any parking spot

	for _, s := range p.Spots {
		if s.GetStatus() == "Free" && ve.GetVehicleType() == string(ve.VehicleType) {

			vehicleId := ve.VehicleNumber + string(time.Now().UnixMilli())
			vehdetails := NewVehicleTicket(vehicleId, ve.GetVehicleNumber(), s.SpotNumber, time.Now())
			p.VehicleTicket[vehicleId] = vehdetails
			s.ChangeStatus(Occupied)
			return true

		}
	}
	return false
}

func (p *ParkingManager) ExitVehicle(id string) bool {
	// get vehicle by id

	vehicle := p.VehicleTicket[id]
	timeParked := time.Now().Sub(vehicle.EntryTime)
	spotNumber := vehicle.SpotNumber
	for _, s := range p.Spots {
		if s.SpotNumber == spotNumber {
			s.ChangeStatus(Free)
		}
	}
	calculateCost := timeParked * 10
	log.Println("Total Time Parked is ", calculateCost)
	delete(p.VehicleTicket, vehicle.VehicleId)
	return true
}

// I can include the method of payment
