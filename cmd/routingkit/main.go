//main.go
package main

import (
	"encoding/json"
	"errors"
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
	params, err := parseFlags()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error parsing flags: %v\n", err)
		os.Exit(1)
	}

	var client Router
	switch params.measure {
	case measureEnum.DISTANCE:
		c, err := routingkit.NewDistanceClient(
			params.mapFile,
			params.profile,
		)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error creating client: %v\n", err)
			os.Exit(1)
		}
		client = c
	case measureEnum.TRAVELTIME:
		if params.profile != routingkit.CarTravelProfile {
			fmt.Fprintf(os.Stderr, "invalid parameter combination. "+
				"This profile can only be used with measure=distance.\n")
			os.Exit(1)
		}
		c, err := routingkit.NewTravelTimeClient(
			params.mapFile,
		)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error creating client: %v\n", err)
			os.Exit(1)
		}
		client = c
	default:
		fmt.Fprintf(os.Stderr, "invalid option for measure "+params.measure+"\n")
		os.Exit(1)
	}

	input, err := read(params.in)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading input: %v\n", err)
		os.Exit(1)
	}

	trips := make([]trip, len(input.Tuples))
	var wg sync.WaitGroup
	wg.Add(len(input.Tuples))
	for i, p := range input.Tuples {
		go func(i int, p pointTuple) {
			defer wg.Done()
			from := []float32{p.From.Lon, p.From.Lat}
			to := []float32{p.To.Lon, p.To.Lat}
			dist, wpTuples := client.Route(from, to)
			waypoints := make([]position, len(wpTuples))
			for w, tuple := range wpTuples {
				waypoints[w] = position{Lon: tuple[0], Lat: tuple[1]}
			}
			trips[i] = trip{Cost: dist, Waypoints: waypoints}
		}(i, p)
	}
	wg.Wait()

	output := output{Trips: trips}
	err = write(params.out, output)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error writing output: %v", err)
		os.Exit(1)
	}
}

func parseFlags() (params parameters, err error) {
	var in, out, profile string
	flag.StringVar(
		&in,
		"input",
		"",
		"path to input file. default is stdin.",
	)
	flag.StringVar(
		&out,
		"output",
		"",
		"path to output file. default is stdout.",
	)
	flag.StringVar(
		&params.mapFile,
		"map",
		"data/map.osm.pbf",
		"path to map file",
	)
	flag.StringVar(
		&profile,
		"profile",
		profileEnum.CAR,
		"car|bike|pedestrian - bike and pedestrian only work with measure=distance",
	)
	flag.StringVar(
		&params.measure,
		"measure",
		measureEnum.DISTANCE,
		"distance|traveltime",
	)
	flag.Parse()
	if in == "" {
		params.in = os.Stdin
	} else {
		params.in, err = os.Open(in)
		if err != nil {
			return parameters{}, err
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
		return parameters{}, errors.New("invalid option for profile" + profile)
	}

	if out == "" {
		params.out = os.Stdout
	} else {
		params.out, err = os.Create(out)
		if err != nil {
			return parameters{}, err
		}
	}
	if err != nil {
		return parameters{}, err
	}
	return params, nil
}

func read(file *os.File) (in input, err error) {
	dat, err := ioutil.ReadAll(file)
	if err != nil {
		return in, err
	}
	err = json.Unmarshal(dat, &in)
	if err != nil {
		return in, err
	}
	return in, nil
}

func write(file *os.File, output output) (err error) {
	b, err := json.Marshal(output)
	if err != nil {
		return err
	}
	_, err = file.Write(b)
	if err != nil {
		return err
	}
	return nil
}

type input struct {
	Tuples []pointTuple `json:"tuples"`
}

type pointTuple struct {
	From position `json:"from"`
	To   position `json:"to"`
}

type position struct {
	Lon float32 `json:"lon"`
	Lat float32 `json:"lat"`
}

type output struct {
	Trips []trip `json:"trips"`
}

type trip struct {
	Waypoints []position `json:"waypoints"`
	Cost      uint32     `json:"cost"`
}
