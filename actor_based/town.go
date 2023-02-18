package main

import (
	"fmt"
	"math/rand"
	"strings"
	"sync"
)

// Town represents a town and the routes to its neighbouring towns.
type Town struct {
	sync.RWMutex
	Name        string // Town's name
	isDestroyed bool
	occupier    *Alien
	routes      map[string]*Town
}

func (town *Town) String() string {
	line := make([]string, 0, len(town.routes))
	line = append(line, town.Name)
	if town.isDestroyed {
		return town.Name + " Destroyed"
	}
	for direction, t := range town.AdjoiningTowns() {
		line = append(line, fmt.Sprintf("%v=%v", direction, t.Name))
	}
	return strings.Join(line, " ")
}

// NewTown creates a new Town given the town's name
func NewTown(name string) *Town {
	return &Town{
		Name:        name,
		routes:      make(map[string]*Town),
		isDestroyed: false,
	}
}

// IsDestroyed  tests if the town is destroyed.
// Puts a read lock on the town struct.
func (town *Town) IsDestroyed() bool {
	town.RLock()
	defer town.RUnlock()
	return town.isDestroyed
}

// Destroy changes the towns state isDestroyed to true and takes occupier for printing to STDOUT.
// Puts a read lock on the town struct.
func (town *Town) Destroy(occupier *Alien) {

	town.Lock()
	defer town.Unlock()
	fmt.Printf("\n%v destroyed by alien %v and alien %v\n", town.Name, town.occupier, occupier)
	town.isDestroyed = true
}

// AdjoiningTowns returns a list of neighbouring towns.
func (town *Town) AdjoiningTowns() []*Town {
	towns := make([]*Town, 0, len(town.routes))
	for _, t := range town.routes {
		towns = AppendTown(towns, t)
	}
	return towns
}

// AppendTown returns a list of towns not destroyed.
// Puts a read lock on each town tested and appended.
func AppendTown(towns []*Town, town *Town) []*Town {
	if town != nil {
		town.RLock()
		if town != nil && town.isDestroyed == false {
			towns = append(towns, town)
		}
		town.RUnlock()
	} else {
		town = nil
	}
	return towns
}

// RandomRouteFrom returns a random neighbouring town that is not destroyed.
func (town *Town) RandomRouteFrom() *Town {
	towns := town.AdjoiningTowns()
	if routeNumber := len(towns); routeNumber > 0 {
		possibleTown := towns[rand.Intn(routeNumber)]
		return possibleTown
	}
	return nil
}

func (town *Town) notifyNeighbours() {
	for _, neighbouringTown := range town.routes {
		neighbouringTown.RemoveRoute(town.Name)
	}
}

// IsAlreadyOccupied checks to see if the town is occupied
// Puts a read lock on each town tested.
func (town *Town) IsAlreadyOccupied() bool {
	town.RLock()
	defer town.RUnlock()
	return town.occupier != nil
}

// SetOccupier sets the occupier
// Puts a  lock on each town tested.
func (town *Town) SetOccupier(alien *Alien) {
	town.Lock()
	defer town.Unlock()
	town.occupier = alien
}

// GetOccupier gets the current  occupier, returns nill if there are none.
// Puts a read lock on each town tested.
func (town *Town) GetOccupier() *Alien {
	town.RLock()
	defer town.RUnlock()
	return town.occupier
}

// RemoveRoute deletes a route from Town route table.
// Puts a read lock on each town tested.
func (town *Town) RemoveRoute(name string) {
	town.Lock()
	defer town.Unlock()
	delete(town.routes, name)
}

func compass(direction string) string {
	switch direction {
	case "north":
		return "south"
	case "east":
		return "west"
	case "south":
		return "north"
	case "west":
		return "east"
	default:
		return "COMPASS ERROR"
	}
}
