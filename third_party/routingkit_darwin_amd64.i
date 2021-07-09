%module routingkit
%insert(cgo_comment_typedefs) %{
#cgo LDFLAGS: ${SRCDIR}/libroutingkit_darwin_amd64.a
#cgo CPPFLAGS: -I${SRCDIR}/../../../routingkit/internal/routingkit/include
#cgo CXXFLAGS: -std=c++11
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

