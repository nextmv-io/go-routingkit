package routingkit_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"reflect"
	"testing"

	"github.com/nextmv-io/go-routingkit/routingkit"
)

// This is a small map file containing data for the boudning box from
// -76.60735000000001,39.28971 to -76.57749,39.31587
var marylandMap string = "testdata/maryland.osm.pbf"

// tempFile returns the location of a temporary file. It uses ioutil.TempFile
// under the hood, but if the file exists (but does not contain a valid
// contraction hierarchy), we'll get an error from routingkit, so we need to
// delete it and allow routingkit to recreate it
func tempFile(dir, pattern string) (string, error) {
	ch, err := ioutil.TempFile("", "routingkit_test.ch")
	if err != nil {
		return "", fmt.Errorf("creating tmp ch: %v", err)
	}
	filename := ch.Name()
	if err := os.Remove(filename); err != nil {
		return "", fmt.Errorf("removing temp file: %v", err)
	}
	return filename, nil
}

func TestCreateCH(t *testing.T) {
	chFile, err := tempFile("", "routingkit-test.ch")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(chFile)

	_, err = routingkit.NewDistanceClient(marylandMap, routingkit.Car(marylandMap))
	if err != nil {
		t.Fatalf("creating Client: %v", err)
	}
	if os.Stat(chFile); err != nil {
		t.Errorf("expected ch file to be created, but got error stating file: %v", err)
	}
}

func TestNearest(t *testing.T) {
	tests := []struct {
		point        []float32
		snap         float32
		expected     []float32
		expectedSnap bool
	}{
		{
			point:        []float32{-76.587490, 39.299710},
			snap:         1000,
			expected:     []float32{-76.58753, 39.29971},
			expectedSnap: true,
		},
		{
			point:        []float32{-76.584897, 39.280774},
			snap:         10,
			expected:     nil,
			expectedSnap: false,
		},
	}

	car := routingkit.Car(marylandMap)

	cli, err := routingkit.NewDistanceClient(marylandMap, car)
	if err != nil {
		t.Fatalf("creating Client: %v", err)
	}
	for i, test := range tests {
		cli.SetSnapRadius(test.snap)
		got, ok := cli.Nearest(test.point)
		if ok != test.expectedSnap {
			t.Errorf("[%d] expected snap=%t, snapped=%t", i, test.expectedSnap, ok)
		}
		if !reflect.DeepEqual(got, test.expected) {
			t.Errorf("[%d] expected %v, got %v", i, test.expected, got)
		}
	}
}

func TestDistances(t *testing.T) {
	tests := []struct {
		source       []float32
		destinations [][]float32
		snap         float32
		profile      routingkit.Profile
		ch           string

		expected []uint32
	}{
		{
			source: []float32{-76.587490, 39.299710},
			destinations: [][]float32{
				{-76.582855, 39.309095},
				{-76.591286, 39.298443},
			},
			snap:    1000,
			profile: routingkit.Car(marylandMap),

			expected: []uint32{1496, 617},
		},
		{
			source: []float32{-76.587490, 39.299710},
			destinations: [][]float32{
				{-76.582855, 39.309095},
				{-76.591286, 39.298443},
			},
			snap:    1000,
			profile: routingkit.Bike(marylandMap),

			expected: []uint32{1440, 617},
		},
		{
			source: []float32{-76.587490, 39.299710},
			destinations: [][]float32{
				{-76.582855, 39.309095},
				{-76.591286, 39.298443},
			},
			snap:    1000,
			profile: routingkit.Pedestrian(marylandMap),

			expected: []uint32{1588, 912},
		},
		{
			// should receive MaxDistance for invalid destinations
			source: []float32{-76.587490, 39.299710},
			destinations: [][]float32{
				{-76.60548, 39.30772},
				{-76.582855, 39.309095},
				{-76.584897, 39.280774},
				{-76.599388, 39.302014},
			},
			snap:    100,
			profile: routingkit.Car(marylandMap),

			expected: []uint32{routingkit.MaxDistance, 1496, routingkit.MaxDistance, 1259},
		},
		{
			// invalid source - should receive all MaxDistance
			source: []float32{-76.60586, 39.30228},
			destinations: [][]float32{
				{-76.60548, 39.30772},
				{-76.584897, 39.280774},
			},
			snap:    10,
			profile: routingkit.Car(marylandMap),

			expected: []uint32{routingkit.MaxDistance, routingkit.MaxDistance},
		},
	}

	for i, test := range tests {
		cli, err := routingkit.NewDistanceClient(marylandMap, test.profile)
		if err != nil {
			t.Fatalf("creating Client: %v", err)
		}
		cli.SetSnapRadius(test.snap)
		got := cli.Distances(test.source, test.destinations)
		if !reflect.DeepEqual(test.expected, got) {
			t.Errorf("[%d] expected %v, got %v", i, test.expected, got)
		}
	}
}

