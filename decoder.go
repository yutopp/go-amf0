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
	"fmt"
	_ "github.com/pkg/errors"
	"io"
	"math"
	"reflect"
	"time"
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

func (dec *Decoder) Reset(r io.Reader) {
	dec.r = r
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

	case MarkerMovieclip:
		return dec.decodeMovieClip(rv)

	case MarkerNull:
		return dec.decodeNull(rv)

	case MarkerUndefined:
		return dec.decodeUndefined(rv)

	case MarkerReference:
		return dec.decodeReference(rv)

	case MarkerEcmaArray:
		return dec.decodeECMAArray(rv)

	case MarkerObjectEnd:
		return ErrObjectEndMarker

	case MarkerStrictArray:
		return dec.decodeStrictArray(rv)

	case MarkerDate:
		return dec.decodeDate(rv)

	case MarkerLongString:
		return dec.decodeLongString(rv)

	case MarkerUnsupported:
		panic("Not implemented: Unsupported") // TODO: returns error

	case MarkerRecordSet:
		return dec.decodeRecordSet(rv)

	case MarkerXMLDocument:
		return dec.decodeXMLDocument(rv)

	case MarkerTypedObject:
		return dec.decodeTypedObject(rv)

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

	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		fallthrough
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		fallthrough
	case reflect.Float32, reflect.Float64:
		fallthrough
	case reflect.Interface:
		rv.Set(reflect.ValueOf(num).Convert(rv.Type()))

	default:
		return &NotAssignableError{
			Message: "Not numeric type",
			Kind:    rv.Kind(),
			Type:    rv.Type(),
		}
	}

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

func (dec *Decoder) decodeObjectProperty(rk *string, rv reflect.Value) (bool, error) {
	key, err := dec.readUTF8()
	if err != nil {
		return false, err
	}
	if key == "" {
		// End object
		marker, err := dec.readU8()
		if err != nil {
			return false, err
		}
		if marker != MarkerObjectEnd {
			return false, &DecodeError{
				Message: "Not ended with object-end",
			}
		}

		return true, nil
	}

	*rk = key
	return false, dec.decode(rv)
}

func (dec *Decoder) decodeMovieClip(rv reflect.Value) error {
	panic("Not implemented: MovieClip")
}

func (dec *Decoder) decodeNull(rv reflect.Value) error {
	rv, err := indirect(rv)
	if err != nil {
		return err
	}

	if rv.Kind() != reflect.Map && rv.Kind() != reflect.Slice && rv.Kind() != reflect.Interface {
		return &NotAssignableError{
			Message: "Not reference type",
			Kind:    rv.Kind(),
			Type:    rv.Type(),
		}
	}

	rv.Set(reflect.Zero(rv.Type()))

	return nil
}

func (dec *Decoder) decodeUndefined(rv reflect.Value) error {
	panic("Not implemented: Undefined")
}

func (dec *Decoder) decodeReference(rv reflect.Value) error {
	panic("Not implemented: Reference")
}

func (dec *Decoder) decodeECMAArray(rv reflect.Value) error {
	rv, err := indirect(rv)
	if err != nil {
		return err
	}

	if rv.Kind() != reflect.Interface && rv.Kind() != reflect.Map {
		return &NotAssignableError{
			Message: "Not map or interface type",
			Kind:    rv.Kind(),
			Type:    rv.Type(),
		}
	}

	if rv.IsNil() {
		switch rv.Kind() {
		case reflect.Interface:
			rv.Set(reflect.MakeMap(reflect.TypeOf(ECMAArray{})))
			rv = rv.Elem()
		case reflect.Map:
			rv.Set(reflect.MakeMap(rv.Type()))
		}
	}

	if rv.Kind() == reflect.Map {
		keyTy := rv.Type().Key()
		if keyTy.Kind() != reflect.String {
			return &NotAssignableError{
				Message: "Key of map is not string type",
				Kind:    keyTy.Kind(),
				Type:    keyTy,
			}
		}
	}

	numElems, err := dec.readU32()
	if err != nil {
		return err
	}
	_ = numElems

	var key string
	value := reflect.New(rv.Type().Elem())

	for {
		isEnd, err := dec.decodeObjectProperty(&key, value)
		if err != nil {
			return err
		}
		if isEnd {
			break
		}

		rv.SetMapIndex(reflect.ValueOf(key), value.Elem())
	}

	return nil
}

