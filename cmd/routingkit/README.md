# routingkit binary

`routingkit` allows the usage of go-routingkit as a standalone executable.

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

Here is a sample input:

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
