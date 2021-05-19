/* ----------------------------------------------------------------------------
 * This file was automatically generated by SWIG (http://www.swig.org).
 * Version 4.0.1
 *
 * This file is not intended to be easily readable and contains a number of
 * coding conventions designed to improve portability and efficiency. Do not make
 * changes to this file unless you know what you are doing--modify the SWIG
 * interface file instead.
 * ----------------------------------------------------------------------------- */

// source: routingkit.i

#define SWIGMODULE routingkit

#ifdef __cplusplus
/* SwigValueWrapper is described in swig.swg */
template<typename T> class SwigValueWrapper {
  struct SwigMovePointer {
    T *ptr;
    SwigMovePointer(T *p) : ptr(p) { }
    ~SwigMovePointer() { delete ptr; }
    SwigMovePointer& operator=(SwigMovePointer& rhs) { T* oldptr = ptr; ptr = 0; delete oldptr; ptr = rhs.ptr; rhs.ptr = 0; return *this; }
  } pointer;
  SwigValueWrapper& operator=(const SwigValueWrapper<T>& rhs);
  SwigValueWrapper(const SwigValueWrapper<T>& rhs);
public:
  SwigValueWrapper() : pointer(0) { }
  SwigValueWrapper& operator=(const T& t) { SwigMovePointer tmp(new T(t)); pointer = tmp; return *this; }
  operator T&() const { return *pointer.ptr; }
  T *operator&() { return pointer.ptr; }
};

template <typename T> T SwigValueInit() {
  return T();
}
#endif

/* -----------------------------------------------------------------------------
 *  This section contains generic SWIG labels for method/variable
 *  declarations/attributes, and other compiler dependent labels.
 * ----------------------------------------------------------------------------- */

/* template workaround for compilers that cannot correctly implement the C++ standard */
#ifndef SWIGTEMPLATEDISAMBIGUATOR
# if defined(__SUNPRO_CC) && (__SUNPRO_CC <= 0x560)
#  define SWIGTEMPLATEDISAMBIGUATOR template
# elif defined(__HP_aCC)
/* Needed even with `aCC -AA' when `aCC -V' reports HP ANSI C++ B3910B A.03.55 */
/* If we find a maximum version that requires this, the test would be __HP_aCC <= 35500 for A.03.55 */
#  define SWIGTEMPLATEDISAMBIGUATOR template
# else
#  define SWIGTEMPLATEDISAMBIGUATOR
# endif
#endif

/* inline attribute */
#ifndef SWIGINLINE
# if defined(__cplusplus) || (defined(__GNUC__) && !defined(__STRICT_ANSI__))
#   define SWIGINLINE inline
# else
#   define SWIGINLINE
# endif
#endif

/* attribute recognised by some compilers to avoid 'unused' warnings */
#ifndef SWIGUNUSED
# if defined(__GNUC__)
#   if !(defined(__cplusplus)) || (__GNUC__ > 3 || (__GNUC__ == 3 && __GNUC_MINOR__ >= 4))
#     define SWIGUNUSED __attribute__ ((__unused__))
#   else
#     define SWIGUNUSED
#   endif
# elif defined(__ICC)
#   define SWIGUNUSED __attribute__ ((__unused__))
# else
#   define SWIGUNUSED
# endif
#endif

#ifndef SWIG_MSC_UNSUPPRESS_4505
# if defined(_MSC_VER)
#   pragma warning(disable : 4505) /* unreferenced local function has been removed */
# endif
#endif

#ifndef SWIGUNUSEDPARM
# ifdef __cplusplus
#   define SWIGUNUSEDPARM(p)
# else
#   define SWIGUNUSEDPARM(p) p SWIGUNUSED
# endif
#endif

/* internal SWIG method */
#ifndef SWIGINTERN
# define SWIGINTERN static SWIGUNUSED
#endif

/* internal inline SWIG method */
#ifndef SWIGINTERNINLINE
# define SWIGINTERNINLINE SWIGINTERN SWIGINLINE
#endif

/* exporting methods */
#if defined(__GNUC__)
#  if (__GNUC__ >= 4) || (__GNUC__ == 3 && __GNUC_MINOR__ >= 4)
#    ifndef GCC_HASCLASSVISIBILITY
#      define GCC_HASCLASSVISIBILITY
#    endif
#  endif
#endif

#ifndef SWIGEXPORT
# if defined(_WIN32) || defined(__WIN32__) || defined(__CYGWIN__)
#   if defined(STATIC_LINKED)
#     define SWIGEXPORT
#   else
#     define SWIGEXPORT __declspec(dllexport)
#   endif
# else
#   if defined(__GNUC__) && defined(GCC_HASCLASSVISIBILITY)
#     define SWIGEXPORT __attribute__ ((visibility("default")))
#   else
#     define SWIGEXPORT
#   endif
# endif
#endif

/* calling conventions for Windows */
#ifndef SWIGSTDCALL
# if defined(_WIN32) || defined(__WIN32__) || defined(__CYGWIN__)
#   define SWIGSTDCALL __stdcall
# else
#   define SWIGSTDCALL
# endif
#endif

