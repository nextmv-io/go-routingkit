#!/bin/bash
HERE="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
pushd "${HERE}" || exit 1
sudo apt-get install -y zlib1g-dev
git clone https://github.com/RoutingKit/RoutingKit.git || (cd RoutingKit ; git pull; cd ..)
cd RoutingKit
make
rm -v lib/libroutingkit.so
popd
