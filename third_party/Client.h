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
        float distance(int i, float from_longitude, float from_latitude, float to_longitude, float to_latitude);
        float threaded(int i, float from_longitude, float from_latitude, float to_longitude, float to_latitude);
        QueryResponse queryrequest(int i, float radius, float from_longitude, float from_latitude, float to_longitude, float to_latitude);
        std::vector<float> table(int i, std::vector<Point> sources, std::vector<Point> targets);
        void build_ch(int conc, char *pbf_file, char *ch_file);
        void load(int conc, char *pbf_file, char *ch_file);
        double average(std::vector<int> v);
};

#endif
