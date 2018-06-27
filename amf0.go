//
// Copyright (c) 2018- yutopp (yutopp@gmail.com)
//
// Distributed under the Boost Software License, Version 1.0. (See accompanying
// file LICENSE_1_0.txt or copy at  https://www.boost.org/LICENSE_1_0.txt)
//

package amf0

type Marker byte

const (
	MarkerNumber      Marker = 0x00
	MarkerBoolean     Marker = 0x01
	MarkerString      Marker = 0x02
	MarkerObject             = 0x03
	MarkerMovieclip          = 0x04 // reserved, not supported
	MarkerNull               = 0x05
	MarkerUndefined          = 0x06
	MarkerReference          = 0x07
	MarkerEcmaArray          = 0x08
	MarkerObjectEnd          = 0x09
	MarkerStrictArray        = 0x0A
	MarkerDate               = 0x0B
	MarkerLongString         = 0x0C
	MarkerUnsupported        = 0x0D
	MarkerRecordSet          = 0x0E // reserved, not supported
	MarkerXMLDocument        = 0x0F
	MarkerTypedObject        = 0x10
)

type ECMAArray map[string]interface{}

var ObjectEnd = &struct{}{}
