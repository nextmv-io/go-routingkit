# routingkit binary

`routingkit` allows the usage of go-routingkit as a standalone executable.

Run `./example.sh` for a quickstart. Prior to executing a sample call the script
will download a matching osm region file.

## Usage

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

Find a sample input below. Each request is given as a tuple of two locations
defining the `from` and `to` part of the trip. The coordinates of the location
are given in `[lon, lat]` order.

```json
{
    "tuples": [
        {
            "from":[-76.73266990511583, 38.88656521737339],
            "to":[-77.0950991, 38.98064919999999]
        },
        {
            "from":[-76.8233293, 39.2142304],
            "to":[-76.738197, 39.392148]
        },
        {
            "from":[-76.6977902, 39.0648349],
            "to":[-77.029122, 39.0220273]
        },
        {
            "from":[-76.8875218, 39.2067232],
            "to":[-76.4852039, 39.3641577]
        }
    ]
}
```

Here is an excerpt of a sample output. There is one 'trip' per tuple in the input
given as an array in input order. Each trip contains the `cost` (distance or
drive-time) and the `waypoints`, which define the shape of the route. All
coordinates are given in `[lon, lat]` order.

```json
{
    "trips": [
        {
            "waypoints": [
                [-76.7316, 38.887352],
                [-76.731415, 38.88516],
                [-76.73108, 38.88498],
                ...
                [-77.09511, 38.98097]
            ],
            "cost": 42480
        },
        {
            "waypoints": [
                [-76.888374, 38.95079],
                [-76.88743, 38.950684],
                [-76.88705, 38.950535],
                ...
                [-76.867935, 38.937477]
            ],
            "cost": 5200
        },
        {
            "waypoints": [
                [-76.69816, 39.064713],
                [-76.697945, 39.065346],
                [-76.69793, 39.065422],
                ...
                [-77.028854, 39.022038]
            ],
            "cost": 38806
        },
        {
            "waypoints": [
                [-76.88802, 39.206516],
                [-76.88725, 39.20655],
                [-76.88608, 39.2066],
                ...
                [-76.48349, 39.363342]
            ],
            "cost": 45401
        }
    ]
}

```
