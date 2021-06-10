# go-routingkit

go-routingkit is a Go wrapper around the [RoutingKit](https://github.com/RoutingKit/RoutingKit) C++ library. It answers queries about the shortest path between points found within a road network.

## Install

```go
go get -u github.com/nextmv-io/go-routingkit
```

go-routingkit is currently only supported on Linux.

## Usage

### Initialization

All queries require first constructing a Client.

```go
cli, err := routingkit.NewClient("philadelphia.osm.pbf", "philadelphia.ch")
```

The Client constructor takes two arguments: the path to an .osm.pbf file and path to a .ch file (which may or may not already exist). The _.osm.pbf_ file contains geographic data about the road network. You can download files that contain the region you're interested in from [Geofabrik](http://download.geofabrik.de/), or you can use another source such as the [Overpass API](http://overpass-api.de/) to generate a file describing a custom bounding box.

The contraction hierarchy file contains indices that allow routing queries to be executed more efficiently. If there is no file at the .ch filepath you pass to `NewClient`, routingkit will build this file for you. Reusing this file between initializations will lead to faster start times.

### Distance Queries

`routingkit.Client` allows a few different types of queries for the shortest paths between points. Points are represented as `[]float32`s where the first element is the longitude and the second is the latitude. Distances are represented as `uint32`s, representing distance in meters.

The simplest query finds the distance between two points:

```go
distance := cli.Distance([]float32{-75.1785585,39.9532349}, []float32{-75.1650723,39.9515036})
```

The `Route` query returns not only the distance between two points, but also a series of waypoints representing turn-by-turn directions from the starting point to the destination.

```go
distance, waypoints := cli.Route([]float32{-75.1785585,39.9532349}, []float32{-75.1650723,39.9515036})
```

The `Distances` method performs a vectorized query for distances from a source to multiple destinations.

```go
distances := cli.Distances([]float32{-75.1785585,39.9532349}, [][]float32{{-75.1650723,39.9515036}, {-75.1524708,39.9496144}},)
```

And `Matrix` creates a matrix containing distances from multiple source points to multiple destination points.

```go
matrix := cli.Matrix([][]float32{{-75.1785585,39.9532349}, {-75.2135608,39.9610131}}, [][]float32{{-75.1650723,39.9515036}, {-75.1524708,39.9496144}})
```

### Snap Radius

The `routingkit.Client` can find routes between points that are located within road networks, but it's often useful to query for points that do not fall exactly on a road, automatically snapping the point to the nearest location on a road. The client's snap radius defaults to 1000 meters and determines the maximum distance a point will be snapped to on the road grid. It can be set with:

```go
cli.SetSnapRadius(100)
```

After being adjusted, this snap radius will be applied to any query done by the client.
