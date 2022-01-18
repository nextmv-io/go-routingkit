package routingkit

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
)

type SpeedMapper func(wayId int, tagMap map[string]string) int

var osmTagWithCountryCode = regexp.MustCompile(`^(\w{2}):(.*)$`)
var maxSpeedAndUnits = regexp.MustCompile(`^([0-9][\.0-9]*?)(?:[ ]?(km/h|kmh|kph|mph|knots))?$`)

func parseMaxspeed(maxspeed string) float64 {
	switch maxspeed {
	case "at:rural":
		return 100
	case "at:trunk":
		return 100
	case "be:motorway":
		return 120
	case "be-bru:rural":
		return 70
	case "be-bru:urban":
		return 30
	case "be-vlg:rural":
		return 70
	case "by:urban":
		return 60
	case "by:motorway":
		return 110
	case "ch:rural":
		return 80
	case "ch:trunk":
		return 100
	case "ch:motorway":
		return 120
	case "cz:trunk":
		return 0
	case "cz:motorway":
		return 0
	case "de:living_street":
		return 7
	case "de:rural":
		return 100
	case "de:motorway":
		return 0
	case "dk:rural":
		return 80
	case "fr:rural":
		return 80
	case "gb:nsl_single":
		return (60 * 1609) / 1000
	case "gb:nsl_dual":
		return (70 * 1609) / 1000
	case "gb:motorway":
		return (70 * 1609) / 1000
	case "nl:rural":
		return 80
	case "nl:trunk":
		return 100
	case "no:rural":
		return 80
	case "no:motorway":
		return 110
	case "pl:rural":
		return 100
	case "pl:trunk":
		return 120
	case "pl:motorway":
		return 140
	case "ro:trunk":
		return 100
	case "ru:living_street":
		return 20
	case "ru:urban":
		return 60
	case "ru:motorway":
		return 110
	case "uk:nsl_single":
		return (60 * 1609) / 1000
	case "uk:nsl_dual":
		return (70 * 1609) / 1000
	case "uk:motorway":
		return (70 * 1609) / 1000
	case "za:urban":
		return 60
	case "za:rural":
		return 100
	case "none":
		return 140
	}

	withoutCountryCode := osmTagWithCountryCode.ReplaceAllString(maxspeed, "${1}")
	switch withoutCountryCode {
	case "urban":
		return 50
	case "rural":
		return 90
	case "trunk":
		return 110
	case "motorway":
		return 130
	}

	if speed, ok := ParseOSMSpeedToKM(maxspeed); ok {
		return speed
	}

	// TODO: logging... we don't have a strategy for how a consumer should inject a logger

	return 0
}

func ParseOSMSpeedToKM(str string) (float64, bool) {
	speedUnitsMatch := maxSpeedAndUnits.FindStringSubmatch(str)
	if len(speedUnitsMatch) != 3 {
		return 0, false
	}
	speedStr, units := speedUnitsMatch[1], speedUnitsMatch[2]
	speed, err := strconv.ParseFloat(speedStr, 64)
	if err != nil {
		// This should not be possible due to the contruction of the regexp
		panic(fmt.Errorf("extracted an invalid integer from maxspeed tag %s: %v", str, err))
	}
	if units == "" || units == "km/h" || units == "kmh" || units == "kph" {
		return speed, true
	}
	if units == "mph" {
		return speed * 1609 / 1000, true
	}
	if units == "knots" {
		return speed * 1852 / 1000, true
	}
	// TODO: logging... we don't have a strategy for how a consumer should inject a logger
	return speed, true
}

// BikeSpeedMapper sets the speed for bikes according to the map tag
func BikeSpeedMapper(_ int, tags map[string]string) int {
	defaultSpeed := 15
	walkingSpeed := 4
	if tags["bridge"] == "movable" {
		return 5
	}
	if tags["route"] == "ferry" {
		return 5
	}
	if tags["public_transport"] == "platform" {
		return walkingSpeed

	}
	switch tags["railway"] {
	case "platform":
		return walkingSpeed
	case "train":
		return 10
	case "railway":
		return 10
	case "subway":
		return 10
	case "light_rail":
		return 10
	case "monorail":
		return 10
	case "tram":
		return 10
	}
	amenity := tags["amenity"]
	if amenity == "parking" || amenity == "parking_entrance" {
		return 10
	}
	switch tags["highway"] {
	case "track":
		return 12
	case "path":
		return 12
	}
	return defaultSpeed
}

