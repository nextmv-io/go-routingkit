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

// This is a small map file containing data for the boudning box from
// -76.60735000000001,39.28971 to -76.57749,39.31587
var marylandMap string = "testdata/maryland.osm.pbf"

func TestDistances(t *testing.T) {
	tests := []struct {
		source       []float64
		destinations [][]float64
		expected     []float64
		snap         float64
	}{
		{
			source: []float64{-76.587490, 39.299710},
			destinations: [][]float64{
				{-76.582855, 39.309095},
				{-76.599388, 39.302014},
			},
			snap:     1000,
			expected: []float64{1496, 1259},
		},
		{
			// should receive -1 for invalid destinations
			source: []float64{-76.587490, 39.299710},
			destinations: [][]float64{
				{-76.60548, 39.30772},
				{-76.582855, 39.309095},
				{-76.584897, 39.280774},
				{-76.599388, 39.302014},
			},
			snap:     100,
			expected: []float64{-1, 1496, -1, 1259},
		},
		{
			// invalid source - should receive all -1
			source: []float64{-76.60586, 39.30228},
			destinations: [][]float64{
				{-76.60548, 39.30772},
				{-76.584897, 39.280774},
			},
			snap:     10,
			expected: []float64{-1, -1},
		},
	}
	chFile, err := tempFile("", "routingkit-test.ch")
	defer os.Remove(chFile)
	cli, err := routingkit.NewClient(marylandMap, chFile)
	if err != nil {
		t.Fatalf("creating Client: %v", err)
	}

	for i, test := range tests {
		cli.SetSnapRadius(test.snap)
		got := cli.Distances(test.source, test.destinations)
		if !reflect.DeepEqual(test.expected, got) {
			t.Errorf("[%d] expected %v, got %v", i, test.expected, got)
		}
	}
}

func TestMatrix(t *testing.T) {
	tests := []struct {
		sources      [][]float64
		destinations [][]float64
		expected     [][]float64
	}{
		{
			sources: [][]float64{
				{-76.587490, 39.299710},
				{-76.594045, 39.300524},
				{-76.586664, 39.290938},
				{-76.598423, 39.289484},
			},
			destinations: [][]float64{
				{-76.582855, 39.309095},
				{-76.599388, 39.302014},
			},
			expected: [][]float64{
				{1496, 1259},
				{1831, 575},
				{2372, 2224},
				{3399, 1548},
			},
		},
	}
	chFile, err := tempFile("", "routingkit-test.ch")
	defer os.Remove(chFile)
	cli, err := routingkit.NewClient(marylandMap, chFile)
	if err != nil {
		t.Fatalf("creating Client: %v", err)
	}

	for i, test := range tests {
		got := cli.Matrix(test.sources, test.destinations)
		if !reflect.DeepEqual(test.expected, got) {
			t.Errorf("[%d] expected %v, got %v", i, test.expected, got)
		}
	}

}

func TestDistance(t *testing.T) {
	tests := []struct {
		source            []float64
		destination       []float64
		snap              float64
		expectedDistance  float64
		expectedWaypoints [][]float64
	}{
		{
			source:           []float64{-76.587490, 39.299710},
			destination:      []float64{-76.584897, 39.280774},
			snap:             1000,
			expectedDistance: 1897,
			expectedWaypoints: [][]float64{
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
		},
		{
			source: []float64{-76.587490, 39.299710},
			// point is in a river so should not snap
			destination:       []float64{-76.584897, 39.280774},
			snap:              10,
			expectedDistance:  -1,
			expectedWaypoints: [][]float64{},
		},
	}
	chFile, err := tempFile("", "routingkit-test.ch")
	defer os.Remove(chFile)
	cli, err := routingkit.NewClient(marylandMap, chFile)
	if err != nil {
		t.Fatalf("creating Client: %v", err)
	}

	for i, test := range tests {
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

var distances [][]float64

func benchmarkMatrix(pointsFile string, nSources int, nDestinations int, b *testing.B) {
	chFile, err := tempFile("", "routingkit-test.ch")
	defer os.Remove(chFile)
	cli, err := routingkit.NewClient(marylandMap, chFile)
	if err != nil {
		b.Fatalf("creating Client: %v", err)
	}
	cli.SetSnapRadius(100000)

	f, err := os.Open(pointsFile)
	if err != nil {
		b.Fatal(err)
	}
	var data struct {
		Points [][]float64
	}
	if err := json.NewDecoder(f).Decode(&data); err != nil {
		b.Fatal(err)
	}
	sources := data.Points[:nSources]
	destinations := data.Points[:nDestinations]

	var d [][]float64
	for n := 0; n < b.N; n++ {
		d = cli.Matrix(sources, destinations)
	}
	distances = d
}

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

//func TestGenerateRandomPoints(t *testing.T) {
//	points := randomPointsInBoundingBox(1000, [2]float64{-76.60735000000001, 39.28971}, [2]float64{-76.57749, 39.31587})
//	f, err := os.OpenFile("testdata/points.json", os.O_CREATE|os.O_RDWR, 0755)
//	if err != nil {
//		t.Fatal(err)
//	}
//	p := map[string]interface{}{
//		"points": points,
//	}
//	if err := json.NewEncoder(f).Encode(p); err != nil {
//		t.Fatal(err)
//	}
//}

func BenchmarkMatrix(b *testing.B) {
	benchmarkMatrix("testdata/points.json", 100, 100, b)
}
