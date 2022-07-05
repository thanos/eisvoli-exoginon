package main

import (
	"fmt"
	"log"
)

type RequestType int

const (
	NEXT_TOWN RequestType = iota
)

func (request RequestType) String() string {
	switch request {
	case NEXT_TOWN:
		return "NEXT_TOWN"
	}
	return "unknown"
}

type ResponseType int

const (
	TOWN_FOUND ResponseType = iota
	TOWN_NOTFOUND
	TOWN_OCCUPIED
	NO_ROUTES
)

func (request ResponseType) String() string {
	switch request {
	case TOWN_FOUND:
		return "TOWN_FOUND"
	case TOWN_NOTFOUND:
		return "TOWN_NOTFOUND"
	case TOWN_OCCUPIED:
		return "TOWN_OCCUPIED"
	case NO_ROUTES:
		return "NO_ROUTES"
	}
	return "unknown"
}

type navRequest struct {
	alienName    string
	respChannel  chan navResponse
	deathChannel chan alienAlert
	currentTown  *Town
	requestType  RequestType
}

type navResponse struct {
	status ResponseType
	town   *Town
}
type PlanetaryNavStat struct {
	routeMap   *RouteMap
	townAlien  map[string]string
	alienTown  map[string]*Town
	navChannel chan navRequest
}

func NewPlanetaryNavStat(mapFile string) *PlanetaryNavStat {
	routeMap := BuildMap(mapFile)
	return &PlanetaryNavStat{
		routeMap:   routeMap,
		townAlien:  make(map[string]string),
		alienTown:  make(map[string]*Town),
		navChannel: make(chan navRequest),
	}

}

func (pNavStat *PlanetaryNavStat) serve() {
	log.Printf("Towns")
	log.Printf("%v", pNavStat.routeMap)
	log.Printf("pNavStat started serving")
	fightClub := make(map[string]*Town)
	for {
		select {
		case request := <-pNavStat.navChannel:
			//log.Printf("pstat recieved %v", request)
			if _, ok := fightClub[request.alienName]; ok {
				request.respChannel <- navResponse{
					status: TOWN_OCCUPIED,
					town:   request.currentTown,
				}
			} else {
				switch request.requestType {
				case NEXT_TOWN:
					resp, nextTown := pNavStat.nextTown(request.currentTown)
					if resp == TOWN_FOUND {
						if pNavStat.townOccupied(nextTown) {
							log.Printf("town occupied %v", nextTown.Name)
							pNavStat.destroyTown(request, nextTown)
							fightClub[request.alienName] = nextTown
							fightClub[pNavStat.townAlien[nextTown.Name]] = nextTown
							request.respChannel <- navResponse{
								status: TOWN_OCCUPIED,
								town:   nextTown,
							}
							//pNavStat.townAlien[nextTown.Name] <- navResponse{
							//	status: TOWN_OCCUPIED,
							//	town:   nextTown,
							//}

						} else {
							pNavStat.invadeTown(request, nextTown)
							request.respChannel <- navResponse{
								status: resp,
								town:   nextTown,
							}
						}
					} else {
						request.respChannel <- navResponse{
							status: resp,
							town:   nextTown,
						}
					}
				}
			}
		default:
			//
		}
	}
	log.Printf("navSat exits")
}

func (pNavStat PlanetaryNavStat) nextTown(currentTown *Town) (ResponseType, *Town) {
	if currentTown == nil {
		if town := pNavStat.routeMap.AnyTown(); town != nil {
			return TOWN_FOUND, town
		} else {
			return TOWN_NOTFOUND, town
		}
	} else {
		if town := currentTown.RandomRouteFrom(); town != nil {
			return TOWN_FOUND, town
		} else {
			return NO_ROUTES, town
		}
	}
}

func (pNavStat PlanetaryNavStat) destroyTown(request navRequest, town *Town) {
	fmt.Printf("twon %v destroyed by aliend %v and alien %v", town.Name, pNavStat.townAlien[town.Name], request.alienName)
	town.Destroy(nil)

}

func (pNavStat *PlanetaryNavStat) invadeTown(request navRequest, town *Town) {
	delete(pNavStat.alienTown, request.alienName)
	pNavStat.alienTown[request.alienName] = town
	if request.currentTown != nil {
		delete(pNavStat.townAlien, request.currentTown.Name)
	}
	pNavStat.townAlien[town.Name] = request.alienName
}

func (pNavStat PlanetaryNavStat) townOccupied(town *Town) bool {
	_, ok := pNavStat.townAlien[town.Name]
	return ok
}
