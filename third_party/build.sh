set -e

# Move to script dir
HERE="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
pushd "${HERE}" || exit 1

# Set OS and ARCH according to go
GOOS="$( go env GOOS )"
GOARCH=$( go env GOARCH )

# alias ar to llvm-ar
case $GOOS in
	darwin)
		export CC="$(brew --prefix llvm@14)/bin/clang"
		export CXX=$(brew --prefix llvm@14)/bin/clang++
		export AR="$(brew --prefix llvm@14)/bin/llvm-ar"
		alias ar="$(brew --prefix llvm@14)/bin/llvm-ar"
	;;
	linux)
		export AR=ar
	;;
esac

# Compile according to platform
case $GOOS in
	linux)
		g++ -IRoutingKit/include -LRoutingKit/lib/libroutingkit.a \
			-std=c++11 -c Client.cpp -lroutingkit -lz -fopenmp -pthread -lm -fPIC -ffast-math -O3
	;;
	darwin)
		$CXX -IRoutingKit/include \
			-std=c++11 -stdlib=libc++ -c Client.cpp -Xpreprocessor -fopenmp \
			-pthread -fPIC -ffast-math -O3 -mmacosx-version-min=10.15
	;;
esac

mkdir -p temp
cd temp

# Extract libz library
case $GOOS in
	linux)
		# Find libz.a in the usual locations.
		libzlocation=($(find /usr ! -readable -prune -o -name "libz.a" -print))
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

# Link everything
cd ..
$AR rvs libroutingkit.a Client.o RoutingKit/build/* temp/*
rm -r temp

# Generate bindings via swig
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

# Copy to final location
cp Client.h ../routingkit/internal/routingkit/
rm -f ../routingkit/internal/routingkit/include/routingkit/*.h
cp -v RoutingKit/include/routingkit/* ../routingkit/internal/routingkit/include/routingkit/

popd
