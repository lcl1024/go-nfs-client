package xdr

import (
	"fmt"
	"io"
	"reflect"
	"encoding/binary"
)

func Uint32(b []byte) (uint32, []byte) {
	return binary.BigEndian.Uint32(b[0:4]), b[4:]
}

func Read(r io.Reader, val interface{}) error {
	if err := read(r, reflect.ValueOf(val)); err != nil {
		if err == io.EOF {
			return nil
		}
		return err
	}
	return nil
}

func read(r io.Reader, v reflect.Value) error {
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	fmt.Println("value:", v)
	switch t := v.Type(); t.Kind() {
	case reflect.Struct:
		for i := 0; i < v.NumField() ; i++ {
			if err := read(r, v.Field(i)); err != nil {
				return err
			}
		}
	case reflect.Uint32:
		var val uint32
		if err := binary.Read(r, binary.BigEndian, &val); err != nil {
			return err
		}
		v.SetUint(uint64(val))
	default:
		return fmt.Errorf("rpc.read: invalid type: %v ", t.String())
	}
	return nil
}
