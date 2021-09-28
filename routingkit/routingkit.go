package routingkit

import (
	"fmt"
	"os"
	"runtime"

	"github.com/nextmv-io/go-routingkit/routingkit/internal/routingkit"
)

// MaxDistance represents the maximum possible route distance.
var MaxDistance uint32

type TravelProfile routingkit.GoRoutingKitTravel_profile

var CarTravelProfile, BikeTravelProfile, PedestrianTravelProfile TravelProfile

type Wayfilter struct {
	// the Tag to be matched
	Tag string
	// either the Tag has to match or the Tag is not allowed to match
	MatchTag bool
	// optional: the value that the Tag has to equal
	Value string
	// either the value has to match or the value is not allowed to match
	MatchValue bool
	// expresses whether this way is allowed or not due to this filter
	Allowed bool
}

func Car() Profile {
	profile := Profile{
		Name: "car",
		Wayfilters: []Wayfilter{
			{
				Tag:      "junction",
				MatchTag: true,
				Allowed:  true,
			},
			{
				Tag:        "route",
				MatchTag:   true,
				Value:      "ferry",
				MatchValue: true,
				Allowed:    true,
			},
			{
				Tag:        "ferry",
				MatchTag:   true,
				Value:      "yes",
				MatchValue: true,
				Allowed:    true,
			},
			{
				Tag:      "highway",
				MatchTag: false,
				Allowed:  false,
			},
			{
				Tag:        "motorcar",
				MatchTag:   true,
				Value:      "no",
				MatchValue: true,
				Allowed:    false,
			},
			{
				Tag:        "motor_vehicle",
				MatchTag:   true,
				Value:      "no",
				MatchValue: true,
				Allowed:    false,
			},
			{
				Tag:        "access",
				MatchTag:   true,
				Value:      "no",
				MatchValue: true,
				Allowed:    false,
			},
			{
				Tag:        "highway",
				MatchTag:   true,
				Value:      "construction",
				MatchValue: true,
				Allowed:    false,
			},
			{
				Tag:        "highway",
				MatchTag:   true,
				Value:      "path",
				MatchValue: true,
				Allowed:    false,
			},
			{
				Tag:        "highway",
				MatchTag:   true,
				Value:      "footway",
				MatchValue: true,
				Allowed:    false,
			},
			{
				Tag:        "highway",
				MatchTag:   true,
				Value:      "cycleway",
				MatchValue: true,
				Allowed:    false,
			},
			{
				Tag:        "highway",
				MatchTag:   true,
				Value:      "bridleway",
				MatchValue: true,
				Allowed:    false,
			},
			{
				Tag:        "highway",
				MatchTag:   true,
				Value:      "pedestrian",
				MatchValue: true,
				Allowed:    false,
			},
			{
				Tag:        "highway",
				MatchTag:   true,
				Value:      "bus_guideway",
				MatchValue: true,
				Allowed:    false,
			},
			{
				Tag:        "highway",
				MatchTag:   true,
				Value:      "raceway",
				MatchValue: true,
				Allowed:    false,
			},
			{
				Tag:        "highway",
				MatchTag:   true,
				Value:      "escape",
				MatchValue: true,
				Allowed:    false,
			},
			{
				Tag:        "highway",
				MatchTag:   true,
				Value:      "steps",
				MatchValue: true,
				Allowed:    false,
			},
			{
				Tag:        "highway",
				MatchTag:   true,
				Value:      "proposed",
				MatchValue: true,
				Allowed:    false,
			},
			{
				Tag:        "highway",
				MatchTag:   true,
				Value:      "conveying",
				MatchValue: true,
				Allowed:    false,
			},
			{
				Tag:        "highway",
				MatchTag:   true,
				Value:      "corridor",
				MatchValue: true,
				Allowed:    false,
			},
			{
				Tag:      "maxspeed",
				MatchTag: false,
				Allowed:  false,
			},
		},
	}
	return profile
}

func Bike() Profile {
	profile := Profile{
		Name: "bike",
		Wayfilters: []Wayfilter{
			{
				Tag: "highway",
			},
		},
	}
	return profile
}

func Pedestrian() Profile {
	profile := Profile{
		Name: "pedestrian",
		Wayfilters: []Wayfilter{
			{
				Tag:      "junction",
				MatchTag: true,
				Allowed:  true,
			},
			{
				Tag:      "route",
				MatchTag: true,
				Value:    "ferry",
				Allowed:  true,
			},
			{
				Tag:      "ferry",
				MatchTag: true,
				Value:    "ferry",
				Allowed:  true,
			},
			{
				Tag:      "public_transport",
				MatchTag: true,
				Value:    "stop_position",
				Allowed:  true,
			},
			{
				Tag:      "public_transport",
				MatchTag: true,
				Value:    "platform",
				Allowed:  true,
			},
			{
				Tag:      "public_transport",
				MatchTag: true,
				Value:    "stop_area",
				Allowed:  true,
			},
			{
				Tag:      "public_transport",
				MatchTag: true,
				Value:    "station",
				Allowed:  true,
			},
			{
				Tag:      "railway",
				MatchTag: true,
				Value:    "halt",
				Allowed:  true,
			},
			{
				Tag:      "railway",
				MatchTag: true,
				Value:    "platform",
				Allowed:  true,
			},
			{
				Tag:      "railway",
				MatchTag: true,
				Value:    "subway_entrance",
				Allowed:  true,
			},
			{
				Tag:      "railway",
				MatchTag: true,
				Value:    "station",
				Allowed:  true,
			},
			{
				Tag:      "railway",
				MatchTag: true,
				Value:    "tram_stop",
				Allowed:  true,
			},
			{
				Tag: "highway",
			},
		},
	}
	return profile
}

