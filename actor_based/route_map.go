package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"regexp"
	"strings"
)

// RouteMap is used to represent a map of towns.
type RouteMap struct {
	towns map[string]*Town
}

// RouteMap is used to represent a map of towns.
type RouteLookupTable map[string]map[string]string

func connectTown(town *Town, townsMap *RouteMap, routes RouteLookupTable) {
	for neighbourName, direction := range routes[town.Name] {
		town.routes[direction] = townsMap.towns[neighbourName]
	}
}

// BuildMap loads a map file as specified by mapFile and generates the RouteMap
func BuildMap(mapFile string) *RouteMap {

	townsMap := &RouteMap{
		towns: make(map[string]*Town),
	}
	file, err := os.Open(mapFile)

	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(file)
	rgx, _ := regexp.Compile(`(^[\w\-]+)|(\w+=\w+)`)
	routeTable := make(RouteLookupTable)
	for scanner.Scan() {
		recordLine := scanner.Text()
		if strings.HasPrefix(recordLine, "#") {
			continue
		}
		results := rgx.FindAllString(recordLine, -1)
		townName := results[0]
		routes := make(map[string]string)
		for _, route := range results[1:] {
			s := strings.Split(route, "=")
			routes[s[1]] = s[0]
		}

		routeTable[townName] = routes
		townsMap.towns[townName] = NewTown(townName)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	for _, town := range townsMap.towns {
		connectTown(town, townsMap, routeTable)
	}
	return townsMap
}

func (routeMap RouteMap) String() string {
	mapStr := ""
	for _, town := range routeMap.towns {
		mapStr += fmt.Sprintf("%v\n", town)
	}
	return mapStr
}

// AnyTown returns a random *Town that is not destroyed.
// If no Town is found it will return nil. This function is used
// to land an alien on the planet.
// Aliens that can be landed will be in Limbo.
func (routeMap RouteMap) AnyTown() *Town {
	towns := make([]*Town, 0, len(routeMap.towns))
	for _, possibleTown := range routeMap.towns {
		towns = AppendTown(towns, possibleTown)
	}
	if townNum := len(towns); townNum > 0 {
		return towns[rand.Intn(townNum)]
	}
	return nil
}

// ExistingTowns returns list of  *Towns that have not been destroyed.
func (routeMap RouteMap) ExistingTowns() []*Town {
	towns := make([]*Town, 0, len(routeMap.towns))
	for _, possibleTown := range routeMap.towns {
		towns = AppendTown(towns, possibleTown)
	}
	return towns
}

//func (routeMap RouteMap) Dump() {
//	for _, town := range routeMap.ExistingTowns() {
//		fmt.Printf("%v ", town.Name)
//		for route, neighbour := range town.routes {
//			fmt.Printf("%v=%v ", route, neighbour.Name)
//		}
//	}
//
//}
