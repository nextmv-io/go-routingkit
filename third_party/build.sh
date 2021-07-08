#!/bin/bash
HERE="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
pushd "${HERE}" || exit 1
g++ -IRoutingKit/include -LRoutingKit/lib/libroutingkit.a -std=c++11 -c Client.cpp -lroutingkit -lz -fopenmp -pthread -lm -fPIC -ffast-math -O3
mkdir temp
cd temp
ar -x /usr/lib/x86_64-linux-gnu/libz.a
cd ..
ar rvs libroutingkit.a Client.o RoutingKit/build/* temp/*
rm -r temp
swig -go -cgo -c++ -IRoutingKit/include/routingkit -intgosize 64 -O routingkit.i
mv routingkit_wrap.cxx routingkit.go libroutingkit.a ../routingkit/internal/routingkit/
cp Client.h ../routingkit/internal/routingkit/
rm -rf ../routingkit/internal/routingkit/include
cp -r RoutingKit/include ../routingkit/internal/routingkit/include
popd
