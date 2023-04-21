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

	"github.com/stretchr/testify/assert"
)

func TestEncodeCommon(t *testing.T) {
	var allTestCases []testCase
	allTestCases = append(testCases, onlyEncodingTestCases...)

	for _, tc := range allTestCases {
		tc := tc // capture

		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			buf := bytes.NewBuffer([]byte{})
			enc := NewEncoder(buf)
			enc.sortKeys = true // for debuging

			err := enc.Encode(tc.Value)
			assert.Nil(t, err)
			assert.Equal(t, tc.Binary, buf.Bytes())
		})
	}
}

func TestEncodeObjectEnd(t *testing.T) {
	buf := bytes.NewBuffer([]byte{})
	enc := NewEncoder(buf)

	err := enc.Encode(ObjectEnd)
	assert.Nil(t, err)
}
