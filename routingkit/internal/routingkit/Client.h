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
        float distance;
        std::vector<Point> waypoints;
};

class Client
{
public:
        QueryResponse queryrequest(int i, float radius, float from_longitude, float from_latitude, float to_longitude, float to_latitude, bool include_waypoints);
        std::vector<unsigned> distances(int i, Point source, std::vector<Point> targets);
        void build_ch(int conc, char *pbf_file, char *ch_file);
        void load(int conc, char *pbf_file, char *ch_file);
};

#endif
