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
)

type Json struct {
	// Value typeOf Object or Array or JSON-Type
	Value      interface{}
	escapeHTML bool
}

func (j *Json) SetEscapeHTML(escapeHTML bool) {
	j.escapeHTML = escapeHTML
}

func (j *Json) UnmarshalJSON(data []byte) error {
	return unmarshalJson(j, json.NewDecoder(bytes.NewReader(data)))
}

func (j *Json) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(j.escapeHTML)

	if err := marshalValue(enc, &buf, j.Value); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type Array struct {
	Elements []interface{}
}

type Object struct {
	Properties []*Property
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
	Key   string
	Value interface{}
}

func unmarshalJson(j *Json, dec *json.Decoder) error {
	// first token
	token, err := dec.Token()
	if err != nil {
		return err
	}

	if delim, ok := token.(json.Delim); ok {
		if delim == '{' {
			obj := &Object{Properties: []*Property{}}
			if err = unmarshalObject(obj, dec); err != nil {
				return err
			}
			j.Value = obj
		} else if delim == '[' {
			arr := &Array{}
			if err = unmarshalArray(arr, dec); err != nil {
				return err
			}
			j.Value = arr
		}
	} else {
		j.Value = token
	}

	return nil
}

func unmarshalObject(obj *Object, dec *json.Decoder) error {
	for {
		token, err := dec.Token()
		if err != nil {
			return err
		}

		if delim, ok := token.(json.Delim); ok && delim == '}' {
			return nil
		}

		prop := &Property{Key: token.(string)}

		token, err = dec.Token()
		if err != nil {
			return err
		}

		if delim, ok := token.(json.Delim); ok {
			switch delim {
			case '{': // this is an object
				nestObj := &Object{
					Properties: []*Property{},
				}
				if err = unmarshalObject(nestObj, dec); err != nil {
					return err
				}
				prop.Value = nestObj
			case '[': // this is an array
				nestArr := &Array{
					Elements: []interface{}{},
				}
				if err = unmarshalArray(nestArr, dec); err != nil {
					return err
				}
				prop.Value = nestArr
			case '}', ']':
				return nil
			}
		} else {
			// this is value
			prop.Value = token
		}

		obj.Properties = append(obj.Properties, prop)
	}
}

func unmarshalArray(arr *Array, dec *json.Decoder) error {
	for index := 0; ; index++ {
		var elem interface{}

		token, err := dec.Token()
		if err != nil {
			return err
		}

		if delim, ok := token.(json.Delim); ok {
			switch delim {
			case '{':
				nestObj := &Object{
					Properties: []*Property{},
				}
				if err = unmarshalObject(nestObj, dec); err != nil {
					return err
				}
				elem = nestObj
			case '[':
				nestArr := &Array{
					Elements: []interface{}{},
				}
				if err = unmarshalArray(nestArr, dec); err != nil {
					return err
				}
				elem = nestArr
			case '}', ']':
				return nil
			}
		} else {
			// this is value
			elem = token
		}

		arr.Elements = append(arr.Elements, elem)
	}
}

func marshalValue(enc *json.Encoder, buf *bytes.Buffer, value interface{}) error {
	if obj, ok := value.(*Object); ok {
		if err := marshalObject(enc, buf, obj); err != nil {
			return err
		}
	} else if arr, ok := value.(*Array); ok {
		if err := marshalArray(enc, buf, arr); err != nil {
			return err
		}
	} else if err := enc.Encode(value); err != nil {
		return err
	}

	return nil
}

func marshalObject(enc *json.Encoder, buf *bytes.Buffer, obj *Object) error {
	buf.WriteByte('{')
	defer buf.WriteByte('}')

	for i, prop := range obj.Properties {
		if i > 0 {
			buf.WriteByte(',')
		}
		if err := enc.Encode(prop.Key); err != nil {
			return err
		}
		buf.WriteByte(':')
		if err := marshalValue(enc, buf, prop.Value); err != nil {
			return err
		}
	}

	return nil
}

func marshalArray(enc *json.Encoder, buf *bytes.Buffer, arr *Array) error {
	buf.WriteByte('[')
	defer buf.WriteByte(']')

	for i, el := range arr.Elements {
		if i > 0 {
			buf.WriteByte(',')
		}
		if err := marshalValue(enc, buf, el); err != nil {
			return err
		}
	}

	return nil
}
