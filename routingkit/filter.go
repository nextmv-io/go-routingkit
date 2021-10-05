package routingkit

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
)

type TagMapFilter func(wayId int, tagMap map[string]string) bool

func carTagMapFilter(_ int, tagMap map[string]string) bool {
	if _, ok := tagMap["junction"]; ok {
		return true
	}
	if val, ok := tagMap["route"]; ok && val == "ferry" {
		return true
	}
	if val, ok := tagMap["ferry"]; ok && val == "yes" {
		return true
	}
	highway, ok := tagMap["highway"]
	if !ok {
		return false
	}
	if val, ok := tagMap["motorcar"]; ok && val == "no" {
		return false
	}
	if val, ok := tagMap["motor_vehicle"]; ok && val == "no" {
		return false
	}

	if val, ok := tagMap["access"]; ok {
		if !(val == "yes" || val == "permissive" || val == "delivery" || val == "designated" || val == "destination") {
			return false
		}
	}

	if highway == "motorway" ||
		highway == "trunk" ||
		highway == "primary" ||
		highway == "secondary" ||
		highway == "tertiary" ||
		highway == "unclassified" ||
		highway == "residential" ||
		highway == "service" ||
		highway == "motorway_link" ||
		highway == "trunk_link" ||
		highway == "primary_link" ||
		highway == "secondary_link" ||
		highway == "tertiary_link" ||
		highway == "motorway_junction" ||
		highway == "living_street" ||
		highway == "track" ||
		highway == "ferry" {
		return true
	}

	if highway == "bicycle_road" {
		if val, ok := tagMap["motorcar"]; ok && val == "yes" {
			return true
		}
		return false
	}

	if highway == "construction" ||
		highway == "path" ||
		highway == "footway" ||
		highway == "cycleway" ||
		highway == "bridleway" ||
		highway == "pedestrian" ||
		highway == "bus_guideway" ||
		highway == "raceway" ||
		highway == "escape" ||
		highway == "steps" ||
		highway == "proposed" ||
		highway == "conveying" {
		return false
	}

	if val, ok := tagMap["oneway"]; ok && val == "reversible" || val == "alternating" {
		return false
	}

	if _, ok := tagMap["maxspeed"]; ok {
		return true
	}

	return false
}

func bikeTagMapFilter(_ int, tagMap map[string]string) bool {
	if _, ok := tagMap["junction"]; ok {
		return true
	}
	if val, ok := tagMap["route"]; ok && val == "ferry" {
		return true
	}
	// TODO: I noticed this is different from cars, where the val is "yes" instead of "ferry".
	// This matches what RoutingKit does but I'd like to double check this
	if val, ok := tagMap["ferry"]; ok && val == "ferry" {
		return true
	}
	highway, ok := tagMap["highway"]
	if !ok {
		return false
	}
	// TODO: proposed highways aren't filtered out until later in the car profile,
	// which seems wrong...
	if highway == "proposed" {
		return false
	}

	if val, ok := tagMap["access"]; ok {
		if !(val == "yes" ||
			val == "permissive" ||
			val == "delivery" ||
			val == "designated" ||
			val == "destination" ||
			val == "agricultural" ||
			val == "forestry" ||
			val == "public") {
			return false
		}
	}

	if val, ok := tagMap["bicycle"]; ok && val == "no" || val == "use_sidepath" {
		return false
	}

	if _, ok := tagMap["cycleway"]; ok {
		return true
	}
	if _, ok := tagMap["cycleway:left"]; ok {
		return true
	}
	if _, ok := tagMap["cycleway:right"]; ok {
		return true
	}
	if _, ok := tagMap["cycleway:both"]; ok {
		return true
	}

	if highway == "secondary" ||
		highway == "tertiary" ||
		highway == "unclassified" ||
		highway == "residential" ||
		highway == "service" ||
		highway == "secondary_link" ||
		highway == "tertiary_link" ||
		highway == "living_street" ||
		highway == "track" ||
		highway == "bicycle_road" ||
		highway == "primary" ||
		highway == "primary_link" ||
		highway == "path" ||
		highway == "footway" ||
		highway == "cycleway" ||
		// TODO: from OSM docs it doesn't seem like bridleways universally permit biking
		highway == "bridleway" ||
		highway == "pedestrian" ||
		highway == "crossing" ||
		highway == "escape" ||
		highway == "steps" ||
		highway == "ferry" {
		return true
	}

	if highway == "motorway" ||
		highway == "motorway_link" ||
		highway == "motorway_junction" ||
		highway == "trunk" ||
		highway == "trunk_link" ||
		highway == "construction" ||
		highway == "bus_guideway" ||
		highway == "raceway" ||
		highway == "conveying" {
		return false
	}

	// TODO: curious about lack of handling for one-way streets

	return false
}

