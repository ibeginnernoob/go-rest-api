package utils

import "unsafe"

func StringToBytes(string_ string) (bytes []byte) {
	return unsafe.Slice(unsafe.StringData(string_), len(string_))
}
