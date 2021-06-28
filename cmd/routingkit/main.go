package main

import (
	"fmt"
	"os"

	"github.com/nextmv-io/go-routingkit/routingkit"
)

func main() {
	osmFile := os.Getenv("OSM_FILE")
	chFile := os.Getenv("CH_FILE")
	if osmFile == "" {
		fmt.Fprintf(os.Stderr, "OSM_FILE env var is required\n")
		os.Exit(1)
	}
	if chFile == "" {
		fmt.Fprintf(os.Stderr, "CH_FILE env var is required\n")
		os.Exit(1)
	}
	_, err := os.Stat(chFile)
	if os.IsNotExist(err) {
		fmt.Println("ch file not found: creating")
	} else if err != nil {
		fmt.Fprintf(os.Stderr, "error inspecting ch file: %v\n", err)
		os.Exit(1)
	} else {
		fmt.Println("ch file found")
	}
	cli, err := routingkit.NewDistanceClient(osmFile, chFile, routingkit.CarTravelProfile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating client: %v\n", err)
		os.Exit(1)
	}
	m := cli.Matrix([][]float32{{-76.587490, 39.299710}},
		[][]float32{{-76.582855, 39.309095},
			{-76.591286, 39.298443}})
	fmt.Println(m)
}
