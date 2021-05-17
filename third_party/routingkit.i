%module routingkit
%{
#include "Client.h"
%}

%include "std_vector.i"
// Instantiate templates used by example
namespace std {
  %template(IntVector) vector<int>;
  %template(FloatVector) vector<float>;
  %template(PointVector) vector<Point>;
}

%include "Client.h"
