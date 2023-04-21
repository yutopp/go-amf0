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

var number20ForAddr = 20

type sampleObject struct {
	A string `amf0:"a"`
	B int    `amf0:"b"`
}

var testCases = []testCase{
	{
		Name:  "Number(Int)",
		Value: float64(10),
		Binary: []byte{
			// Number Marker
			0x00,
			// Value(10: double) BigEndian
			0x40, 0x24, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		},
	},

	{
		Name:  "Boolean (false)",
		Value: false,
		Binary: []byte{
			// Boolean Marker
			0x01,
			// False
			0x00,
		},
	},
	{
		Name:  "Boolean (true)",
		Value: true,
		Binary: []byte{
			// Boolean Marker
			0x01,
			// True
			0x01,
		},
	},
	{
		Name:  "String",
		Value: "abc",
		// 0x02: String Marker
		// 0x00, 0x03: Length(3: u16) BigEndian
		// 0x61, 0x62, 0x63: Value(abc: []byte)
		Binary: []byte{0x02, 0x00, 0x03, 0x61, 0x62, 0x63},
	},
	{
		Name:  "Nil",
		Value: nil,
		Binary: []byte{
			// Null Marker
			0x05,
		},
	},
	{
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
	{
		Name: "Strict Array",
		Value: []interface{}{
			"str",
			float64(10),
		},
		Binary: []byte{
			// Strict Array Marker
			0x0a,
			// Array length (2: u32) BigEndian
			0x00, 0x00, 0x00, 0x02,
			// Elem 0 (string)
			0x2, 0x00, 0x03, 0x73, 0x74, 0x72,
			// Elem 1 (number)
			0x00, 0x40, 0x24, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		},
	},
	{
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

var onlyEncodingTestCases = []testCase{
	ptrNestedNumberTest,
	objectTest,
}

var ptrNestedNumberTest = testCase{
	Name:  "Number(Int ptr)",
	Value: &number20ForAddr,
	Binary: []byte{
		// Number Marker
		0x00,
		// Value(20: double) BigEndian
		0x40, 0x34, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	},
}

var objectTest = testCase{
	Name: "Object",
	Value: sampleObject{
		A: "s",
		B: 42,
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
}
