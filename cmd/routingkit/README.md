# go-routingkit standalone

`routingkit` allows the usage of go-routingkit as a standalone executable.

## Example

1. Download an _.osm.pbf_ file for Maryland from [Geofabrik](http://download.geofabrik.de/):

    ```bash
    wget -P "data/" "http://download.geofabrik.de/north-america/us/maryland-latest.osm.pbf"
    ```

1. Execute the sample request from `data/maryland-points.json`

    ```bash
    go run . \
        -map "data/maryland-latest.osm.pbf" \
        -ch "data/maryland-latest.ch" \
        -input "data/maryland-points.json" \
        -measure "distance"
    ```

Note that [contraction hierarchies][ch] are built and saved to the _.ch_ file,
if the file is not yet present. This process takes a while. Subsequent calls
(with a _.ch_ file present) will be faster.

## Usage

Usage of _standalone go-routingkit_ is described below.

### Build

Build a binary by simply invoking

```bash
go build
```

### Install

Install binary to `PATH`. Requires `$(go env GOPATH)/bin` to be included in
`PATH`. After successful installation `routingkit` can be invoked as a command.

```bash
go install github.com/nextmv-io/go-routingkit/cmd/routingkit@latest
```

### Arguments

```go
Usage of ./routingkit:
  -ch string
        path to ch file (default "data/map.ch")
  -input string
        path to input file. default is stdin.
  -map string
        path to map file (default "data/map.osm.pbf")
  -measure string
        distance|traveltime (default "distance")
  -output string
        path to output file. default is stdout.
  -profile string
        car|bike|pedestrian - bike and pedestrian only work with measure=distance
        (default "car")
```

### Input / output

Find a sample `--input` below. Each request is given as a tuple of two locations
defining the `from` and `to` part of the trip. The position of the location is
given as longitude (`lon`) and latitude (`lat`).

```json
{
    "tuples": [
        {
            "from": { "lon": -76.733, "lat": 38.887 },
            "to": { "lon": -77.095, "lat": 38.981 }
        },
        {
            "from": { "lon": -76.888, "lat": 38.951 },
            "to": { "lon": -76.868, "lat": 38.938 }
        },
        {
            "from": { "lon": -76.698, "lat": 39.065 },
            "to": { "lon": -77.029, "lat": 39.022 }
        },
        {
            "from": { "lon": -76.888, "lat": 39.207 },
            "to": { "lon": -76.485, "lat": 39.364 }
        }
    ]
}
```

Here is an excerpt of a sample output. There is one 'trip' per tuple in the
input given as an array in input order. Each trip contains the `cost` (distance
or drive-time) and the `waypoints`, which define the shape of the route. All
positions are again given as longitude (`lon`) & latitude (`lat`).

```jsonc
{
    "trips": [
        {
            "waypoints": [
                { "lon": -76.7316, "lat": 38.887352 },
                { "lon": -76.731415, "lat": 38.88516 },
                { "lon": -76.73108, "lat": 38.88498 },
                // ...
                { "lon": -77.09511, "lat": 38.98097 }
            ],
            "cost": 42480
        },
        {
            "waypoints": [
                { "lon": -76.888374, "lat": 38.95079 },
                { "lon": -76.88743, "lat": 38.950684 },
                { "lon": -76.88705, "lat": 38.950535 },
                // ...
                { "lon": -76.867935, "lat": 38.937477 }
            ],
            "cost": 5200
        },
        {
            "waypoints": [
                { "lon": -76.69816, "lat": 39.064713 },
                { "lon": -76.697945, "lat": 39.065346 },
                { "lon": -76.69793, "lat": 39.065422 },
                // ...
                { "lon": -77.028854, "lat": 39.022038 }
            ],
            "cost": 38806
        },
        {
            "waypoints": [
                { "lon": -76.88802, "lat": 39.206516 },
                { "lon": -76.88725, "lat": 39.20655 },
                { "lon": -76.88608, "lat": 39.2066 },
                // ...
                { "lon": -76.48349, "lat": 39.363342 }
            ],
            "cost": 45401
        }
    ]
}
```

[ch]: https://en.wikipedia.org/wiki/Contraction_hierarchies
