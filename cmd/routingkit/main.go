//main.go
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"sync"

	"github.com/nextmv-io/go-routingkit/routingkit"
)

type Router interface {
	Route(from []float32, to []float32) (uint32, [][]float32)
}

type parameters struct {
	in      *os.File
	out     *os.File
	mapFile string
	chFile  string
	measure string
	profile routingkit.TravelProfile
}

var measureEnum = struct {
	DISTANCE   string
	TRAVELTIME string
}{
	DISTANCE:   "distance",
	TRAVELTIME: "traveltime",
}

var profileEnum = struct {
	CAR        string
	BIKE       string
	PEDESTRIAN string
}{
	CAR:        "car",
	BIKE:       "bike",
	PEDESTRIAN: "pedestrian",
}

func main() {
	params := parseFlags()

	var client Router
	switch params.measure {
	case measureEnum.DISTANCE:
		c, err := routingkit.NewDistanceClient(
			params.mapFile,
			params.chFile,
			params.profile,
		)
		if err != nil {
			fmt.Println("Error creating client %", err)
		}
		client = c
	case measureEnum.TRAVELTIME:
		if params.profile != routingkit.CarTravelProfile {
			panic(`Invalid parameter combination.
			This profile can only be used with measure=distance.`)
		}
		c, err := routingkit.NewTravelTimeClient(
			params.mapFile,
			params.chFile,
		)
		if err != nil {
			fmt.Println("Error creating client %", err)
		}
		client = c
	default:
		panic("Invalid option for measure" + params.measure)
	}

	input := read(params.in)

	trips := make([]trip, len(input.Tuples))
	var wg sync.WaitGroup
	wg.Add(len(input.Tuples))
	for i, p := range input.Tuples {
		go func(i int, p pointTuple) {
			defer wg.Done()
			dist, waypoints := client.Route(p.From, p.To)
			trips[i] = trip{Cost: dist, Waypoints: waypoints}
		}(i, p)
	}
	wg.Wait()

	output := output{Trips: trips}
	write(params.out, output)
}

func parseFlags() (params parameters) {
	var in, out, profile string
	var err error
	flag.StringVar(
		&in,
		"input",
		"",
		"path to input file",
	)
	flag.StringVar(
		&out,
		"output",
		"",
		"path to output file",
	)
	flag.StringVar(
		&params.mapFile,
		"map",
		"data/map.osm.pbf",
		"path to map file",
	)
	flag.StringVar(
		&params.chFile,
		"ch",
		"data/map.ch",
		"path to ch file",
	)
	flag.StringVar(
		&profile,
		"profile",
		"car",
		"car|bike|pedestrian - bike and pedestrian only work with measure=distance",
	)
	flag.StringVar(
		&params.measure,
		"measure",
		"distance",
		"distance|traveltime",
	)
	flag.Parse()
	if in == "" {
		params.in = os.Stdin
	} else {
		params.in, err = os.Open(in)
		if err != nil {
			panic(err)
		}
	}

	switch profile {
	case profileEnum.CAR:
		params.profile = routingkit.CarTravelProfile
	case profileEnum.BIKE:
		params.profile = routingkit.BikeTravelProfile
	case profileEnum.PEDESTRIAN:
		params.profile = routingkit.PedestrianTravelProfile
	default:
		panic("Invalid option for profile: " + profile)
	}

	if out == "" {
		params.out = os.Stdout
	} else {
		params.out, err = os.Open(out)
		if err != nil {
			panic(err)
		}
	}
	if err != nil {
		panic(err)
	}
	return params
}

func read(file *os.File) input {
	dat, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}
	var in input
	err = json.Unmarshal(dat, &in)
	if err != nil {
		panic(err)
	}
	return in
}

func write(file *os.File, output output) {
	b, err := json.Marshal(output)
	if err != nil {
		panic(err)
	}
	_, err = file.Write(b)
	if err != nil {
		panic(err)
	}
}

type input struct {
	Tuples []pointTuple `json:"tuples"`
}

type pointTuple struct {
	From []float32 `json:"from"`
	To   []float32 `json:"to"`
}

type output struct {
	Trips []trip `json:"trips"`
}

type trip struct {
	Waypoints [][]float32 `json:"waypoints"`
	Cost      uint32      `json:"cost"`
}
