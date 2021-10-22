package routingkit

import "math"

type Point struct {
	Lon float64 `json:"lon"`
	Lat float64 `json:"lat"`
}

func polylineSimilarity(first, second []Point) float64 {
	return 0
}

func haversine(p1, p2 Point) float64 {
	x1, y1 := degToRad(p1.Lon), degToRad(p1.Lat)
	x2, y2 := degToRad(p2.Lon), degToRad(p2.Lat)

	dx := x1 - x2
	dy := y1 - y2

	a := math.Pow(math.Sin(dy/2), 2) +
		math.Cos(y1)*math.Cos(y2)*
			math.Pow(math.Sin(dx/2), 2)

	return 2 * radius * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
}

const radius = 6371 * 1000 // radius of the earth in meters

// degToRad converts a degree value to radians.
func degToRad(d float64) float64 {
	return d * math.Pi / 180.0
}
