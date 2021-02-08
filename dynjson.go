/*
 * Copyright (c) 2021 Nikita Krasnikov
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package dynjson

import (
	"bytes"
	"encoding/json"
	"errors"
	"strconv"
)

var (
	ErrFirstTokenJson = errors.New("incorrect json first token")
)

type Json struct {
	// Value typeOf Object or Array
	Value interface{} `json:"value"`
}

func (j *Json) UnmarshalJSON(data []byte) error {
	return parseJson(j, json.NewDecoder(bytes.NewReader(data)))
}

type Array struct {
	Elements []interface{} `json:"elements"`
}

type Object struct {
	Key        string      `json:"key"`
	Properties []*Property `json:"properties"`
}

func (o Object) GetProperty(key string) (*Property, bool) {
	for _, p := range o.Properties {
		if p.Key == key {
			return p, true
		}
	}

	return nil, false
}

type Property struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

func parseJson(j *Json, dec *json.Decoder) error {
	// first token
	token, err := dec.Token()
	if err != nil {
		return err
	}

	if delim, ok := token.(json.Delim); ok {
		if delim == '{' {
			obj := &Object{Key: "root", Properties: []*Property{}}
			parseMap(obj, dec)
			j.Value = obj
		} else if delim == '[' {
			arr := &Array{}
			parseArray(arr, dec)
			j.Value = arr
		}
	} else {
		return ErrFirstTokenJson
	}

	return nil
}

func parseMap(obj *Object, dec *json.Decoder) {
	for {
		prop := &Property{
			Key:   "",
			Value: nil,
		}

		token, err := dec.Token()
		if err != nil {
			return
		}

		if delim, ok := token.(json.Delim); ok && delim == '}' {
			return
		}

		key := token.(string)
		prop.Key = key

		token, err = dec.Token()
		if err != nil {
			return
		}

		if delim, ok := token.(json.Delim); ok {
			switch delim {
			case '{': // this a object
				nestObj := &Object{
					Key:        prop.Key,
					Properties: []*Property{},
				}
				parseMap(nestObj, dec)
				prop.Value = nestObj
			case '[': // this a array
				nestArr := &Array{
					Elements: []interface{}{},
				}
				parseArray(nestArr, dec)
				prop.Value = nestArr
			case '}', ']':
				return
			}
		} else {
			// this a value
			prop.Value = token
		}

		obj.Properties = append(obj.Properties, prop)
	}
}

func parseArray(arr *Array, dec *json.Decoder) {
	for index := 0; ; index++ {
		var elem interface{}

		token, err := dec.Token()
		if err != nil {
			return
		}

		if delim, ok := token.(json.Delim); ok {
			switch delim {
			case '{':
				nestObj := &Object{
					Key:        strconv.Itoa(index),
					Properties: []*Property{},
				}
				parseMap(nestObj, dec)
				elem = nestObj
			case '[':
				nestArr := &Array{
					Elements: []interface{}{},
				}
				parseArray(nestArr, dec)
				elem = nestArr
			case '}', ']':
				return
			}
		} else {
			// this a value
			elem = token
		}

		arr.Elements = append(arr.Elements, elem)
	}
}
