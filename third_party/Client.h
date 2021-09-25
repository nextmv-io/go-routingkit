#ifndef __MYCLASS_H
#define __MYCLASS_H
#include <vector>
#include <routingkit/osm_simple.h>
#include <routingkit/contraction_hierarchy.h>
#include <routingkit/geo_position_to_node.h>

struct Point
{
        float lon;
        float lat;
};

struct QueryResponse
{
        unsigned distance;
        std::vector<Point> waypoints;
};

struct WayFilter
{
        // the tag to be matched
        const char*tag;
        // optional: the value that the tag has to equal
        const char*value;
        // either the tag has to match or the tag is not allowed to match
        bool matchTag;
        // expresses whether this way is allowed or not due to this filter
        bool allowed;
};

namespace GoRoutingKit
{
        extern const unsigned max_distance;
        enum travel_profile
        {
                car,
                bike,
                pedestrian
        };
        class Client
        {
                Point point(int i);
                RoutingKit::ContractionHierarchy ch;
                RoutingKit::GeoPositionToNode map;
                std::vector<RoutingKit::ContractionHierarchyQuery> queries;
                travel_profile profile;
                RoutingKit::SimpleOSMCarRoutingGraph car_graph;
                RoutingKit::SimpleOSMBicycleRoutingGraph bike_graph;
                RoutingKit::SimpleOSMPedestrianRoutingGraph pedestrian_graph;

        public:
                QueryResponse query(int i, float radius, float from_longitude, float from_latitude,
                                    float to_longitude, float to_latitude, bool include_waypoints);
                std::vector<unsigned> distances(int i, float radius, Point source, std::vector<Point> targets);
                Point *nearest(int i, float radius, float lon, float lat);
                Client(int conc, char *pbf_file, char *ch_file, travel_profile profile, bool use_travel_time);
        };
}

#endif
