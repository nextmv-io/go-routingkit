#ifndef __MYCLASS_H
#define __MYCLASS_H
#include <vector>

struct Point {
	float lon;
	float lat;
};

class Client {
        public:
                float distance(float from_longitude, float from_latitude, float to_longitude, float to_latitude);
                std::vector<float> table(std::vector<Point> sources, std::vector<Point> targets);
                void build_ch(char* pbf_file, char* ch_file);
                void load(char *pbf_file, char *ch_file);
                double average(std::vector<int> v);
};

#endif