func TestMatrix(t *testing.T) {
	tests := []struct {
		sources      [][]float32
		destinations [][]float32
		profile      routingkit.Profile
		ch           string

		expected [][]uint32
	}{
		{
			sources: [][]float32{
				{-76.587490, 39.299710},
				{-76.594045, 39.300524},
				{-76.586664, 39.290938},
				{-76.598423, 39.289484},
			},
			destinations: [][]float32{
				{-76.582855, 39.309095},
				{-76.599388, 39.302014},
			},
			expected: [][]uint32{
				{1496, 1259},
				{1831, 575},
				{2372, 2224},
				{3399, 1548},
			},
			profile: routingkit.Car(marylandMap),
		},
		{
			sources: [][]float32{
				{-76.587490, 39.299710},
				{-76.594045, 39.300524},
				{-76.586664, 39.290938},
				{-76.598423, 39.289484},
			},
			destinations: [][]float32{
				{-76.582855, 39.309095},
				{-76.599388, 39.302014},
			},
			expected: [][]uint32{
				{1440, 1242},
				{1792, 558},
				{2370, 2192},
				{3354, 1547},
			},
			profile: routingkit.Bike(marylandMap),
		},
		{
			sources: [][]float32{
				{-76.587490, 39.299710},
				{-76.594045, 39.300524},
				{-76.586664, 39.290938},
				{-76.598423, 39.289484},
			},
			destinations: [][]float32{
				{-76.582855, 39.309095},
				{-76.599388, 39.302014},
			},
			expected: [][]uint32{
				{1588, 1404},
				{1589, 558},
				{2368, 2209},
				{3151, 1544},
			},
			profile: routingkit.Pedestrian(marylandMap),
		},
	}

	for i, test := range tests {
		cli, err := routingkit.NewDistanceClient(marylandMap, test.profile)
		if err != nil {
			t.Fatalf("creating Client: %v", err)
		}
		got := cli.Matrix(test.sources, test.destinations)
		if !reflect.DeepEqual(test.expected, got) {
			t.Errorf("[%d] expected %v, got %v", i, test.expected, got)
		}
	}
}