// skip ObjectEnd

func (dec *Decoder) decodeStrictArray(rv reflect.Value) error {
	rv, err := indirect(rv)
	if err != nil {
		return err
	}

	if rv.Kind() != reflect.Interface && rv.Kind() != reflect.Slice && rv.Kind() != reflect.Array {
		return &NotAssignableError{
			Message: "Not array or slice or interface type",
			Kind:    rv.Kind(),
			Type:    rv.Type(),
		}
	}

	length, err := dec.readU32()
	if err != nil {
		return err
	}
	if length > math.MaxInt32 {
		// specification said "maximum 4294967295", however we cannot support that... TODO: support if possible
		return fmt.Errorf("Unsupported array length: Expected <= %d, Actual = %d", math.MaxInt32, length)
	}

	if rv.Kind() == reflect.Interface || rv.Kind() == reflect.Slice {
		if rv.IsNil() {
			switch rv.Kind() {
			case reflect.Interface:
				rv.Set(reflect.ValueOf(make([]interface{}, int(length), int(length))))
				rv = rv.Elem()
			case reflect.Slice:
				rv.Set(reflect.MakeSlice(rv.Type(), int(length), int(length)))
			}
		}
	}

	if rv.Kind() == reflect.Slice || rv.Kind() == reflect.Array {
		if rv.Len() != int(length) {
			return fmt.Errorf("Length of array/slice is different: Expected = %d, Actual = %d", int(length), rv.Len())
		}
	}

	for i := 0; i < int(length); i++ {
		if err := dec.decode(rv.Index(i).Addr()); err != nil {
			return err
		}
	}

	return nil
}

func (dec *Decoder) decodeDate(rv reflect.Value) error {
	rv, err := indirect(rv)
	if err != nil {
		return err
	}

	if rv.Kind() != reflect.Interface && rv.Kind() != reflect.Struct {
		return &NotAssignableError{
			Message: "Not struct or interface type",
			Kind:    rv.Kind(),
			Type:    rv.Type(),
		}
	}

	if rv.Kind() == reflect.Struct && rv.Type() != reflect.TypeOf(time.Time{}) {
		return &NotAssignableError{
			Message: "Not time.Time type",
			Kind:    rv.Kind(),
			Type:    rv.Type(),
		}
	}

	unixMs, err := dec.readDouble()
	if err != nil {
		return err
	}

	tz, err := dec.readS16()
	if err != nil {
		return err
	}

	t := time.Unix(int64(unixMs)/1000, int64(unixMs)%1000*int64(time.Nanosecond)).In(time.UTC)

	if tz != 0x00 {
		// Timezone is specified
		// TODO: support
	}

	rv.Set(reflect.ValueOf(t))

	return nil
}

func (dec *Decoder) decodeLongString(rv reflect.Value) error {
	panic("Not implemented: LongString")
}

// skip Unsupported

func (dec *Decoder) decodeRecordSet(rv reflect.Value) error {
	panic("Not implemented: RecordSet")
}

func (dec *Decoder) decodeXMLDocument(rv reflect.Value) error {
	panic("Not implemented: XMLDocument")
}

func (dec *Decoder) decodeTypedObject(rv reflect.Value) error {
	panic("Not implemented: TypedObject")
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

func (dec *Decoder) readS16() (int16, error) {
	n, err := dec.readU16()
	if err != nil {
		return 0, err
	}

	return int16(n), nil
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
			Message: "Not pointer",
			Kind:    rv.Kind(),
			Type:    rv.Type(),
		}
	}
	if rv.IsNil() {
		return reflect.Value{}, &NotAssignableError{
			Message: "Nil",
			Kind:    rv.Kind(),
			Type:    rv.Type(),
		}
	}

	return reflect.Indirect(rv), nil
}
