#!/usr/bin/env bash

# Specify instance file
instance_file="data/maryland-points.json"
# Specify osm data file to download and use
osm_link="http://download.geofabrik.de/north-america/us/maryland-latest.osm.pbf"
# Specify some depending variables
osm_file_name=$(basename "${osm_link}")
osm_path="data/${osm_file_name}"

# Run script from script's directory
HERE="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
cd "${HERE}" || (echo "cannot switch to script's dir"; exit 1)

# Download osm file, if not already present
if [ ! -f "${osm_path}" ]; then
    echo "Downloading osm file via wget ..."
    wget -P "data/" "${osm_link}"
    echo "Please note: first build of contraction hierarchies may take a while"
else
    echo "Found osm file at ${osm_path}"
fi

# Run routingkit standalone
echo "Running routingkit standalone ..."
go run . \
    -ch "data/${osm_file_name}.ch" \
    -input "${instance_file}" \
    -map "${osm_path}" \
    -measure "distance" \
