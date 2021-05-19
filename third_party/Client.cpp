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
static GeoPositionToNode map;
static ContractionHierarchyQuery query;

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
	ContractionHierarchyQuery ch_query(ch);
	query = ch_query;
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
	ContractionHierarchyQuery ch_query(ch);
	query = ch_query;
}

double Client::average(std::vector<int> v) {
  return std::accumulate(v.begin(), v.end(), 0.0)/v.size();
}

float Client::distance(float from_longitude, float from_latitude, float to_longitude, float to_latitude) {

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
	query.reset().add_source(from).add_target(to).run();
	auto distance = query.get_distance();
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

QueryResponse Client::queryrequest(float radius, float from_longitude, float from_latitude, float to_longitude, float to_latitude) {
	// Use the query object to answer queries from stdin to stdout
	unsigned from = map.find_nearest_neighbor_within_radius(from_latitude, from_longitude, radius).id;
	unsigned to = map.find_nearest_neighbor_within_radius(to_latitude, to_longitude, radius).id;
	QueryResponse response;
	if(from == invalid_id || to == invalid_id){
		response.distance=-1.0;
		return response;
	}

	query.reset().add_source(from).add_target(to).run();
	auto distance = query.get_distance();
	response.distance=distance;
	auto path = query.get_node_path();
	for(auto x:path)
	    response.waypoints.push_back(Point{lon: graph.longitude[x], lat: graph.latitude[x]});

	return response;
}

std::vector<float> Client::table(std::vector<Point> sources, std::vector<Point> targets){
	vector<float> vect;

	for (auto &source : sources)
	{  
		for (auto &target : targets)
		{  
			vect.push_back(distance(source.lon, source.lat, target.lon, target.lat));
		}
	}
 
    return vect;
}
