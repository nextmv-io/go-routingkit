%module routingkit
%insert(cgo_comment_typedefs) %{
#cgo LDFLAGS: ${SRCDIR}/libroutingkit.a -static -Wl,--whole-archive -lpthread -Wl,--no-whole-archive
#cgo CPPFLAGS: -I${SRCDIR}/../../../third_party/RoutingKit/include/
%}

%{
#include "Client.h"
%}

%include <typemaps.i>
%include "std_vector.i"
// Instantiate templates used by example
namespace std {
  %template(IntVector) vector<int>;
  %template(FloatVector) vector<float>;
  %template(PointVector) vector<Point>;
  %template(UnsignedVector) vector<unsigned>;
  %template(LongIntVector) vector<long int>;
}

%include "Client.h"
