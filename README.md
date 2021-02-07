# dynjson

[![Godoc Reference](https://godoc.org/github.com/nikitaksv/strcase?status.svg)](http://godoc.org/github.com/nikitaksv/strcase)
[![Coverage Status](https://coveralls.io/repos/github/nikitaksv/strcase/badge.svg?branch=main)](https://coveralls.io/github/nikitaksv/strcase?branch=main)

Parsing JSON to dynamic json structure with sorted fields

## Features

* Save json order fields
* Parse data to dynamic json structure

## Install

```sh
go get -u github.com/nikitaksv/dynjson
```

## Usage

```go
package main

import (
	"fmt"
	"log"

	"github.com/nikitaksv/dynjson"
)

func main() {
	bs := []byte(`
{
  "a1": "one",
  "a2": -10,
  "a3": "ahead",
  "a4": false,
  "a5": 17,
  "a6": "wet",
  "a7": {
    "b1": "one",
    "b2": "two",
    "b3": false,
    "b4": 4,
    "b5": [
      "five"
    ],
    "b6": true
  },
  "a8": null
}
`)

	j := &dynjson.Json{}
	err := j.UnmarshalJSON(bs)
	if err != nil {
		log.Fatal(err)
	}

	if rootObj, ok := j.Value.(dynjson.Object); ok {
		for originJsonIndex, prop := range rootObj.Properties {
			fmt.Printf("RootObj - Index: %d, Key: %s, Value: %+v", originJsonIndex, prop.Key, prop.Value)

			if nestObj, ok := prop.Value.(dynjson.Object); ok {
				for nestedOriginJsonIndex, propNestObj := range nestObj.Properties {
					fmt.Printf("NestObj - Index: %d, Key: %s, Value: %+v", nestedOriginJsonIndex, propNestObj.Key, propNestObj.Value)
				}
			}
		}
	}
}
```
 