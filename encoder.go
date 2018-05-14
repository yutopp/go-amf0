//
// Copyright (c) 2018- yutopp (yutopp@gmail.com)
//
// Distributed under the Boost Software License, Version 1.0. (See accompanying
// file LICENSE_1_0.txt or copy at  https://www.boost.org/LICENSE_1_0.txt)
//

package amf0

import (
	"encoding/binary"
	"io"
	"math"
	"reflect"
)

type Encoder struct {
	w io.Writer
}

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{
		w: w,
	}
}

func (enc *Encoder) Encode(v interface{}) error {
	rv := reflect.ValueOf(v)
	return enc.encode(rv)
}

func (enc *Encoder) encode(rv reflect.Value) error {
	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		fallthrough
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		fallthrough
	case reflect.Float32, reflect.Float64:
		return enc.encodeNumber(rv)
	case reflect.String:
		return enc.encodeString(rv)
	case reflect.Map:
		return enc.encodeMapAsObject(rv)
	case reflect.Array, reflect.Slice:
		return enc.encodeArray(rv)
	case reflect.Interface:
		if rv.IsNil() {
			return enc.encodeNull()
		}
		return enc.Encode(rv.Interface())
	case reflect.Invalid:
		return enc.encodeNull()
	}

	return &UnsupportedKindError{
		Kind: rv.Kind(),
	}
}

func (enc *Encoder) encodeNumber(rv reflect.Value) error {
	if err := enc.writeU8(uint8(MarkerNumber)); err != nil {
		return err
	}

	var d float64
	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		fallthrough
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		d = float64(rv.Int())
	case reflect.Float32, reflect.Float64:
		d = rv.Float()
	default:
		return &UnsupportedKindError{
			Kind: rv.Kind(),
		}
	}

	return enc.writeDouble(d)
}

func (enc *Encoder) encodeString(rv reflect.Value) error {
	s := rv.String()
	if len(s) > 65535 {
		// TODO: use long string
		panic("not implemented")
	}

	if err := enc.writeU8(uint8(MarkerString)); err != nil {
		return err
	}
	return enc.writeUTF8(s)
}

func (enc *Encoder) encodeMapAsObject(rv reflect.Value) error {
	if err := enc.writeU8(uint8(MarkerObject)); err != nil {
		return err
	}

	for _, key := range rv.MapKeys() {
		if key.Kind() != reflect.String {
			return &UnexpectedKeyTypeError{
				ActualKind: key.Kind(),
				ExpectKind: reflect.String,
			}
		}

		if err := enc.writeUTF8(key.String()); err != nil {
			return err
		}

		value := rv.MapIndex(key)
		if err := enc.encode(value); err != nil {
			return err
		}
	}

	if err := enc.writeUTF8(""); err != nil { // utf-8-empty
		return err
	}

	return enc.encodeObjectEnd()
}

func (enc *Encoder) encodeObjectEnd() error {
	return enc.writeU8(uint8(MarkerObjectEnd))
}

func (enc *Encoder) encodeArray(rv reflect.Value) error {
	if rv.Len() >= 1 {
		re := rv.Index(0)
		if re.Kind() == reflect.Ptr {
			re = reflect.Indirect(re)
		}
		if isECMAArrayElem(re) {
			return enc.encodeArrayAsECMAArray(rv)
		}
	}

	panic("not implemented") // TODO
	//return enc.encodeArrayAsStrictArray(rv)
}

func (enc *Encoder) encodeArrayAsECMAArray(rv reflect.Value) error {
	if err := enc.writeU8(uint8(MarkerEcmaArray)); err != nil {
		return err
	}

	l := rv.Len()
	if err := enc.writeU32(uint32(l)); err != nil {
		return err
	}

	for i := 0; i < rv.Len(); i++ {
		re := rv.Index(i)
		if re.Kind() == reflect.Ptr {
			re = reflect.Indirect(re)
		}
		key := re.Field(0).String()
		value := re.Field(1)

		if err := enc.writeUTF8(key); err != nil {
			return err
		}
		if err := enc.encode(value); err != nil {
			return err
		}
	}

	return nil
}

func (enc *Encoder) encodeNull() error {
	return enc.writeU8(MarkerNull)
}

func (enc *Encoder) writeU8(num uint8) error {
	_, err := enc.w.Write([]byte{num}) // TODO: optimize
	return err
}

func (enc *Encoder) writeU16(num uint16) error {
	buf := make([]byte, 2) // TODO: optimize
	binary.BigEndian.PutUint16(buf, num)

	_, err := enc.w.Write(buf)
	return err
}

func (enc *Encoder) writeU32(num uint32) error {
	buf := make([]byte, 4) // TODO: optimize
	binary.BigEndian.PutUint32(buf, num)

	_, err := enc.w.Write(buf)
	return err
}

func (enc *Encoder) writeDouble(f64 float64) error {
	buf := make([]byte, 8) // TODO: optimize
	u64 := math.Float64bits(f64)
	binary.BigEndian.PutUint64(buf, u64)

	_, err := enc.w.Write(buf)
	return err
}

func (enc *Encoder) writeUTF8(str string) error {
	l := uint16(len(str))
	if err := enc.writeU16(l); err != nil {
		return err
	}
	_, err := enc.w.Write([]byte(str))
	return err
}

func isECMAArrayElem(rv reflect.Value) bool {
	if rv.Kind() != reflect.Struct {
		return false
	}

	ty := rv.Type()
	if ty.NumField() != 2 {
		return false
	}

	keyField := ty.Field(0)
	if keyField.Type.Kind() != reflect.String {
		return false
	}

	return keyField.Tag == `amf0:"ecma"`
}
