#include <routingkit/osm_simple.h>
#include <routingkit/contraction_hierarchy.h>
#include <routingkit/inverse_vector.h>
#include <routingkit/timer.h>
#include <routingkit/geo_position_to_node.h>
#include "Client.h"
#include <fstream>
#include <iostream>
#include <numeric>
#include <thread>
#include <future>

using namespace RoutingKit;
using namespace GoRoutingKit;
using namespace std;

namespace GoRoutingKit
{
	const unsigned max_distance = RoutingKit::inf_weight;
}

bool file_exists(char *file)
{
	ifstream f;
	f.open(file);
	return !!f;
}

Client::Client(int conc, char *pbf_file, char *ch_file, travel_profile prof)
{
	vector<unsigned int> tail;
	profile = prof;

	bool ch_exists = file_exists(ch_file);

	// Load a car routing graph from OpenStreetMap-based data
	switch (profile)
	{
	case car:
		car_graph = simple_load_osm_car_routing_graph_from_pbf(pbf_file);
		tail = invert_inverse_vector(car_graph.first_out);
		if (ch_exists)
		{
			ch = ContractionHierarchy::load_file(ch_file);
		}
		else
		{
			ch = ContractionHierarchy::build(
			    car_graph.node_count(),
			    tail, car_graph.head,
			    car_graph.geo_distance);
			ch.save_file(ch_file);
		}
		map = GeoPositionToNode{car_graph.latitude, car_graph.longitude};
		break;
	case pedestrian:
		pedestrian_graph = simple_load_osm_pedestrian_routing_graph_from_pbf(pbf_file);
		tail = invert_inverse_vector(pedestrian_graph.first_out);
		if (ch_exists)
		{
			ch = ContractionHierarchy::load_file(ch_file);
		}
		else
		{
			ch = ContractionHierarchy::build(
			    pedestrian_graph.node_count(),
			    tail, pedestrian_graph.head,
			    pedestrian_graph.geo_distance);
			ch.save_file(ch_file);
		}
		map = GeoPositionToNode{pedestrian_graph.latitude, pedestrian_graph.longitude};
		break;
	case bike:
		bike_graph = simple_load_osm_bicycle_routing_graph_from_pbf(pbf_file);
		tail = invert_inverse_vector(bike_graph.first_out);
		if (ch_exists)
		{
			ch = ContractionHierarchy::load_file(ch_file);
		}
		else
		{
			ch = ContractionHierarchy::build(
			    bike_graph.node_count(),
			    tail, bike_graph.head,
			    bike_graph.geo_distance);
			ch.save_file(ch_file);
		}
		map = GeoPositionToNode{bike_graph.latitude, bike_graph.longitude};
		break;
	}

	// Besides the CH itself we need a query object.
	for (int i = 0; i < conc; i++)
	{
		ContractionHierarchyQuery ch_query(ch);
		queries.push_back(ch_query);
	}
}

Point Client::point(int i)
{
	switch (profile)
	{
	case car:
		return Point{
			lon :
			    car_graph.longitude[i],
			lat : car_graph.latitude[i]
		};
	case pedestrian:
		return Point{
			lon :
			    pedestrian_graph.longitude[i],
			lat : pedestrian_graph.latitude[i]
		};
	case bike:
		return Point{
			lon :
			    bike_graph.longitude[i],
			lat : bike_graph.latitude[i]
		};
	};
	return Point{};
}

Point *Client::nearest(int i, float radius, float lon, float lat)
{
	auto n = [this, i, lon, lat, radius]() -> Point *
	{
		unsigned neighbor = map.find_nearest_neighbor_within_radius(lat, lon, radius).id;
		if (neighbor == invalid_id)
			return NULL;
		Point p = point(neighbor);
		return new Point(p);
	};

	return async(launch::deferred, n).get();
}

std::vector<unsigned> Client::distances(int i, float radius, Point source, std::vector<struct Point> targets)
{
	auto tbl = [this, i, radius, source, targets]() -> vector<unsigned int>
	{
		vector<unsigned> results;
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
				results[i] = RoutingKit::inf_weight;
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
				results[i] = RoutingKit::inf_weight;
				invalid_id++;
			}
			else
			{
				results[i] = *distance;
				distance++;
			}
		}
		return results;
	};

	auto future = async(launch::deferred, tbl);
	auto result = future.get();
	return result;
}

QueryResponse Client::query(int i, float radius, float from_longitude, float from_latitude, float to_longitude, float to_latitude, bool include_waypoints)
{
	auto query = [this, i, radius, from_longitude, from_latitude, to_longitude, to_latitude, include_waypoints]()
	{
		unsigned from = map.find_nearest_neighbor_within_radius(from_latitude, from_longitude, radius).id;
		unsigned to = map.find_nearest_neighbor_within_radius(to_latitude, to_longitude, radius).id;

		QueryResponse response;
		if (from == invalid_id || to == invalid_id)
		{
			response.distance = RoutingKit::inf_weight;
			return response;
		}

		queries[i].reset().add_source(from).add_target(to).run();
		auto distance = queries[i].get_distance();

		response.distance = distance;
		if (include_waypoints)
		{
			auto path = queries[i].get_node_path();
			for (auto x : path)
				response.waypoints.push_back(point(x));
		}

		return response;
	};

	auto future = std::async(launch::deferred, query);
	auto result = future.get();
	return result;
}
