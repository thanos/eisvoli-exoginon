# Aisvoli Exoginon
A simulation example in golang


## Problem Outline

1. Read in  map file of towns
1. Each town has upto four routes (north, south, east, or west) to other towns
1. N aliens land at random towns on the map and wander around randomly.
1. if two *(or more ?)* should meet they destroy the town along with themselves.
1. On a town's destruction the simulation should print "t has been destroyed by alien x and alien y!"
1. When a town is destroyed its removed from the map along with the routes to it.
1. When an alien finds itself in a town with all it's neighbours destroyed it becomes stuck.
1. Finnaly print out a map file of the remaining towns.
1. language to use: Golang

Some simple questions:
1. Can an alien start it's invasion in a isolated town ? This is a possible edgae case, unless we assume the `landing` town is randomly chosen from the connected twon set.
1. When aliens meet in a town do the kill each other and in the process destroy the town or when the town is desrotyed they die with it. The requirements states that they *kill each other and then destroy the town.* The interpretation will effect how this is coded. 


## Basic Needs

1. a map file reader that builds an internal map representation.
```
Foo north=Bar west=Baz south=Qu-ux
Bar south=Foo west=Bee
.
.
```

to 

```
// non language specific maps of structs ?
{
  "Foo": {
    "routes": { "north": &Bar, "west": &Baz, "south": &Quux}},
    "occupiers": set{"Alien1"},
    }
    
   "Bar": {
    "routes":{ "west": &Bee, "south": &Foo},
      "occupiers": set{},
    }
 .
 .
 }
```


### Actor Solution

As possibly hinted in the requirements with *"Each iteration..."*,  one solution  would be to cose the simulation without goroutes as an O(N**2) solution 
The main loop being the 10k iterations:


```
map = buildMap(map_file)
dead_or_trapped = {}

for iteration in 10k:
  alien_active_set = aliens(N) - dead_or_trapped
  for alien in alien_living_set:
    {killed, trapped, rampaging}, town = alien.pillage(map)
    if killed:
          map.remove(town)
          dead_or_trapped.add(alien)
     if trapped:
           dead_or_trapped.add(alien)
```

#### PROS
  1. no race conditions
  1. no deadlocks
  1. deterministic so easy to debug (use constant randon seeds)
  1. easy to test
  1. performance predictable and linear
  1. predictable dev cost and fast to market
  1. any developer from any language.
   
#### CONS
 1. can't scale to utilizes available CPU cores
 1. simulation is forced - __deterministic__
 1. boring
 1. doesnt showoff the programmer's skill
 

### Actor/Sprite Solution

Define one actor Alien and Town

an Alien's responibility is to 
  1. get a route from the Map
  1. invade the new town
  1. if the town is already invaded kill the other alien,  destroy the town and die
  1. if the town is vacant goto the first step 
  
#### PROS
  1. elegant
  1. simulation is non-deterministic
  1. will scale use available CPU cores

   
#### CONS
  1. non-deterministc so it's hard to test and debug
  1. performance unpredictable and linear
  1. unpredictable dev cost
  1. race and dead lock prone.
  1. needs an expienced golang developer
  
  
  #### A FEW POSSIBLE RACE CONDITIONS
  
 1. an alien needs to get an available route and then occupy the town - but the at the same time could have beed destroyed. for instance the code look like:
```golang 
if alien.location == nil {
	next_town = route_map.AnyTown() 
} else {
	next_town = alien.location.RandomRouteFrom()
}
if next_town != nill {
	if next_town.IsAlreadyOccupied() {
		occupier = next_town.GetOccupier()
		next_town.Destroy()
		occupier.Kill()
	} else {
		next_town.SetOccupier(alien)
		alien.location = next_town
	}
} else {
	alien.Trapped()
}
```
 *  `location` is `*Town`

 1. there are race conditions on the next_town's state
 1. there are race conditions on the next_town's routes - it's neighbours could have been detroyed between `AnyTown()`, `RandomRouteFrom()` and  `IsAlreadyOccupied()`
 1. between the test `IsAlreadyOccupied() == false`  and `next_town.SetOccupier(alien)` it could have been already occupied.
 1. And there are more in `if next_town.IsAlreadyOccupied() {}` condition.
 

  So the next question is `channels` or `Mutex` and where ?
   


