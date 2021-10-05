package routingkit_test

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/nextmv-io/go-routingkit/routingkit"
)

// This is a small map file containing data for the boudning box from
// -76.60735000000001,39.28971 to -76.57749,39.31587
var marylandMap string = "testdata/maryland.osm.pbf"

// This is another small map of Maryland that contains data for the
// bounding box from -76.663640,39.240043 to -76.605023,39.269623.
// A low overpass with a max height of 13'11" is found at -76.638449,39.254932
// (way ID 456490563), where Annapolis Rd. passes under a train line
var marylandMapWithHeightRestriction = "testdata/maryland_height_restriction.osm.pbf"

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

// plotWaypoints uses nextplot to render two paths, the expected and received
// set of waypoints. It returns the paths of the two plots
func plotWaypoints(
	i int,
	expectedWaypoints [][]float32,
	gotWaypoints [][]float32,
) (expectedPlot string, gotPlot string, err error) {
	tempDir := os.TempDir()
	expectedPath := filepath.Join(tempDir, fmt.Sprintf("routingkit_debug_expected_%d.json", i))
	gotPath := filepath.Join(tempDir, fmt.Sprintf("routingkit_debug_got_%d.json", i))
	expectedFile, err := os.OpenFile(expectedPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0655)
	if err != nil {
		return "", "", fmt.Errorf("opening expected waypoints file: %v", err)
	}
	defer expectedFile.Close()
	gotFile, err := os.OpenFile(gotPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0655)
	if err != nil {
		return "", "", fmt.Errorf("opening got waypoints file: %v", err)
	}
	defer gotFile.Close()

	type waypoints struct {
		Waypoints [][]float32 `json:"waypoints"`
	}
	if err := json.NewEncoder(expectedFile).Encode(waypoints{expectedWaypoints}); err != nil {
		return "", "", fmt.Errorf("writing expected points: %v", err)
	}
	if err := json.NewEncoder(gotFile).Encode(waypoints{gotWaypoints}); err != nil {
		return "", "", fmt.Errorf("writing expected points: %v", err)
	}
	for _, path := range []string{expectedPath, gotPath} {
		out, err := exec.Command(
			"nextplot",
			"route",
			"--input_route",
			path,
			"--jpath_route",
			"waypoints",
		).CombinedOutput()
		if err != nil {
			return "", "", fmt.Errorf("nextplot error: %v, stdout: %s", err, string(out))
		}
	}

	return expectedPath + ".html", gotPath + ".html", nil
}

func TestCreateCH(t *testing.T) {
	chFile, err := tempFile("", "routingkit-test.ch")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(chFile)

	_, err = routingkit.NewDistanceClient(marylandMap, routingkit.Car())
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

	car := routingkit.Car()

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
			profile: routingkit.Car(),

			expected: []uint32{1496, 617},
		},
		{
			source: []float32{-76.587490, 39.299710},
			destinations: [][]float32{
				{-76.582855, 39.309095},
				{-76.591286, 39.298443},
			},
			snap:    1000,
			profile: routingkit.Bike(),

			expected: []uint32{1496, 617},
		},
		{
			source: []float32{-76.587490, 39.299710},
			destinations: [][]float32{
				{-76.582855, 39.309095},
				{-76.591286, 39.298443},
			},
			snap:    1000,
			profile: routingkit.Pedestrian(),

			expected: []uint32{1429, 428},
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
			profile: routingkit.Car(),

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
			profile: routingkit.Car(),

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
			profile: routingkit.Car(),
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
				{1496, 1259},
				{1831, 575},
				{2372, 2224},
				{3399, 1548},
			},
			profile: routingkit.Bike(),
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
				{1429, 1259},
				{1589, 575},
				{2367, 2221},
				{3157, 1535},
			},
			profile: routingkit.Pedestrian(),
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

var update *bool

func init() {
	update = flag.Bool("update", false, "update text fixtures")
}