func TestDistance(t *testing.T) {
	tests := []struct {
		source      []float32
		destination []float32
		snap        float32
		profile     routingkit.Profile
		ch          string

		expectedDistance  uint32
		expectedWaypoints [][]float32
	}{
		{
			source:           []float32{-76.587490, 39.299710},
			destination:      []float32{-76.584897, 39.280774},
			snap:             1000,
			expectedDistance: 1897,
			expectedWaypoints: [][]float32{
				{-76.58753204345703, 39.29970932006836},
				{-76.58747863769531, 39.29899978637695},
				{-76.58726501464844, 39.29899978637695},
				{-76.58705139160156, 39.299007415771484},
				{-76.58668518066406, 39.29902267456055},
				{-76.58667755126953, 39.29899215698242},
				{-76.58666229248047, 39.298675537109375},
				{-76.58663940429688, 39.29836654663086},
				{-76.58662414550781, 39.29810333251953},
				{-76.58661651611328, 39.29795455932617},
				{-76.58660125732422, 39.297767639160156},
				{-76.58659362792969, 39.29757308959961},
				{-76.5865707397461, 39.29726028442383},
				{-76.5865478515625, 39.296871185302734},
				{-76.5865249633789, 39.296566009521484},
				{-76.58650970458984, 39.29627227783203},
				{-76.58650970458984, 39.296241760253906},
				{-76.58647918701172, 39.2957763671875},
				{-76.58645629882812, 39.29545593261719},
				{-76.58644104003906, 39.29514694213867},
				{-76.58643341064453, 39.29507827758789},
				{-76.58641815185547, 39.29477310180664},
				{-76.58641052246094, 39.29462814331055},
				{-76.5864028930664, 39.294586181640625},
				{-76.58638000488281, 39.294246673583984},
				{-76.58562469482422, 39.29427719116211},
				{-76.58485412597656, 39.29430389404297},
				{-76.5848388671875, 39.293941497802734},
				{-76.5848159790039, 39.29353332519531},
				{-76.58477783203125, 39.293006896972656},
				{-76.58473205566406, 39.292274475097656},
				{-76.58470916748047, 39.291893005371094},
				{-76.58467864990234, 39.291358947753906},
				{-76.58463287353516, 39.29073715209961},
				{-76.58534240722656, 39.290714263916016},
				{-76.58531951904297, 39.29029846191406},
				{-76.58529663085938, 39.290008544921875},
				{-76.58494567871094, 39.284912109375},
			},
			profile: routingkit.Car(marylandMap),
		},
		{
			source:           []float32{-76.587490, 39.299710},
			destination:      []float32{-76.584897, 39.280774},
			snap:             1000,
			expectedDistance: 1893,
			expectedWaypoints: [][]float32{
				{-76.58753, 39.29971},
				{-76.58748, 39.299},
				{-76.587265, 39.299},
				{-76.58705, 39.299007},
				{-76.586685, 39.299023},
				{-76.58668, 39.298992},
				{-76.58666, 39.298676},
				{-76.58664, 39.298367},
				{-76.586624, 39.298103},
				{-76.58662, 39.297955},
				{-76.5866, 39.297768},
				{-76.58659, 39.297573},
				{-76.58657, 39.29726},
				{-76.58655, 39.29687},
				{-76.586525, 39.296566},
				{-76.58651, 39.296272},
				{-76.58651, 39.29624},
				{-76.58648, 39.295776},
				{-76.58646, 39.295456}, {-76.58644, 39.295147}, {-76.58643, 39.29508}, {-76.58642, 39.294773}, {-76.58641, 39.29463}, {-76.5864, 39.294586}, {-76.58638, 39.294247}, {-76.585625, 39.294277}, {-76.584854, 39.294304}, {-76.58484, 39.29394}, {-76.584816, 39.293533}, {-76.58478, 39.293007}, {-76.58473, 39.292274}, {-76.58472, 39.291973}, {-76.58471, 39.291893}, {-76.58471, 39.291817}, {-76.58468, 39.29136}, {-76.58464, 39.2908}, {-76.584724, 39.290737}, {-76.58534, 39.290714},
				{-76.58532, 39.2903},
				{-76.5853, 39.29001},
				{-76.584946, 39.284912},
			},
			profile: routingkit.Bike(marylandMap),
		},
		{
			source:           []float32{-76.587490, 39.299710},
			destination:      []float32{-76.584897, 39.280774},
			snap:             1000,
			expectedDistance: 1777,
			expectedWaypoints: [][]float32{
				{-76.58753, 39.29971},
				{-76.58748, 39.299},
				{-76.587265, 39.299},
				{-76.58725, 39.298653},
				{-76.587234, 39.298347},
				{-76.58718, 39.29755},
				{-76.58716, 39.297253},
				{-76.58716, 39.297234},
				{-76.587135, 39.296844},
				{-76.58712, 39.296543},
				{-76.5871, 39.296238},
				{-76.58707, 39.295742},
				{-76.587036, 39.295124},
				{-76.58703, 39.29506},
				{-76.587, 39.294594},
				{-76.586975, 39.294224},
				{-76.586945, 39.293785},
				{-76.58693, 39.293446},
				{-76.5869, 39.29292},
				{-76.586586, 39.292934},
				{-76.58655, 39.292404},
				{-76.58652, 39.29182},
				{-76.58625, 39.291832},
				{-76.58621, 39.291294},
				{-76.586205, 39.291046},
				{-76.5862, 39.290985},
				{-76.58618, 39.290684},
				{-76.58615, 39.29031},
				{-76.58614, 39.290264},
				{-76.58612, 39.289978},
				{-76.5853, 39.29001},
				{-76.584946, 39.284912},
			},
			profile: routingkit.Pedestrian(marylandMap),
		},
		{
			source: []float32{-76.587490, 39.299710},
			// point is in a river so should not snap
			destination:       []float32{-76.584897, 39.280774},
			snap:              10,
			expectedDistance:  routingkit.MaxDistance,
			expectedWaypoints: [][]float32{},
			profile:           routingkit.Car(marylandMap),
		},
	}

	for i, test := range tests {
		cli, err := routingkit.NewDistanceClient(marylandMap, test.profile)
		if err != nil {
			t.Fatalf("creating Client: %v", err)
		}
		cli.SetSnapRadius(test.snap)
		distance, waypoints := cli.Route(test.source, test.destination)
		if test.expectedDistance != distance {
			t.Errorf("[%d] expected distance %v, got %v", i, test.expectedDistance, distance)
		}
		if !reflect.DeepEqual(test.expectedWaypoints, waypoints) {
			t.Errorf("[%d] expected waypoints %v, got %v", i, test.expectedWaypoints, waypoints)
		}
		distance = cli.Distance(test.source, test.destination)
		if test.expectedDistance != distance {
			t.Errorf("[%d] expected distance %v, got %v", i, test.expectedDistance, distance)
		}
	}
}

