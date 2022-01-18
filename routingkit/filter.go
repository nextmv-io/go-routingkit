package routingkit

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type TagMapFilter func(wayId int, tagMap map[string]string) bool

func allLanesAreHOV(tags map[string]string) bool {
	lanes := strings.Split(tags["hov:lanes"], "|")
	for _, lane := range lanes {
		if lane != "designated" {
			return false
		}
	}
	return true
}

// CarTagMapFilter filters the map for map tags usable by a car
func CarTagMapFilter(id int, tags map[string]string) bool {
	highway := tags["highway"]
	route := tags["route"]
	if highway == "" && route == "" {
		return false
	}

	if tags["area"] == "yes" {
		return false
	}

	if val, ok := tags["route"]; ok && val == "shuttle_train" {
		return true
	}
	// TODO: make this configurable
	if tags["toll"] == "yes" {
		return false
	}
	if highway == "steps" {
		return false
	}
	if highway == "construction" || tags["railway"] == "construction" {
		return false
	}
	if construction, ok := tags["construction"]; ok && !(construction == "no" ||
		construction == "widening" ||
		construction == "minor") {
		return false
	}
	if tags["proposed"] != "" {
		return false
	}
	// TODO: if we are time-aware we may be able to handle this better
	if tags["oneway"] == "reversible" {
		return false
	}
	if tags["impassable"] == "yes" || tags["status"] == "impassable" {
		return false
	}

	if highway == "area" ||
		highway == "reversible" ||
		highway == "impassable" ||
		highway == "hov_lanes" ||
		highway == "steps" ||
		highway == "construction" ||
		highway == "proposed" {
		return false
	}

	var access string
	if motorcar, ok := tags["motorcar"]; ok {
		access = motorcar
	} else if motorVehicle, ok := tags["motor_vehicle"]; ok {
		access = motorVehicle
	} else if vehicle, ok := tags["vehicle"]; ok {
		access = vehicle
	} else if accessVal, ok := tags["access"]; ok {
		access = accessVal
	}
	if access == "no" ||
		access == "agricultural" ||
		access == "forestry" ||
		access == "emergency" ||
		access == "psv" ||
		access == "private" {
		return false
	}
	// TODO: the following access tags should only be usable if the way is needed to reach
	// the destination: "customers", "private", "delivery", "destination".
	// To implement this, we can set the edge weight to be very high.
	if highway == "service" {
		return true
	}
	if highway == "motorway" ||
		highway == "motorway_link" ||
		highway == "trunk" ||
		highway == "trunk_link" ||
		highway == "primary" ||
		highway == "primary_link" ||
		highway == "secondary" ||
		highway == "secondary_link" ||
		highway == "tertiary" ||
		highway == "tertiary_link" ||
		highway == "unclassified" ||
		highway == "residential" ||
		highway == "living_street" ||
		highway == "service" {
		return true
	}
	if val, ok := tags["route"]; ok && val == "ferry" {
		return true
	}
	if tags["bridge"] == "movable" && tags["capacity:car"] != "0" {
		return true
	}

	if tags["service"] == "emergency_access" {
		return false
	}

	//TODO: should be configurable whether the vehicle can travel in HOV lanes
	if allLanesAreHOV(tags) {
		return false
	}
	if access == "yes" ||
		access == "motorcar" ||
		access == "motor_vehicle" ||
		access == "vehicle" ||
		access == "permissive" ||
		access == "designated" ||
		access == "hov" {
		return true
	}

	return false
}