// CarSpeedMapper sets the speed for cars according to allowed speed, map tag and surface
func CarSpeedMapper(_ int, tags map[string]string) int {
	route := tags["route"]
	if route == "ferry" {
		return 5
	}
	if route == "shuttle_train" {
		return 10
	}
	if tags["bridge"] == "movable" {
		if tags["capacity:car"] == "0" {
			return 0
		}
		return 5
	}
	speed := map[string]int{
		"motorway":       90,
		"motorway_link":  45,
		"trunk":          85,
		"trunk_link":     40,
		"primary":        65,
		"primary_link":   30,
		"secondary":      55,
		"secondary_link": 25,
		"tertiary":       40,
		"tertiary_link":  20,
		"unclassified":   25,
		"residential":    25,
		"living_street":  10,
		"service":        15,
	}[tags["highway"]]
	if speed == 0 {
		speed = 10
	}

	{
		var speedStr string
		if speed, ok := tags["maxspeed:advisory"]; ok {
			speedStr = speed
		} else if speed, ok := tags["maxspeed"]; ok {
			speedStr = speed
		} else if speed, ok := tags["source:maxspeed"]; ok {
			speedStr = speed
		} else if speed, ok := tags["maxspeed:type"]; ok {
			speedStr = speed
		}
		if speedStr != "" {
			speed = int(math.Ceil((parseMaxspeed(strings.TrimLeft(speedStr, " ")))))
		}
	}

	surface := map[string]int{
		"cement":        80,
		"compacted":     80,
		"fine_gravel":   80,
		"paving_stones": 60,
		"metal":         60,
		"bricks":        60,
		"grass":         40,
		"wood":          40,
		"sett":          40,
		"grass_paver":   40,
		"gravel":        40,
		"unpaved":       40,
		"ground":        40,
		"dirt":          40,
		"pebblestone":   40,
		"tartan":        40,
		"cobblestone":   30,
		"clay":          30,
		"earth":         20,
		"stone":         20,
		"rocky":         20,
		"sand":          20,
		"mud":           10,
	}[tags["surface"]]
	if surface != 0 && surface < speed {
		speed = surface
	}
	trackType := map[string]int{
		"grade1": 60,
		"grade2": 40,
		"grade3": 30,
		"grade4": 25,
		"grade5": 20,
	}[tags["tracktype"]]
	if trackType != 0 && trackType < speed {
		speed = trackType
	}
	if smoothness, ok := map[string]int{
		"intermediate":  80,
		"bad":           40,
		"very_bad":      20,
		"horrible":      10,
		"very_horrible": 5,
		"impassable":    0,
	}[tags["smoothness"]]; ok && smoothness < speed {
		speed = smoothness
	}

	return speed
}

// PedestrianSpeedMapper sets to 5km/h and reduces the speed according to the
// surface of the underlying way
func PedestrianSpeedMapper(_ int, tags map[string]string) int {
	speed := 5.0
	multiplier := map[string]float64{
		"fine_gravel": 0.75,
		"gravel":      0.75,
		"pebblestone": 0.75,
		"mud":         0.5,
		"sand":        0.5,
	}[tags["surface"]]
	if multiplier != 0 {
		return int(speed * multiplier)
	}
	return int(speed)
}

// MaxSpeedMapper caps the allowed speed at the given value
func MaxSpeedMapper(maxSpeed int) SpeedMapper {
	return func(_ int, tagMap map[string]string) int {
		speed := CarSpeedMapper(0, tagMap)
		if speed > maxSpeed {
			speed = maxSpeed
		}
		return speed
	}
}
