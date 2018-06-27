//
// Copyright (c) 2018- yutopp (yutopp@gmail.com)
//
// Distributed under the Boost Software License, Version 1.0. (See accompanying
// file LICENSE_1_0.txt or copy at  https://www.boost.org/LICENSE_1_0.txt)
//

package amf0

import (
	"time"
)

type testCase struct {
	Name   string
	Value  interface{}
	Binary []byte
}

var testCases = []testCase{
	testCase{
		Name:  "Number(Int)",
		Value: float64(10),
		Binary: []byte{
			// Number Marker
			0x00,
			// Value(10: double) BigEndian
			0x40, 0x24, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		},
	},
	testCase{
		Name:  "String",
		Value: "abc",
		// 0x02: String Marker
		// 0x00, 0x03: Length(3: u16) BigEndian
		// 0x61, 0x62, 0x63: Value(abc: []byte)
		Binary: []byte{0x02, 0x00, 0x03, 0x61, 0x62, 0x63},
	},
	testCase{
		Name: "Map",
		Value: map[string]interface{}{
			"a": "s",
			"b": float64(42),
		},
		Binary: []byte{
			// Object Marker
			0x03,
			// - Length(1: u16) BigEndian
			0x00, 0x01,
			//   Key(a: []byte)
			0x61,
			//   - String Marker
			0x02,
			//     Length(1: u16) BigEndian
			0x00, 0x01,
			//     Value(s: []byte)
			0x73,
			// - Length(1: u16) BigEndian
			0x00, 0x01,
			//   Key(b: []byte)
			0x62,
			//   - Number Marker
			0x00,
			//     Value(42: double) BigEndian
			0x40, 0x45, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			// - Length(0: u16) BigEndian
			0x00, 0x00,
			//   Key(empty)
			//   - ObjectEndMarker
			0x09,
		},
	},
	testCase{
		Name:  "Nil",
		Value: nil,
		Binary: []byte{
			// Null Marker
			0x05,
		},
	},
	testCase{
		Name: "ECMA Array",
		Value: ECMAArray{
			"a": "str",
			"b": float64(10), // all decoded numbers become float64 type
		},
		Binary: []byte{
			// ECMA Array Marker
			0x08,
			// Associative count(2: u32) BigEndian
			0x00, 0x00, 0x00, 0x02,
			// - Length(1: u16)
			0x00, 0x01,
			//   Key(a: []byte)
			0x61,
			//   - String Marker
			0x02,
			//     Length(3: u16) BigEndian
			0x00, 0x03,
			//     Value(abc: []byte)
			0x73, 0x74, 0x72,
			// - Length(1: u16)
			0x00, 0x01,
			//   Key(b: []byte)
			0x62,
			//   - Number Marker
			0x00,
			//     Value(10: double) BigEndian
			0x40, 0x24, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			// - Length(0: u16) BigEndian
			0x00, 0x00,
			//   Key(empty)
			//   - ObjectEndMarker
			0x09,
		},
	},
	testCase{
		Name:  "Date",
		Value: time.Unix(0x1234, 0).In(time.UTC),
		Binary: []byte{
			// Date Marker
			0x0b,
			// Unix time[ms]
			0x41, 0x51, 0xc6, 0xc8, 0x00, 0x00, 0x00, 0x00,
			// Time zone
			0x00, 0x00,
		},
	},
}