type Profile struct {
	Wayfilters []Wayfilter
	Name       string
}

func (p Profile) swigProfile() routingkit.Profile {
	customProfile := routingkit.NewProfile()
	wayFilterVector := routingkit.NewWayFilterVector()
	for _, wayFilter := range p.Wayfilters {
		wf := routingkit.NewWayFilter()
		wf.SetTag(wayFilter.Tag)
		wf.SetMatchTag(wayFilter.MatchTag)
		wf.SetValue(wayFilter.Value)
		wf.SetMatchValue(wayFilter.MatchValue)
		wf.SetAllowed(wayFilter.Allowed)
		wayFilterVector.Add(wf)
	}
	customProfile.SetWayfilters(wayFilterVector)
	return customProfile
}

func init() {
	MaxDistance = uint32(routingkit.GetMax_distance())
	CarTravelProfile = TravelProfile(routingkit.Car)
	BikeTravelProfile = TravelProfile(routingkit.Bike)
	PedestrianTravelProfile = TravelProfile(routingkit.Pedestrian)
}

// NewDistanceClient initializes a DistanceClient using the provided .osm.pbf file and
// .ch file. The .ch file will be created if it does not already exist. It is the caller's
// responsibility to call Delete on the client when it is no longer needed.
func NewDistanceClient(mapFile string, p Profile) (DistanceClient, error) {
	if _, err := os.Stat(mapFile); os.IsNotExist(err) {
		return DistanceClient{}, fmt.Errorf("could not find map file at %v", mapFile)
	}

	chFile, err := chFileName(mapFile, p, false)
	if err != nil {
		return DistanceClient{}, err
	}

	concurrentQueries := runtime.GOMAXPROCS(0)
	customProfile := p.swigProfile()
	defer func() {
		routingkit.DeleteProfile(customProfile)
	}()
	c := routingkit.NewClient(concurrentQueries, mapFile, chFile, routingkit.GoRoutingKitTravel_profile(CarTravelProfile), customProfile)

	channel := make(chan int, concurrentQueries)
	for i := 0; i < concurrentQueries; i++ {
		channel <- i
	}

	return DistanceClient{
		client: client{
			client:     c,
			channel:    channel,
			snapRadius: 1000,
		}}, nil
}

func chFileName(mapFile string, profile Profile, duration bool) (string, error) {
	extension := profile.Name
	if profile.Name == "" {
		return "", fmt.Errorf("profile name was empty")
	}

	distOrDuration := "distance"
	if duration {
		distOrDuration = "duration"
	}
	return mapFile + "_" + extension + "_" + distOrDuration + ".ch", nil
}

// Delete deletes the client, releasing memory allocated for C++ routing data structures
func (c *client) Delete() {
	routingkit.DeleteClient(c.client)
}

// SetSnapRadius updates Client so that all queries will snap points to the nearest
// street network point within the given radius in meters.
func (c *client) SetSnapRadius(n float32) {
	c.snapRadius = n
}

type DistanceClient struct {
	client
}

// client allows routing queries to be executed against a particular region.
type client struct {
	client     routingkit.Client
	channel    chan int
	snapRadius float32
}

// Route finds the fastest route between the two points, returning the total route
// distance and the waypoints describing the route.
func (c client) Route(from []float32, to []float32) (uint32, [][]float32) {
	counter := <-c.channel
	defer func() {
		c.channel <- counter
	}()
	resp := c.client.Query(
		int(counter),
		float32(c.snapRadius),
		float32(from[0]),
		float32(from[1]),
		float32(to[0]),
		float32(to[1]),
		true,
	)
	defer routingkit.DeleteQueryResponse(resp)
	wp := resp.GetWaypoints()
	waypoints := make([][]float32, wp.Size())
	for i := 0; i < len(waypoints); i++ {
		p := wp.Get(i)
		waypoints[i] = []float32{float32(p.GetLon()), float32(p.GetLat())}
	}

	return uint32(resp.GetDistance()), waypoints
}

// Distance returns the length of the shortest possible route between the points
func (c client) Distance(from []float32, to []float32) uint32 {
	counter := <-c.channel
	defer func() {
		c.channel <- counter
	}()
	resp := c.client.Query(
		int(counter),
		c.snapRadius,
		from[0],
		from[1],
		to[0],
		to[1],
		false,
	)
	defer routingkit.DeleteQueryResponse(resp)

	return uint32(resp.GetDistance())
}

type distanceMatrixRow struct {
	i         int
	distances []uint32
}