// BikeTagMapFilter filters the map for map tags usable by bikes
func BikeTagMapFilter(_ int, tags map[string]string) bool {
	// TODO: make it configurable whether to use public transit - now we assume yes
	if tags["impassable"] == "yes" {
		return false
	}
	if construction, ok := tags["construction"]; ok && !(construction == "no" ||
		construction == "widening" ||
		construction == "minor") {
		return false
	}
	highway := tags["highway"]
	route := tags["route"]
	railway := tags["railway"]
	amenity := tags["amenity"]
	manMade := tags["man_made"]
	publicTransport := tags["public_transport"]
	bridge := tags["bridge"]
	if highway == "" &&
		route == "" &&
		railway == "" &&
		amenity == "" &&
		manMade == "" &&
		publicTransport == "" &&
		bridge == "" {
		return false
	}

	var access string
	if bicycle, ok := tags["bicycle"]; ok {
		access = bicycle
	} else if vehicle, ok := tags["vehicle"]; ok {
		access = vehicle
	} else if accessTag, ok := tags["access"]; ok {
		access = accessTag
	}
	if access == "no" ||
		access == "private" ||
		access == "agricultural" ||
		access == "forestry" ||
		access == "delivery" ||
		access == "use_sidepath" {
		return false
	}
	if access == "yes" ||
		access == "permissive" ||
		access == "designated" {
		return true
	}
	if val, ok := tags["route"]; ok && val == "ferry" {
		return true
	}
	if tags["bridge"] == "movable" {
		return true
	}
	if railway == "platform" || publicTransport == "platform" {
		return true
	}
	if railway == "train" ||
		railway == "railway" ||
		railway == "subway" ||
		railway == "light_rail" ||
		railway == "monorail" ||
		railway == "tram" {
		return true
	}
	// TODO: as with pedestrian, this should include amenity=parking
	// and amenity=parking_entrance, but somehow including these ways
	// results in the query being unroutable
	if _, ok := tags["cycleway"]; ok {
		return true
	}
	if _, ok := tags["cycleway:left"]; ok {
		return true
	}
	if _, ok := tags["cycleway:right"]; ok {
		return true
	}
	if _, ok := tags["cycleway:both"]; ok {
		return true
	}
	if highway == "cycleway" ||
		highway == "primary" ||
		highway == "primary_link" ||
		highway == "secondary" ||
		highway == "secondary_link" ||
		highway == "tertiary" ||
		highway == "tertiary_link" ||
		highway == "residential" ||
		highway == "unclassified" ||
		highway == "living_street" ||
		highway == "road" ||
		highway == "service" ||
		highway == "track" ||
		highway == "path" {
		return true
	}

	return false
}

