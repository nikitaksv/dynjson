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
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

func BenchmarkJson_UnmarshalJSON(b *testing.B) {
	data := []byte(`{"a1":[{"b1":"one","b2":"two","b3":false,"b4":4,"b5":["five"],"b6":true}],"a2":-10,"a3":"ahead","a4":false,"a5":17,"a6":"wet","a7":null}`)
	b.Run("Unmarshall all data types", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			j := &Json{}
			err := j.UnmarshalJSON(data)
			if err != nil {
				b.Errorf("UnmarshalJSON error = %v", err)

			}
		}
	})
}

func TestJson_UnmarshalJSON(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		want    *Json
		args    args
		wantErr bool
	}{
		{
			name: "one field",
			want: &Json{
				Value: &Object{
					Key: "root",
					Properties: []*Property{
						{
							Key:   "a1",
							Value: true,
						},
					},
				},
			},
			args:    args{data: []byte(`{"a1":true}`)},
			wantErr: false,
		},
		{
			name: "nested object",
			want: &Json{
				Value: &Object{
					Key: "root",
					Properties: []*Property{
						{
							Key: "a1",
							Value: &Object{
								Key: "a1",
								Properties: []*Property{
									{
										Key:   "b1",
										Value: true,
									},
								},
							},
						},
					},
				},
			},
			args:    args{data: []byte(`{"a1":{"b1":true}}`)},
			wantErr: false,
		},
		{
			name: "nested array",
			want: &Json{
				Value: &Object{
					Key: "root",
					Properties: []*Property{
						{
							Key: "a1",
							Value: &Array{
								Elements: []interface{}{
									&Object{
										Key: "0",
										Properties: []*Property{
											{
												Key:   "b1",
												Value: true,
											},
										},
									},
								},
							},
						},
					},
				},
			},
			args:    args{data: []byte(`{"a1":[{"b1":true}]}`)},
			wantErr: false,
		},
		{
			name: "root array",
			want: &Json{
				Value: &Array{
					Elements: []interface{}{
						&Object{
							Key: "0",
							Properties: []*Property{
								{
									Key:   "b1",
									Value: true,
								},
							},
						},
					},
				},
			},
			args:    args{data: []byte(`[{"b1":true}]`)},
			wantErr: false,
		},
		{
			name: "check error equal",
			want: &Json{
				Value: &Array{
					Elements: []interface{}{
						&Object{
							Key: "0",
							Properties: []*Property{
								{
									Key:   "b1",
									Value: "error value",
								},
							},
						},
					},
				},
			},
			args:    args{data: []byte(`[{"b1":true}]`)},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &Json{}
			if err := j.UnmarshalJSON(tt.args.data); err != nil {
				t.Errorf("UnmarshalJSON() error = %v", err)
			}

			if !reflect.DeepEqual(j, tt.want) && !tt.wantErr {
				t.Errorf("DeepEqual json = %s not equal want = %s, wantErr %v", j, tt.want, tt.wantErr)
			}
		})
	}
}

func TestJson(t *testing.T) {

	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		want    []byte
		args    args
		wantErr bool
	}{
		{
			name: "check object sorted json",
			want: []byte(`{"value":{"key":"root","properties":[{"key":"a1","value":{"elements":[{"key":"0","properties":[{"key":"b1","value":"one"},{"key":"b2","value":"two"},{"key":"b3","value":false},{"key":"b4","value":4},{"key":"b5","value":{"elements":["five"]}},{"key":"b6","value":true}]}]}},{"key":"a2","value":-10},{"key":"a3","value":"ahead"},{"key":"a4","value":false},{"key":"a5","value":17},{"key":"a6","value":"wet"},{"key":"a7","value":null}]}}`),
			args: args{
				data: []byte(`{"a1":[{"b1":"one","b2":"two","b3":false,"b4":4,"b5":["five"],"b6":true}],"a2":-10,"a3":"ahead","a4":false,"a5":17,"a6":"wet","a7":null}`),
			},
			wantErr: false,
		},
		{
			name: "check array sorted json",
			want: []byte(`{"value":{"elements":[{"key":"0","properties":[{"key":"a1","value":{"elements":[{"key":"0","properties":[{"key":"b1","value":"one"},{"key":"b2","value":"two"},{"key":"b3","value":false},{"key":"b4","value":4},{"key":"b5","value":{"elements":["five"]}},{"key":"b6","value":true}]}]}},{"key":"a2","value":-10},{"key":"a3","value":"ahead"},{"key":"a4","value":false},{"key":"a5","value":17},{"key":"a6","value":"wet"},{"key":"a7","value":null}]}]}}`),
			args: args{
				data: []byte(`[{"a1":[{"b1":"one","b2":"two","b3":false,"b4":4,"b5":["five"],"b6":true}],"a2":-10,"a3":"ahead","a4":false,"a5":17,"a6":"wet","a7":null}]`),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &Json{}
			if err := j.UnmarshalJSON(tt.args.data); err != nil {
				t.Errorf("UnmarshalJSON() error = %v", err)
			}
			bs, err := json.Marshal(j)
			if err != nil {
				t.Errorf("MarshalJSON() error = %v", err)
			}

			actualStr := string(bs)
			wantStr := string(tt.want)
			if strings.Compare(actualStr, wantStr) != 0 && !tt.wantErr {
				t.Errorf("Compare() actual = %v not equal want = %v, wantErr %v", actualStr, wantStr, tt.wantErr)
			}
		})
	}
}
