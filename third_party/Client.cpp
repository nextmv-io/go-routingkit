#include <routingkit/osm_simple.h>
#include <routingkit/contraction_hierarchy.h>
#include <routingkit/inverse_vector.h>
#include <routingkit/timer.h>
#include <routingkit/geo_position_to_node.h>
#include <routingkit/osm_graph_builder.h>
#include <routingkit/osm_profile.h>
#include "Client.h"
#include <fstream>
#include <iostream>
#include <numeric>
#include <thread>
#include <future>
#include <vector>

using namespace RoutingKit;
using namespace GoRoutingKit;
using namespace std;

namespace GoRoutingKit
{
	const unsigned max_distance = inf_weight;

	void log_message(const std::string &msg)
	{
		// cout << msg << endl;
	}

	bool str_eq(const char *l, const char *r)
	{
		return !strcmp(l, r);
	}

	bool is_osm_way_used_by_cars_custom(uint64_t osm_way_id, const TagMap &tags, std::function<void(const std::string &)> log_message)
	{
		const char *junction = tags["junction"];
		if (junction != nullptr)
			return true;

		const char *route = tags["route"];
		if (route && str_eq(route, "ferry"))
			return true;

		const char *ferry = tags["ferry"];
		if (ferry && str_eq(ferry, "yes"))
			return true;

		const char *highway = tags["highway"];
		if (highway == nullptr)
			return false;

		const char *motorcar = tags["motorcar"];
		if (motorcar && str_eq(motorcar, "no"))
			return false;

		const char *motor_vehicle = tags["motor_vehicle"];
		if (motor_vehicle && str_eq(motor_vehicle, "no"))
			return false;

		const char *access = tags["access"];
		if (access)
		{
			if (!(str_eq(access, "yes") || str_eq(access, "permissive") || str_eq(access, "delivery") || str_eq(access, "designated") || str_eq(access, "destination")))
				return false;
		}

		if (
			str_eq(highway, "motorway") ||
			str_eq(highway, "trunk") ||
			str_eq(highway, "primary") ||
			str_eq(highway, "secondary") ||
			str_eq(highway, "tertiary") ||
			str_eq(highway, "unclassified") ||
			str_eq(highway, "residential") ||
			str_eq(highway, "service") ||
			str_eq(highway, "motorway_link") ||
			str_eq(highway, "trunk_link") ||
			str_eq(highway, "primary_link") ||
			str_eq(highway, "secondary_link") ||
			str_eq(highway, "tertiary_link") ||
			str_eq(highway, "motorway_junction") ||
			str_eq(highway, "living_street") ||
			str_eq(highway, "residential") ||
			str_eq(highway, "track") ||
			str_eq(highway, "ferry"))
			return true;

		if (str_eq(highway, "bicycle_road"))
		{
			auto motorcar = tags["motorcar"];
			if (motorcar != nullptr)
				if (str_eq(motorcar, "yes"))
					return true;
			return false;
		}

		if (
			str_eq(highway, "construction") ||
			str_eq(highway, "path") ||
			str_eq(highway, "footway") ||
			str_eq(highway, "cycleway") ||
			str_eq(highway, "bridleway") ||
			str_eq(highway, "pedestrian") ||
			str_eq(highway, "bus_guideway") ||
			str_eq(highway, "raceway") ||
			str_eq(highway, "escape") ||
			str_eq(highway, "steps") ||
			str_eq(highway, "proposed") ||
			str_eq(highway, "conveying"))
			return false;

		const char *oneway = tags["oneway"];
		if (oneway != nullptr)
		{
			if (str_eq(oneway, "reversible") || str_eq(oneway, "alternating"))
			{
				return false;
			}
		}

		const char *maxspeed = tags["maxspeed"];
		if (maxspeed != nullptr)
			return true;

		return false;
	}

