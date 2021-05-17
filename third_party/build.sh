#!/bin/bash
HERE="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
pushd "${HERE}" || exit 1
g++ -I../../RoutingKit/include -L../../RoutingKit/lib -std=c++11 Client.cpp -shared -o routingkit.so -lroutingkit -lz -fopenmp -pthread -lm -fPIC
swig -go -cgo -c++ -I../../RoutingKit/include/routingkit -intgosize 64 -O routingkit.i
mv routingkit_wrap.cxx routingkit.go ../routingkit/
cp Client.h ../routingkit/
popd
