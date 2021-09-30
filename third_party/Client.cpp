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
#include <unordered_set>
#include <vector>
#include <execinfo.h>
#include <signal.h>
#include <unistd.h>
#include <stdexcept>

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

	bool starts_with(const char *prefix, const char *str)
	{
		while (*prefix != '\0' && *str == *prefix)
		{
			++prefix;
			++str;
		}
		return *prefix == '\0';
	}

	void decode_osm_car_turn_restrictions_custom(
		uint64_t osm_relation_id, const std::vector<OSMRelationMember> &member_list,
		const TagMap &tags,
		std::function<void(OSMTurnRestriction)> on_new_turn_restriction,
		std::function<void(const std::string &)> log_message, bool prevent_left_turns)
	{
		const char *restriction = tags["restriction"];
		if (restriction == nullptr)
			return;

		OSMTurnRestrictionCategory restriction_type;

		int direction_offset;

		if (starts_with("only_", restriction))
		{
			restriction_type = OSMTurnRestrictionCategory::mandatory;
			direction_offset = 5;
		}
		else if (starts_with("no_", restriction))
		{
			restriction_type = OSMTurnRestrictionCategory::prohibitive;
			direction_offset = 3;
		}
		else
		{
			if (log_message)
				log_message("Unknown OSM turn restriction with ID " + std::to_string(osm_relation_id) + " and value \"" + restriction + "\", ignoring restriction");
			return;
		}

		OSMTurnDirection turn_direction;

		if (str_eq("left_turn", restriction + direction_offset))
		{
			turn_direction = OSMTurnDirection::left_turn;
		}
		else if (str_eq("right_turn", restriction + direction_offset))
		{
			turn_direction = OSMTurnDirection::right_turn;
		}
		else if (str_eq("straight_on", restriction + direction_offset))
		{
			turn_direction = OSMTurnDirection::straight_on;
		}
		else if (str_eq("u_turn", restriction + direction_offset))
		{
			turn_direction = OSMTurnDirection::u_turn;
		}
		else
		{
			if (log_message)
				log_message("Unknown OSM turn restriction with ID " + std::to_string(osm_relation_id) + " and value \"" + restriction + "\", ignoring restriction");
			return;
		}

		std::vector<unsigned> from_member_list;
		std::vector<unsigned> to_member_list;
		unsigned via_member = invalid_id;

		for (unsigned i = 0; i < member_list.size(); ++i)
		{
			if (str_eq(member_list[i].role, "via"))
			{
				if (via_member != invalid_id)
				{
					if (log_message)
						log_message("OSM turn restriction with ID " + std::to_string(osm_relation_id) + " has several \"via\" roles, ignoring restriction");
					return;
				}
				via_member = i;
			}
			else if (str_eq(member_list[i].role, "from"))
			{
				from_member_list.push_back(i);
			}
			else if (str_eq(member_list[i].role, "to"))
			{
				to_member_list.push_back(i);
			}
			else if (str_eq(member_list[i].role, "location_hint"))
			{
				// ignore
			}
			else
			{
				if (log_message)
					log_message("OSM turn restriction with ID " + std::to_string(osm_relation_id) + " and unknown role \"" + member_list[i].role + "\", ignoring role");
			}
		}

		if (via_member != invalid_id && member_list[via_member].type == OSMIDType::relation)
		{
			if (log_message)
				log_message("OSM turn restriction with ID " + std::to_string(osm_relation_id) + " has a relation as \"via\"-role, this is invalid, ignoring restriction");
			return;
		}

		if (via_member != invalid_id && member_list[via_member].type == OSMIDType::way)
		{
			//log_message("OSM turn restriction with ID "+std::to_string(osm_relation_id)+" and name \""+restriction+"\" has a relation as \"way\"-role, this feature is not supported, ignoring restriction");
			return;
		}

		uint64_t via_node = (uint64_t)-1;
		if (via_member != invalid_id)
		{
			via_node = member_list[via_member].id;
		}

		from_member_list.erase(
			std::remove_if(
				from_member_list.begin(), from_member_list.end(),
				[&](unsigned member)
				{
					if (member_list[member].type != OSMIDType::way)
					{
						if (log_message)
							log_message("OSM turn restriction with ID " + std::to_string(osm_relation_id) + " has \"from\"-role that is not a way, ignoring role");
						return true;
					}
					else
					{
						return false;
					}
				}),
			from_member_list.end());

		to_member_list.erase(
			std::remove_if(
				to_member_list.begin(), to_member_list.end(),
				[&](unsigned member)
				{
					if (member_list[member].type != OSMIDType::way)
					{
						if (log_message)
							log_message("OSM turn restriction with ID " + std::to_string(osm_relation_id) + " has \"to\"-role that is not a way, ignoring role");
						return true;
					}
					else
					{
						return false;
					}
				}),
			to_member_list.end());

		if (to_member_list.empty())
		{
			if (log_message)
				log_message("OSM turn restriction with ID " + std::to_string(osm_relation_id) + " is missing \"to\" role, ignoring restriction");
			return;
		}

		if (from_member_list.empty())
		{
			if (log_message)
				log_message("OSM turn restriction with ID " + std::to_string(osm_relation_id) + " is missing \"from\" role, ignoring restriction");
			return;
		}

		if (restriction_type == OSMTurnRestrictionCategory::mandatory && to_member_list.size() != 1)
		{
			if (log_message)
				log_message("OSM turn restriction with ID " + std::to_string(osm_relation_id) + " is mandatory but has several \"to\" roles, ignoring restriction");
			return;
		}

		if (restriction_type == OSMTurnRestrictionCategory::mandatory && from_member_list.size() != 1)
		{
			if (log_message)
				log_message("OSM turn restriction with ID " + std::to_string(osm_relation_id) + " is mandatory but has several \"to\" roles, ignoring restriction");
			return;
		}

		// override the restriction_type in case we want to forbid left turns
		if (prevent_left_turns == true && turn_direction == OSMTurnDirection::left_turn)
		{
			restriction_type = OSMTurnRestrictionCategory::prohibitive;
		}

		for (unsigned from_member : from_member_list)
			for (unsigned to_member : to_member_list)
				on_new_turn_restriction(OSMTurnRestriction{osm_relation_id, restriction_type, turn_direction, member_list[from_member].id, via_node, member_list[to_member].id});
	}

	RoutingGraph load_custom_osm_routing_graph_from_pbf(
		const std::string &pbf_file, Profile profile)
	{
		bool all_modelling_nodes_are_routing_nodes = false;
		bool file_is_ordered_even_though_file_header_says_that_it_is_unordered = false;

		std::unordered_set<uint64_t> allowedWayIds;
		for (int i = 0; i < profile.allowedWayIds.size(); i++)
		{
			auto wayId = profile.allowedWayIds[i];
			allowedWayIds.insert(wayId);
		}

		auto mapping = load_osm_id_mapping_from_pbf(
			pbf_file,
			nullptr,
			[&](uint64_t osm_way_id, const TagMap &tags)
			{
				if (allowedWayIds.find(osm_way_id) != allowedWayIds.end())
				{
					return true;
				}
				return false;
			},
			log_message,
			all_modelling_nodes_are_routing_nodes);

		unsigned routing_way_count = mapping.is_routing_way.population_count();

		auto waySpeeds = std::vector<unsigned>(routing_way_count);

		std::function<
			void(
				uint64_t osm_relation_id,
				const std::vector<OSMRelationMember> &member_list,
				const TagMap &tags,
				std::function<void(OSMTurnRestriction)>)>
			turn_restriction_decoder = nullptr;

		if (profile.transportMode == vehicle)
		{
			turn_restriction_decoder = [&](uint64_t osm_relation_id, const std::vector<OSMRelationMember> &member_list, const TagMap &tags, std::function<void(OSMTurnRestriction)> on_new_restriction)
			{
				return decode_osm_car_turn_restrictions_custom(osm_relation_id, member_list, tags, on_new_restriction, log_message, profile.prevent_left_turns);
			};
		}

		auto routing_graph = load_osm_routing_graph_from_pbf(
			pbf_file,
			mapping,
			[&](uint64_t osm_way_id, unsigned routing_way_id, const TagMap &way_tags)
			{
				if (profile.waySpeeds.find(osm_way_id) == profile.waySpeeds.end())
				{
					waySpeeds[routing_way_id] = get_osm_way_speed(osm_way_id, way_tags, log_message);
				}
				else
				{
					waySpeeds[routing_way_id] = profile.waySpeeds[osm_way_id];
				}
				switch (profile.transportMode)
				{
				case vehicle:
					return get_osm_car_direction_category(osm_way_id, way_tags, log_message);
				case bike:
					return get_osm_bicycle_direction_category(osm_way_id, way_tags, log_message);
				}
				return OSMWayDirectionCategory::open_in_both;
			},
			turn_restriction_decoder,
			log_message);

		mapping = OSMRoutingIDMapping(); // release memory

		RoutingGraph ret;
		ret.first_out = std::move(routing_graph.first_out);
		ret.head = std::move(routing_graph.head);
		ret.geo_distance = std::move(routing_graph.geo_distance);
		ret.latitude = std::move(routing_graph.latitude);
		ret.longitude = std::move(routing_graph.longitude);

		ret.travel_time = ret.geo_distance;
		for (unsigned a = 0; a < ret.travel_time.size(); ++a)
		{
			ret.travel_time[a] *= 18000;
			ret.travel_time[a] /= waySpeeds[routing_graph.way[a]];
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

namespace ErrorHandler
{
	void dump_stack(int sig)
	{
		fprintf(stderr, "Error: signal %d:\n", sig);

		void *array[20];
		size_t size;

		// get void*'s for all entries on the stack
		size = backtrace(array, sizeof(array));

		// print out all the frames to stderr
		backtrace_symbols_fd(array, size, STDERR_FILENO);
	}

	void exception_handler(int sig)
	{
		dump_stack(sig);
		exit(1);
	}

	void install_exception_handlers()
	{
		signal(SIGSEGV, exception_handler);
		signal(SIGBUS, exception_handler);
		signal(SIGINT, exception_handler);
		signal(SIGQUIT, exception_handler);
		signal(SIGILL, exception_handler);
		signal(SIGABRT, exception_handler);
		signal(SIGFPE, exception_handler);
		signal(SIGTERM, exception_handler);
		signal(SIGSYS, exception_handler);

		signal(SIGUSR1, dump_stack);
	}
}

Client::Client(int conc, char *pbf_file, char *ch_file, Profile profile)
{
	ErrorHandler::install_exception_handlers();
	vector<unsigned int> tail;

	bool ch_exists = file_exists(ch_file);

	// Load a routing graph from OpenStreetMap-based data
	graph = load_custom_osm_routing_graph_from_pbf(pbf_file, profile);
	tail = invert_inverse_vector(graph.first_out);
	if (ch_exists)
	{
		ch = ContractionHierarchy::load_file(ch_file);
	}
	else
	{
		vector<unsigned> weight = profile.travel_time ? graph.travel_time : graph.geo_distance;
		ch = ContractionHierarchy::build(graph.node_count(), tail, graph.head, weight);
		ch.save_file(ch_file);
	}
	map = GeoPositionToNode{graph.latitude, graph.longitude};
	// Besides the CH itself we need a query object.
	for (int i = 0; i < conc; i++)
	{
		ContractionHierarchyQuery ch_query(ch);
		queries.push_back(ch_query);
	}
}

Point Client::point(int i)
{
	return Point{
		lon :
			graph.longitude[i],
		lat : graph.latitude[i]
	};
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
