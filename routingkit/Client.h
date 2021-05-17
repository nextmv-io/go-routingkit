#ifndef __MYCLASS_H
#define __MYCLASS_H
#include <vector>

class Client {
        public:
                float int_get(float from_latitude,  float from_longitude,  float to_latitude,  float to_longitude);
                void build_ch(char* pbf_file, char* ch_file);
                void load(char *pbf_file, char *ch_file);
                double average(std::vector<int> v);
};

#endif
