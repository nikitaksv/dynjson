# dynjson

[![Godoc Reference](https://godoc.org/github.com/nikitaksv/dynjson?status.svg)](http://godoc.org/github.com/nikitaksv/dynjson)
[![Coverage Status](https://coveralls.io/repos/github/nikitaksv/dynjson/badge.svg?branch=main)](https://coveralls.io/github/nikitaksv/dynjson?branch=main)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fnikitaksv%2Fdynjson.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Fnikitaksv%2Fdynjson?ref=badge_shield)

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

	if rootObj, ok := j.Value.(*dynjson.Object); ok {
		for originJsonIndex, prop := range rootObj.Properties {
			fmt.Printf("root[%d] Key: %s, Value: %+v\n", originJsonIndex, prop.Key, prop.Value)

			if nestObj, ok := prop.Value.(*dynjson.Object); ok {
				for nestedOriginJsonIndex, propNestObj := range nestObj.Properties {
					fmt.Printf("\t%s[%d] Key: %s, Value: %+v\n", prop.Key, nestedOriginJsonIndex, propNestObj.Key, propNestObj.Value)
				}
			}
		}
	}
	/*
		root[0] Key: a1, Value: one
		root[1] Key: a2, Value: -10
		root[2] Key: a3, Value: ahead
		root[3] Key: a4, Value: false
		root[4] Key: a5, Value: 17
		root[5] Key: a6, Value: wet
		root[6] Key: a7, Value: &{Key:a7 Properties:[0xc000004720 0xc000004780 0xc0000047e0 0xc000004860 0xc0000048c0 0xc000004920]}
			a7[0] Key: b1, Value: one
			a7[1] Key: b2, Value: two
			a7[2] Key: b3, Value: false
			a7[3] Key: b4, Value: 4
			a7[4] Key: b5, Value: &{Elements:[five]}
			a7[5] Key: b6, Value: true
		root[7] Key: a8, Value: <nil>

		Process finished with exit code 0
	*/
}
```
 

## License
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fnikitaksv%2Fdynjson.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Fnikitaksv%2Fdynjson?ref=badge_large)