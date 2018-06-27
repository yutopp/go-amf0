//
// Copyright (c) 2018- yutopp (yutopp@gmail.com)
//
// Distributed under the Boost Software License, Version 1.0. (See accompanying
// file LICENSE_1_0.txt or copy at  https://www.boost.org/LICENSE_1_0.txt)
//

package amf0

import (
	"encoding/binary"
	"encoding/hex"
	_ "github.com/pkg/errors"
	"io"
	"math"
	"reflect"
	"unicode/utf8"
)

type Decoder struct {
	r io.Reader
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{
		r: r,
	}
}

func (dec *Decoder) Decode(v interface{}) error {
	rv := reflect.ValueOf(v)
	return dec.decode(rv)
}

func (dec *Decoder) decode(rv reflect.Value) error {
	marker, err := dec.readU8()
	if err != nil {
		return err
	}

	switch Marker(marker) {
	case MarkerNumber:
		return dec.decodeNumber(rv)
	case MarkerBoolean:
		return dec.decodeBoolean(rv)
	case MarkerString:
		return dec.decodeString(rv)
	case MarkerObject:
		return dec.decodeObject(rv)
	case MarkerNull:
		return dec.decodeNull(rv)
	case MarkerEcmaArray:
		return dec.decodeECMAArray(rv)
	case MarkerObjectEnd:
		return dec.decodeObjectEnd(rv)
	default:
		return &UnsupportedMarkerError{
			Marker: marker,
		}
	}
}

func (dec *Decoder) decodeNumber(rv reflect.Value) error {
	num, err := dec.readDouble()
	if err != nil {
		return err
	}

	rv, err = indirect(rv)
	if err != nil {
		return err
	}

	rv.Set(reflect.ValueOf(num).Convert(rv.Type()))

	return nil
}

func (dec *Decoder) decodeBoolean(rv reflect.Value) error {
	num, err := dec.readU8()
	if err != nil {
		return err
	}

	tf := false
	if num != 0 {
		tf = true
	}

	rv, err = indirect(rv)
	if err != nil {
		return err
	}

	rv.Set(reflect.ValueOf(tf))

	return nil
}

func (dec *Decoder) decodeString(rv reflect.Value) error {
	str, err := dec.readUTF8()
	if err != nil {
		return err
	}

	rv, err = indirect(rv)
	if err != nil {
		return err
	}

	rv.Set(reflect.ValueOf(str))

	return nil
}

func (dec *Decoder) decodeObject(rv reflect.Value) error {
	obj := make(map[string]interface{}) // TODO: optimize

	for {
		key, err := dec.readUTF8()
		if err != nil {
			return err
		}

		if key == "" {
			marker, err := dec.readU8()
			if err != nil {
				return err
			}
			if marker != MarkerObjectEnd {
				return &DecodeError{
					Message: "Not ended with object-end",
				}
			}
			break
		}

		var value interface{}
		if err := dec.Decode(&value); err != nil {
			return err
		}

		obj[key] = value
	}

	rv, err := indirect(rv)
	if err != nil {
		return err
	}

	rv.Set(reflect.ValueOf(obj))

	return nil
}

func (dec *Decoder) decodeObjectProperty(rk *string, rv reflect.Value) error {
	key, err := dec.readUTF8()
	if err != nil {
		return err
	}
	*rk = key

	return dec.decode(rv)
}

func (dec *Decoder) decodeNull(rv reflect.Value) error {
	rv, err := indirect(rv)
	if err != nil {
		return err
	}

	if rv.Kind() != reflect.Map && rv.Kind() != reflect.Slice && rv.Kind() != reflect.Interface {
		return &NotAssignableError{
			Message:      "Not reference type",
			ReceiverKind: rv.Kind(),
		}
	}

	rv.Set(reflect.Zero(rv.Type()))

	return nil
}

func (dec *Decoder) decodeECMAArray(rv reflect.Value) error {
	rv, err := indirect(rv)
	if err != nil {
		return err
	}

	if rv.Kind() != reflect.Interface && rv.Kind() != reflect.Map {
		return &NotAssignableError{
			Message:      "Not map or interface type",
			ReceiverKind: rv.Kind(),
		}
	}

	if rv.IsNil() {
		switch rv.Kind() {
		case reflect.Interface:
			rv.Set(reflect.MakeMap(reflect.TypeOf(ECMAArray{})))
		case reflect.Map:
			rv.Set(reflect.MakeMap(rv.Type()))
		}
		rv = rv.Elem()
	}

	numElems, err := dec.readU32()
	if err != nil {
		return err
	}

	for i := uint32(0); i < numElems; i++ {
		var key string
		var value interface{}

		rvM := reflect.ValueOf(&value)
		if err := dec.decodeObjectProperty(&key, rvM); err != nil {
			return err
		}

		rv.SetMapIndex(reflect.ValueOf(key), rvM.Elem())
	}

	return nil
}

func (dec *Decoder) decodeObjectEnd(rv reflect.Value) error {
	rv, err := indirect(rv)
	if err != nil {
		return err
	}

	rv.Set(reflect.ValueOf(ObjectEnd))

	return nil
}

func (dec *Decoder) readU8() (uint8, error) {
	u8 := make([]byte, 1) // TODO: optimize
	_, err := io.ReadAtLeast(dec.r, u8, 1)
	if err != nil {
		return 0, err
	}

	return u8[0], nil
}

func (dec *Decoder) readU16() (uint16, error) {
	u16 := make([]byte, 2) // TODO: optimize
	_, err := io.ReadAtLeast(dec.r, u16, 2)
	if err != nil {
		return 0, err
	}

	return binary.BigEndian.Uint16(u16), nil
}

func (dec *Decoder) readU32() (uint32, error) {
	bin := make([]byte, 4)
	_, err := io.ReadAtLeast(dec.r, bin, len(bin))
	if err != nil {
		return 0, err
	}

	return binary.BigEndian.Uint32(bin), nil
}

func (dec *Decoder) readDouble() (float64, error) {
	d := make([]byte, 8) // TODO: optimize
	_, err := io.ReadAtLeast(dec.r, d, 8)
	if err != nil {
		return 0, err
	}

	bits := binary.BigEndian.Uint64(d)
	return math.Float64frombits(bits), nil
}

func (dec *Decoder) readUTF8Chars(len int) (string, error) {
	str := make([]byte, len) // TODO: optimize
	_, err := io.ReadAtLeast(dec.r, str, len)
	if err != nil {
		return "", err
	}

	if !utf8.Valid(str) {
		return "", &DecodeError{
			Message: "Invalid utf8 sequence",
			Dump:    hex.Dump(str),
		}
	}

	return string(str), nil
}

func (dec *Decoder) readUTF8() (string, error) {
	len, err := dec.readU16()
	if err != nil {
		return "", err
	}
	if len == 0 {
		return "", nil // empty
	}

	str, err := dec.readUTF8Chars(int(len))
	if err != nil {
		return "", err
	}

	return str, nil
}

func indirect(rv reflect.Value) (reflect.Value, error) {
	if rv.Kind() != reflect.Ptr {
		return reflect.Value{}, &NotAssignableError{
			Message:      "Not pointer",
			ReceiverKind: rv.Kind(),
		}
	}
	if rv.IsNil() {
		return reflect.Value{}, &NotAssignableError{
			Message:      "Nil",
			ReceiverKind: rv.Kind(),
		}
	}

	return reflect.Indirect(rv), nil
}
