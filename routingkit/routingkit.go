package routingkit

import (
	"errors"
	"fmt"
	"os"
	"runtime"

	"github.com/nextmv-io/go-routingkit/routingkit/internal/routingkit"
	rk "github.com/nextmv-io/go-routingkit/routingkit/internal/routingkit"
)

type Response struct {
	WayPoints [][]float64
	Distance  float64
}

type Client interface {
	Query(float64, []float64, []float64) Response
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
	}, nil
}

type client struct {
	client rk.Client
}

func (c client) Query(radius float64, from []float64, to []float64) Response {
	resp := c.client.Queryrequest(
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

	reponse := Response{
		Distance:  float64(resp.GetDistance()),
		WayPoints: waypoints,
	}
	return reponse
}
