//
// Copyright (c) 2018- yutopp (yutopp@gmail.com)
//
// Distributed under the Boost Software License, Version 1.0. (See accompanying
// file LICENSE_1_0.txt or copy at  https://www.boost.org/LICENSE_1_0.txt)
//

package amf0

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDecodeCommon(t *testing.T) {
	for _, tc := range testCases {
		tc := tc // capture

		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			buf := bytes.NewBuffer(tc.Binary)
			dec := NewDecoder(buf)

			var v interface{}
			err := dec.Decode(&v)
			assert.Nil(t, err)
			assert.Equal(t, tc.Value, v)
		})
	}
}

func TestDecodeNumber(t *testing.T) {
	bin := []byte{0x00, 0x40, 0x24, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00} // Number: 10

	t.Run("int", func(t *testing.T) {
		buf := bytes.NewBuffer(bin)
		dec := NewDecoder(buf)

		var v int
		err := dec.Decode(&v)
		assert.Nil(t, err)
		assert.Equal(t, 10, v)
	})

	t.Run("float64", func(t *testing.T) {
		buf := bytes.NewBuffer(bin)
		dec := NewDecoder(buf)

		var v float64
		err := dec.Decode(&v)
		assert.Nil(t, err)
		assert.Equal(t, float64(10), v)
	})
}

func TestDecodeNil(t *testing.T) {
	bin := []byte{0x05} // Null

	t.Run("assignable to interface{}", func(t *testing.T) {
		buf := bytes.NewBuffer(bin)
		dec := NewDecoder(buf)

		var v interface{}
		err := dec.Decode(&v)
		assert.Nil(t, err)
		assert.Equal(t, nil, v)
	})

	t.Run("assignable to map", func(t *testing.T) {
		buf := bytes.NewBuffer(bin)
		dec := NewDecoder(buf)

		var v map[int]int
		err := dec.Decode(&v)
		assert.Nil(t, err)
		assert.Equal(t, map[int]int(nil), v)
	})

	t.Run("assignable to int (set to not reference value will fail)", func(t *testing.T) {
		buf := bytes.NewBuffer(bin)
		dec := NewDecoder(buf)

		var v int
		err := dec.Decode(&v)
		assert.NotNil(t, err)
	})
}

func TestDecodeStrictArraySame(t *testing.T) {
	bin := []byte{
		// Strict Array Marker
		0x0a,
		// Array length (2: u32) BigEndian
		0x00, 0x00, 0x00, 0x02,
		// Elem 0 (string)
		0x2, 0x00, 0x03, 0x73, 0x74, 0x72,
		// Elem 1 (string)
		0x2, 0x00, 0x03, 0x73, 0x74, 0x72,
	}

	t.Run("assignable to typed slice", func(t *testing.T) {
		buf := bytes.NewBuffer(bin)
		dec := NewDecoder(buf)

		var v []string
		err := dec.Decode(&v)
		assert.Nil(t, err)
		assert.Equal(t, []string{"str", "str"}, v)
	})

	t.Run("assignable to typed array (same length)", func(t *testing.T) {
		buf := bytes.NewBuffer(bin)
		dec := NewDecoder(buf)

		var v [2]string
		err := dec.Decode(&v)
		assert.Nil(t, err)
		assert.Equal(t, [2]string{"str", "str"}, v)
	})

	t.Run("assignable to typed array (different length)", func(t *testing.T) {
		buf := bytes.NewBuffer(bin)
		dec := NewDecoder(buf)

		var v [10]string
		err := dec.Decode(&v)
		assert.NotNil(t, err) // should support an array which has length that more than of equals to length of an encoded strict array?
	})
}

func TestDecodeStrictArrayHetero(t *testing.T) {
	bin := []byte{
		// Strict Array Marker
		0x0a,
		// Array length (2: u32) BigEndian
		0x00, 0x00, 0x00, 0x02,
		// Elem 0 (string)
		0x2, 0x00, 0x03, 0x73, 0x74, 0x72,
		// Elem 1 (number)
		0x00, 0x40, 0x24, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}

	t.Run("assignable to typed slice", func(t *testing.T) {
		buf := bytes.NewBuffer(bin)
		dec := NewDecoder(buf)

		var v []string
		err := dec.Decode(&v)
		assert.NotNil(t, err)
	})
}