// Nearest returns the nearest point in the road network within the radius configured on
// the Client. The second argument will be false if no point could be found.
func (c client) Nearest(point []float32) ([]float32, bool) {
	counter := <-c.channel
	defer func() {
		c.channel <- counter
	}()
	res := c.client.Nearest(counter, c.snapRadius, point[0], point[1])
	if res.Swigcptr() == 0 {
		return nil, false
	}
	defer routingkit.DeletePoint(res)
	return []float32{res.GetLon(), res.GetLat()}, true
}

// Matrix creates a matrix representing the minimum distances from the points in
// sources to the points in targets.
func (c client) Matrix(sources [][]float32, targets [][]float32) [][]uint32 {
	matrix := make([][]uint32, len(sources))

	workers := make(chan struct{}, runtime.GOMAXPROCS(0))
	results := make(chan distanceMatrixRow)

	go func() {
		for i, source := range sources {
			workers <- struct{}{}
			go func(i int, source []float32) {
				distances := c.Distances(source, targets)
				results <- distanceMatrixRow{i, distances}
				<-workers
			}(i, source)
		}
	}()

	for range sources {
		matrixRow := <-results
		matrix[matrixRow.i] = matrixRow.distances
	}

	return matrix
}

// Distances returns a slice containing the minimum distances from the source to the
// points in targets.
func (c client) Distances(source []float32, targets [][]float32) []uint32 {
	counter := <-c.channel
	defer func() {
		c.channel <- counter
	}()

	s := routingkit.NewPoint()
	defer routingkit.DeletePoint(s)
	s.SetLon(float32(source[0]))
	s.SetLat(float32(source[1]))

	targetsVector := routingkit.NewPointVector(int64(len(targets)))
	defer routingkit.DeletePointVector(targetsVector)

	for i := 0; i < len(targets); i++ {
		t := routingkit.NewPoint()
		t.SetLon(float32(targets[i][0]))
		t.SetLat(float32(targets[i][1]))
		targetsVector.Set(i, t)
	}

	distanceVec := c.client.Distances(counter, float32(c.snapRadius), s, targetsVector)
	defer routingkit.DeleteUnsignedVector(distanceVec)
	numDistances := distanceVec.Size()
	distances := make([]uint32, numDistances)
	for i := 0; i < int(numDistances); i++ {
		col := uint32(distanceVec.Get(i))
		distances[i] = col
	}

	return distances
}

type TravelTimeClient struct {
	client client
}

// NewTravelTimeClient initializes a TravelTimeClient using the provided .osm.pbf file and
// .ch file. The .ch file will be created if it does not already exist. It is the caller's
// responsibility to call Delete on the client when it is no longer needed.
func NewTravelTimeClient(mapFile string, profile Profile) (TravelTimeClient, error) {
	if _, err := os.Stat(mapFile); os.IsNotExist(err) {
		return TravelTimeClient{}, fmt.Errorf("could not find map file at %v", mapFile)
	}
	chFile, err := chFileName(mapFile, profile, true)
	if err != nil {
		return TravelTimeClient{}, err
	}
	concurrentQueries := runtime.GOMAXPROCS(0)
	customProfile := profile.swigProfile()
	// sets that we are interested in the travel time rather than the distance
	customProfile.SetTravel_time(true)
	defer func() {
		routingkit.DeleteProfile(customProfile)
	}()
	c := routingkit.NewClient(concurrentQueries, mapFile, chFile, routingkit.Car, customProfile)

	channel := make(chan int, concurrentQueries)
	for i := 0; i < concurrentQueries; i++ {
		channel <- i
	}

	return TravelTimeClient{
		client: client{
			client:     c,
			channel:    channel,
			snapRadius: 1000,
		}}, nil
}

// Route finds the fastest route between the two points, returning the total route
// travel time by car and the waypoints describing the route.
func (c TravelTimeClient) Route(from []float32, to []float32) (uint32, [][]float32) {
	return c.client.Route(from, to)
}

// TravelTime returns the travel time by car for the shortest possible route between
// the points.
func (c TravelTimeClient) TravelTime(from []float32, to []float32) uint32 {
	return c.client.Distance(from, to)
}

// Nearest returns the nearest point in the road network within the radius configured on
// the Client. The second argument will be false if no point could be found.
func (c TravelTimeClient) Nearest(point []float32) ([]float32, bool) {
	return c.client.Nearest(point)
}

// Matrix creates a matrix representing the minimum travel times (by car) from the
// points in sources to the points in targets.
func (c TravelTimeClient) Matrix(sources [][]float32, targets [][]float32) [][]uint32 {
	return c.client.Matrix(sources, targets)
}

// TravelTimes returns a slice containing the minimum car travel times from the source
// to the points in targets.
func (c TravelTimeClient) TravelTimes(source []float32, targets [][]float32) []uint32 {
	return c.client.Distances(source, targets)
}

// SetSnapRadius updates Client so that all queries will snap points to the nearest
// street network point within the given radius in meters.
func (c *TravelTimeClient) SetSnapRadius(n float32) {
	c.client.SetSnapRadius(n)
}