// PedestrianTagMapFilter filters the map for map tags usable by pedestrians
func PedestrianTagMapFilter(id int, tags map[string]string) bool {
	{
		var routable bool
		for _, tag := range []string{
			"highway",
			"bridge",
			"route",
			"leisure",
			"man_made",
			"railway",
			"platform",
			"amenity",
			"public_transport",
		} {
			if tags[tag] != "" {
				routable = true
				break
			}
		}
		if !routable {
			return false
		}
	}

	if tags["impassable"] == "yes" {
		return false
	}
	if tags["status"] == "impassable" {
		return false
	}

	var access string
	if foot, ok := tags["foot"]; ok {
		access = foot
	} else if accessTag, ok := tags["access"]; ok {
		access = accessTag
	}

	if access == "no" ||
		access == "agricultural" ||
		access == "forestry" ||
		access == "private" {
		return false
	}
	if access == "yes" ||
		access == "foot" ||
		access == "permissive" ||
		access == "destination" ||
		access == "delivery" ||
		access == "designated" {
		return true
	}
	if tags["bridge"] == "movable" {
		return true
	}
	railway := tags["railway"]
	publicTransport := tags["public_transport"]
	if railway == "platform" || publicTransport == "platform" {
		return true
	}
	if railway == "train" ||
		railway == "railway" ||
		railway == "subway" ||
		railway == "light_rail" ||
		railway == "monorail" ||
		railway == "tram" {
		return true
	}

	if val := tags["route"]; val == "ferry" {
		return true
	}
	highway := tags["highway"]
	manMade := tags["man_made"]

	// TODO: OSRM includes ways where leisure=track, but I found that this
	// sometimes caused no valid route to be found. I have not been able
	// to unpack why this type of way is problematic, but my theory is it
	// has something to do with tracks being circular (for an example, see
	// way 352240450), and that OSRM must be able to handle this in a way
	// that routingkit is not. Nevertheless, I don't expect it should be
	// a huge problem to omit this type of way.
	// The same issue exists with amenity=parking/amenity=parking_entrance
	if highway == "primary" ||
		highway == "primary_link" ||
		highway == "secondary" ||
		highway == "secondary_link" ||
		highway == "tertiary" ||
		highway == "tertiary_link" ||
		highway == "residential" ||
		highway == "unclassified" ||
		highway == "living_street" ||
		highway == "road" ||
		highway == "service" ||
		highway == "track" ||
		highway == "path" ||
		highway == "steps" ||
		highway == "pedestrian" ||
		highway == "footway" ||
		highway == "pier" ||
		railway == "platform" ||
		manMade == "pier" {
		return true
	}

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

var weightMeasure = regexp.MustCompile(`^(\d+(?:\.\d*)?)(?: (t|kg|st|lt|lbs|cwt))?$`)

const tonnesPerShortTon = .9071847
const tonnesPerLbs = 0.00045359237
const tonnesPerKG = 0.001
const tonnesPerLongTon = 1.016047
const tonnesPerLongHundredWeights = .05080

func parseAsTonnes(val string) (float64, error) {
	weight := weightMeasure.FindStringSubmatch(val)
	if len(weight) == 3 {
		f, err := strconv.ParseFloat(weight[1], 64)
		if err != nil {
			return 0, fmt.Errorf("invalid weight value %s: %v", val, err)
		}
		switch weight[2] {
		case "":
			return f, nil
		case "t":
			return f, nil
		case "st":
			return f * tonnesPerShortTon, nil
		case "lbs":
			return f * tonnesPerLbs, nil
		case "kg":
			return f * tonnesPerKG, nil
		case "lt":
			return f * tonnesPerLongTon, nil
		case "cwt":
			return f * tonnesPerLongHundredWeights, nil
		}
	}
	return 0, fmt.Errorf("could not parse %s as tonnes value", val)
}

// TruckTagMapFilter filters the map for map tags usable by trucks
func TruckTagMapFilter(truckHeight, truckWidth, truckLength, truckWeight float64) TagMapFilter {
	// to handle the units and filter by actual values, we need to mirror this:
	// https://github.com/Project-OSRM/osrm-profiles-contrib/blob/master/5/21/truck-soft/lib/measure.lua
	return func(wayId int, tagMap map[string]string) bool {
		parseMeterTag := func(tag string) float64 {
			str, ok := tagMap[tag]
			if !ok {
				return 0.0
			}
			if str == "default" ||
				str == "below_default" ||
				str == "no_indications" ||
				str == "no_sign" ||
				str == "none" ||
				str == "unsigned" {
				return 0.0
			}
			meters, err := parseAsMeters(str)
			if err != nil {
				// TODO: decide on a real logging strategy
				fmt.Fprintf(os.Stderr, "invalid %s tag %s: %v\n", tag, str, err)
				return 0.0
			}
			return meters
		}
		min := func(a, b float64) float64 {
			if b > a {
				return b
			}
			return a
		}

		//see https://wiki.openstreetmap.org/wiki/Key:maxheight
		heightLimit := min(parseMeterTag("maxheight"), parseMeterTag("maxheight:physical"))
		if heightLimit > 0.0 && truckHeight > heightLimit {
			return false
		}

		widthLimit := min(parseMeterTag("maxwidth"), parseMeterTag("maxwidth:physical"))
		if widthLimit > 0.0 && truckWidth > widthLimit {
			return false
		}

		//// there is also maxlength:hgv_articulated, see:
		//// https://wiki.openstreetmap.org/wiki/Key:maxlength
		lengthLimit := parseMeterTag("maxlength")
		if lengthLimit > 0.0 && truckLength > lengthLimit {
			return false
		}

		//// see https://wiki.openstreetmap.org/wiki/Key:maxweight
		if maxweight, ok := tagMap["maxweight"]; ok {
			weightLimit, err := parseAsTonnes(maxweight)
			if err != nil {
				fmt.Fprintf(os.Stderr, "invalid maxweight tag %s: %v\n", tagMap["maxweight"], err)
			}
			if weightLimit > 0.0 && truckWeight > weightLimit {
				return false
			}
		}

		// TODO: there are more things to filter out ways for certain trucks:
		// https://wiki.openstreetmap.org/wiki/Key:maxweightrating
		// https://wiki.openstreetmap.org/wiki/Key:hgv_articulated
		// https://wiki.openstreetmap.org/wiki/Key:maxaxleload
		// https://wiki.openstreetmap.org/wiki/Key:maxbogieweight

		// car is the default for trucks
		return CarTagMapFilter(wayId, tagMap)
	}
}
