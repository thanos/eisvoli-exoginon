/*
simulate, simulates an alien invasion of a planet with towns.

Usage:

	go run main.go [flags] N map_file

where required arguments are:

	N is the number of aliens to land on the por planet
	map_file the path to a map file in the described format below

where the optional flags are:

	-energy=int
	    The energy level the aliens starts with
	-debug=bool
	   	set to ture to turn on logging.

examples:

			go run main.go  10 maps/map_01.txt
			go run main.go -energy= 10 3 maps/map_01.txt
	        go run main.go -debug=true -energy= 10 3 maps/map_01.txt
			Map file format is like this:

			#Foo north=Bar west=Baz south=Qu-ux
			#Bar south=Foo west=Bee
			#Qu-ux north=Foo west=Bee
			#Bee west=Bar east=Qu-ux
			#Baz west=Foo south=Daz
			Daz north=Kaz
			Naz west=Zaz south=Daz
			Zaz north=Kaz east=Naz
			Kaz south=Zaz

lines starting with # are ignored
*/
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"sync"
	"time"
)

const DEFAULT_ENERGY = 10_000

var alienGroup sync.WaitGroup

func land(navChannel chan navRequest, alien *Alien) {
	defer alienGroup.Done()
	alien.Rampage(navChannel)
}

type argsGiven struct {
	numberAliens   int    // number of aliens to land on the planet
	startingEnergy int    // amount of energy they start with
	mapFile        string // the path to the map file that will be loaded on startup
	logOn          bool   // used to turn on or off logging
}

func onArgErr() {
	flagSet := flag.CommandLine
	fmt.Println(`usage:
		go run main.go -energy=E  N map_file

		where required arguments are:
		N is the number of aliens to land on the por planet
		map_file the path to a map file in the described format below

		where the optional arguments are:`)
	order := []string{"energy"}
	for _, name := range order {
		option_flag := flagSet.Lookup(name)
		fmt.Printf("-%s\t%s\n", option_flag.Name, option_flag.Usage)
	}
	fmt.Println(`examples:

		go run main.go  10 maps/map_01.txt
		go run main.go -energy= 10 3 maps/map_01.txt
        go run main.go -debug=true -energy= 10 3 maps/map_01.txt
		Map file format is like this:

		#Foo north=Bar west=Baz south=Qu-ux
		#Bar south=Foo west=Bee
		#Qu-ux north=Foo west=Bee
		#Bee west=Bar east=Qu-ux
		#Baz west=Foo south=Daz
		Daz north=Kaz
		Naz west=Zaz south=Daz
		Zaz north=Kaz east=Naz
		Kaz south=Zaz

		lines starting with # are ignored`)
	os.Exit(1)
}

func processCmdArgs() *argsGiven {
	args := argsGiven{
		numberAliens:   1,
		startingEnergy: DEFAULT_ENERGY,
		mapFile:        "maps/map_01.txt",
		logOn:          false,
	}
	flag.IntVar(&args.startingEnergy, "energy", args.startingEnergy, "energy the aliens starts with")
	flag.BoolVar(&args.logOn, "debug", args.logOn, "set to ture to turn on logging")
	flag.Parse()
	var err error
	if flag.NArg() == 2 {
		if args.numberAliens, err = strconv.Atoi(flag.Arg(0)); err == nil {
			args.mapFile = flag.Arg(1)
			if _, err = os.Stat(args.mapFile); err == nil {
				return &args
			}
		}

	}
	if err != nil {
		fmt.Errorf("%v", err)
	}
	onArgErr()
	return nil
}

func main() {
	args := processCmdArgs()
	if args.logOn == false {
		log.SetOutput(ioutil.Discard)
	}
	log.Printf("starting")
	navStat := NewPlanetaryNavStat(args.mapFile)

	alienGroup.Add(args.numberAliens)
	for alienNumber := 0; alienNumber < args.numberAliens; alienNumber++ {
		alien := NewAlien(fmt.Sprintf("alien_%03d", alienNumber+1), args.startingEnergy)
		go land(navStat.navChannel, alien)
		time.Sleep(1 * time.Millisecond)
	}
	time.Sleep(time.Second)
	go navStat.serve()
	alienGroup.Wait()

	fmt.Println("\nRemaining Towns ")
	for _, town := range navStat.routeMap.ExistingTowns() {
		fmt.Println(town)
	}

}
