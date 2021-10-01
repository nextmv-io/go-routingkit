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
var maxSpeedAndUnits = regexp.MustCompile(`^([0-9][\.0-9]+?)(?:[ ]?(km/h|kmh|kph|mph|knots))?$`)

func parseMaxspeed(maxspeed string) int {
	if maxspeed == "signals" || maxspeed == "variable" {
		return math.MaxInt64
	}
	if maxspeed == "none" || maxspeed == "unlimited" {
		return 130
	}
	withoutCountryCode := osmTagWithCountryCode.ReplaceAllString(maxspeed, "${1}")
	if withoutCountryCode == "walk" || maxspeed == "foot" {
		return 5
	}
	if withoutCountryCode == "urban" {
		return 40
	}
	if withoutCountryCode == "living_street" {
		return 10
	}
	if maxspeed == "rural" || maxspeed == "de:rural" || maxspeed == "at:rural" || maxspeed == "ro:rural" {
		return 100
	}
	if maxspeed == "ru:rural" || maxspeed == "ua:rural" {
		return 90
	}
	if maxspeed == "ru:motorway" {
		return 110
	}
	if maxspeed == "at:motorway" || maxspeed == "ro:motorway" {
		return 130
	}
	if maxspeed == "national" {
		return 100
	}
	if maxspeed == "ro:trunk" {
		return 100
	}
	if maxspeed == "dk:rural" || maxspeed == "ch:rural" || maxspeed == "fr:rural" {
		return 80
	}
	if maxspeed == "it:rural" || maxspeed == "hu:rural" {
		return 90
	}
	if maxspeed == "de:zone:30" || maxspeed == "de:zone30" {
		return 30
	}

	speedUnitsMatch := maxSpeedAndUnits.FindStringSubmatch(maxspeed)
	if len(speedUnitsMatch) == 3 {
		speedStr, units := speedUnitsMatch[1], speedUnitsMatch[2]
		speed, err := strconv.Atoi(speedStr)
		if err != nil {
			// This should not be possible due to the contruction of the regexp
			panic(fmt.Errorf("extracted an invalid integer from maxspeed tag %s: %v", maxspeed, err))
		}
		if units == "" || units == "km/h" || units == "kmh" || units == "kph" {
			return speed
		}
		if units == "mph" {
			return speed * 1609 / 1000
		}
		if units == "knots" {
			return speed * 1852 / 1000
		}
		// TODO: logging... we don't have a strategy for how a consumer should inject a logger
		return speed
	}
	// TODO: logging... we don't have a strategy for how a consumer should inject a logger

	return math.MaxInt64
}

func carSpeedMapper(_ int, tagMap map[string]string) int {
	maxspeed, maxspeedOk := tagMap["maxspeed"]
	if maxspeedOk && maxspeed != "unposted" {
		entries := strings.Split(maxspeed, ";")
		minSpeed := math.MaxInt64
		for _, entry := range entries {
			speed := parseMaxspeed(strings.TrimLeft(entry, " "))
			if speed < minSpeed {
				minSpeed = speed
			}
		}

		if minSpeed == math.MaxInt64 {
			return 1
		}
		return minSpeed
	}
	highway, highwayOk := tagMap["highway"]
	if highwayOk {
		switch highway {
		case "motorway":
			return 90
		case "motorway_link":
			return 45
		case "trunk":
			return 85
		case "trunk_link":
			return 40
		case "primary":
			return 65
		case "primary_link":
			return 30
		case "secondary":
			return 55
		case "secondary_link":
			return 25
		case "tertiary":
			return 40
		case "tertiary_link":
			return 20
		case "unclassified":
			return 25
		case "residential":
			return 25
		case "living_street":
			return 10
		case "service":
			return 8
		case "track":
			return 8
		case "ferry":
			return 5
		}
	}

	if _, ok := tagMap["junction"]; ok {
		return 20
	}

	if val, ok := tagMap["route"]; ok && val == "ferry" {
		return 5
	}
	if _, ok := tagMap["ferry"]; ok {
		return 5
	}

	return 50
}

func pedestrianSpeedMapper(_ int, tagMap map[string]string) int {
	return 5
}

func maxSpeedMapper(maxSpeed int) SpeedMapper {
	return func(_ int, tagMap map[string]string) int {
		speed := carSpeedMapper(0, tagMap)
		if speed > maxSpeed {
			speed = maxSpeed
		}
		return speed
	}
}
