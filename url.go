package slack

import (
	"encoding"
	"fmt"
	"log"
	"net/url"
	"reflect"
)

// Encode marshals the input as a series of URL query parameters for an http.Request
// All zero values are ommited from the URL.
//
// Example:
//
//  type Options{
//		Name string `url:"name"`
//	}
//
func Encode(o interface{}) url.Values {
	if o == nil {
		return nil
	}
	t := reflect.TypeOf(o)
	rv := reflect.ValueOf(o)
	if t.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return nil
		}
		return Encode(rv.Elem().Interface())
	}

	param := make(url.Values)
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		fv := rv.Field(i)
		if fv.Interface() == reflect.Zero(f.Type).Interface() {
			continue
		}
		tag := f.Tag.Get("url")
		if tag == "" {
			continue
		}
		switch t := fv.Interface().(type) {
		case encoding.TextMarshaler:
			s, err := t.MarshalText()
			if err != nil {
				log.Printf("errors with text marshaler: %v", err)
				param.Add(tag, fmt.Sprintf("%v", fv.Interface()))
			} else {
				param.Add(tag, string(s))
			}
		case interface{}:
			param.Add(tag, fmt.Sprintf("%v", fv.Interface()))
		}
	}

	return param
}
