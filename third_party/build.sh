#!/bin/bash
HERE="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
pushd "${HERE}" || exit 1
GOOS="$( go env GOOS )"

case $GOOS in
	linux)
		g++ -IRoutingKit/include -LRoutingKit/lib/libroutingkit.a \
			-std=c++11 -c Client.cpp -lroutingkit -lz -fopenmp -pthread -lm -fPIC -ffast-math -O3
	;;
	darwin)
		clang++ -IRoutingKit/include -LRoutingKit/lib/libroutingkit.a \
			-std=c++11 -stdlib=libc++ -c Client.cpp -lroutingkit -lz -Xpreprocessor -fopenmp -lomp \
			-pthread -lm -fPIC -ffast-math -O3
	;;
esac

mkdir temp
cd temp

case $GOOS in
	linux)
		ar -x /usr/lib/x86_64-linux-gnu/libz.a
	;;
	darwin)
		ar -x /opt/homebrew/opt/zlib/lib/libz.a
	;;
esac

cd ..
ar rvs libroutingkit.a Client.o RoutingKit/build/* temp/*
rm -r temp

case $GOOS in
	linux)
		swig -go -cgo -c++ -IRoutingKit/include/routingkit -intgosize 64 -O routingkit_linux.i
		mv routingkit_linux_wrap.cxx libroutingkit.a ../routingkit/internal/routingkit/
		mv routingkit.go ../routingkit/internal/routingkit/routingkit_linux.go
	;;
	darwin)
		swig -go -cgo -c++ -IRoutingKit/include/routingkit -intgosize 64 -O routingkit_darwin.i
		mv routingkit_darwin_wrap.cxx libroutingkit.a ../routingkit/internal/routingkit/
		mv routingkit.go ../routingkit/internal/routingkit/routingkit_darwin.go
	;;
esac


cp Client.h ../routingkit/internal/routingkit/
rm -rf ../routingkit/internal/routingkit/include
cp -r RoutingKit/include ../routingkit/internal/routingkit/include
popd
