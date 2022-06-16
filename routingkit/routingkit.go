package routingkit

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"runtime"
	"sort"
	"strconv"

	routingkit "github.com/nextmv-io/go-routingkit/routingkit/bindings"
	"github.com/paulmach/osm"
	"github.com/paulmach/osm/osmpbf"
)

// MaxDistance represents the maximum possible route distance.
var MaxDistance uint32

func parsePBF(osmFile string, tagMapFilter TagMapFilter, speedMapper SpeedMapper) (map[int]bool, map[int]int) {
	file, err := os.Open(osmFile)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// The third parameter is the number of parallel decoders to use.
	scanner := osmpbf.New(context.Background(), file, runtime.GOMAXPROCS(0))
	scanner.SkipNodes = true
	scanner.SkipRelations = true
	defer scanner.Close()

	allowed := map[int]bool{}
	waySpeeds := map[int]int{}

	for scanner.Scan() {
		switch o := scanner.Object().(type) {
		case *osm.Way:
			id := int(o.ID)
			tagMap := o.Tags.Map()
			if tagMapFilter != nil && tagMapFilter(id, tagMap) {
				allowed[id] = true
				// we only need to write the speed into ways that are actually
				// allowed
				if speedMapper != nil {
					waySpeeds[id] = speedMapper(id, tagMap)
					// if the speed is returned as 0, that means the way is not
					// usable
					if waySpeeds[id] == 0 {
						waySpeeds[id] = 1 // prevents routingkit crash
						allowed[id] = false
					}
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	return allowed, waySpeeds
}

func Car() Profile {
	return NewProfile("car", VehicleMode, false, false, CarTagMapFilter, CarSpeedMapper)
}

func Bike() Profile {
	return NewProfile("bike", BikeMode, false, false, BikeTagMapFilter, BikeSpeedMapper)
}

func Pedestrian() Profile {
	return NewProfile(
		"pedestrian",
		PedestrianMode,
		false,
		false,
		PedestrianTagMapFilter,
		PedestrianSpeedMapper,
	)
}

func Truck(height, width, length, weight float64, speed int) Profile {
	return NewProfile(
		"truck",
		VehicleMode,
		false,
		true,
		TruckTagMapFilter(height, width, length, weight),
		// TODO: instead of reusing the car max speed and capping it at a given
		// maxSpeed, we could/should also react to maxspeed:hgv
		MaxSpeedMapper(speed),
	)
}

type TransportMode routingkit.Transport_mode

var (
	VehicleMode    TransportMode = TransportMode(routingkit.Vehicle)
	BikeMode       TransportMode = TransportMode(routingkit.Bike)
	PedestrianMode TransportMode = TransportMode(routingkit.Pedestrian)
)

type Profile struct {
	Name             string
	TransportMode    TransportMode
	PreventLeftTurns bool
	PreventUTurns    bool
	Filter           TagMapFilter
	SpeedMapper      SpeedMapper
}

func NewProfile(
	name string,
	transportMode TransportMode,
	preventLeftTurns bool,
	preventUTurns bool,
	filter TagMapFilter,
	speedMapper SpeedMapper,
) Profile {
	return Profile{
		Name:             name,
		TransportMode:    transportMode,
		PreventLeftTurns: preventLeftTurns,
		PreventUTurns:    preventUTurns,
		Filter:           filter,
		SpeedMapper:      speedMapper,
	}
}

func withSwigProfile(p Profile, allowedWayIDs map[int]bool, waySpeeds map[int]int, f func(routingkit.Profile)) {
	customProfile := routingkit.NewProfile()
	customProfile.SetName(p.Name)
	customProfile.SetTransportMode(routingkit.Transport_mode(p.TransportMode))
	customProfile.SetPrevent_left_turns(p.PreventLeftTurns)
	customProfile.SetPrevent_u_turns(p.PreventUTurns)

	allowedWayIds := routingkit.NewIntVector()
	for wayId := range allowedWayIDs {
		allowedWayIds.Add(wayId)
	}
	customProfile.SetAllowedWayIds(allowedWayIds)

	rkWaySpeeds := routingkit.NewIntIntMap()
	for wayId, speed := range waySpeeds {
		rkWaySpeeds.Set(uint64(wayId), uint(speed))
	}
	customProfile.SetWaySpeeds(rkWaySpeeds)

	defer func() {
		routingkit.DeleteIntVector(allowedWayIds)
		routingkit.DeleteIntIntMap(rkWaySpeeds)
		routingkit.DeleteProfile(customProfile)
	}()

	f(customProfile)
}

func init() {
	MaxDistance = uint32(routingkit.GetMax_distance())
}

// NewDistanceClient initializes a DistanceClient using the provided .osm.pbf file and
// .ch file. The .ch file will be created if it does not already exist. It is the caller's
// responsibility to call Delete on the client when it is no longer needed.
func NewDistanceClient(mapFile string, profile Profile) (DistanceClient, error) {
	if _, err := os.Stat(mapFile); os.IsNotExist(err) {
		return DistanceClient{}, fmt.Errorf("could not find map file at %v", mapFile)
	}

	allowedWayIDs, waySpeeds := parsePBF(mapFile, profile.Filter, profile.SpeedMapper)

	chFile, err := chFileName(mapFile, profile, allowedWayIDs, waySpeeds, false)
	if err != nil {
		return DistanceClient{}, err
	}

	concurrentQueries := runtime.GOMAXPROCS(0)
	var c routingkit.Client
	withSwigProfile(profile, allowedWayIDs, waySpeeds, func(customProfile routingkit.Profile) {
		c = routingkit.NewClient(concurrentQueries, mapFile, chFile, customProfile)
	})

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

func chFileName(mapFile string, profile Profile, allowedWayIDs map[int]bool, waySpeeds map[int]int, duration bool) (string, error) {
	extension := profile.Name
	if profile.Name == "" {
		return "", fmt.Errorf("profile name was empty")
	}

	distOrDuration := "distance"
	if duration {
		distOrDuration = "duration"
	}

	// compute a hash based on the contents of the profile
	h := sha1.New()

	// iterate over profile.AllowedWayIds in order
	keys := make([]int, 0)
	for k := range allowedWayIDs {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	for _, s := range keys {
		_, _ = io.WriteString(h, "-")
		_, _ = io.WriteString(h, strconv.Itoa(s))
	}

	// iterate over profile.WaySpeeds in order
	keys = make([]int, 0)
	for k := range waySpeeds {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	for _, s := range keys {
		_, _ = io.WriteString(h, "-")
		_, _ = io.WriteString(h, strconv.Itoa(waySpeeds[s]))
		_, _ = io.WriteString(h, "-")
		_, _ = io.WriteString(h, strconv.Itoa(s))
	}

	// add simple fields
	_, _ = io.WriteString(h, "-")
	_, _ = io.WriteString(h, strconv.FormatBool(profile.PreventLeftTurns))
	_, _ = io.WriteString(h, "-")
	_, _ = io.WriteString(h, strconv.Itoa(int(profile.TransportMode)))
	hash := hex.EncodeToString(h.Sum(nil))

	return mapFile + "_" + extension + "_" + distOrDuration + "_" + hash + ".ch", nil
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

	allowedWayIDs, waySpeeds := parsePBF(mapFile, profile.Filter, profile.SpeedMapper)
	chFile, err := chFileName(mapFile, profile, allowedWayIDs, waySpeeds, true)
	if err != nil {
		return TravelTimeClient{}, err
	}

	concurrentQueries := runtime.GOMAXPROCS(0)
	var c routingkit.Client
	withSwigProfile(profile, allowedWayIDs, waySpeeds, func(swigProfile routingkit.Profile) {
		// sets that we are interested in the travel time rather than the distance
		swigProfile.SetTravel_time(true)
		c = routingkit.NewClient(concurrentQueries, mapFile, chFile, swigProfile)
	})

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
