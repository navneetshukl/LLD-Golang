package main

// =====================
// Vehicle
// =====================

type VehicleType int

const (
	Car VehicleType = iota + 1
	Bike
	Truck
)

type Vehicle struct {
	Number string
	Type   VehicleType
}

// =====================
// Parking Spot
// =====================

type ParkingSpot struct {
	ID      int
	Type    VehicleType
	Vehicle *Vehicle
}

func (s *ParkingSpot) IsFree() bool {
	return s.Vehicle == nil
}

func (s *ParkingSpot) Park(v *Vehicle) bool {
	if !s.IsFree() || s.Type != v.Type {
		return false
	}
	s.Vehicle = v
	return true
}

func (s *ParkingSpot) UnPark() {
	s.Vehicle = nil
}

// =====================
// Floor
// =====================

type Floor struct {
	Number    int
	Spots     []*ParkingSpot
	freeCount map[VehicleType]int
}

func NewFloor(number, car, bike, truck int) *Floor {
	spots := []*ParkingSpot{}
	id := 0

	for i := 0; i < car; i++ {
		spots = append(spots, &ParkingSpot{ID: id, Type: Car})
		id++
	}
	for i := 0; i < bike; i++ {
		spots = append(spots, &ParkingSpot{ID: id, Type: Bike})
		id++
	}
	for i := 0; i < truck; i++ {
		spots = append(spots, &ParkingSpot{ID: id, Type: Truck})
		id++
	}

	return &Floor{
		Number: number,
		Spots:  spots,
		freeCount: map[VehicleType]int{
			Car:   car,
			Bike:  bike,
			Truck: truck,
		},
	}
}

func (f *Floor) ParkVehicle(v *Vehicle) (*ParkingSpot, bool) {
	if f.freeCount[v.Type] == 0 {
		return nil, false
	}

	for _, spot := range f.Spots {
		if spot.Park(v) {
			f.freeCount[v.Type]--
			return spot, true
		}
	}
	return nil, false
}

func (f *Floor) UnParkVehicle(spotID int) {
	spot := f.Spots[spotID]
	if spot.Vehicle != nil {
		f.freeCount[spot.Type]++
		spot.UnPark()
	}
}

func (f *Floor) GetFreeCount(t VehicleType) int {
	return f.freeCount[t]
}

// =====================
// Parking Lot (Manager)
// =====================

type ParkingLocation struct {
	Floor int
	Spot  int
}

type ParkingLot struct {
	Floors []*Floor
	parked map[string]*ParkingLocation
}

func NewParkingLot(numFloors int) *ParkingLot {
	floors := make([]*Floor, numFloors)
	for i := 0; i < numFloors; i++ {
		floors[i] = NewFloor(i, 50, 40, 10)
	}

	return &ParkingLot{
		Floors: floors,
		parked: make(map[string]*ParkingLocation),
	}
}

func (p *ParkingLot) ParkVehicle(v *Vehicle) bool {
	for _, floor := range p.Floors {
		spot, ok := floor.ParkVehicle(v)
		if ok {
			p.parked[v.Number] = &ParkingLocation{
				Floor: floor.Number,
				Spot:  spot.ID,
			}
			return true
		}
	}
	return false
}

func (p *ParkingLot) UnParkVehicle(vehicleNumber string) bool {
	location, ok := p.parked[vehicleNumber]
	if !ok {
		return false
	}

	floor := p.Floors[location.Floor]
	floor.UnParkVehicle(location.Spot)
	delete(p.parked, vehicleNumber)
	return true
}

func (p *ParkingLot) GetFreeSpots(floor int, t VehicleType) int {
	return p.Floors[floor].GetFreeCount(t)
}

// =====================
// main intentionally omitted
// =====================
