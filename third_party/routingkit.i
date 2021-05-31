%module routingkit
%{
#include "Client.h"
%}

%insert(cgo_comment_typedefs) %{
#cgo LDFLAGS: ${SRCDIR}/libroutingkit.a
%}

%include <typemaps.i>
%include "std_vector.i"
// Instantiate templates used by example
namespace std {
  %template(IntVector) vector<int>;
  %template(FloatVector) vector<float>;
  %template(PointVector) vector<Point>;
  %template(UnsignedVector) vector<unsigned>;
}

%include "Client.h"