func TestTravelTimes(t *testing.T) {
	tests := []struct {
		source       []float32
		destinations [][]float32
		snap         float32
		ch           string

		expected []uint32
	}{
		{
			source: []float32{-76.587490, 39.299710},
			destinations: [][]float32{
				{-76.582855, 39.309095},
				{-76.591286, 39.298443},
			},
			snap: 1000,

			expected: []uint32{134910, 55530},
		},
	}

	for i, test := range tests {
		cli, err := routingkit.NewTravelTimeClient(marylandMap, routingkit.Car(marylandMap))
		if err != nil {
			t.Fatalf("creating Client: %v", err)
		}
		cli.SetSnapRadius(test.snap)
		got := cli.TravelTimes(test.source, test.destinations)
		if !reflect.DeepEqual(test.expected, got) {
			t.Errorf("[%d] expected %v, got %v", i, test.expected, got)
		}
	}
}

func TestTravelTimeMatrix(t *testing.T) {
	tests := []struct {
		sources      [][]float32
		destinations [][]float32
		ch           string

		expected [][]uint32
	}{
		{
			sources: [][]float32{
				{-76.587490, 39.299710},
				{-76.594045, 39.300524},
				{-76.586664, 39.290938},
				{-76.598423, 39.289484},
			},
			destinations: [][]float32{
				{-76.582855, 39.309095},
				{-76.599388, 39.302014},
			},
			expected: [][]uint32{
				{134910, 113043},
				{157170, 53118},
				{210945, 190230},
				{295634, 129178},
			},
		},
	}

	for i, test := range tests {
		cli, err := routingkit.NewTravelTimeClient(marylandMap, routingkit.Car(marylandMap))
		if err != nil {
			t.Fatalf("creating Client: %v", err)
		}
		got := cli.Matrix(test.sources, test.destinations)
		if !reflect.DeepEqual(test.expected, got) {
			t.Errorf("[%d] expected %v, got %v", i, test.expected, got)
		}
	}
}

