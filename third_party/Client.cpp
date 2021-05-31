#include "Client.h"
#include <routingkit/osm_simple.h>
#include <routingkit/contraction_hierarchy.h>
#include <routingkit/inverse_vector.h>
#include <routingkit/timer.h>
#include <routingkit/geo_position_to_node.h>
#include <iostream>
#include <numeric>
#include <thread>
#include <future>
using namespace RoutingKit;
using namespace std;

static ContractionHierarchy ch;
static SimpleOSMCarRoutingGraph graph;
static GeoPositionToNode map;
static std::vector<ContractionHierarchyQuery> queries;

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

	// Store contraction hierarchy
	ch.save_file(ch_file);

	// Build the index to quickly map latitudes and longitudes
	GeoPositionToNode map_geo_position(graph.latitude, graph.longitude);
	map = map_geo_position;

	// Besides the CH itself we need a query object. 
	for (int i = 0; i < 100; i++) {
		ContractionHierarchyQuery ch_query(ch);
		queries.push_back(ch_query);
	}
}

void Client::load(char* pbf_file, char* ch_file){
	// Load a car routing graph from OpenStreetMap-based data
    graph = simple_load_osm_car_routing_graph_from_pbf(pbf_file);

	// Load corresponding contraction hierarchy
	ch = ContractionHierarchy::load_file(ch_file);

	// Build the index to quickly map latitudes and longitudes
	GeoPositionToNode map_geo_position(graph.latitude, graph.longitude);
	map = map_geo_position;

	// Besides the CH itself we need a query object. 
	for (int i = 0; i < 100; i++) {
		ContractionHierarchyQuery ch_query(ch);
		queries.push_back(ch_query);
	}
}

double Client::average(std::vector<int> v) {
  return std::accumulate(v.begin(), v.end(), 0.0)/v.size();
}

float Client::threaded(int i, float from_longitude, float from_latitude, float to_longitude, float to_latitude){
	auto dist= [](int i, float from_longitude, float from_latitude, float to_longitude, float to_latitude) {
		// Use the query object to answer queries from stdin to stdout
		unsigned from = map.find_nearest_neighbor_within_radius(from_latitude, from_longitude, 1000).id;
		unsigned to = map.find_nearest_neighbor_within_radius(to_latitude, to_longitude, 1000).id;
		if(from == invalid_id || to == invalid_id){
			return -1.0;
		}
		queries[i].reset().add_source(from).add_target(to).run();
		float distance = queries[i].get_distance();
		return static_cast<double>(distance);
	};

	auto future = std::async(launch::deferred, dist, i, from_longitude, from_latitude, to_longitude, to_latitude);
	float simple = future.get();
	return simple;
}


float Client::distance(int i, float from_longitude, float from_latitude, float to_longitude, float to_latitude) {
	// cout << "distance lon: " << from_longitude << " lat: " << from_latitude <<  ", lon: " << to_longitude << " lat: " << to_latitude << endl;

	long long start_time = get_micro_time();
	// Use the query object to answer queries from stdin to stdout
	unsigned from = map.find_nearest_neighbor_within_radius(from_latitude, from_longitude, 1000).id;
	// if(from == invalid_id){
	// 	//cout << "No node within 1000m from source position" << endl;
	// 	return -1.0;
	// }
	unsigned to = map.find_nearest_neighbor_within_radius(to_latitude, to_longitude, 1000).id;
	// if(to == invalid_id){
	// 	// cout << "No node within 1000m from target position" << endl;
	// 	return -1.0;
	// }
	long long end_time = get_micro_time();
	// cout << "find_nearest_neighbor_within_radius; took " << (end_time - start_time) << " microseconds." << endl;
	if(from == invalid_id || to == invalid_id){
		// cout << "No node within 1000m from target position" << endl;
		return -1.0;
	}

	start_time = get_micro_time();
	queries[i].reset().add_source(from).add_target(to).run();
	auto distance = queries[i].get_distance();
	// auto path = query.get_node_path();
	end_time = get_micro_time();

	// cout << "To get from "<< from << " to "<< to << " one needs " << distance << " meters." << endl;
	// cout << "This query was answered in " << (end_time - start_time) << " microseconds." << endl;
	// cout << "The path is";
	// for(auto x:path)
	// 	cout << " " << graph.longitude[x] << "," << graph.latitude[x];
	// cout << endl;

	return distance;
}

std::vector<unsigned> Client::table(int i, Point source, std::vector<struct Point> targets){
	std::vector<unsigned>target_list;
	for (auto &target : targets)
	{  
		unsigned to = map.find_nearest_neighbor_within_radius(target.lat, target.lon, 1000).id;
		target_list.push_back(to);
	}

	queries[i].reset().pin_targets(target_list);

	unsigned from = map.find_nearest_neighbor_within_radius(source.lat, source.lon, 1000).id;
	return queries[i].reset_source().add_source(from).run_to_pinned_targets().get_distances_to_targets();
}
