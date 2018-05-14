//
// Copyright (c) 2018- yutopp (yutopp@gmail.com)
//
// Distributed under the Boost Software License, Version 1.0. (See accompanying
// file LICENSE_1_0.txt or copy at  https://www.boost.org/LICENSE_1_0.txt)
//

package amf0

import (
	"fmt"
	"reflect"
)

// UnsupportedKindError returned by Encode when the value which has an unsupported kind is passed to the function.
type UnsupportedKindError struct {
	Kind reflect.Kind
}

// Error ...
func (e *UnsupportedKindError) Error() string {
	return fmt.Sprintf("Unsupported kind: %+v", e.Kind.String())
}

// UnexpectedKeyTypeError ...
type UnexpectedKeyTypeError struct {
	ActualKind reflect.Kind
	ExpectKind reflect.Kind
}

// Error ...
func (e *UnexpectedKeyTypeError) Error() string {
	return fmt.Sprintf("Unsupported key kind: %+v should be %+v", e.ActualKind.String(), e.ExpectKind.String())
}

type UnsupportedMarkerError struct {
	Marker uint8
}

// Error ...
func (e *UnsupportedMarkerError) Error() string {
	return fmt.Sprintf("Unsupported marker: %+v", e.Marker)
}

type DecodeError struct {
	Message string
}

// Error ...
func (e *DecodeError) Error() string {
	return e.Message
}

// NotAssignableError ...
type NotAssignableError struct {
	Message      string
	ReceiverKind reflect.Kind
}

// Error ...
func (e *NotAssignableError) Error() string {
	return fmt.Sprintf("Not assignable to receiver value: Message=%+v, Kind=%+v", e.Message, e.ReceiverKind)
}
