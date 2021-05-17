#include "Client.h"
#include <routingkit/osm_simple.h>
#include <routingkit/contraction_hierarchy.h>
#include <routingkit/inverse_vector.h>
#include <routingkit/timer.h>
#include <routingkit/geo_position_to_node.h>
#include <iostream>
#include <numeric>
using namespace RoutingKit;
using namespace std;

static ContractionHierarchy ch;
static SimpleOSMCarRoutingGraph graph;

void Client::build_ch(char* pbf_file, char* ch_file){
	// Load a car routing graph from OpenStreetMap-based data
    graph = simple_load_osm_car_routing_graph_from_pbf(pbf_file);
	auto tail = invert_inverse_vector(graph.first_out);

	// Build the shortest path index
	ch = ContractionHierarchy::build(
		graph.node_count(), 
		tail, graph.head, 
		graph.geo_distance
	);

	ch.save_file(ch_file);
}

void Client::load(char* pbf_file, char* ch_file){
	// Load a car routing graph from OpenStreetMap-based data
    graph = simple_load_osm_car_routing_graph_from_pbf(pbf_file);

	ch = ContractionHierarchy::load_file(ch_file);
}

double Client::average(std::vector<int> v) {
  return std::accumulate(v.begin(), v.end(), 0.0)/v.size();
}

float Client::int_get(float from_latitude,  float from_longitude,  float to_latitude,  float to_longitude) {
	// Build the index to quickly map latitudes and longitudes
	GeoPositionToNode map_geo_position(graph.latitude, graph.longitude);

	// Besides the CH itself we need a query object. 
	ContractionHierarchyQuery ch_query(ch);

	// Use the query object to answer queries from stdin to stdout
	unsigned from = map_geo_position.find_nearest_neighbor_within_radius(from_latitude, from_longitude, 1000).id;
	if(from == invalid_id){
		//cout << "No node within 1000m from source position" << endl;
		return -1.0;
	}
	unsigned to = map_geo_position.find_nearest_neighbor_within_radius(to_latitude, to_longitude, 1000).id;
	if(to == invalid_id){
		// cout << "No node within 1000m from target position" << endl;
		return -1.0;
	}

	// long long start_time = get_micro_time();
	ch_query.reset().add_source(from).add_target(to).run();
	auto distance = ch_query.get_distance();
	// auto path = ch_query.get_node_path();
	// long long end_time = get_micro_time();

	// cout << "To get from "<< from << " to "<< to << " one needs " << distance << " meters." << endl;
	// cout << "This query was answered in " << (end_time - start_time) << " microseconds." << endl;
	// cout << "The path is";
	// for(auto x:path)
	// 	cout << " " << graph.longitude[x] << "," << graph.latitude[x];
	// cout << endl;

	return distance;
}