/* Deal with Microsoft's attempt at deprecating C standard runtime functions */
#if !defined(SWIG_NO_CRT_SECURE_NO_DEPRECATE) && defined(_MSC_VER) && !defined(_CRT_SECURE_NO_DEPRECATE)
# define _CRT_SECURE_NO_DEPRECATE
#endif

/* Deal with Microsoft's attempt at deprecating methods in the standard C++ library */
#if !defined(SWIG_NO_SCL_SECURE_NO_DEPRECATE) && defined(_MSC_VER) && !defined(_SCL_SECURE_NO_DEPRECATE)
# define _SCL_SECURE_NO_DEPRECATE
#endif

/* Deal with Apple's deprecated 'AssertMacros.h' from Carbon-framework */
#if defined(__APPLE__) && !defined(__ASSERT_MACROS_DEFINE_VERSIONS_WITHOUT_UNDERSCORES)
# define __ASSERT_MACROS_DEFINE_VERSIONS_WITHOUT_UNDERSCORES 0
#endif

/* Intel's compiler complains if a variable which was never initialised is
 * cast to void, which is a common idiom which we use to indicate that we
 * are aware a variable isn't used.  So we just silence that warning.
 * See: https://github.com/swig/swig/issues/192 for more discussion.
 */
#ifdef __INTEL_COMPILER
# pragma warning disable 592
#endif


#include <stddef.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/types.h>



typedef long long intgo;
typedef unsigned long long uintgo;


# if !defined(__clang__) && (defined(__i386__) || defined(__x86_64__))
#   define SWIGSTRUCTPACKED __attribute__((__packed__, __gcc_struct__))
# else
#   define SWIGSTRUCTPACKED __attribute__((__packed__))
# endif



typedef struct { char *p; intgo n; } _gostring_;
typedef struct { void* array; intgo len; intgo cap; } _goslice_;




#define swiggo_size_assert_eq(x, y, name) typedef char name[(x-y)*(x-y)*-2+1];
#define swiggo_size_assert(t, n) swiggo_size_assert_eq(sizeof(t), n, swiggo_sizeof_##t##_is_not_##n)

swiggo_size_assert(char, 1)
swiggo_size_assert(short, 2)
swiggo_size_assert(int, 4)
typedef long long swiggo_long_long;
swiggo_size_assert(swiggo_long_long, 8)
swiggo_size_assert(float, 4)
swiggo_size_assert(double, 8)

#ifdef __cplusplus
extern "C" {
#endif
extern void crosscall2(void (*fn)(void *, int), void *, int);
extern char* _cgo_topofstack(void) __attribute__ ((weak));
extern void _cgo_allocate(void *, int);
extern void _cgo_panic(void *, int);
#ifdef __cplusplus
}
#endif

static char *_swig_topofstack() {
  if (_cgo_topofstack) {
    return _cgo_topofstack();
  } else {
    return 0;
  }
}

static void _swig_gopanic(const char *p) {
  struct {
    const char *p;
  } SWIGSTRUCTPACKED a;
  a.p = p;
  crosscall2(_cgo_panic, &a, (int) sizeof a);
}




#define SWIG_contract_assert(expr, msg) \
  if (!(expr)) { _swig_gopanic(msg); } else


static void Swig_free(void* p) {
  free(p);
}

static void* Swig_malloc(int c) {
  return malloc(c);
}


#include "Client.h"


#include <vector>
#include <stdexcept>

SWIGINTERN std::vector< int >::const_reference std_vector_Sl_int_Sg__get(std::vector< int > *self,int i){
                int size = int(self->size());
                if (i>=0 && i<size)
                    return (*self)[i];
                else
                    throw std::out_of_range("vector index out of range");
            }
SWIGINTERN void std_vector_Sl_int_Sg__set(std::vector< int > *self,int i,std::vector< int >::value_type const &val){
                int size = int(self->size());
                if (i>=0 && i<size)
                    (*self)[i] = val;
                else
                    throw std::out_of_range("vector index out of range");
            }
SWIGINTERN std::vector< float >::const_reference std_vector_Sl_float_Sg__get(std::vector< float > *self,int i){
                int size = int(self->size());
                if (i>=0 && i<size)
                    return (*self)[i];
                else
                    throw std::out_of_range("vector index out of range");
            }
SWIGINTERN void std_vector_Sl_float_Sg__set(std::vector< float > *self,int i,std::vector< float >::value_type const &val){
                int size = int(self->size());
                if (i>=0 && i<size)
                    (*self)[i] = val;
                else
                    throw std::out_of_range("vector index out of range");
            }
SWIGINTERN std::vector< Point >::const_reference std_vector_Sl_Point_Sg__get(std::vector< Point > *self,int i){
                int size = int(self->size());
                if (i>=0 && i<size)
                    return (*self)[i];
                else
                    throw std::out_of_range("vector index out of range");
            }
