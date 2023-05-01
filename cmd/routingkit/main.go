// main.go
package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"sync"

	"github.com/nextmv-io/go-routingkit/routingkit"
)

type Router interface {
	Route(from []float32, to []float32) (uint32, [][]float32)
	Matrix(sources [][]float32, targets [][]float32) [][]uint32
}

type parameters struct {
	in      *os.File
	out     *os.File
	mapFile string
	measure string
	mode    string
	width   float64
	height  float64
	length  float64
	weight  float64
	speed   int
	profile routingkit.Profile
}

var measureEnum = struct {
	DISTANCE   string
	TRAVELTIME string
}{
	DISTANCE:   "distance",
	TRAVELTIME: "traveltime",
}

var modeEnum = struct {
	TUPLES string
	MATRIX string
}{
	TUPLES: "tuples",
	MATRIX: "matrix",
}

var profileEnum = struct {
	CAR        string
	BIKE       string
	PEDESTRIAN string
	TRUCK      string
}{
	CAR:        "car",
	BIKE:       "bike",
	PEDESTRIAN: "pedestrian",
	TRUCK:      "truck",
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
		c, err := routingkit.NewTravelTimeClient(
			params.mapFile,
			params.profile,
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

	switch params.mode {
	case modeEnum.TUPLES:
		input, err := readTuples(params.in)
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

		output := outputTuples{Trips: trips}
		err = writeTuples(params.out, output)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error writing output: %v", err)
			os.Exit(1)
		}
	case modeEnum.MATRIX:
		input, err := readMatrix(params.in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error reading input: %v\n", err)
			os.Exit(1)
		}

		sources := make([][]float32, len(input.Points))
		for i, p := range input.Points {
			sources[i] = []float32{p.Lon, p.Lat}
		}
		targets := make([][]float32, len(input.Points))
		for i, p := range input.Points {
			targets[i] = []float32{p.Lon, p.Lat}
		}
		distances := client.Matrix(sources, targets)

		output := outputMatrix{Matrix: distances}
		err = writeMatrix(params.out, output)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error writing output: %v", err)
			os.Exit(1)
		}
	default:
		fmt.Fprintf(os.Stderr, "invalid option for mode "+params.mode+"\n")
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
		"car|truck|bike|pedestrian",
	)
	flag.StringVar(
		&params.measure,
		"measure",
		measureEnum.DISTANCE,
		"distance|traveltime",
	)
	flag.StringVar(
		&params.mode,
		"mode",
		modeEnum.TUPLES,
		"tuples|matrix",
	)
	flag.Float64Var(
		&params.length,
		"length",
		math.MaxFloat64,
		"truck length",
	)
	flag.Float64Var(
		&params.width,
		"width",
		math.MaxFloat64,
		"truck width",
	)
	flag.Float64Var(
		&params.height,
		"height",
		math.MaxFloat64,
		"truck height",
	)
	flag.Float64Var(
		&params.weight,
		"weight",
		math.MaxFloat64,
		"truck weight",
	)
	flag.IntVar(
		&params.speed,
		"speed",
		27,
		"truck speed in m/s (default=27)",
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
		params.profile = routingkit.Car()
	case profileEnum.BIKE:
		params.profile = routingkit.Bike()
	case profileEnum.PEDESTRIAN:
		params.profile = routingkit.Pedestrian()
	case profileEnum.TRUCK:
		params.profile = routingkit.Truck(params.height, params.width, params.length, params.weight, params.speed)
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

func readTuples(file *os.File) (in inputTuples, err error) {
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

func writeTuples(file *os.File, output outputTuples) (err error) {
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

func readMatrix(file *os.File) (in inputMatrix, err error) {
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

func writeMatrix(file *os.File, output outputMatrix) (err error) {
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

type inputTuples struct {
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

type outputTuples struct {
	Trips []trip `json:"trips"`
}

type trip struct {
	Waypoints []position `json:"waypoints"`
	Cost      uint32     `json:"cost"`
}

type inputMatrix struct {
	Points []position `json:"points"`
}

type outputMatrix struct {
	Matrix [][]uint32 `json:"matrix"`
}
