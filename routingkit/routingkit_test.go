package routingkit_test

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/nextmv-io/go-routingkit/routingkit"
)

// This is a small map file containing data for the bounding box from
// -76.60735000000001,39.28971 to -76.57749,39.31587
var marylandMap string = "testdata/maryland.osm.pbf"

// This is another small map of Maryland that contains data for the
// bounding box from -76.663640,39.240043 to -76.605023,39.269623.
// A low overpass with a max height of 13'11" is found at -76.638449,39.254932
// (way ID 456490563), where Annapolis Rd. passes under a train line
var marylandMapWithHeightRestriction = "testdata/maryland_height_restriction.osm.pbf"

// This is a small area of England, covering data from 0.301445,51.363350 to
// 0.371563,51.392494. There is a narrow pass on Harley Bottom Road that a
// vehicle with a width larger than 6'6" will need to go around.
var englandMapWithWidthRestriction = "testdata/england_width_restriction.osm.pbf"

// This is a part of London, covering data from -0.121322,51.508732 to -0.088926,51.525289.
// There is a 40' length restriction on Fleet and Farringdon Streets.
var englandMapWithLengthRestriction = "testdata/england_length_restriction.osm.pbf"

// This is a section of London, covering -0.100959,51.487403 to -0.083680,51.496247.
// Larcom Street has a weight limit of 7.5 tonnes
var englandMapWithWeightRestriction = "testdata/england_weight_restriction.osm.pbf"