// TestDistance not only tests that the distance between two points is the expected value, but also that
// the provided waypoints match test fixtures located in testdata/fixtures. These fixtures can automatically
// be updated when they don't match if you pass the -update flag when running the tests. If you're using
// "go test", you'll have to pass this as e.g. "go test ./... -args -update."
// If the test fails, nextplot will be used to create a plot of the expected and received cases in the temporary
// directory. Before updating a case or adding a new case, you should look at these plots and confirm them
// to other sources (OSRM, Google Maps) to ensure they look reasonable - if anything looks off please note somewhere.
func TestDistance(t *testing.T) {
	tests := []struct {
		source      []float32
		destination []float32
		snap        float32
		profile     routingkit.Profile
		ch          string
		osmFile     string

		expectedDistance uint32
		waypointsFile    string
	}{
		// The destination has a strange way of snapping, it snaps to -76.58494567871094, 39.284912109375
		// which is a few blocks inland even though the point is in the water. May want to look into
		// this more
		// Also, I noticed the route we get takes one strange turn, on the block between E. Monument and
		// E. Madison. This street is tagged as highway=service and service=alley but on Google Street
		// View it really doesn't seem like something you should be driving through (narrow alleyway).
		// I wonder if we should think about forbidding streets tagged like this unless they're being
		// used to reach a destination
		{
			source:           []float32{-76.587490, 39.299710},
			destination:      []float32{-76.584897, 39.280774},
			snap:             1000,
			expectedDistance: 1897,
			osmFile:          marylandMap,
			waypointsFile:    "waypoints_0.json",
			profile:          routingkit.Car(),
		},
		{
			source:           []float32{-76.587490, 39.299710},
			destination:      []float32{-76.584897, 39.280774},
			snap:             1000,
			expectedDistance: 1897,
			osmFile:          marylandMap,
			waypointsFile:    "waypoints_1.json",
			profile:          routingkit.Bike(),
		},
		{
			source:           []float32{-76.587490, 39.299710},
			destination:      []float32{-76.584897, 39.280774},
			snap:             1000,
			expectedDistance: 1777,
			osmFile:          marylandMap,
			waypointsFile:    "waypoints_2.json",
			profile:          routingkit.Pedestrian(),
		},
		{
			source:           []float32{-76.587490, 39.299710},
			destination:      []float32{-76.582855, 39.309095},
			snap:             1000,
			profile:          routingkit.Car(),
			osmFile:          marylandMap,
			expectedDistance: 1496,
			waypointsFile:    "waypoints_3.json",
		},
		{
			source:           []float32{-76.587490, 39.299710},
			destination:      []float32{-76.582855, 39.309095},
			snap:             1000,
			profile:          routingkit.Bike(),
			osmFile:          marylandMap,
			expectedDistance: 1496,
			waypointsFile:    "waypoints_4.json",
		},
		{
			source:           []float32{-76.587490, 39.299710},
			destination:      []float32{-76.582855, 39.309095},
			snap:             1000,
			profile:          routingkit.Pedestrian(),
			osmFile:          marylandMap,
			expectedDistance: 1429,
			waypointsFile:    "waypoints_5.json",
		},
		{
			source:           []float32{-76.587490, 39.299710},
			destination:      []float32{-76.591286, 39.298443},
			snap:             1000,
			profile:          routingkit.Car(),
			osmFile:          marylandMap,
			expectedDistance: 617,
			waypointsFile:    "waypoints_6.json",
		},
		{
			source:           []float32{-76.587490, 39.299710},
			destination:      []float32{-76.591286, 39.298443},
			snap:             1000,
			profile:          routingkit.Bike(),
			osmFile:          marylandMap,
			expectedDistance: 617,
			waypointsFile:    "waypoints_7.json",
		},
		{
			source:           []float32{-76.587490, 39.299710},
			destination:      []float32{-76.591286, 39.298443},
			snap:             912,
			profile:          routingkit.Pedestrian(),
			osmFile:          marylandMap,
			expectedDistance: 428,
			waypointsFile:    "waypoints_8.json",
		},
		{
			source: []float32{-76.587490, 39.299710},
			// point is in a river so should not snap
			destination:      []float32{-76.584897, 39.280774},
			snap:             10,
			osmFile:          marylandMap,
			expectedDistance: routingkit.MaxDistance,
			waypointsFile:    "waypoints_9.json",
			profile:          routingkit.Car(),
		},
		// a truck with this height will need to go around the train overpass
		{
			source:           []float32{-76.638843, 39.254254},
			destination:      []float32{-76.637647, 39.256933},
			snap:             1000,
			osmFile:          marylandMapWithHeightRestriction,
			expectedDistance: 1972,
			waypointsFile:    "waypoints_10.json",
			profile:          routingkit.Truck(4.25, 0, 0, 0, 100),
		},
	}

	for i, test := range tests {
		cli, err := routingkit.NewDistanceClient(test.osmFile, test.profile)
		if err != nil {
			t.Fatalf("creating Client: %v", err)
		}
		cli.SetSnapRadius(test.snap)
		distance, waypoints := cli.Route(test.source, test.destination)
		if test.expectedDistance != distance {
			t.Errorf("[%d] expected distance %v, got %v", i, test.expectedDistance, distance)
		}

		waypointsFile, err := os.OpenFile(
			filepath.Join("testdata/fixtures", test.waypointsFile),
			os.O_RDONLY,
			0655,
		)
		if err != nil {
			t.Errorf("[%d] opening test fixture: %v", i, err)
			continue
		}
		var expectedWaypoints [][]float32
		if err := json.NewDecoder(waypointsFile).Decode(&expectedWaypoints); err != nil && err != io.EOF {
			t.Errorf("[%d] loading test fixture: %v", i, err)
			continue
		}
		if diff := cmp.Diff(expectedWaypoints, waypoints); diff != "" {
			if *update {
				waypointsFileW, err := os.OpenFile(
					filepath.Join("testdata/fixtures", test.waypointsFile),
					os.O_CREATE|os.O_WRONLY|os.O_TRUNC,
					0655,
				)
				if err != nil {
					t.Errorf("[%d] opening test fixture for update: %v", i, err)
					continue
				}
				if err := json.NewEncoder(waypointsFileW).Encode(waypoints); err != nil {
					t.Errorf("[%d updating fixtures: %v", i, err)
				}
			}
			expectedPlot, gotPlot, err := plotWaypoints(i, expectedWaypoints, waypoints)
			if err == nil {
				t.Errorf(
					"[%d] waypoints mismatch (-want +got):\n%s\nsee expected path at %s and actual path at %s",
					i,
					diff,
					expectedPlot,
					gotPlot,
				)
			} else {
				t.Errorf(
					"[%d] expected waypoints %v, got %v\nerror plotting paths: %v",
					i,
					expectedWaypoints,
					waypoints,
					err,
				)
			}
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
		cli, err := routingkit.NewTravelTimeClient(marylandMap, routingkit.Car())
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
		cli, err := routingkit.NewTravelTimeClient(marylandMap, routingkit.Car())
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
		profile     routingkit.Profile

		expectedTravelTime uint32
		waypointsFile      string
	}{
		{
			source:             []float32{-76.587490, 39.299710},
			destination:        []float32{-76.584897, 39.280774},
			snap:               1000,
			expectedTravelTime: 213405,
			waypointsFile:      "travel_time_waypoints_0.json",
			profile:            routingkit.Car(),
		},
		{
			source:             []float32{-76.587490, 39.299710},
			destination:        []float32{-76.591286, 39.298443},
			snap:               1000,
			expectedTravelTime: 53712,
			waypointsFile:      "travel_time_waypoints_1.json",
			profile:            routingkit.Bike(),
		},
	}

	for i, test := range tests {
		cli, err := routingkit.NewTravelTimeClient(marylandMap, test.profile)
		if err != nil {
			t.Fatalf("creating Client: %v", err)
		}
		cli.SetSnapRadius(test.snap)
		travelTime, waypoints := cli.Route(test.source, test.destination)
		if test.expectedTravelTime != travelTime {
			t.Errorf("[%d] expected travel time %v, got %v", i, test.expectedTravelTime, travelTime)
		}

		waypointsFile, err := os.OpenFile(
			filepath.Join("testdata/fixtures", test.waypointsFile),
			os.O_RDONLY,
			0655,
		)
		var expectedWaypoints [][]float32
		if err := json.NewDecoder(waypointsFile).Decode(&expectedWaypoints); err != nil && err != io.EOF {
			t.Errorf("[%d] loading test fixture: %v", i, err)
			continue
		}
		if err != nil {
			t.Errorf("[%d] opening test fixture: %v", i, err)
			continue
		}
		if diff := cmp.Diff(expectedWaypoints, waypoints); diff != "" {
			if *update {
				waypointsFileW, err := os.OpenFile(
					filepath.Join("testdata/fixtures", test.waypointsFile),
					os.O_CREATE|os.O_WRONLY|os.O_TRUNC,
					0655,
				)
				if err != nil {
					t.Errorf("[%d] opening test fixture for update: %v", i, err)
					continue
				}
				if err := json.NewEncoder(waypointsFileW).Encode(waypoints); err != nil {
					t.Errorf("[%d updating fixtures: %v", i, err)
				}
			}
			expectedPlot, gotPlot, err := plotWaypoints(i, expectedWaypoints, waypoints)
			if err == nil {
				t.Errorf(
					"[%d] waypoints mismatch (-want +got):\n%s\nsee expected path at %s and actual path at %s",
					i,
					diff,
					expectedPlot,
					gotPlot,
				)
			} else {
				t.Errorf(
					"[%d] expected waypoints %v, got %v\nerror plotting paths: %v",
					i,
					expectedWaypoints,
					waypoints,
					err,
				)
			}
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
	cli, err := routingkit.NewDistanceClient(marylandMap, routingkit.Car())
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
	cli, err := routingkit.NewDistanceClient(marylandMap, routingkit.Car())
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
