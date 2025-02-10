package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

// DumpToString interface to string
func DumpToString(v interface{}) string {

	switch v.(type) {
	case string:
		return v.(string)
	case []byte:
		return string(v.([]byte))
	case bytes.Buffer:
		x := v.(bytes.Buffer)
		return x.String()
	case io.ReadCloser:
		return IoReadCloserToBuffer(v.(io.ReadCloser)).String()
	}

	buff := &bytes.Buffer{}
	json.NewEncoder(buff).Encode(v)
	return buff.String()
}

// DebugPrint for debug print on terminal
func DebugPrint(v ...interface{}) {
	for _, x := range v {
		fmt.Println(DumpToString(x))
	}
}

// Data To json Bytes
func ToJSONByte(v interface{}) []byte {

	switch v.(type) {
	case []byte:
		return v.([]byte)
	case string:
		return []byte(v.(string))
	}

	buff := &bytes.Buffer{}
	json.NewEncoder(buff).Encode(v)
	return buff.Bytes()
}

func ToBuffer(v interface{}) *bytes.Buffer {
	buff := &bytes.Buffer{}
	switch v.(type) {
	case []byte:
		buff.Write(v.([]byte))
		return buff
	case string:
		buff.WriteString(v.(string))
		return buff
	}

	json.NewEncoder(buff).Encode(v)
	return buff
}

func IoReadCloserToBuffer(closer io.ReadCloser) *bytes.Buffer {
	buf := new(bytes.Buffer)
	buf.ReadFrom(closer)
	return buf
}

func PrettyToString(v interface{}) string {
	str, ok := v.(string)
	if ok {
		return str
	}

	buff := &bytes.Buffer{}
	enc := json.NewEncoder(buff)
	enc.SetIndent("", " ")
	enc.Encode(v)
	return buff.String()
}

func PrettyPrint(v interface{}) {
	fmt.Println(PrettyToString(v))
}

// ToJSONMinify for insignificant space characters elided
func ToJSONMinify(v interface{}) []byte {
	buff := bytes.Buffer{}
	b := ToJSONByte(v)

	json.Compact(&buff, b)
	return buff.Bytes()
}
