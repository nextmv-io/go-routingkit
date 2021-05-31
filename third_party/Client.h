#ifndef __MYCLASS_H
#define __MYCLASS_H
#include <vector>

struct Point {
	float lon;
	float lat;
};

class Client {
        public:
                float distance(int i, float from_longitude, float from_latitude, float to_longitude, float to_latitude);
                float threaded(int i, float from_longitude, float from_latitude, float to_longitude, float to_latitude);
                std::vector<unsigned> table(int i, Point source, std::vector<struct Point> targets);
                void build_ch(char* pbf_file, char* ch_file);
                void load(char *pbf_file, char *ch_file);
                double average(std::vector<int> v);
};

#endif
