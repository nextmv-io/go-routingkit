#!/bin/bash
set -eu
HERE="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
pushd "${HERE}" || exit 1
GOOS="$( go env GOOS )"
case $GOOS in
	linux)
		sudo apt-get install -y zlib1g-dev
	;;
	darwin)
		sudo brew install zlib
	;;
esac
git clone https://github.com/RoutingKit/RoutingKit.git || (cd RoutingKit ; git pull; cd ..)
cd RoutingKit
if [ "$GOOS" = "darwin" ]; then
	sed -i 's/CC=g++/CC=clang++/' Makefile
	sed -i 's/\(CFLAGS=.*-std=c++11\) \(.*\)/\1 -stdlib=libc++ \2/' Makefile
	sed -i 's/OMP_CFLAGS=-fopenmp/OMP_CFLAGS=-Xpreprocessor -fopenmp -lomp/' Makefile
	sed -i 's/OMP_LDFLAGS=-fopenmp/OMP_LDFLAGS=-Xpreprocessor -fopenmp -lomp/' Makefile
fi
make
rm -v lib/libroutingkit.so
popd
