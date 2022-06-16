#!/bin/bash
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
		# Find libz.a in the usual locations.
		libzlocation=($(find /usr/lib ! -readable -prune -o -name "libz.a" -print))
		if [ ${#libzlocation[@]} -eq 0 ]; then
			echo "libz.a not found"
			exit 1
		else
			echo "libz.a found at ${libzlocation[0]}"
			ar -x "${libzlocation[0]}"
		fi
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
			mv routingkit_linux_amd64_wrap.cxx ../routingkit/bindings/routingkit/
			mv libroutingkit.a ../routingkit/bindings/routingkit/libroutingkit_linux_amd64.a
			mv routingkit.go ../routingkit/bindings/routingkit/routingkit_linux_amd64.go
		;;
		esac
	;;
	darwin)
		case $GOARCH in
		amd64)
			swig -go -cgo -c++ -IRoutingKit/include/routingkit -intgosize 64 -O routingkit_darwin_amd64.i
			mv routingkit_darwin_amd64_wrap.cxx ../routingkit/bindings/routingkit/
			mv libroutingkit.a ../routingkit/bindings/routingkit/libroutingkit_darwin_amd64.a
			mv routingkit.go ../routingkit/bindings/routingkit/routingkit_darwin_amd64.go
		;;
		arm64)
			swig -go -cgo -c++ -IRoutingKit/include/routingkit -intgosize 64 -O routingkit_darwin_arm64.i
			mv routingkit_darwin_arm64_wrap.cxx ../routingkit/bindings/routingkit/
			mv libroutingkit.a ../routingkit/bindings/routingkit/libroutingkit_darwin_arm64.a
			mv routingkit.go ../routingkit/bindings/routingkit/routingkit_darwin_arm64.go
		;;
		esac
	;;
esac


cp Client.h ../routingkit/bindings/routingkit/
rm -f ../routingkit/bindings/routingkit/include/routingkit/*.h
cp -v RoutingKit/include/routingkit/* ../routingkit/bindings/routingkit/include/routingkit/
popd
