#ifndef __MYCLASS_H
#define __MYCLASS_H
#include <vector>
#include <map>
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

enum transport_mode
{
        vehicle = 1,
        bike = 2,
        pedestrian = 3
};

struct Profile
{
        std::vector<int> allowedWayIds;
        std::map<uint64_t, unsigned int> waySpeeds;
        transport_mode transportMode;
        const char *name;
        bool prevent_left_turns;
        bool prevent_u_turns;
        bool travel_time;
};

namespace GoRoutingKit
{
        extern const unsigned max_distance;

        struct RoutingGraph
        {
                std::vector<unsigned> first_out;
                std::vector<unsigned> head;
                std::vector<unsigned> travel_time;
                std::vector<unsigned> geo_distance;
                std::vector<float> latitude;
                std::vector<float> longitude;
                std::vector<unsigned> forbidden_turn_from_arc;
                std::vector<unsigned> forbidden_turn_to_arc;

                unsigned node_count() const
                {
                        return first_out.size() - 1;
                }

                unsigned arc_count() const
                {
                        return head.size();
                }
        };

        class Client
        {
                Point point(int i);
                RoutingKit::ContractionHierarchy ch;
                RoutingKit::GeoPositionToNode map;
                std::vector<RoutingKit::ContractionHierarchyQuery> queries;
                RoutingGraph graph;

        public:
                QueryResponse query(int i, float radius, float from_longitude, float from_latitude,
                                    float to_longitude, float to_latitude, bool include_waypoints);
                std::vector<unsigned> distances(int i, float radius, Point source, std::vector<Point> targets);
                Point *nearest(int i, float radius, float lon, float lat);
                Client(int conc, char *pbf_file, char *ch_file, Profile customProfile);
        };
}

#endif
