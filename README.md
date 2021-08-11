# go-routingkit

Go-routingkit is a Go wrapper around the [RoutingKit][rk] C++ library. It
answers queries about the shortest path between points found within a road
network.

## Installing and Building

```go
go get -u github.com/nextmv-io/go-routingkit
```

Go-routingkit is currently supported on Linux, MacOS (both Intel and Apple Silicon).

## Deployment

As go-routingkit uses cgo, any programs that use it should ensure that at
runtime they can dynamically link against a C standard library version that is
compatible with the version the program was built with. If using glibc, version
2.26 or higher is required.

The default AWS Lambda image does not meet the version requirements for glibc.
However, the amazonlinux 2 image provides a more recent version of glibc that
is compatible with routingkit. To use this image, simply use the dropdown under
`Runtime` to select `Provide your own bootstrap on Amazon Linux 2` when creating
your lambda function. If creating your Lambda function with SAM, enter
`provided.al2` under the `Runtime` setting.

## Usage

### Initialization

All queries require first constructing a Client. Go-routingkit provides two
clients, `DistanceClient` and `TravelTimeClient`, which measure and minimize
route lengths using distance and total travel time, respectively.

```go
distanceCli, err := routingkit.NewDistanceClient("philadelphia.osm.pbf", routingkit.CarTravelProfile)
timeCli, err := routingkit.NewTravelTimeClient("philadelphia.osm.pbf")
```

The `DistanceClient` constructs different routes depending on the mode of
transportation, which is specified in its final argument. The `TravelTimeClient`
only applies to cars, since they travel at substantially different speeds along
different paths in the road network.

The constructors for DistanceClient and TravelTimeClient share their first two
arguments: these are the path to an .osm.pbf file and the path to a .ch file
(which may or may not already exist). The _.osm.pbf_ file contains geographic
data about the road network. You can download files that contain the region
you're interested in from [Geofabrik](http://download.geofabrik.de/), or you can
use another source such as the [Overpass API](http://overpass-api.de/) to
generate a file describing a custom bounding box.

The contraction hierarchy (.ch) file contains indices that allow routing queries
to be executed more efficiently. If there is no file at the .ch filepath you
pass to `NewClient`, routingkit will build this file for you. Reusing this file
between initializations will lead to faster start times. Contraction hierarchies
are specific to the travel profile and the route measurement, and should not be
reused between different types of clients.

### Distance and Travel Time Queries

`routingkit.DistanceClient` and `routingkit.TravelTimeClient` allow a few
different types of queries for the shortest paths between points. Points are
represented as `[]float32`s where the first element is the longitude and the
second is the latitude. Distances are represented as `uint32`s, representing
distance in meters for `DistanceClient` and in milliseconds for
`TravelTimeClient`.

The simplest query finds the distance between two points:

```go
distance := distanceCli.Distance([]float32{-75.1785585,39.9532349}, []float32{-75.1650723,39.9515036})
time := timeCli.TravelTime([]float32{-75.1785585,39.9532349}, []float32{-75.1650723,39.9515036})
```

The `Route` query returns not only the distance between two points, but also a
series of waypoints along the path from the starting point to the destination.

```go
distance, waypoints := distanceCli.Route([]float32{-75.1785585,39.9532349}, []float32{-75.1650723,39.9515036})
time, waypoints := timeCli.Route([]float32{-75.1785585,39.9532349}, []float32{-75.1650723,39.9515036})
```

The `Distances` and `TravelTimes` methods perform a vectorized query for
distances or travel times from a source to multiple destinations.

```go
distances := distanceCli.Distances(
    []float32{-75.1785585,39.9532349},
    [][]float32{{-75.1650723,39.9515036}, {-75.1524708,39.9496144}},
)
times := timeCli.TravelTimes(
    []float32{-75.1785585,39.9532349},
    [][]float32{{-75.1650723,39.9515036}, {-75.1524708,39.9496144}},
)
```

And `Matrix` creates a matrix containing distances (or travel times) from
multiple source points to multiple destination points.

```go
matrix := distanceCli.Matrix(
    [][]float32{{-75.1785585,39.9532349}, {-75.2135608,39.9610131}},
    [][]float32{{-75.1650723,39.9515036}, {-75.1524708,39.9496144}},
)
matrix := timeCli.Matrix(
    [][]float32{{-75.1785585,39.9532349}, {-75.2135608,39.9610131}},
    [][]float32{{-75.1650723,39.9515036}, {-75.1524708,39.9496144}},
)
```

### Snap Radius

The clients can find routes between points that are located within road
networks, but it's often useful to query for points that do not fall exactly on
a road, automatically snapping the point to the nearest location on a road. The
client's snap radius defaults to 1000 meters and determines the maximum distance
a point will be snapped to on the road grid. It can be set with:

```go
cli.SetSnapRadius(100)
```

After being adjusted, this snap radius will be applied to any query done by the
client.

[rk]: https://github.com/RoutingKit/RoutingKit
