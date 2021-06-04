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

func NewClient(mapFile, chFile string) (Client, error) {
	if _, err := os.Stat(mapFile); os.IsNotExist(err) {
		return Client{}, errors.New(fmt.Sprintf("could not find map file at %v", mapFile))
	}

	c := rk.NewClient()
	runtime.SetFinalizer(&c, finalizer)

	if _, err := os.Stat(chFile); os.IsNotExist(err) {
		c.Build_ch(mapFile, chFile)
	} else {
		c.Load(mapFile, chFile)
	}

	count := 100
	channel := make(chan int, count)
	for i := 0; i < count; i++ {
		channel <- i
	}

	return Client{
		client:  c,
		channel: channel,
	}, nil
}

type Client struct {
	client  rk.Client
	channel chan int
}

func (c Client) Average() float64 {
	ints := make([]int, 10)
	vector := rk.NewIntVector(int64(len(ints)))
	defer rk.DeleteIntVector(vector)
	return c.client.Average(vector)
}

func (c Client) Query(radius float64, from []float64, to []float64) (float64, [][]float64) {
	counter := <-c.channel
	defer func() {
		c.channel <- counter
	}()
	resp := c.client.Queryrequest(
		int(counter),
		float32(radius),
		float32(from[0]),
		float32(from[1]),
		float32(to[0]),
		float32(to[1]),
	)
	wp := resp.GetWaypoints()
	waypoints := make([][]float64, wp.Size())
	for i := 0; i < len(waypoints); i++ {
		p := wp.Get(i)
		waypoints[i] = []float64{float64(p.GetLon()), float64(p.GetLat())}
	}

	return float64(resp.GetDistance()), waypoints
}

func (c Client) Threaded(from []float64, to []float64) float64 {
	counter := <-c.channel
	defer func() {
		c.channel <- counter
	}()
	return float64(c.client.Threaded(
		int(counter),
		float32(from[0]),
		float32(from[1]),
		float32(to[0]),
		float32(to[1]),
	))
}

func (c Client) Table(sources [][]float64, targets [][]float64) [][]float64 {
	counter := <-c.channel
	defer func() {
		c.channel <- counter
	}()

	sourcesVector := rk.NewPointVector(int64(len(sources)))
	defer routingkit.DeletePointVector(sourcesVector)
	for i := 0; i < len(sources); i++ {
		s := rk.NewPoint()
		s.SetLon(float32(sources[i][0]))
		s.SetLat(float32(sources[i][1]))
		sourcesVector.Set(i, s)
	}

	targetsVector := rk.NewPointVector(int64(len(targets)))
	defer routingkit.DeletePointVector(targetsVector)
	for i := 0; i < len(targets); i++ {
		t := rk.NewPoint()
		t.SetLon(float32(targets[i][0]))
		t.SetLat(float32(targets[i][1]))
		targetsVector.Set(i, t)
	}

	matrix := c.client.Table(counter, sourcesVector, targetsVector)
	defer routingkit.DeleteFloatVector(matrix)
	numRows := len(sources)
	rows := make([][]float64, len(sources))
	pos := 0
	for i := 0; i < int(numRows); i++ {
		numCols := len(targets)
		cols := make([]float64, numCols)
		for j := 0; j < numCols; j++ {
			cols[j] = float64(matrix.Get(pos))
			pos++
		}
		rows[i] = cols
	}

	return rows
}
