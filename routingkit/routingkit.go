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
