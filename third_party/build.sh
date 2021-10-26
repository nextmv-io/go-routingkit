#!/bin/bash
set -e

HERE="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
pushd "${HERE}" || exit 1
GOOS="$( go env GOOS )"
GOARCH=$( go env GOARCH )

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
		case $GOARCH in
		amd64)
			ar -x /usr/lib/x86_64-linux-gnu/libz.a
		;;
		arm64)
			ar -x /usr/lib64/libz.a
		;;
		esac
	;;
	darwin)
		ar -x "$(brew --prefix zlib)/lib/libz.a"
	;;
esac

cd ..
ar rvs libroutingkit.a Client.o RoutingKit/build/* temp/*
rm -r temp

case $GOOS in
	linux)
		case $GOARCH in
		amd64)
			swig -go -cgo -c++ -IRoutingKit/include/routingkit -intgosize 64 -O routingkit_linux_amd64.i
			mv routingkit_linux_amd64_wrap.cxx ../routingkit/internal/routingkit/
			mv libroutingkit.a ../routingkit/internal/routingkit/libroutingkit_linux_amd64.a
			mv routingkit.go ../routingkit/internal/routingkit/routingkit_linux_amd64.go
		;;
		arm64)
			swig -go -cgo -c++ -IRoutingKit/include/routingkit -intgosize 64 -O routingkit_linux_arm64.i
			mv routingkit_linux_arm64_wrap.cxx ../routingkit/internal/routingkit/
			mv libroutingkit.a ../routingkit/internal/routingkit/libroutingkit_linux_arm64.a
			mv routingkit.go ../routingkit/internal/routingkit/routingkit_linux_arm64.go
		;;
		esac
	;;
	darwin)
		case $GOARCH in
		amd64)
			swig -go -cgo -c++ -IRoutingKit/include/routingkit -intgosize 64 -O routingkit_darwin_amd64.i
			mv routingkit_darwin_amd64_wrap.cxx ../routingkit/internal/routingkit/
			mv libroutingkit.a ../routingkit/internal/routingkit/libroutingkit_darwin_amd64.a
			mv routingkit.go ../routingkit/internal/routingkit/routingkit_darwin_amd64.go
		;;
		arm64)
			swig -go -cgo -c++ -IRoutingKit/include/routingkit -intgosize 64 -O routingkit_darwin_arm64.i
			mv routingkit_darwin_arm64_wrap.cxx ../routingkit/internal/routingkit/
			mv libroutingkit.a ../routingkit/internal/routingkit/libroutingkit_darwin_arm64.a
			mv routingkit.go ../routingkit/internal/routingkit/routingkit_darwin_arm64.go
		;;
		esac
	;;
esac


cp Client.h ../routingkit/internal/routingkit/
rm -rf ../routingkit/internal/routingkit/include
cp -r RoutingKit/include ../routingkit/internal/routingkit/include
popd
