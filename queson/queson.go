// Package queson implements QUESON for Go.
//
// QUESON is an alternative notation for the JSON data model, designed to look
// less noisy in URLs. See https://tcard.github.io/queson
package queson

import (
	"encoding/json"

	"github.com/tcard/queson-go/internal/json2queson"
	"github.com/tcard/queson-go/internal/queson2json"
)

// ToJSONBytes compiles QUESON source code into JSON.
func ToJSONBytes(src []byte) ([]byte, error) {
	v, err := queson2json.Parse("", src)
	if err != nil {
		return nil, err
	}
	return v.([]byte), nil
}

// FromJSONBytes compiles JSON source code into QUESON.
func FromJSONBytes(src []byte) ([]byte, error) {
	v, err := json2queson.Parse("", src, json2queson.Recover(false))
	if err != nil {
		return nil, err
	}
	return v.([]byte), nil
}

// ToJSON compiles QUESON source code into JSON.
func ToJSON(src string) (string, error) {
	v, err := ToJSONBytes([]byte(src))
	return string(v), err
}

// FromJSON compiles JSON source code into QUESON.
func FromJSON(src string) (string, error) {
	v, err := FromJSONBytes([]byte(src))
	return string(v), err
}

// Marshal marshals a value into QUESON via json.Marshal.
func Marshal(v interface{}) ([]byte, error) {
	js, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	return FromJSONBytes(js)
}

// Unmarshal unmarshals a value from QUESON via json.Unmarshal.
func Unmarshal(data []byte, v interface{}) error {
	src, err := ToJSONBytes(data)
	if err != nil {
		return err
	}
	return json.Unmarshal(src, v)
}
