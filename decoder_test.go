//
// Copyright (c) 2018- yutopp (yutopp@gmail.com)
//
// Distributed under the Boost Software License, Version 1.0. (See accompanying
// file LICENSE_1_0.txt or copy at  https://www.boost.org/LICENSE_1_0.txt)
//

package amf0

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDecodeCommon(t *testing.T) {
	for _, tc := range testCases {
		tc := tc // capture

		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			r := bytes.NewReader(tc.Binary)
			dec := NewDecoder(r)

			var v interface{}
			err := dec.Decode(&v)
			require.NoError(t, err)
			if tc.AssumeNil {
				require.Nil(t, v)
			} else {
				require.Equal(t, tc.Value, v)
			}

			require.Equal(t, 0, r.Len()) // Assure that all bytes are consumed
		})
	}
}

func TestDecodeNumber(t *testing.T) {
	bin := []byte{0x00, 0x40, 0x24, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00} // Number: 10

	t.Run("int", func(t *testing.T) {
		r := bytes.NewReader(bin)
		dec := NewDecoder(r)

		var v int
		err := dec.Decode(&v)
		require.Nil(t, err)
		require.Equal(t, 10, v)
	})

	t.Run("float64", func(t *testing.T) {
		r := bytes.NewReader(bin)
		dec := NewDecoder(r)

		var v float64
		err := dec.Decode(&v)
		require.Nil(t, err)
		require.Equal(t, float64(10), v)
	})

	t.Run(ptrNestedNumberTest.Name, func(t *testing.T) {
		r := bytes.NewReader(ptrNestedNumberTest.Binary)
		dec := NewDecoder(r)

		var v float64
		err := dec.Decode(&v)
		require.Nil(t, err)
		require.Equal(t, float64(20), v)
	})
}

func TestDecodePartialNumber(t *testing.T) {
	{
		bin := []byte{0x00}

		r := bytes.NewReader(bin)
		dec := NewDecoder(r)

		var v int
		err := dec.Decode(&v)
		require.EqualError(t, err, "unexpected EOF")
	}

	{
		bin := []byte{0x00, 0x00}

		r := bytes.NewReader(bin)
		dec := NewDecoder(r)

		var v int
		err := dec.Decode(&v)
		require.EqualError(t, err, "unexpected EOF")
	}
}

func TestDecodeNil(t *testing.T) {
	bin := []byte{0x05} // Null

	t.Run("assignable to interface{}", func(t *testing.T) {
		r := bytes.NewReader(bin)
		dec := NewDecoder(r)

		var v interface{}
		err := dec.Decode(&v)
		require.NoError(t, err)
		require.Equal(t, nil, v)
	})

	t.Run("assignable to map", func(t *testing.T) {
		r := bytes.NewReader(bin)
		dec := NewDecoder(r)

		var v map[int]int
		err := dec.Decode(&v)
		require.NoError(t, err)
		require.Equal(t, map[int]int(nil), v)
	})

	t.Run("assignable to slice (interface)", func(t *testing.T) {
		r := bytes.NewReader(bin)
		dec := NewDecoder(r)

		var v []interface{}
		err := dec.Decode(&v)
		require.NoError(t, err)
		require.Equal(t, []interface{}(nil), v)
	})

	t.Run("assignable to slice (non-nil primitive)", func(t *testing.T) {
		r := bytes.NewReader(bin)
		dec := NewDecoder(r)

		var v []int
		err := dec.Decode(&v)
		require.NoError(t, err)
		require.Equal(t, []int(nil), v)
	})

	t.Run("NOT assignable to array (despite the length)", func(t *testing.T) {
		r := bytes.NewReader(bin)
		dec := NewDecoder(r)

		var v [42]interface{}
		err := dec.Decode(&v)
		require.Error(t, err)
	})

	t.Run("NOT assignable to int (set to not reference value will fail)", func(t *testing.T) {
		r := bytes.NewReader(bin)
		dec := NewDecoder(r)

		var v int
		err := dec.Decode(&v)
		require.Error(t, err)
	})
}

