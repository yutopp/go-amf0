//
// Copyright (c) 2018- yutopp (yutopp@gmail.com)
//
// Distributed under the Boost Software License, Version 1.0. (See accompanying
// file LICENSE_1_0.txt or copy at  https://www.boost.org/LICENSE_1_0.txt)
//

package amf0

type testCase struct {
	Name   string
	Value  interface{}
	Binary []byte
}

var testCases = []testCase{
	testCase{
		Name:  "Number(Int)",
		Value: float64(10),
		// 0x00: Number Marker
		// 0x40, 0x24, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00: Value(10: double) BigEndian
		Binary: []byte{0x00, 0x40, 0x24, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
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
		// 0x03: Object Marker
		//   0x00, 0x01: Length(1: u16) BigEndian
		//   0x61: Key(a: []byte)
		//     0x02: String Marker
		//     0x00, 0x01: Length(1: u16) BigEndian
		//     0x73: Value(s: []byte)
		//   0x00, 0x01: Length(1: u16) BigEndian
		//   0x62: Key(b: []byte)
		//     0x00: Number Marker
		//     0x40, 0x45, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00: Value(42: double) BigEndian
		//   0x00, 0x00: Length(0: u16) BigEndian
		//   (empty)
		//     0x09: ObjectEndMarker
		Binary: []byte{
			0x03,
			0x00, 0x01, 0x61,
			0x02, 0x00, 0x01, 0x73,
			0x00, 0x01, 0x62,
			0x00, 0x40, 0x45, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00,
			0x09,
		},
	},
	testCase{
		Name:  "Nil",
		Value: nil,
		// 0x05: Null Marker
		Binary: []byte{0x05},
	},
}
