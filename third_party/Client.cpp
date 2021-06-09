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
#include <algorithm>

using namespace RoutingKit;
using namespace std;

static ContractionHierarchy ch;
static SimpleOSMCarRoutingGraph graph;
static GeoPositionToNode map;
static std::vector<ContractionHierarchyQuery> queries;

void Client::build_ch(int conc, char *pbf_file, char *ch_file)
{
	// Load a car routing graph from OpenStreetMap-based data
	graph = simple_load_osm_car_routing_graph_from_pbf(pbf_file);
	auto tail = invert_inverse_vector(graph.first_out);

	// Build the shortest path index
	ch = ContractionHierarchy::build(
		graph.node_count(),
		tail, graph.head,
		graph.geo_distance);

	// Store contraction hierarchy
	ch.save_file(ch_file);

	// Build the index to quickly map latitudes and longitudes
	GeoPositionToNode map_geo_position(graph.latitude, graph.longitude);
	map = map_geo_position;

	// Besides the CH itself we need a query object.
	for (int i = 0; i < conc; i++)
	{
		ContractionHierarchyQuery ch_query(ch);
		queries.push_back(ch_query);
	}
}

void Client::load(int conc, char *pbf_file, char *ch_file)
{
	// Load a car routing graph from OpenStreetMap-based data
	graph = simple_load_osm_car_routing_graph_from_pbf(pbf_file);

	// Load corresponding contraction hierarchy
	ch = ContractionHierarchy::load_file(ch_file);

	// Build the index to quickly map latitudes and longitudes
	GeoPositionToNode map_geo_position(graph.latitude, graph.longitude);
	map = map_geo_position;

	// Besides the CH itself we need a query object.
	for (int i = 0; i < conc; i++)
	{
		ContractionHierarchyQuery ch_query(ch);
		queries.push_back(ch_query);
	}
}

Point Client::nearest(int i, float radius, Point p)
{
	auto n = [i, p, radius]()
	{
		unsigned neighbor = map.find_nearest_neighbor_within_radius(p.lat, p.lon, radius).id;
		return Point{lon : graph.longitude[neighbor], lat : graph.latitude[neighbor]};
	};

	auto future = async(launch::deferred, n);
	auto result = future.get();
	return result;
}

std::vector<long int> Client::distances(int i, float radius, Point source, std::vector<struct Point> targets)
{
	auto tbl = [](int i, float radius, Point source, std::vector<struct Point> targets) -> vector<long int>
	{
		vector<long int> results;
		results.resize(targets.size());

		vector<unsigned> target_list;
		vector<int> invalid_ids;
		for (int i = 0; i < targets.size(); i++)
		{
			auto target = targets[i];
			unsigned to = map.find_nearest_neighbor_within_radius(target.lat, target.lon, radius).id;
			if (to == invalid_id)
			{
				invalid_ids.push_back(i);
			}
			else
			{
				target_list.push_back(to);
			}
		}
		queries[i].reset().pin_targets(target_list);

		unsigned from = map.find_nearest_neighbor_within_radius(source.lat, source.lon, radius).id;
		if (from == invalid_id)
		{
			for (int i = 0; i < targets.size(); i++)
			{
				results[i] = -1;
			}
			return results;
		}
		vector<unsigned> distances = queries[i].reset_source().add_source(from).run_to_pinned_targets().get_distances_to_targets();

		auto invalid_id = invalid_ids.begin();
		auto distance = distances.begin();
		for (int i = 0; i < targets.size(); i++)
		{
			if (invalid_id != invalid_ids.end() && i == *invalid_id)
			{
				results[i] = -1;
				invalid_id++;
			}
			else if (*distance == INT_MAX)
			{
				results[i] = -1;
				distance++;
			}
			else
			{
				results[i] = *distance;
				distance++;
			}
		}
		return results;
	};

	auto future = std::async(launch::deferred, tbl, i, radius, source, targets);
	auto result = future.get();
	return result;
}

QueryResponse Client::query(int i, float radius, float from_longitude, float from_latitude, float to_longitude, float to_latitude, bool include_waypoints)
{
	auto query = [](int i, float radius, float from_longitude, float from_latitude, float to_longitude, float to_latitude, bool include_waypoints)
	{
		unsigned from = map.find_nearest_neighbor_within_radius(from_latitude, from_longitude, radius).id;
		unsigned to = map.find_nearest_neighbor_within_radius(to_latitude, to_longitude, radius).id;

		QueryResponse response;
		if (from == invalid_id || to == invalid_id)
		{
			response.distance = -1.0;
			return response;
		}

		queries[i].reset().add_source(from).add_target(to).run();
		auto distance = queries[i].get_distance();
		if (distance == INT_MAX)
		{
			response.distance = -1.0;
			return response;
		}

		response.distance = distance;
		if (include_waypoints)
		{
			auto path = queries[i].get_node_path();
			for (auto x : path)
				response.waypoints.push_back(Point{lon : graph.longitude[x], lat : graph.latitude[x]});
		}

		return response;
	};

	auto future = std::async(launch::deferred, query, i, radius, from_longitude, from_latitude, to_longitude, to_latitude, include_waypoints);
	auto result = future.get();
	return result;
}