func TestDecodeObject(t *testing.T) {
	t.Run("assignable to interface", func(t *testing.T) {
		r := bytes.NewReader(objectTest.Binary)
		dec := NewDecoder(r)

		var v interface{}
		err := dec.Decode(&v)
		require.Nil(t, err)
		require.Equal(t, map[string]interface{}{
			"a": "s",
			"b": float64(42),
		}, v)
	})

	t.Run("assignable to map", func(t *testing.T) {
		r := bytes.NewReader(objectTest.Binary)
		dec := NewDecoder(r)

		var v map[string]interface{}
		err := dec.Decode(&v)
		require.Nil(t, err)
		require.Equal(t, map[string]interface{}{
			"a": "s",
			"b": float64(42),
		}, v)
	})

	t.Run("assignable to map which has invalid type", func(t *testing.T) {
		r := bytes.NewReader(objectTest.Binary)
		dec := NewDecoder(r)

		var v map[string]int
		err := dec.Decode(&v)
		require.NotNil(t, err)
	})

	t.Run("assignable to struct", func(t *testing.T) {
		r := bytes.NewReader(objectTest.Binary)
		dec := NewDecoder(r)

		var v sampleObject
		err := dec.Decode(&v)
		require.Nil(t, err)
		require.Equal(t, sampleObject{
			A: "s",
			B: 42,
		}, v)
	})

	t.Run("assignable to struct which keys are not exists", func(t *testing.T) {
		r := bytes.NewReader(objectTest.Binary)
		dec := NewDecoder(r)

		type empty struct{}
		var v empty
		err := dec.Decode(&v)
		require.Nil(t, err)
		require.Equal(t, empty{}, v)
	})
}

func TestDecodeECMAArray(t *testing.T) {
	bin := []byte{
		0x08,
		0x00, 0x00, 0x00, 0x01,
		// Key a
		0x00, 0x01, 0x61, // a
		// Value nil
		0x05,
		// End
		0x00, 0x00, 0x09,
	}

	t.Run("assignable to interface{}", func(t *testing.T) {
		r := bytes.NewReader(bin)
		dec := NewDecoder(r)

		var v interface{}
		err := dec.Decode(&v)
		require.Nil(t, err)
		require.Equal(t, ECMAArray{"a": nil}, v)
	})

	t.Run("assignable to map", func(t *testing.T) {
		r := bytes.NewReader(bin)
		dec := NewDecoder(r)

		var v map[string]interface{}
		err := dec.Decode(&v)
		require.Nil(t, err)
		require.Equal(t, map[string]interface{}{"a": nil}, v)
	})

	t.Run("assignable to map which has not string type key will fail", func(t *testing.T) {
		r := bytes.NewReader(bin)
		dec := NewDecoder(r)

		var v map[int]interface{}
		err := dec.Decode(&v)
		require.NotNil(t, err)
	})

	t.Run("assignable to map which has unmatched value type will fail", func(t *testing.T) {
		r := bytes.NewReader(bin)
		dec := NewDecoder(r)

		var v map[string]int
		err := dec.Decode(&v)
		require.NotNil(t, err)
	})

	t.Run("assignable to int (set to not reference value will fail)", func(t *testing.T) {
		r := bytes.NewReader(bin)
		dec := NewDecoder(r)

		var v int
		err := dec.Decode(&v)
		require.NotNil(t, err)
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
		r := bytes.NewReader(bin)
		dec := NewDecoder(r)

		var v []string
		err := dec.Decode(&v)
		require.NoError(t, err)
		require.Equal(t, []string{"str", "str"}, v)
	})

	t.Run("assignable to typed array (same length)", func(t *testing.T) {
		r := bytes.NewReader(bin)
		dec := NewDecoder(r)

		var v [2]string
		err := dec.Decode(&v)
		require.NoError(t, err)
		require.Equal(t, [2]string{"str", "str"}, v)
	})

	t.Run("assignable to typed array (different length)", func(t *testing.T) {
		r := bytes.NewReader(bin)
		dec := NewDecoder(r)

		var v [10]string
		err := dec.Decode(&v)
		require.Error(t, err) // should support an array which has length that more than of equals to length of an encoded strict array?
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
		r := bytes.NewReader(bin)
		dec := NewDecoder(r)

		var v []string
		err := dec.Decode(&v)
		require.NotNil(t, err)
	})
}

func TestDecodeUnknownMarker(t *testing.T) {
	bin := []byte{
		// Unknown Marker
		0xff,
	}

	t.Run("Cannot decode", func(t *testing.T) {
		r := bytes.NewReader(bin)
		dec := NewDecoder(r)

		var v interface{}
		err := dec.Decode(&v)
		require.Error(t, err)
	})
}
