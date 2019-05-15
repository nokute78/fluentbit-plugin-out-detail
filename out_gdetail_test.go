/*
   Copyright 2019 Takahiro Yamashita

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package main

import (
	"bytes"
	"encoding/json"
	"github.com/nokute78/msgpack-microscope/pkg/msgpack"
	"testing"
)

type MPBase struct {
	Format string `json:"format"`
	Byte   string `json:"header"`
	Raw    string `json:"raw"`
}

type MPString struct {
	MPBase
	Value string `json:"value"`
}

func TestVerboseJSONString(t *testing.T) {
	type testcase struct {
		casename string
		bytes    []byte
		expected string
	}

	cases := []testcase{
		{"fixstr", []byte{0xa2, 0x41, 0x42}, `AB`},
		{"str8", []byte{0xd9, 0x0f, 0xe3, 0x81, 0x93, 0xe3, 0x82, 0x93, 0xe3, 0x81, 0xab, 0xe3, 0x81, 0xa1, 0xe3, 0x81, 0xaf}, `こんにちは`},
		{"bin8", []byte{0xc4, 0x04, 0xde, 0xad, 0xbe, 0xef}, "0xdeadbeef"},
	}

	buf := bytes.Buffer{}
	for _, v := range cases {
		buf.Reset()
		ret, err := msgpack.Decode(bytes.NewBuffer(v.bytes))
		outputVerboseJSON(ret, &buf, 0)

		p := MPString{}
		err = json.Unmarshal(buf.Bytes(), &p)
		if err != nil {
			t.Errorf("%s: Unmarshal Error %s", v.casename, err)
		}
		if v.expected != p.Value {
			t.Errorf("%s: mismatch. given: %s. expected: %s", v.casename, p.Value, v.expected)
		}
	}
}

type MPExt struct {
	MPBase
	Type  int8   `json:"type"`
	Value string `json:"value"`
}

func TestVerboseJSONExt(t *testing.T) {
	type testcase struct {
		casename string
		bytes    []byte
		expected string
	}

	cases := []testcase{
		{"fixext1", []byte{0xd4, 0x01, 0xff}, "0xff"},
		{"fixext2", []byte{0xd5, 0x01, 0xfe, 0xed}, "0xfeed"},
		{"fixext4", []byte{0xd6, 0x01, 0xde, 0xad, 0xbe, 0xef}, "0xdeadbeef"},
		{"fixext8", []byte{0xd7, 0x01, 0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef}, "0xdeadbeefdeadbeef"},
	}

	buf := bytes.Buffer{}
	for _, v := range cases {
		buf.Reset()
		ret, err := msgpack.Decode(bytes.NewBuffer(v.bytes))
		outputVerboseJSON(ret, &buf, 0)

		p := MPExt{}
		err = json.Unmarshal(buf.Bytes(), &p)
		if err != nil {
			t.Errorf("%s: Unmarshal Error %s", v.casename, err)
		}
		if v.expected != p.Value {
			t.Errorf("%s: mismatch. given: %s. expected: %s", v.casename, p.Value, v.expected)
		}
	}
}

type MPBool struct {
	MPBase
	Value bool `json:"value"`
}

func TestVerboseJSONBool(t *testing.T) {
	type testcase struct {
		casename string
		bytes    []byte
		expected bool
	}

	cases := []testcase{
		{"true", []byte{0xc3}, true},
		{"false", []byte{0xc2}, false},
	}

	buf := bytes.Buffer{}
	for _, v := range cases {
		buf.Reset()
		ret, err := msgpack.Decode(bytes.NewBuffer(v.bytes))
		outputVerboseJSON(ret, &buf, 0)

		p := MPBool{}
		err = json.Unmarshal(buf.Bytes(), &p)
		if err != nil {
			t.Errorf("%s: Unmarshal Error %s", v.casename, err)
		}
		if v.expected != p.Value {
			t.Errorf("%s: mismatch. given: %t. expected: %t", v.casename, p.Value, v.expected)
		}
	}
}

type MPNil struct {
	MPBase
	Value *bool `json:"value"`
}

func TestVerboseJSONNil(t *testing.T) {
	b := []byte{0xc0}
	buf := bytes.Buffer{}

	ret, err := msgpack.Decode(bytes.NewBuffer(b))
	outputVerboseJSON(ret, &buf, 0)

	p := MPNil{}
	err = json.Unmarshal(buf.Bytes(), &p)
	if err != nil {
		t.Errorf("Nil: Unmarshal Error %s", err)
	}
	if p.Value != nil {
		t.Errorf("Nil: Value is not nil")
	}
}

type MPInt struct {
	MPBase
	Value int64 `json:"value"`
}

func TestVerboseJSONInt(t *testing.T) {
	type testcase struct {
		casename string
		bytes    []byte
		expected int64
	}

	cases := []testcase{
		{"p fixint", []byte{0x01}, 1},
		{"n fixint", []byte{0xff}, -1},
		{"int8", []byte{0xd0, 0xff}, -1},
		{"int16", []byte{0xd1, 0xff, 0x00}, -256},
		{"int32", []byte{0xd2, 0xff, 0x00, 0xff, 0x00}, -16711936},
		{"int64", []byte{0xd3, 0xff, 0x00, 0xff, 0x00, 0xff, 0x00, 0xff, 0x00}, -71777214294589696},
	}

	buf := bytes.Buffer{}
	for _, v := range cases {
		buf.Reset()
		ret, err := msgpack.Decode(bytes.NewBuffer(v.bytes))
		outputVerboseJSON(ret, &buf, 0)

		p := MPInt{}
		err = json.Unmarshal(buf.Bytes(), &p)
		if err != nil {
			t.Errorf("%s: Unmarshal Error %s", v.casename, err)
		}
		if v.expected != p.Value {
			t.Errorf("%s: mismatch. given: %d. expected: %d", v.casename, p.Value, v.expected)
		}
	}
}

type MPUint struct {
	MPBase
	Value uint64 `json:"value"`
}

func TestVerboseJSONUint(t *testing.T) {
	type testcase struct {
		casename string
		bytes    []byte
		expected uint64
	}

	cases := []testcase{
		{"uint8", []byte{0xcc, 0xff}, 255},
		{"uint16", []byte{0xcd, 0xff, 0x00}, 65280},
		{"uint32", []byte{0xce, 0xff, 0x00, 0xff, 0x00}, 4278255360},
		{"uint64", []byte{0xcf, 0xff, 0x00, 0xff, 0x00, 0xff, 0x00, 0xff, 0x00}, 18374966859414961920},
	}

	buf := bytes.Buffer{}
	for _, v := range cases {
		buf.Reset()
		ret, err := msgpack.Decode(bytes.NewBuffer(v.bytes))
		outputVerboseJSON(ret, &buf, 0)

		p := MPUint{}
		err = json.Unmarshal(buf.Bytes(), &p)
		if err != nil {
			t.Errorf("%s: Unmarshal Error %s", v.casename, err)
		}
		if v.expected != p.Value {
			t.Errorf("%s: mismatch. given: %d. expected: %d", v.casename, p.Value, v.expected)
		}
	}
}

type MPFloat struct {
	MPBase
	Value float64 `json:"value"`
}

func TestVerboseJSONFloat(t *testing.T) {
	type testcase struct {
		casename string
		bytes    []byte
		expected float64
	}

	cases := []testcase{
		{"float32", []byte{0xca, 0x80, 0x00, 0x00, 0x00}, -0.000000},
		{"float64", []byte{0xcb, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, -0.000000},
	}

	buf := bytes.Buffer{}
	for _, v := range cases {
		buf.Reset()
		ret, err := msgpack.Decode(bytes.NewBuffer(v.bytes))
		outputVerboseJSON(ret, &buf, 0)

		p := MPFloat{}
		err = json.Unmarshal(buf.Bytes(), &p)
		if err != nil {
			t.Errorf("%s: Unmarshal Error %s", v.casename, err)
		}
		if v.expected != p.Value {
			t.Errorf("%s: mismatch. given: %f. expected: %f", v.casename, p.Value, v.expected)
		}
	}
}

type MPArray struct {
	MPBase
	Value []MPInt `json:"value"`
}

func TestVerboseJSONArray(t *testing.T) {
	type testcase struct {
		casename string
		bytes    []byte
		length   int
	}

	cases := []testcase{
		{"fixarray len2", []byte{0x92, 0x00, 0x01}, 2},
		{"array16", []byte{0xdc, 0x00, 0x0f, 0x00, 0x01, 0x00, 0x01, 0x00, 0x01, 0x00, 0x01, 0x00, 0x01, 0x00, 0x01, 0x00, 0x01, 0x00}, 15},
	}

	zb := bytes.Buffer{}
	z, err := msgpack.Decode(bytes.NewBuffer([]byte{0x00}))
	outputVerboseJSON(z, &zb, 0)
	zero := MPInt{}
	err = json.Unmarshal(zb.Bytes(), &zero)
	if err != nil {
		t.Errorf("json.Unmarshal Error")
	}

	ob := bytes.Buffer{}
	o, err := msgpack.Decode(bytes.NewBuffer([]byte{0x01}))
	outputVerboseJSON(o, &ob, 0)
	one := MPInt{}
	err = json.Unmarshal(ob.Bytes(), &one)
	if err != nil {
		t.Errorf("json.Unmarshal Error")
	}

	buf := bytes.Buffer{}
	for _, v := range cases {
		buf.Reset()
		ret, err := msgpack.Decode(bytes.NewBuffer(v.bytes))
		outputVerboseJSON(ret, &buf, 0)

		p := MPArray{}
		err = json.Unmarshal(buf.Bytes(), &p)
		if err != nil {
			t.Errorf("%s: Unmarshal Error %s", v.casename, err)
		}
		if len(p.Value) != v.length {
			t.Errorf("%s: Length Error given: %d. expected: %d", v.casename, len(p.Value), v.length)
		}
		for i, c := range p.Value {
			if i%2 == 0 && c != zero {
				t.Errorf("%s:mismatch given: %v. expected: %v", v.casename, c, zero)
			} else if i%2 != 0 && c != one {
				t.Errorf("%s:mismatch given: %v. expected: %v", v.casename, c, one)
			}
		}
	}
}

type MPMap struct {
	MPBase
	Value []map[string]interface{} `json:"value"`
}

/*
type MPKey struct {
	MPBase
	Value MPString `json:key`
}

type MPValue struct {
	MPBase
	Value MPInt `json:value`
}
*/

func TestVerboseJSONMap(t *testing.T) {
	type testcase struct {
		casename string
		bytes    []byte
		length   int
	}

	cases := []testcase{
		{"fixmap len2", []byte{0x82, 0xa1, 0x30, 0x00, 0xa1, 0x31, 0x01}, 2},
	}

	buf := bytes.Buffer{}
	for _, v := range cases {
		buf.Reset()
		ret, err := msgpack.Decode(bytes.NewBuffer(v.bytes))
		outputVerboseJSON(ret, &buf, 0)

		p := MPMap{}
		err = json.Unmarshal(buf.Bytes(), &p)
		if err != nil {
			t.Errorf("%s: Unmarshal Error %s", v.casename, err)
		}
		if len(p.Value) != v.length {
			t.Errorf("%s: Length Error given: %d. expected: %d", v.casename, len(p.Value), v.length)
		}
		// TODO: check p.Value
	}
}