func TestTravelTime(t *testing.T) {
	tests := []struct {
		source      []float32
		destination []float32
		snap        float32
		ch          string

		expectedTravelTime uint32
		expectedWaypoints  [][]float32
	}{
		{
			source:             []float32{-76.587490, 39.299710},
			destination:        []float32{-76.584897, 39.280774},
			snap:               1000,
			expectedTravelTime: 213405,
			expectedWaypoints: [][]float32{
				{-76.58753, 39.29971},
				{-76.587906, 39.299694},
				{-76.58786, 39.298977},
				{-76.587845, 39.29863},
				{-76.58725, 39.298653},
				{-76.58666, 39.298676},
				{-76.58664, 39.298367},
				{-76.586624, 39.298103},
				{-76.58662, 39.297955},
				{-76.5866, 39.297768},
				{-76.58659, 39.297573},
				{-76.58657, 39.29726},
				{-76.58655, 39.29687},
				{-76.586525, 39.296566},
				{-76.58651, 39.296272},
				{-76.58651, 39.29624},
				{-76.58648, 39.295776},
				{-76.58646, 39.295456},
				{-76.58644, 39.295147},
				{-76.58643, 39.29508},
				{-76.58642, 39.294773},
				{-76.58641, 39.29463},
				{-76.5864, 39.294586},
				{-76.58638, 39.294247},
				{-76.585625, 39.294277},
				{-76.584854, 39.294304},
				{-76.58484, 39.29394},
				{-76.584816, 39.293533},
				{-76.58478, 39.293007},
				{-76.58473, 39.292274},
				{-76.58471, 39.291893},
				{-76.58468, 39.29136},
				{-76.58463, 39.290737},
				{-76.58534, 39.290714},
				{-76.58532, 39.2903},
				{-76.5853, 39.29001},
				{-76.584946, 39.284912},
			},
		},
	}

	for i, test := range tests {
		cli, err := routingkit.NewTravelTimeClient(marylandMap, routingkit.Car(marylandMap))
		if err != nil {
			t.Fatalf("creating Client: %v", err)
		}
		cli.SetSnapRadius(test.snap)
		travelTime, waypoints := cli.Route(test.source, test.destination)
		if test.expectedTravelTime != travelTime {
			t.Errorf("[%d] expected travel time %v, got %v", i, test.expectedTravelTime, travelTime)
		}
		if !reflect.DeepEqual(test.expectedWaypoints, waypoints) {
			t.Errorf("[%d] expected waypoints %v, got %v", i, test.expectedWaypoints, waypoints)
		}
		travelTime = cli.TravelTime(test.source, test.destination)
		if test.expectedTravelTime != travelTime {
			t.Errorf("[%d] expected travel time %v, got %v", i, test.expectedTravelTime, travelTime)
		}
	}
}

var distance uint32
var distances [][]uint32

// These two functions are utilities for generating random points within a bounding box for benchmarking
// Keeping them around even though they aren't used now. Note that even though the points will lie within
// the bounding it box, it may not be possible to route between them if the route requires leaving the box

func pointInRange(low float64, high float64) float64 {
	var mult float64 = 100000
	lowInt := int(low * mult)
	highInt := int(high * mult)
	return float64(rand.Intn(highInt-lowInt)+lowInt) / mult
}

func randomPointsInBoundingBox(n int, bottomLeft [2]float64, topRight [2]float64) [][]float64 {
	points := make([][]float64, n)
	for i := 0; i < n; i++ {
		points[i] = []float64{pointInRange(bottomLeft[0], topRight[0]), pointInRange(bottomLeft[1], topRight[1])}
	}
	return points
}

func BenchmarkDistance(b *testing.B) {
	cli, err := routingkit.NewDistanceClient(marylandMap, routingkit.Car(marylandMap))
	if err != nil {
		b.Fatalf("creating Client: %v", err)
	}
	cli.SetSnapRadius(1000)

	f, err := os.Open("testdata/points.json")
	if err != nil {
		b.Fatal(err)
	}
	var data struct {
		Points [][]float32
	}
	if err := json.NewDecoder(f).Decode(&data); err != nil {
		b.Fatal(err)
	}

	b.Run("distance", func(b *testing.B) {
		var d uint32
		for n := 0; n < b.N; n++ {
			d = cli.Distance(data.Points[n%len(data.Points)], data.Points[(n+1)%len(data.Points)])
		}
		distance = d
	})
}

func BenchmarkMatrix(b *testing.B) {
	cli, err := routingkit.NewDistanceClient(marylandMap, routingkit.Car(marylandMap))
	if err != nil {
		b.Fatalf("creating Client: %v", err)
	}
	cli.SetSnapRadius(1000)

	f, err := os.Open("testdata/points.json")
	if err != nil {
		b.Fatal(err)
	}
	var data struct {
		Points [][]float32
	}
	if err := json.NewDecoder(f).Decode(&data); err != nil {
		b.Fatal(err)
	}

	b.Run("matrix", func(b *testing.B) {
		var d [][]uint32
		for n := 0; n < b.N; n++ {
			d = cli.Matrix(data.Points, data.Points)
		}
		distances = d
	})
}
