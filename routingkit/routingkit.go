package routingkit

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"sync/atomic"

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

	return client{
		client: c,
		sem:    NewSemaphore(100),
	}, nil
}

type client struct {
	client  rk.Client
	sem     Semaphore
	counter int64
}

func (c client) Distance(from []float64, to []float64) float64 {
	c.sem.Lock()
	defer func() {
		atomic.AddInt64(&c.counter, -1)
		c.sem.Unlock()
	}()
	counter := atomic.AddInt64(&c.counter, 1)
	return float64(c.client.Distance(
		int(counter-1),
		float32(from[0]),
		float32(from[1]),
		float32(to[0]),
		float32(to[1]),
	))
}

func (c client) Threaded(from []float64, to []float64) float64 {
	c.sem.Lock()
	defer func() {
		atomic.AddInt64(&c.counter, -1)
		c.sem.Unlock()
	}()
	counter := atomic.AddInt64(&c.counter, 1)
	return float64(c.client.Threaded(
		int(counter-1),
		float32(from[0]),
		float32(from[1]),
		float32(to[0]),
		float32(to[1]),
	))
}

type Semaphore chan struct {
}

func NewSemaphore(size int) Semaphore {
	return make(Semaphore, size)
}

func (s Semaphore) Lock() {
	// Writes will only succeed if there is room in s.
	s <- struct{}{}
}

func (s Semaphore) Unlock() {
	// Make room for other users of the semaphore.
	<-s
}