func cleanCH() error {
	files, err := filepath.Glob("testdata/*.ch")
	if err != nil {
		return fmt.Errorf("removing testdata ch files: %v", err)
	}
	for _, f := range files {
		if err := os.Remove(f); err != nil {
			return fmt.Errorf("removing file %s: %v", f, err)
		}
	}
	return nil
}

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
	var pathsToPlot []string

	if len(expectedWaypoints) > 0 {
		if err := json.NewEncoder(expectedFile).Encode(waypoints{expectedWaypoints}); err != nil {
			return "", "", fmt.Errorf("writing expected points: %v", err)
		}
		pathsToPlot = append(pathsToPlot, expectedPath)
		expectedPath = expectedPath + ".html"
	} else {
		expectedPath = ""
	}
	if len(gotWaypoints) > 0 {
		if err := json.NewEncoder(gotFile).Encode(waypoints{gotWaypoints}); err != nil {
			return "", "", fmt.Errorf("writing expected points: %v", err)
		}
		pathsToPlot = append(pathsToPlot, gotPath)
		gotPath = gotPath + ".html"
	} else {
		gotPath = ""
	}

	for _, path := range pathsToPlot {
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

	return expectedPath, gotPath, nil
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

			expected: []uint32{1440, 617},
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
				{1440, 1259},
				{1796, 575},
				{2372, 2216},
				{3364, 1548},
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
				{1429, 1242},
				{1589, 558},
				{2367, 2189},
				{3151, 1533},
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
var cleanCHFiles *bool

func init() {
	update = flag.Bool("update", false, "update text fixtures")
	cleanCHFiles = flag.Bool("clean_ch", true, "clean CH files in testdata")
}

func TestMain(m *testing.M) {
	flag.Parse()
	if *cleanCHFiles {
		if err := cleanCH(); err != nil {
			fmt.Fprintf(os.Stderr, "error cleaning CH files: %v", err)
			os.Exit(1)
		}
	}
	exit := m.Run()
	os.Exit(exit)
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
			expectedDistance: 1440,
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
			snap:             1000,
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
		// a truck with this width will need to go around the narrow pass
		{
			source:           []float32{0.328562, 51.387527},
			destination:      []float32{0.328830, 51.389174},
			snap:             1000,
			osmFile:          englandMapWithWidthRestriction,
			expectedDistance: 9303,
			waypointsFile:    "waypoints_11.json",
			profile:          routingkit.Truck(4.25, 2.0, 0, 0, 100),
		},
		// Truck should avoid going down Fleet St. due to the length restriction
		{
			source:           []float32{-0.106210, 51.514208},
			destination:      []float32{-0.103678, 51.514181},
			snap:             1000,
			osmFile:          englandMapWithLengthRestriction,
			expectedDistance: 527,
			waypointsFile:    "waypoints_12.json",
			profile:          routingkit.Truck(4.25, 2.0, 13.0, 0, 100),
		},
		// Truck should avoid going down Larcom St. due to the weight restriction
		{
			source:           []float32{-0.096975, 51.490302},
			destination:      []float32{-0.093553, 51.491785},
			snap:             1000,
			osmFile:          englandMapWithWeightRestriction,
			expectedDistance: 942,
			waypointsFile:    "waypoints_13.json",
			profile:          routingkit.Truck(4.25, 2.0, 13.0, 8.0, 100),
		},
		{
			source:           []float32{-76.594045, 39.300524},
			destination:      []float32{-76.582855, 39.309095},
			snap:             1000,
			osmFile:          marylandMap,
			expectedDistance: 1589,
			waypointsFile:    "waypoints_14.json",
			profile:          routingkit.Pedestrian(),
		},
		{
			source:           []float32{-76.598423, 39.289484},
			destination:      []float32{-76.599388, 39.302014},
			snap:             1000,
			osmFile:          marylandMap,
			expectedDistance: 1533,
			waypointsFile:    "waypoints_15.json",
			profile:          routingkit.Pedestrian(),
		},
		{
			source:           []float32{-76.58749, 39.29971},
			destination:      []float32{-76.59735, 39.30587},
			snap:             1000,
			osmFile:          marylandMap,
			expectedDistance: 1555,
			waypointsFile:    "waypoints_16.json",
			profile:          routingkit.Pedestrian(),
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
			t.Errorf("[%d] route: expected distance %v, got %v", i, test.expectedDistance, distance)
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
			msg := fmt.Sprintf("[%d] waypoints mismatch (-want +got)\n%s\n", i, diff)
			if err == nil {
				if expectedPlot != "" {
					msg += fmt.Sprintf("see expected path at %s\n", expectedPlot)
				}
				if gotPlot != "" {
					msg += fmt.Sprintf("see actual path at %s\n", gotPlot)
				}
			} else {
				msg += fmt.Sprintf("error plotting paths: %v", err)
			}
			t.Errorf(msg)
		}
		distance = cli.Distance(test.source, test.destination)
		if test.expectedDistance != distance {
			t.Errorf("[%d] distance: expected distance %v, got %v", i, test.expectedDistance, distance)
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

			expected: []uint32{131599, 54169},
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
				{131599, 110399},
				{153578, 52230},
				{205841, 187148},
				{289818, 127151},
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
			expectedTravelTime: 206713,
			waypointsFile:      "travel_time_waypoints_0.json",
			profile:            routingkit.Car(),
		},
		{
			source:             []float32{-76.587490, 39.299710},
			destination:        []float32{-76.591286, 39.298443},
			snap:               1000,
			expectedTravelTime: 148080,
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
		const tolerance = .0001
		opt := cmp.Comparer(func(x, y float64) bool {
			diff := math.Abs(x - y)
			mean := math.Abs(x+y) / 2.0
			return (diff / mean) < tolerance
		})
		if diff := cmp.Diff(expectedWaypoints, waypoints, opt); diff != "" {
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
			msg := fmt.Sprintf("[%d] waypoints mismatch (-want +got)\n%s\n", i, diff)
			if err == nil {
				msg += fmt.Sprintf("error plotting paths: %v", err)
				continue
			}
			if expectedPlot != "" {
				msg += fmt.Sprintf("see expected path at %s.\n", expectedPlot)
			}
			if gotPlot != "" {
				msg += fmt.Sprintf("see actual path at %s.\n", gotPlot)
			}
			t.Errorf(msg)
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
