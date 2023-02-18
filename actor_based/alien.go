package main

import (
	"log"
	"sync"
)

type AlienState int

const (
	DEAD AlienState = iota
	ALIVE
	TRAPPED
	CANNOT_LAND
	EXHAUSTED
)

func (request AlienState) String() string {
	switch request {
	case DEAD:
		return "DEAD"
	case ALIVE:
		return "ALIVE"
	case TRAPPED:
		return "TRAPPED"
	case CANNOT_LAND:
		return "CANNOT_LAND"
	case EXHAUSTED:
		return "EXHAUSTED"
	}
	return "unknown"
}

type Alien struct {
	sync.RWMutex
	Name         string
	state        AlienState
	strength     int
	location     *Town
	deathChannel chan alienAlert
	navChannel   chan navResponse
}

type alienAlert struct {
	alien *Alien
}

func NewAlien(name string, startingStrength int) *Alien {
	return &Alien{
		Name:         name, //alienname.Generate(2, ""),
		state:        ALIVE,
		strength:     startingStrength,
		deathChannel: make(chan alienAlert),
		navChannel:   make(chan navResponse),
	}
}

func (alien *Alien) Rampage(navChannel chan navRequest) {
	log.Println("alien", alien.Name, "starting rampaging")
	var currentTowm *Town
	for alien.ShouldContinue() {
		if alien.state == ALIVE {
			log.Printf("alien %v asking  NEXT_TOWN", alien)
			navChannel <- navRequest{
				requestType:  NEXT_TOWN,
				deathChannel: alien.deathChannel,
				respChannel:  alien.navChannel,
				alienName:    alien.Name,
				currentTown:  currentTowm,
			}
		}
		select {
		case <-alien.deathChannel:
			log.Printf("alien %v killed", alien.Name)
			alien.Die()
		case response := <-alien.navChannel:
			log.Printf("alien %v recieved  %v ", alien, response.status)
			switch response.status {
			case TOWN_FOUND:
				log.Printf("alien %v invading %v", alien.Name, response.town.Name)
				alien.location = response.town
				currentTowm = response.town
				alien.consumeEnergy()
			case TOWN_OCCUPIED:
				alien.Die()
			case TOWN_NOTFOUND:
				fallthrough
			case NO_ROUTES:
				if alien.location != nil {
					alien.state = CANNOT_LAND
					return
				} else {
					alien.state = TRAPPED
				}
			}
		} // select
	} // for alien.ShouldContinue()
	log.Printf("alien ends %v with %v", alien.Name, alien.state)
}

func (alien *Alien) consumeEnergy() {
	alien.Lock()
	defer alien.Unlock()
	alien.strength--
	if alien.strength == 0 {
		alien.state = EXHAUSTED
	}
}
func (alien *Alien) changeLocation(nextTown *Town) {
	alien.Lock()
	defer alien.Unlock()
	if alien.location != nil {
		alien.location.SetOccupier(nil)
	}
	alien.location = nextTown
	alien.location.SetOccupier(alien)
}
func (alien *Alien) die() {
	alien.state = DEAD
	//log.Printf("%v dead", alien)
}

func (alien *Alien) Die() {
	//log.Printf("%v Die", alien)
	alien.Lock()
	defer alien.Unlock()
	alien.die()
}

func (alien *Alien) Kill(killer *Alien) {
	log.Printf("%v killed", alien.Name)
	alien.Lock()
	defer alien.Unlock()
	alien.die()
}

func (alien *Alien) Trapped() {
	alien.Lock()
	defer alien.Unlock()
	alien.state = TRAPPED
}

func (alien *Alien) ShouldContinue() bool {
	alien.RLock()
	defer alien.RUnlock()
	return alien.state == ALIVE && alien.strength > 0
}