	SimpleOSMCarRoutingGraph simple_load_osm_car_routing_graph_from_pbf_custom(
		const std::string &pbf_file, std::vector<WayFilter> wayfilters)
	{
		bool all_modelling_nodes_are_routing_nodes = false;
		bool file_is_ordered_even_though_file_header_says_that_it_is_unordered = false;

		auto mapping = load_osm_id_mapping_from_pbf(
			pbf_file,
			nullptr,
			[&](uint64_t osm_way_id, const TagMap &tags)
			{
				for (int i = 0; i < wayfilters.size(); i++)
				{
					auto filter = wayfilters[i];
					const char *route = tags[filter.tag];
					// we want that the tag matches
					if (filter.matchTag)
					{
						// if it did match
						if (route)
						{
							// we want the value to match
							if (filter.matchValue)
							{
								if (str_eq(route, filter.value))
								{
									return filter.allowed;
								}
							}
							else
							{
								if (filter.value == nullptr || !str_eq(route, filter.value))
								{
									return filter.allowed;
								}
							}
						}
					}
					else
					{
						// tag did not match and that is what the filter defined
						if (route == nullptr)
						{
							return filter.allowed;
						}
					}
				}
				return true;
				// return is_osm_way_used_by_cars_custom(osm_way_id, tags, log_message);
			},
			log_message,
			all_modelling_nodes_are_routing_nodes);

		unsigned routing_way_count = mapping.is_routing_way.population_count();
		std::vector<unsigned> way_speed(routing_way_count);

		auto routing_graph = load_osm_routing_graph_from_pbf(
			pbf_file,
			mapping,
			[&](uint64_t osm_way_id, unsigned routing_way_id, const TagMap &way_tags)
			{
				way_speed[routing_way_id] = get_osm_way_speed(osm_way_id, way_tags, log_message);
				return get_osm_car_direction_category(osm_way_id, way_tags, log_message);
			},
			[&](uint64_t osm_relation_id, const std::vector<OSMRelationMember> &member_list, const TagMap &tags, std::function<void(OSMTurnRestriction)> on_new_restriction)
			{
				return decode_osm_car_turn_restrictions(osm_relation_id, member_list, tags, on_new_restriction, log_message);
			},
			log_message);

		mapping = OSMRoutingIDMapping(); // release memory

		SimpleOSMCarRoutingGraph ret;
		ret.first_out = std::move(routing_graph.first_out);
		ret.head = std::move(routing_graph.head);
		ret.geo_distance = std::move(routing_graph.geo_distance);
		ret.latitude = std::move(routing_graph.latitude);
		ret.longitude = std::move(routing_graph.longitude);

		ret.travel_time = ret.geo_distance;
		for (unsigned a = 0; a < ret.travel_time.size(); ++a)
		{
			ret.travel_time[a] *= 18000;
			ret.travel_time[a] /= way_speed[routing_graph.way[a]];
			ret.travel_time[a] /= 5;
		}

		ret.forbidden_turn_from_arc = std::move(routing_graph.forbidden_turn_from_arc);
		assert(is_sorted_using_less(ret.forbidden_turn_from_arc));
		ret.forbidden_turn_to_arc = std::move(routing_graph.forbidden_turn_to_arc);

		return ret;
	}
}

bool file_exists(char *file)
{
	ifstream f;
	f.open(file);
	return !!f;
}

Client::Client(int conc, char *pbf_file, char *ch_file, travel_profile prof, bool use_travel_time, Profile customProfile)
{
	vector<unsigned int> tail;
	profile = prof;

	bool ch_exists = file_exists(ch_file);

	// Load a car routing graph from OpenStreetMap-based data
	switch (profile)
	{
	case car:
		car_graph = simple_load_osm_car_routing_graph_from_pbf_custom(pbf_file, customProfile.wayfilters);
		tail = invert_inverse_vector(car_graph.first_out);
		if (ch_exists)
		{
			ch = ContractionHierarchy::load_file(ch_file);
		}
		else
		{
			vector<unsigned> weight = use_travel_time ? car_graph.travel_time : car_graph.geo_distance;
			ch = ContractionHierarchy::build(car_graph.node_count(), tail, car_graph.head, weight);
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