SWIGINTERN void std_vector_Sl_Point_Sg__set(std::vector< Point > *self,int i,std::vector< Point >::value_type const &val){
                int size = int(self->size());
                if (i>=0 && i<size)
                    (*self)[i] = val;
                else
                    throw std::out_of_range("vector index out of range");
            }
#ifdef __cplusplus
extern "C" {
#endif

void _wrap_Swig_free_routingkit_007dcf72057aca2d(void *_swig_go_0) {
  void *arg1 = (void *) 0 ;
  
  arg1 = *(void **)&_swig_go_0; 
  
  Swig_free(arg1);
  
}


void *_wrap_Swig_malloc_routingkit_007dcf72057aca2d(intgo _swig_go_0) {
  int arg1 ;
  void *result = 0 ;
  void *_swig_go_result;
  
  arg1 = (int)_swig_go_0; 
  
  result = (void *)Swig_malloc(arg1);
  *(void **)&_swig_go_result = (void *)result; 
  return _swig_go_result;
}


std::vector< int > *_wrap_new_IntVector__SWIG_0_routingkit_007dcf72057aca2d() {
  std::vector< int > *result = 0 ;
  std::vector< int > *_swig_go_result;
  
  
  result = (std::vector< int > *)new std::vector< int >();
  *(std::vector< int > **)&_swig_go_result = (std::vector< int > *)result; 
  return _swig_go_result;
}


std::vector< int > *_wrap_new_IntVector__SWIG_1_routingkit_007dcf72057aca2d(long long _swig_go_0) {
  std::vector< int >::size_type arg1 ;
  std::vector< int > *result = 0 ;
  std::vector< int > *_swig_go_result;
  
  arg1 = (size_t)_swig_go_0; 
  
  result = (std::vector< int > *)new std::vector< int >(arg1);
  *(std::vector< int > **)&_swig_go_result = (std::vector< int > *)result; 
  return _swig_go_result;
}


std::vector< int > *_wrap_new_IntVector__SWIG_2_routingkit_007dcf72057aca2d(std::vector< int > *_swig_go_0) {
  std::vector< int > *arg1 = 0 ;
  std::vector< int > *result = 0 ;
  std::vector< int > *_swig_go_result;
  
  arg1 = *(std::vector< int > **)&_swig_go_0; 
  
  result = (std::vector< int > *)new std::vector< int >((std::vector< int > const &)*arg1);
  *(std::vector< int > **)&_swig_go_result = (std::vector< int > *)result; 
  return _swig_go_result;
}


long long _wrap_IntVector_size_routingkit_007dcf72057aca2d(std::vector< int > *_swig_go_0) {
  std::vector< int > *arg1 = (std::vector< int > *) 0 ;
  std::vector< int >::size_type result;
  long long _swig_go_result;
  
  arg1 = *(std::vector< int > **)&_swig_go_0; 
  
  result = ((std::vector< int > const *)arg1)->size();
  _swig_go_result = result; 
  return _swig_go_result;
}


long long _wrap_IntVector_capacity_routingkit_007dcf72057aca2d(std::vector< int > *_swig_go_0) {
  std::vector< int > *arg1 = (std::vector< int > *) 0 ;
  std::vector< int >::size_type result;
  long long _swig_go_result;
  
  arg1 = *(std::vector< int > **)&_swig_go_0; 
  
  result = ((std::vector< int > const *)arg1)->capacity();
  _swig_go_result = result; 
  return _swig_go_result;
}


void _wrap_IntVector_reserve_routingkit_007dcf72057aca2d(std::vector< int > *_swig_go_0, long long _swig_go_1) {
  std::vector< int > *arg1 = (std::vector< int > *) 0 ;
  std::vector< int >::size_type arg2 ;
  
  arg1 = *(std::vector< int > **)&_swig_go_0; 
  arg2 = (size_t)_swig_go_1; 
  
  (arg1)->reserve(arg2);
  
}


bool _wrap_IntVector_isEmpty_routingkit_007dcf72057aca2d(std::vector< int > *_swig_go_0) {
  std::vector< int > *arg1 = (std::vector< int > *) 0 ;
  bool result;
  bool _swig_go_result;
  
  arg1 = *(std::vector< int > **)&_swig_go_0; 
  
  result = (bool)((std::vector< int > const *)arg1)->empty();
  _swig_go_result = result; 
  return _swig_go_result;
}


void _wrap_IntVector_clear_routingkit_007dcf72057aca2d(std::vector< int > *_swig_go_0) {
  std::vector< int > *arg1 = (std::vector< int > *) 0 ;
  
  arg1 = *(std::vector< int > **)&_swig_go_0; 
  
  (arg1)->clear();
  
}


void _wrap_IntVector_add_routingkit_007dcf72057aca2d(std::vector< int > *_swig_go_0, intgo _swig_go_1) {
  std::vector< int > *arg1 = (std::vector< int > *) 0 ;
  std::vector< int >::value_type *arg2 = 0 ;
  
  arg1 = *(std::vector< int > **)&_swig_go_0; 
  arg2 = (std::vector< int >::value_type *)&_swig_go_1; 
  
  (arg1)->push_back((std::vector< int >::value_type const &)*arg2);
  
}


intgo _wrap_IntVector_get_routingkit_007dcf72057aca2d(std::vector< int > *_swig_go_0, intgo _swig_go_1) {
  std::vector< int > *arg1 = (std::vector< int > *) 0 ;
  int arg2 ;
  std::vector< int >::value_type *result = 0 ;
  intgo _swig_go_result;
  
  arg1 = *(std::vector< int > **)&_swig_go_0; 
  arg2 = (int)_swig_go_1; 
  
  try {
    result = (std::vector< int >::value_type *) &std_vector_Sl_int_Sg__get(arg1,arg2);
  } catch(std::out_of_range &_e) {
    (void)_e;
    _swig_gopanic("C++ std::out_of_range exception thrown");
    
  }
  _swig_go_result = (int)*result; 
  return _swig_go_result;
}


void _wrap_IntVector_set_routingkit_007dcf72057aca2d(std::vector< int > *_swig_go_0, intgo _swig_go_1, intgo _swig_go_2) {
  std::vector< int > *arg1 = (std::vector< int > *) 0 ;
  int arg2 ;
  std::vector< int >::value_type *arg3 = 0 ;
  
  arg1 = *(std::vector< int > **)&_swig_go_0; 
  arg2 = (int)_swig_go_1; 
  arg3 = (std::vector< int >::value_type *)&_swig_go_2; 
  
  try {
    std_vector_Sl_int_Sg__set(arg1,arg2,(int const &)*arg3);
  } catch(std::out_of_range &_e) {
    (void)_e;
    _swig_gopanic("C++ std::out_of_range exception thrown");
    
  }
  
}


void _wrap_delete_IntVector_routingkit_007dcf72057aca2d(std::vector< int > *_swig_go_0) {
  std::vector< int > *arg1 = (std::vector< int > *) 0 ;
  
  arg1 = *(std::vector< int > **)&_swig_go_0; 
  
  delete arg1;
  
}


std::vector< float > *_wrap_new_FloatVector__SWIG_0_routingkit_007dcf72057aca2d() {
  std::vector< float > *result = 0 ;
  std::vector< float > *_swig_go_result;
  
  
  result = (std::vector< float > *)new std::vector< float >();
  *(std::vector< float > **)&_swig_go_result = (std::vector< float > *)result; 
  return _swig_go_result;
}


std::vector< float > *_wrap_new_FloatVector__SWIG_1_routingkit_007dcf72057aca2d(long long _swig_go_0) {
  std::vector< float >::size_type arg1 ;
  std::vector< float > *result = 0 ;
  std::vector< float > *_swig_go_result;
  
  arg1 = (size_t)_swig_go_0; 
  
  result = (std::vector< float > *)new std::vector< float >(arg1);
  *(std::vector< float > **)&_swig_go_result = (std::vector< float > *)result; 
  return _swig_go_result;
}


std::vector< float > *_wrap_new_FloatVector__SWIG_2_routingkit_007dcf72057aca2d(std::vector< float > *_swig_go_0) {
  std::vector< float > *arg1 = 0 ;
  std::vector< float > *result = 0 ;
  std::vector< float > *_swig_go_result;
  
  arg1 = *(std::vector< float > **)&_swig_go_0; 
  
  result = (std::vector< float > *)new std::vector< float >((std::vector< float > const &)*arg1);
  *(std::vector< float > **)&_swig_go_result = (std::vector< float > *)result; 
  return _swig_go_result;
}


long long _wrap_FloatVector_size_routingkit_007dcf72057aca2d(std::vector< float > *_swig_go_0) {
  std::vector< float > *arg1 = (std::vector< float > *) 0 ;
  std::vector< float >::size_type result;
  long long _swig_go_result;
  
  arg1 = *(std::vector< float > **)&_swig_go_0; 
  
  result = ((std::vector< float > const *)arg1)->size();
  _swig_go_result = result; 
  return _swig_go_result;
}


long long _wrap_FloatVector_capacity_routingkit_007dcf72057aca2d(std::vector< float > *_swig_go_0) {
  std::vector< float > *arg1 = (std::vector< float > *) 0 ;
  std::vector< float >::size_type result;
  long long _swig_go_result;
  
  arg1 = *(std::vector< float > **)&_swig_go_0; 
  
  result = ((std::vector< float > const *)arg1)->capacity();
  _swig_go_result = result; 
  return _swig_go_result;
}


void _wrap_FloatVector_reserve_routingkit_007dcf72057aca2d(std::vector< float > *_swig_go_0, long long _swig_go_1) {
  std::vector< float > *arg1 = (std::vector< float > *) 0 ;
  std::vector< float >::size_type arg2 ;
  
  arg1 = *(std::vector< float > **)&_swig_go_0; 
  arg2 = (size_t)_swig_go_1; 
  
  (arg1)->reserve(arg2);
  
}


bool _wrap_FloatVector_isEmpty_routingkit_007dcf72057aca2d(std::vector< float > *_swig_go_0) {
  std::vector< float > *arg1 = (std::vector< float > *) 0 ;
  bool result;
  bool _swig_go_result;
  
  arg1 = *(std::vector< float > **)&_swig_go_0; 
  
  result = (bool)((std::vector< float > const *)arg1)->empty();
  _swig_go_result = result; 
  return _swig_go_result;
}


void _wrap_FloatVector_clear_routingkit_007dcf72057aca2d(std::vector< float > *_swig_go_0) {
  std::vector< float > *arg1 = (std::vector< float > *) 0 ;
  
  arg1 = *(std::vector< float > **)&_swig_go_0; 
  
  (arg1)->clear();
  
}


void _wrap_FloatVector_add_routingkit_007dcf72057aca2d(std::vector< float > *_swig_go_0, float _swig_go_1) {
  std::vector< float > *arg1 = (std::vector< float > *) 0 ;
  std::vector< float >::value_type *arg2 = 0 ;
  
  arg1 = *(std::vector< float > **)&_swig_go_0; 
  arg2 = (std::vector< float >::value_type *)&_swig_go_1; 
  
  (arg1)->push_back((std::vector< float >::value_type const &)*arg2);
  
}


float _wrap_FloatVector_get_routingkit_007dcf72057aca2d(std::vector< float > *_swig_go_0, intgo _swig_go_1) {
  std::vector< float > *arg1 = (std::vector< float > *) 0 ;
  int arg2 ;
  std::vector< float >::value_type *result = 0 ;
  float _swig_go_result;
  
  arg1 = *(std::vector< float > **)&_swig_go_0; 
  arg2 = (int)_swig_go_1; 
  
  try {
    result = (std::vector< float >::value_type *) &std_vector_Sl_float_Sg__get(arg1,arg2);
  } catch(std::out_of_range &_e) {
    (void)_e;
    _swig_gopanic("C++ std::out_of_range exception thrown");
    
  }
  _swig_go_result = (float)*result; 
  return _swig_go_result;
}


void _wrap_FloatVector_set_routingkit_007dcf72057aca2d(std::vector< float > *_swig_go_0, intgo _swig_go_1, float _swig_go_2) {
  std::vector< float > *arg1 = (std::vector< float > *) 0 ;
  int arg2 ;
  std::vector< float >::value_type *arg3 = 0 ;
  
  arg1 = *(std::vector< float > **)&_swig_go_0; 
  arg2 = (int)_swig_go_1; 
  arg3 = (std::vector< float >::value_type *)&_swig_go_2; 
  
  try {
    std_vector_Sl_float_Sg__set(arg1,arg2,(float const &)*arg3);
  } catch(std::out_of_range &_e) {
    (void)_e;
    _swig_gopanic("C++ std::out_of_range exception thrown");
    
  }
  
}


void _wrap_delete_FloatVector_routingkit_007dcf72057aca2d(std::vector< float > *_swig_go_0) {
  std::vector< float > *arg1 = (std::vector< float > *) 0 ;
  
  arg1 = *(std::vector< float > **)&_swig_go_0; 
  
  delete arg1;
  
}


std::vector< Point > *_wrap_new_PointVector__SWIG_0_routingkit_007dcf72057aca2d() {
  std::vector< Point > *result = 0 ;
  std::vector< Point > *_swig_go_result;
  
  
  result = (std::vector< Point > *)new std::vector< Point >();
  *(std::vector< Point > **)&_swig_go_result = (std::vector< Point > *)result; 
  return _swig_go_result;
}


std::vector< Point > *_wrap_new_PointVector__SWIG_1_routingkit_007dcf72057aca2d(long long _swig_go_0) {
  std::vector< Point >::size_type arg1 ;
  std::vector< Point > *result = 0 ;
  std::vector< Point > *_swig_go_result;
  
  arg1 = (size_t)_swig_go_0; 
  
  result = (std::vector< Point > *)new std::vector< Point >(arg1);
  *(std::vector< Point > **)&_swig_go_result = (std::vector< Point > *)result; 
  return _swig_go_result;
}


std::vector< Point > *_wrap_new_PointVector__SWIG_2_routingkit_007dcf72057aca2d(std::vector< Point > *_swig_go_0) {
  std::vector< Point > *arg1 = 0 ;
  std::vector< Point > *result = 0 ;
  std::vector< Point > *_swig_go_result;
  
  arg1 = *(std::vector< Point > **)&_swig_go_0; 
  
  result = (std::vector< Point > *)new std::vector< Point >((std::vector< Point > const &)*arg1);
  *(std::vector< Point > **)&_swig_go_result = (std::vector< Point > *)result; 
  return _swig_go_result;
}


long long _wrap_PointVector_size_routingkit_007dcf72057aca2d(std::vector< Point > *_swig_go_0) {
  std::vector< Point > *arg1 = (std::vector< Point > *) 0 ;
  std::vector< Point >::size_type result;
  long long _swig_go_result;
  
  arg1 = *(std::vector< Point > **)&_swig_go_0; 
  
  result = ((std::vector< Point > const *)arg1)->size();
  _swig_go_result = result; 
  return _swig_go_result;
}


long long _wrap_PointVector_capacity_routingkit_007dcf72057aca2d(std::vector< Point > *_swig_go_0) {
  std::vector< Point > *arg1 = (std::vector< Point > *) 0 ;
  std::vector< Point >::size_type result;
  long long _swig_go_result;
  
  arg1 = *(std::vector< Point > **)&_swig_go_0; 
  
  result = ((std::vector< Point > const *)arg1)->capacity();
  _swig_go_result = result; 
  return _swig_go_result;
}


void _wrap_PointVector_reserve_routingkit_007dcf72057aca2d(std::vector< Point > *_swig_go_0, long long _swig_go_1) {
  std::vector< Point > *arg1 = (std::vector< Point > *) 0 ;
  std::vector< Point >::size_type arg2 ;
  
  arg1 = *(std::vector< Point > **)&_swig_go_0; 
  arg2 = (size_t)_swig_go_1; 
  
  (arg1)->reserve(arg2);
  
}


bool _wrap_PointVector_isEmpty_routingkit_007dcf72057aca2d(std::vector< Point > *_swig_go_0) {
  std::vector< Point > *arg1 = (std::vector< Point > *) 0 ;
  bool result;
  bool _swig_go_result;
  
  arg1 = *(std::vector< Point > **)&_swig_go_0; 
  
  result = (bool)((std::vector< Point > const *)arg1)->empty();
  _swig_go_result = result; 
  return _swig_go_result;
}


void _wrap_PointVector_clear_routingkit_007dcf72057aca2d(std::vector< Point > *_swig_go_0) {
  std::vector< Point > *arg1 = (std::vector< Point > *) 0 ;
  
  arg1 = *(std::vector< Point > **)&_swig_go_0; 
  
  (arg1)->clear();
  
}


void _wrap_PointVector_add_routingkit_007dcf72057aca2d(std::vector< Point > *_swig_go_0, Point *_swig_go_1) {
  std::vector< Point > *arg1 = (std::vector< Point > *) 0 ;
  std::vector< Point >::value_type *arg2 = 0 ;
  
  arg1 = *(std::vector< Point > **)&_swig_go_0; 
  arg2 = *(std::vector< Point >::value_type **)&_swig_go_1; 
  
  (arg1)->push_back((std::vector< Point >::value_type const &)*arg2);
  
}


Point *_wrap_PointVector_get_routingkit_007dcf72057aca2d(std::vector< Point > *_swig_go_0, intgo _swig_go_1) {
  std::vector< Point > *arg1 = (std::vector< Point > *) 0 ;
  int arg2 ;
  std::vector< Point >::value_type *result = 0 ;
  Point *_swig_go_result;
  
  arg1 = *(std::vector< Point > **)&_swig_go_0; 
  arg2 = (int)_swig_go_1; 
  
  try {
    result = (std::vector< Point >::value_type *) &std_vector_Sl_Point_Sg__get(arg1,arg2);
  } catch(std::out_of_range &_e) {
    (void)_e;
    _swig_gopanic("C++ std::out_of_range exception thrown");
    
  }
  *(std::vector< Point >::value_type **)&_swig_go_result = result; 
  return _swig_go_result;
}


void _wrap_PointVector_set_routingkit_007dcf72057aca2d(std::vector< Point > *_swig_go_0, intgo _swig_go_1, Point *_swig_go_2) {
  std::vector< Point > *arg1 = (std::vector< Point > *) 0 ;
  int arg2 ;
  std::vector< Point >::value_type *arg3 = 0 ;
  
  arg1 = *(std::vector< Point > **)&_swig_go_0; 
  arg2 = (int)_swig_go_1; 
  arg3 = *(std::vector< Point >::value_type **)&_swig_go_2; 
  
  try {
    std_vector_Sl_Point_Sg__set(arg1,arg2,(Point const &)*arg3);
  } catch(std::out_of_range &_e) {
    (void)_e;
    _swig_gopanic("C++ std::out_of_range exception thrown");
    
  }
  
}


void _wrap_delete_PointVector_routingkit_007dcf72057aca2d(std::vector< Point > *_swig_go_0) {
  std::vector< Point > *arg1 = (std::vector< Point > *) 0 ;
  
  arg1 = *(std::vector< Point > **)&_swig_go_0; 
  
  delete arg1;
  
}


void _wrap_Point_lon_set_routingkit_007dcf72057aca2d(Point *_swig_go_0, float _swig_go_1) {
  Point *arg1 = (Point *) 0 ;
  float arg2 ;
  
  arg1 = *(Point **)&_swig_go_0; 
  arg2 = (float)_swig_go_1; 
  
  if (arg1) (arg1)->lon = arg2;
  
}


float _wrap_Point_lon_get_routingkit_007dcf72057aca2d(Point *_swig_go_0) {
  Point *arg1 = (Point *) 0 ;
  float result;
  float _swig_go_result;
  
  arg1 = *(Point **)&_swig_go_0; 
  
  result = (float) ((arg1)->lon);
  _swig_go_result = result; 
  return _swig_go_result;
}


void _wrap_Point_lat_set_routingkit_007dcf72057aca2d(Point *_swig_go_0, float _swig_go_1) {
  Point *arg1 = (Point *) 0 ;
  float arg2 ;
  
  arg1 = *(Point **)&_swig_go_0; 
  arg2 = (float)_swig_go_1; 
  
  if (arg1) (arg1)->lat = arg2;
  
}


float _wrap_Point_lat_get_routingkit_007dcf72057aca2d(Point *_swig_go_0) {
  Point *arg1 = (Point *) 0 ;
  float result;
  float _swig_go_result;
  
  arg1 = *(Point **)&_swig_go_0; 
  
  result = (float) ((arg1)->lat);
  _swig_go_result = result; 
  return _swig_go_result;
}


Point *_wrap_new_Point_routingkit_007dcf72057aca2d() {
  Point *result = 0 ;
  Point *_swig_go_result;
  
  
  result = (Point *)new Point();
  *(Point **)&_swig_go_result = (Point *)result; 
  return _swig_go_result;
}


void _wrap_delete_Point_routingkit_007dcf72057aca2d(Point *_swig_go_0) {
  Point *arg1 = (Point *) 0 ;
  
  arg1 = *(Point **)&_swig_go_0; 
  
  delete arg1;
  
}


void _wrap_QueryResponse_distance_set_routingkit_007dcf72057aca2d(QueryResponse *_swig_go_0, float _swig_go_1) {
  QueryResponse *arg1 = (QueryResponse *) 0 ;
  float arg2 ;
  
  arg1 = *(QueryResponse **)&_swig_go_0; 
  arg2 = (float)_swig_go_1; 
  
  if (arg1) (arg1)->distance = arg2;
  
}


float _wrap_QueryResponse_distance_get_routingkit_007dcf72057aca2d(QueryResponse *_swig_go_0) {
  QueryResponse *arg1 = (QueryResponse *) 0 ;
  float result;
  float _swig_go_result;
  
  arg1 = *(QueryResponse **)&_swig_go_0; 
  
  result = (float) ((arg1)->distance);
  _swig_go_result = result; 
  return _swig_go_result;
}


void _wrap_QueryResponse_waypoints_set_routingkit_007dcf72057aca2d(QueryResponse *_swig_go_0, std::vector< Point > *_swig_go_1) {
  QueryResponse *arg1 = (QueryResponse *) 0 ;
  std::vector< Point > *arg2 = (std::vector< Point > *) 0 ;
  
  arg1 = *(QueryResponse **)&_swig_go_0; 
  arg2 = *(std::vector< Point > **)&_swig_go_1; 
  
  if (arg1) (arg1)->waypoints = *arg2;
  
}


std::vector< Point > *_wrap_QueryResponse_waypoints_get_routingkit_007dcf72057aca2d(QueryResponse *_swig_go_0) {
  QueryResponse *arg1 = (QueryResponse *) 0 ;
  std::vector< Point > *result = 0 ;
  std::vector< Point > *_swig_go_result;
  
  arg1 = *(QueryResponse **)&_swig_go_0; 
  
  result = (std::vector< Point > *)& ((arg1)->waypoints);
  *(std::vector< Point > **)&_swig_go_result = (std::vector< Point > *)result; 
  return _swig_go_result;
}


QueryResponse *_wrap_new_QueryResponse_routingkit_007dcf72057aca2d() {
  QueryResponse *result = 0 ;
  QueryResponse *_swig_go_result;
  
  
  result = (QueryResponse *)new QueryResponse();
  *(QueryResponse **)&_swig_go_result = (QueryResponse *)result; 
  return _swig_go_result;
}


void _wrap_delete_QueryResponse_routingkit_007dcf72057aca2d(QueryResponse *_swig_go_0) {
  QueryResponse *arg1 = (QueryResponse *) 0 ;
  
  arg1 = *(QueryResponse **)&_swig_go_0; 
  
  delete arg1;
  
}


float _wrap_Client_distance_routingkit_007dcf72057aca2d(Client *_swig_go_0, float _swig_go_1, float _swig_go_2, float _swig_go_3, float _swig_go_4) {
  Client *arg1 = (Client *) 0 ;
  float arg2 ;
  float arg3 ;
  float arg4 ;
  float arg5 ;
  float result;
  float _swig_go_result;
  
  arg1 = *(Client **)&_swig_go_0; 
  arg2 = (float)_swig_go_1; 
  arg3 = (float)_swig_go_2; 
  arg4 = (float)_swig_go_3; 
  arg5 = (float)_swig_go_4; 
  
  result = (float)(arg1)->distance(arg2,arg3,arg4,arg5);
  _swig_go_result = result; 
  return _swig_go_result;
}


QueryResponse *_wrap_Client_queryrequest_routingkit_007dcf72057aca2d(Client *_swig_go_0, float _swig_go_1, float _swig_go_2, float _swig_go_3, float _swig_go_4, float _swig_go_5) {
  Client *arg1 = (Client *) 0 ;
  float arg2 ;
  float arg3 ;
  float arg4 ;
  float arg5 ;
  float arg6 ;
  QueryResponse result;
  QueryResponse *_swig_go_result;
  
  arg1 = *(Client **)&_swig_go_0; 
  arg2 = (float)_swig_go_1; 
  arg3 = (float)_swig_go_2; 
  arg4 = (float)_swig_go_3; 
  arg5 = (float)_swig_go_4; 
  arg6 = (float)_swig_go_5; 
  
  result = (arg1)->queryrequest(arg2,arg3,arg4,arg5,arg6);
  *(QueryResponse **)&_swig_go_result = new QueryResponse(result); 
  return _swig_go_result;
}


std::vector< float > *_wrap_Client_table_routingkit_007dcf72057aca2d(Client *_swig_go_0, std::vector< Point > *_swig_go_1, std::vector< Point > *_swig_go_2) {
  Client *arg1 = (Client *) 0 ;
  std::vector< Point > arg2 ;
  std::vector< Point > arg3 ;
  std::vector< Point > *argp2 ;
  std::vector< Point > *argp3 ;
  std::vector< float > result;
  std::vector< float > *_swig_go_result;
  
  arg1 = *(Client **)&_swig_go_0; 
  
  argp2 = (std::vector< Point > *)_swig_go_1;
  if (argp2 == NULL) {
    _swig_gopanic("Attempt to dereference null std::vector< Point >");
  }
  arg2 = (std::vector< Point >)*argp2;
  
  
  argp3 = (std::vector< Point > *)_swig_go_2;
  if (argp3 == NULL) {
    _swig_gopanic("Attempt to dereference null std::vector< Point >");
  }
  arg3 = (std::vector< Point >)*argp3;
  
  
  result = (arg1)->table(arg2,arg3);
  *(std::vector< float > **)&_swig_go_result = new std::vector< float >(result); 
  return _swig_go_result;
}


void _wrap_Client_build_ch_routingkit_007dcf72057aca2d(Client *_swig_go_0, _gostring_ _swig_go_1, _gostring_ _swig_go_2) {
  Client *arg1 = (Client *) 0 ;
  char *arg2 = (char *) 0 ;
  char *arg3 = (char *) 0 ;
  
  arg1 = *(Client **)&_swig_go_0; 
  
  arg2 = (char *)malloc(_swig_go_1.n + 1);
  memcpy(arg2, _swig_go_1.p, _swig_go_1.n);
  arg2[_swig_go_1.n] = '\0';
  
  
  arg3 = (char *)malloc(_swig_go_2.n + 1);
  memcpy(arg3, _swig_go_2.p, _swig_go_2.n);
  arg3[_swig_go_2.n] = '\0';
  
  
  (arg1)->build_ch(arg2,arg3);
  
  free(arg2); 
  free(arg3); 
}


void _wrap_Client_load_routingkit_007dcf72057aca2d(Client *_swig_go_0, _gostring_ _swig_go_1, _gostring_ _swig_go_2) {
  Client *arg1 = (Client *) 0 ;
  char *arg2 = (char *) 0 ;
  char *arg3 = (char *) 0 ;
  
  arg1 = *(Client **)&_swig_go_0; 
  
  arg2 = (char *)malloc(_swig_go_1.n + 1);
  memcpy(arg2, _swig_go_1.p, _swig_go_1.n);
  arg2[_swig_go_1.n] = '\0';
  
  
  arg3 = (char *)malloc(_swig_go_2.n + 1);
  memcpy(arg3, _swig_go_2.p, _swig_go_2.n);
  arg3[_swig_go_2.n] = '\0';
  
  
  (arg1)->load(arg2,arg3);
  
  free(arg2); 
  free(arg3); 
}


double _wrap_Client_average_routingkit_007dcf72057aca2d(Client *_swig_go_0, std::vector< int > *_swig_go_1) {
  Client *arg1 = (Client *) 0 ;
  std::vector< int > arg2 ;
  std::vector< int > *argp2 ;
  double result;
  double _swig_go_result;
  
  arg1 = *(Client **)&_swig_go_0; 
  
  argp2 = (std::vector< int > *)_swig_go_1;
  if (argp2 == NULL) {
    _swig_gopanic("Attempt to dereference null std::vector< int >");
  }
  arg2 = (std::vector< int >)*argp2;
  
  
  result = (double)(arg1)->average(arg2);
  _swig_go_result = result; 
  return _swig_go_result;
}


Client *_wrap_new_Client_routingkit_007dcf72057aca2d() {
  Client *result = 0 ;
  Client *_swig_go_result;
  
  
  result = (Client *)new Client();
  *(Client **)&_swig_go_result = (Client *)result; 
  return _swig_go_result;
}


void _wrap_delete_Client_routingkit_007dcf72057aca2d(Client *_swig_go_0) {
  Client *arg1 = (Client *) 0 ;
  
  arg1 = *(Client **)&_swig_go_0; 
  
  delete arg1;
  
}


#ifdef __cplusplus
}
#endif

