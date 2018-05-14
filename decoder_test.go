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
