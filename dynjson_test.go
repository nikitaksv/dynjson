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
	"strings"
	"testing"
)

func BenchmarkJson_UnmarshalJSON(b *testing.B) {
	data := []byte(`{"a1":[{"b1":"one","b2":"two","b3":false,"b4":4,"b5":["five"],"b6":true}],"a2":-10,"a3":"ahead","a4":false,"a5":17.2,"a6":"wet","a7":null,"a8":{"c1":"cc"}}`)
	b.Run("Unmarshall all data types", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			j := &Json{}
			err := json.Unmarshal(data, j)
			if err != nil {
				b.Fatalf("UnmarshalJSON error = %v", err)
			}
		}
	})
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
			want: []byte(`{"a1":[{"b1":"one","b2":"two","b3":false,"b4":4,"b5":["five"],"b6":true}],"a2":-10,"a3":"ahead","a4":false,"a5":17.2,"a6":"wet","a7":null,"a8":{"c1":"cc"}}`),
			args: args{
				data: []byte(`{"a1":[{"b1":"one","b2":"two","b3":false,"b4":4,"b5":["five"],"b6":true}],"a2":-10,"a3":"ahead","a4":false,"a5":17.2,"a6":"wet","a7":null,"a8":{"c1":"cc"}}`),
			},
			wantErr: false,
		},
		{
			name: "check array sorted json",
			want: []byte(`[{"a1":[{"b1":"one","b2":"two","b3":false,"b4":4,"b5":["five"],"b6":true}],"a2":-10,"a3":"ahead","a4":false,"a5":17,"a6":"wet","a7":null}]`),
			args: args{
				data: []byte(`[{"a1":[{"b1":"one","b2":"two","b3":false,"b4":4,"b5":["five"],"b6":true}],"a2":-10,"a3":"ahead","a4":false,"a5":17,"a6":"wet","a7":null}]`),
			},
			wantErr: false,
		},
		{
			name: "first string",
			want: []byte(`"qwerty"`),
			args: args{
				data: []byte(`"qwerty"`),
			},
			wantErr: false,
		},
		{
			name: "first array of string",
			want: []byte(`["qwerty"]`),
			args: args{
				data: []byte(`["qwerty"]`),
			},
			wantErr: false,
		},
		{
			name: "incorrect json",
			args: args{
				data: []byte(`qwerty`),
			},
			wantErr: true,
		},
		{
			name: "empty json",
			args: args{
				data: []byte(``),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &Json{}
			if err := json.Unmarshal(tt.args.data, j); err != nil && !tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v", err)
			}
			bs, err := json.Marshal(j)
			if err != nil && !tt.wantErr {
				t.Errorf("MarshalJSON() error = %v", err)
			}

			actualStr := string(bs)
			wantStr := string(tt.want)
			if strings.Compare(actualStr, wantStr) != 0 && !tt.wantErr {
				t.Errorf("Compare() actual = %s not equal want = %s, wantErr %v", actualStr, wantStr, tt.wantErr)
			}
		})
	}
}
