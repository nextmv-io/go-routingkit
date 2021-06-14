package routingkit

import (
	"errors"
	"fmt"
	"os"
	"runtime"

	"github.com/nextmv-io/go-routingkit/routingkit/internal/routingkit"
	rk "github.com/nextmv-io/go-routingkit/routingkit/internal/routingkit"
)

func finalizer(client *rk.Client) {
	routingkit.DeleteClient(*client)
}

// MaxDistance represents the maximum possible route distance.
var MaxDistance uint32

type TravelProfile routingkit.GoRoutingKitTravel_profile

var CarTravelProfile, BikeTravelProfile, PedestrianTravelProfile TravelProfile

func init() {
	MaxDistance = uint32(routingkit.GetMax_distance())
	CarTravelProfile = TravelProfile(routingkit.Car)
	BikeTravelProfile = TravelProfile(routingkit.Bike)
	PedestrianTravelProfile = TravelProfile(routingkit.Pedestrian)
}

// NewDistanceClient initializes a DistanceClient using the provided .osm.pbf file and
// .ch file. The .ch file will be created if it does not already exist.
func NewDistanceClient(mapFile, chFile string, profile TravelProfile) (DistanceClient, error) {
	if _, err := os.Stat(mapFile); os.IsNotExist(err) {
		return DistanceClient{}, errors.New(fmt.Sprintf("could not find map file at %v", mapFile))
	}

	concurrentQueries := runtime.GOMAXPROCS(0)
	c := rk.NewClient(concurrentQueries, mapFile, chFile, routingkit.GoRoutingKitTravel_profile(profile))
	runtime.SetFinalizer(&c, finalizer)

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
	client     rk.Client
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
	defer rk.DeleteQueryResponse(resp)
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
	defer rk.DeleteQueryResponse(resp)

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
	defer rk.DeletePoint(res)
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

	s := rk.NewPoint()
	defer routingkit.DeletePoint(s)
	s.SetLon(float32(source[0]))
	s.SetLat(float32(source[1]))

	targetsVector := rk.NewPointVector(int64(len(targets)))
	defer routingkit.DeletePointVector(targetsVector)

	for i := 0; i < len(targets); i++ {
		t := rk.NewPoint()
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
