#ifndef __MYCLASS_H
#define __MYCLASS_H
#include <vector>

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

namespace RoutingKit
{
        extern const unsigned max_distance;
}

class Client
{
public:
        QueryResponse query(int i, float radius, float from_longitude, float from_latitude, float to_longitude, float to_latitude, bool include_waypoints);
        std::vector<unsigned> distances(int i, float radius, Point source, std::vector<Point> targets);
        Point *nearest(int i, float radius, float lon, float lat);
        void build_ch(int conc, char *pbf_file, char *ch_file);
        void load(int conc, char *pbf_file, char *ch_file);
};

#endif
