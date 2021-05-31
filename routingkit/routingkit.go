package routingkit

import (
	"errors"
	"fmt"
	"os"
	"runtime"

	"github.com/nextmv-io/go-routingkit/routingkit/internal/routingkit"
	rk "github.com/nextmv-io/go-routingkit/routingkit/internal/routingkit"
)

type Client interface {
	Distance([]float64, []float64) float64
	Threaded([]float64, []float64) float64
	Table([]float64, [][]float64) []float64
	Average() float64
}

func finalizer(client *rk.Client) {
	routingkit.DeleteClient(*client)
}

func Wrapper(mapFile, chFile string) (Client, error) {
	if _, err := os.Stat(mapFile); os.IsNotExist(err) {
		return nil, errors.New(fmt.Sprintf("could not find map file at %v", mapFile))
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

	return client{
		client:  c,
		channel: channel,
	}, nil
}

type client struct {
	client  rk.Client
	channel chan int
}

func (c client) Average() float64 {
	ints := make([]int, 10)
	vector := rk.NewIntVector(int64(len(ints)))
	defer rk.DeleteIntVector(vector)
	return c.client.Average(vector)
}

func (c client) Distance(from []float64, to []float64) float64 {
	counter := <-c.channel
	defer func() {
		c.channel <- counter
	}()
	return float64(c.client.Distance(
		int(counter),
		float32(from[0]),
		float32(from[1]),
		float32(to[0]),
		float32(to[1]),
	))
}

func (c client) Threaded(from []float64, to []float64) float64 {
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

func (c client) Table(source []float64, targets [][]float64) []float64 {
	counter := <-c.channel
	defer func() {
		c.channel <- counter
	}()
	targetsVector := rk.NewPointVector(int64(len(targets)))
	s := rk.NewPoint()
	s.SetLon(float32(source[0]))
	s.SetLat(float32(source[1]))

	for i := 0; i < len(targets); i++ {
		t := rk.NewPoint()
		t.SetLon(float32(targets[i][0]))
		t.SetLat(float32(targets[i][1]))
		targetsVector.Set(i, t)
	}

	matrix := c.client.Table(counter, s, targetsVector)
	numRows := matrix.Size()
	rows := make([]float64, numRows)
	for i := 0; i < int(numRows); i++ {
		col := matrix.Get(i)
		rows[i] = float64(col)
	}

	// defer func() {
	// 	routingkit.DeletePointVector(sourcesVector)
	// 	routingkit.DeletePointVector(targetsVector)
	// 	routingkit.DeleteMatrix(matrix)
	// }()

	return rows
}
