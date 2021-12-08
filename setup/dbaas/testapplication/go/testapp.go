//   Copyright (c) 2019 AT&T Intellectual Property.
//   Copyright (c) 2019 Nokia.
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.

//
//   This source code is part of the near-RT RIC (RAN Intelligent Controller)
//   platform project (RICP).
//

package main

import (
    "fmt"
    "./sdl"
)

func main() {
    sdl1 := sdl.Create("test1")

    var err error

    err = sdl1.Set("key1", "data1", "key2", "data2")
    if err != nil {
        fmt.Printf("unable to write to DB\n")
    }

    err = sdl1.Set("num1", 1, "num2", 2)
    if err != nil {
        fmt.Printf("unable to write to DB\n")
    }

    d := make([]byte, 3)
    d[0] = 1
    d[1] = 2
    d[2] = 3
    err = sdl1.Set("arr1", d)
    if err != nil {
        fmt.Printf("unable to write to DB\n")
    }

    p := []string{"pair1", "data1", "pair2", "data2"}
    err = sdl1.Set(p)
    if err != nil {
        fmt.Printf("unable to write to DB\n")
    }

    a := [4]string{"array1", "adata1", "array2", "adata2"}
    err = sdl1.Set(a)
    if err != nil {
        fmt.Printf("unable to write to DB\n")
    }

    mix1 := []interface{}{"mix1", "data1", "mix2", 2}
    err = sdl1.Set(mix1)
    if err != nil {
        fmt.Printf("unable to write to DB\n")
    }

    mix2 := [4]interface{}{"mix3", "data3", "mix4", 4}
    err = sdl1.Set(mix2)
    if err != nil {
        fmt.Printf("unable to write to DB\n")
    }

    retDataMap, err := sdl1.Get([]string{"key1", "key3", "key2"})
    if err != nil {
        fmt.Printf("Unable to read from DB\n")
    } else {
        for i, v := range retDataMap {
            fmt.Printf("%s:%s\n", i, v)
        }
    }

    retDataMap2, err := sdl1.Get([]string{"num1", "num2"})
    if err != nil {
        fmt.Printf("Unable to read from DB\n")
    } else {
        for i, v := range retDataMap2 {
            fmt.Printf("%s:%s\n", i, v)
        }
    }

    fmt.Println("-------------")
    allKeys := []string{"key1", "key2", "num1", "num2", "pair1", "pair2", "array1", "array2", "mix1", "mix2", "mix3", "mix4", "arr1"}
    retDataMap3, err := sdl1.Get(allKeys)
    if err != nil {
        fmt.Printf("Unable to read from DB\n")
    } else {
        for i3, v3 := range retDataMap3 {
            fmt.Printf("%s:%s\n", i3, v3)
        }
    }
}
