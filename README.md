[![Go Reference](https://pkg.go.dev/badge/github.com/goradd/maps.svg)](https://pkg.go.dev/github.com/goradd/maps)
![Build Status](https://img.shields.io/github/workflow/status/goradd/maps/Go)
[![Go Report Card](https://goreportcard.com/badge/github.com/goradd/maps)](https://goreportcard.com/report/github.com/goradd/maps)
[![codecov](https://codecov.io/gh/goradd/maps/branch/main/graph/badge.svg?token=LZNNI26H3L)](https://codecov.io/gh/goradd/maps)

# maps
maps is a library using Go generics that offers a standard interface for manipulating 
different kinds of maps. 

Using the same interface, you can create and use a standard Go map, a map
that is safe for concurrency and/or a map that lets you order the keys in the map.

## Example

```go
package main

import . "github.com/goradd/maps"
import "fmt"

type myMap Map[string,int]
type myStdMap StdMap[string, int]

func main() {
	m := new(Map[string, int])
	
	m.Merge(myStdMap{"b":2, "c":3})
	m.Set("a",1)

	sum := 0
	m.Range(func(k string, v int) bool {
		sum += v
		return true
    })
	fmt.Print(sum)
}

```

By simply changing myMap to a SafeMap, you can make the map safe for concurrent use.
Or, you can change myMap to a SliceMap, or a SafeSliceMap to also be able to iterate
the map in the order it was created, similar to a PHP map.
