package routingkit_test

import (
	"fmt"
	"io/ioutil"
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

func TestTables(t *testing.T) {
	tests := []struct {
		source       []float64
		destinations [][]float64
		expected     []float64
	}{
		{
			source: []float64{-76.587490, 39.299710},
			destinations: [][]float64{
				{-76.582855, 39.309095},
				{-76.599388, 39.302014},
			},
			expected: []float64{1496, 1259},
		},
	}
	chFile, err := tempFile("", "routingkit-test.ch")
	defer os.Remove(chFile)
	cli, err := routingkit.NewClient("testdata/maryland.osm.pbf", chFile)
	if err != nil {
		t.Fatalf("creating Client: %v", err)
	}

	for i, test := range tests {
		got := cli.Table(test.source, test.destinations)
		if !reflect.DeepEqual(test.expected, got) {
			t.Errorf("[%d] expected %v, got %v", i, test.expected, got)
		}
	}
}

func TestRoutingKit(t *testing.T) {
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
	cli, err := routingkit.NewClient("testdata/maryland.osm.pbf", chFile)
	if err != nil {
		t.Fatalf("creating Client: %v", err)
	}

	for i, test := range tests {
		cli.SetSnapRadius(test.snap)
		distance, waypoints := cli.Query(test.source, test.destination)
		if test.expectedDistance != distance {
			t.Errorf("[%d] expected distance %v, got %v", i, test.expectedDistance, distance)
		}
		if !reflect.DeepEqual(test.expectedWaypoints, waypoints) {
			t.Errorf("[%d] expected waypoints %v, got %v", i, test.expectedWaypoints, waypoints)
		}
	}

}
