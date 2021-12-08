/*
==================================================================================
        Copyright (c) 2018-2019 AT&T Intellectual Property.

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
==================================================================================
*/


#pragma once
#ifndef GENERIC_HELPERS
#define GENERIC_HELPERS

#include <cstddef>

/* Utilities */

class octet_helper {

public:
  octet_helper(void):_ref(NULL), _size(0){};
  octet_helper(const void *ref, int size):_ref(ref), _size(size){};
  void set_ref(const void *ref){
    _ref = ref;
  }
  
  void set_size(size_t size){
    _size = size;
  }
  
  const void * get_ref(void){return _ref ; };
  size_t get_size(void) const {return _size ; } ;

private:
  const void *_ref;
  size_t _size;
};
    
#endif
