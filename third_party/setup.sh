#!/bin/bash
set -eu

# Move to script dir
HERE="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
pushd "${HERE}" || exit 1

# Set OS and ARCH according to go
GOOS="$( go env GOOS )"
GOARCH="$( go env GOARCH )"

# Install dependencies
case $GOOS in
	linux)
		if command -v apt &> /dev/null
		then
			sudo apt-get install -y zlib1g-dev
		elif command -v yum &> /dev/null
		then
			sudo yum install -y zlib-devel zlib-static
		elif command -v pacman &> /dev/null
		then
			sudo pacman -S --noconfirm zlib
		else
			echo "cannot find package manager for zlib installation"
			exit 1
		fi
	;;
	darwin)
		brew install zlib
		brew install libomp
		# this version needs to be compatible to the Xcode version installed on
		# the machine that runs build.sh as per
		# https://en.wikipedia.org/wiki/Xcode#14.x_series
		brew install llvm@14
	;;
esac

# Clone routingkit at specific revision
git clone https://github.com/RoutingKit/RoutingKit.git || (cd RoutingKit ; git pull; cd ..)
cd RoutingKit && git reset --hard && git checkout f7d7d14042268123cf778e6129b99eb2249f7f4d

# Make necessary adjustments for some platforms
if [ "$GOOS" = "darwin" ]; then
	sed -i '' "s/CC=g++/CC=clang++/" Makefile
	sed -i '' "s/\(CFLAGS=.*-std=c++11\) \(.*\)/\1 -stdlib=libc++ \2/" Makefile
	sed -i '' "s/OMP_CFLAGS=-fopenmp/OMP_CFLAGS=-Xpreprocessor -fopenmp -lomp/" Makefile
	sed -i '' "s/OMP_LDFLAGS=-fopenmp/OMP_LDFLAGS=-Xpreprocessor -fopenmp -lomp/" Makefile

	if [ "$GOARCH" = "arm64" ]; then
		sed -i '' "s/-march=native/-mcpu=apple-m1/" Makefile
		sed -i '' "s/-Iinclude/-Iinclude -I\/opt\/homebrew\/opt\/libomp\/include/" Makefile
		sed -i '' "s/^LDFLAGS=/LDFLAGS=-L\/opt\/homebrew\/opt\/libomp\/lib/" Makefile
	else
		sed -i '' "s/-DNDEBUG/-DNDEBUG -DROUTING_KIT_NO_ALIGNED_ALLOC/" Makefile
	fi
fi

# Cleanup
rm -rv build || echo "no build directory"

# Build
make
rm -v lib/libroutingkit.so

popd