func pedestrianTagMapFilter(_ int, tagMap map[string]string) bool {
	if _, ok := tagMap["junction"]; ok {
		return true
	}
	if val, ok := tagMap["route"]; ok && val == "ferry" {
		return true
	}
	// TOOD: same question here as with bikes
	if val, ok := tagMap["ferry"]; ok && val == "ferry" {
		return true
	}

	publicTransport, ok := tagMap["public_transport"]
	if ok && (publicTransport == "stop_position" ||
		publicTransport == "platform" ||
		publicTransport == "stop_area" ||
		publicTransport == "station") {
		return true
	}

	railway, ok := tagMap["railway"]
	if ok && (railway == "halt" ||
		railway == "platform" ||
		railway == "subway_entrance" ||
		railway == "station" ||
		railway == "tram_stop") {
		return true
	}

	highway, ok := tagMap["highway"]
	if !ok {
		return false
	}

	if val, ok := tagMap["access"]; ok {
		if !(val == "yes" ||
			val == "permissive" ||
			val == "delivery" ||
			val == "designated" ||
			val == "destination" ||
			val == "agricultural" ||
			val == "forestry" ||
			val == "public") {
			return false
		}
	}

	if val, ok := tagMap["crossing"]; ok && val == "no" {
		return false
	}

	if highway == "secondary" ||
		highway == "tertiary" ||
		highway == "unclassified" ||
		highway == "residential" ||
		highway == "service" ||
		highway == "secondary_link" ||
		highway == "tertiary_link" ||
		highway == "living_street" ||
		highway == "track" ||
		highway == "bicycle_road" ||
		highway == "path" ||
		highway == "footway" ||
		highway == "cycleway" ||
		highway == "bridleway" ||
		highway == "pedestrian" ||
		highway == "escape" ||
		highway == "steps" ||
		highway == "crossing" ||
		highway == "escalator" ||
		highway == "elevator" ||
		highway == "platform" ||
		highway == "ferry" {
		return true
	}

	if highway == "motorway" ||
		highway == "motorway_link" ||
		highway == "motorway_junction" ||
		highway == "trunk" ||
		highway == "trunk_link" ||
		highway == "primary" ||
		highway == "primary_link" ||
		highway == "construction" ||
		highway == "bus_guideway" ||
		highway == "raceway" ||
		// TODO: again, strikes me as wrong that proposed isn't given higher precedence
		// but maybe there's a reason for this
		highway == "proposed" ||
		highway == "conveying" {
		return false
	}

	// TODO: curious about lack of handling for one-way streets

	return false
}

var imperialMeasure = regexp.MustCompile(`^(\d+)'(?:(\d+)")?$`)
var decimalMeasure = regexp.MustCompile(`^(\d+(.\d*)?)(?: m)?$`)

const inchesPerMeter = 0.0254
const feetPerMeters = 0.3048

func parseAsMeters(val string) (float64, error) {
	imperial := imperialMeasure.FindStringSubmatch(val)
	if len(imperial) == 3 {
		var total float64
		feet, err := strconv.Atoi(imperial[1])
		if err != nil {
			return 0, fmt.Errorf("invalid feet value %s in %s: %v", imperial[1], val, err)
		}
		total += float64(feet) * feetPerMeters
		if imperial[2] != "" {
			inches, err := strconv.Atoi(imperial[2])
			if err != nil {
				return 0, fmt.Errorf("invalid inch value %s in %s: %v", imperial[2], val, err)
			}
			total += float64(inches) * inchesPerMeter
		}
		return total, nil
	}
	decimal := decimalMeasure.FindStringSubmatch(val)
	if len(decimal) == 3 {
		f, err := strconv.ParseFloat(decimal[1], 64)
		if err != nil {
			return 0, fmt.Errorf("invalid meter value %s: %v", val, err)
		}
		return f, nil
	}
	return 0, fmt.Errorf("could not parse %s as meter value", val)
}

// to handle the units and filter by actual values, we need to mirror this:
// https://github.com/Project-OSRM/osrm-profiles-contrib/blob/master/5/21/truck-soft/lib/measure.lua
func truckTagMapFilter(truckHeight, truckWidth, truckLength, truckWeight float64) TagMapFilter {
	return func(wayId int, tagMap map[string]string) bool {
		//see https://wiki.openstreetmap.org/wiki/Key:maxheight
		maxHeightStr, maxHeightOk := tagMap["maxheight"]
		var heightLimit float64
		if maxHeightOk {
			maxHeight, err := parseAsMeters(maxHeightStr)
			if err != nil {
				// TODO: decide on a real logging strategy
				fmt.Fprintf(os.Stderr, "invalid maxheight tag: %v", err)
				maxHeightOk = false
			} else {
				heightLimit = maxHeight
			}
		}

		maxPhysicalHeightStr, maxPhysicalHeightOk := tagMap["maxheight:physical"]
		if maxPhysicalHeightOk {
			maxPhysicalHeight, err := parseAsMeters(maxPhysicalHeightStr)
			if err != nil {
				// TODO: decide on a real logging strategy
				fmt.Fprintf(os.Stderr, "invalid maxheight:physical tag: %v", err)
			} else {
				if !maxHeightOk || maxPhysicalHeight < heightLimit {
					heightLimit = maxPhysicalHeight
				}
			}
		}

		if heightLimit > 0.0 && truckHeight > heightLimit {
			return false
		}

		//// see https://wiki.openstreetmap.org/wiki/Key:maxwidth
		//_, maxWidthOk := tagMap["maxwidth"]
		//_, maxPhysicalWidthOk := tagMap["maxwidth:physical"]
		//_, widthOk := tagMap["width"]

		//if maxWidthOk || maxPhysicalWidthOk || widthOk {
		//	return false
		//}

		//// there is also maxlength:hgv_articulated, see:
		//// https://wiki.openstreetmap.org/wiki/Key:maxlength
		//_, maxLengthOk := tagMap["maxlength"]
		//if maxLengthOk {
		//	return false
		//}

		//// see https://wiki.openstreetmap.org/wiki/Key:maxweight
		//_, maxWeightOk := tagMap["maxweight"]
		//if maxWeightOk {
		//	return false
		//}

		// TODO: there are more things to filter out ways for certain trucks:
		// https://wiki.openstreetmap.org/wiki/Key:maxweightrating
		// https://wiki.openstreetmap.org/wiki/Key:hgv_articulated
		// https://wiki.openstreetmap.org/wiki/Key:maxaxleload
		// https://wiki.openstreetmap.org/wiki/Key:maxbogieweight

		// car is the default for trucks
		return carTagMapFilter(wayId, tagMap)
	}
}
