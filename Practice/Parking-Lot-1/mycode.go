package main

import "fmt"

type spot struct {
	free bool
	kind int
}

func (s *spot) changeStatus(status bool) bool {
	if s.getStatus() == status {
		return false
	}
	s.free = status
	return true
}

func (s *spot) getStatus() bool {
	return s.free
}

// 1->car 2->bike 3->truck
type vehicle struct {
	number string
	kind   int
}

type floor struct {
	parkingSpots  []spot
	freeCarSpot   int
	freeBikeSpot  int
	freeTruckSpot int
	totalSpot     int
}

func newFloor(car, bike, truck int) *floor {
	spots := make([]spot, car+bike+truck)
	idx := 0
	for idx = 0; idx < car; idx++ {
		spots[idx].free = true
		spots[idx].kind = 1

	}
	for idx = car; idx < car+bike; idx++ {
		spots[idx].free = true
		spots[idx].kind = 2
	}
	for idx = car + bike; idx < car+bike+truck; idx++ {
		spots[idx].free = true
		spots[idx].kind = 3
	}
	return &floor{
		parkingSpots:  spots,
		freeCarSpot:   car,
		freeBikeSpot:  bike,
		freeTruckSpot: truck,
		totalSpot:     car + truck + bike,
	}
}

func (f *floor) ParkVehicle(v vehicle) int {
	if v.kind == 1 && f.freeCarSpot > 0 {
		for i := 0; i < f.totalSpot; i++ {
			if f.parkingSpots[i].kind == 1 && f.parkingSpots[i].free {
				f.parkingSpots[i].changeStatus(false)
				f.freeCarSpot--
				return i
			}

		}
		return -1

	} else if v.kind == 2 && f.freeBikeSpot > 0 {
		for i := 0; i < f.totalSpot; i++ {
			if f.parkingSpots[i].kind == 2 && f.parkingSpots[i].free {
				f.parkingSpots[i].changeStatus(false)
				f.freeBikeSpot--
				return i
			}

		}
		return -1

	} else if v.kind == 3 && f.freeTruckSpot > 0 {
		for i := 0; i < f.totalSpot; i++ {
			if f.parkingSpots[i].kind == 3 && f.parkingSpots[i].free {
				f.parkingSpots[i].changeStatus(false)
				f.freeTruckSpot--
				return i
			}

		}
		return -1
	} else {
		return -1
	}
}

func (f *floor) UnParkVehicle(spotNum int, v vehicle) {
	f.parkingSpots[spotNum].changeStatus(true)
	if v.kind == 1 {
		f.freeCarSpot++
	} else if v.kind == 2 {
		f.freeBikeSpot++
	} else {
		f.freeTruckSpot++
	}

}

type parkingDescription struct {
	floorNum int
	spotNum  int
}

type parkingManager struct {
	allFloors     []*floor
	parkedVehicle map[string]*parkingDescription
	numFloor      int
}

func newParkingManager(numFloor int) *parkingManager {
	floors := make([]*floor, numFloor)
	for i := 0; i < numFloor; i++ {
		floors[i] = newFloor(50, 40, 10)
	}
	return &parkingManager{
		allFloors:     floors,
		parkedVehicle: make(map[string]*parkingDescription),
		numFloor:      numFloor,
	}
}

func (p *parkingManager) ParkVehicle(v vehicle) {
	for i := 0; i < p.numFloor; i++ {
		spotNum := p.allFloors[i].ParkVehicle(v)
		if spotNum != -1 {
			p.parkedVehicle[v.number] = &parkingDescription{
				spotNum:  spotNum,
				floorNum: i,
			}
		} else {
			fmt.Println("No Spot is empty")
			return
		}

	}
}

func (p *parkingManager) UnParkVehicle(v vehicle) {
	dis := p.parkedVehicle[v.number]
	p.allFloors[dis.floorNum].UnParkVehicle(dis.spotNum, v)
	delete(p.parkedVehicle, v.number)
}

func (p *parkingManager) GetFreeSpot(f, v int) int {
	floor := p.allFloors[f]
	switch v {
	case 1:
		return floor.freeCarSpot
	case 2:
		return floor.freeBikeSpot
	case 3:
		return floor.freeTruckSpot

	}
	return 0

}
