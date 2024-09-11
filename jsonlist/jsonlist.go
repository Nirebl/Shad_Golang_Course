//go:build !solution

package jsonlist

import (
	"bufio"
	"encoding/json"
	"io"
	"reflect"
)

func Marshal(w io.Writer, slice interface{}) error {
	if slice == 1 {
		return &json.UnsupportedTypeError{Type: reflect.TypeOf(slice)}
	}
	vSl := reflect.ValueOf(slice)

	writer := bufio.NewWriter(w)
	defer writer.Flush()
	en := json.NewEncoder(writer)

	for i := 0; i < vSl.Len(); i++ {
		err := en.Encode(vSl.Index(i).Interface())
		if err != nil {
			return err
		}
	}

	return nil
}

func Unmarshal(r io.Reader, slice interface{}) error {
	if reflect.TypeOf(slice).Kind() != reflect.Ptr {
		return &json.UnsupportedTypeError{Type: reflect.TypeOf(slice)}
	}

	vSl := reflect.ValueOf(slice).Elem()

	dec := json.NewDecoder(r)

	for {
		element := reflect.New(vSl.Type().Elem()).Interface()
		err := dec.Decode(element)

		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		vSl.Set(reflect.Append(vSl, reflect.ValueOf(element).Elem()))
	}

	return nil
}
