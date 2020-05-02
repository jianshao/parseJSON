package parse

import (
	"fmt"
)

const (
	InvalidType = 0
	Int = 1
	String = 2
	Array = 3
	Object = 4

	ValueEndSep = ','
	ValueStringEndChar = '"'
	ValueStringStartChar = '"'
	ValueObjectStartChar = '{'
	ValueObjectEndChar = '}'
	ValueArrayStartChar = '['
	ValueArrayEndChar = ']'
	KeyValueDelimiter = ':'
)


type Value interface {
	//Print(key string)
	get(key string) (Value, error)
	parse([]byte, int) (int, error)
	getType() int
}


func skipInvalidChars(buf []byte, s int) int {

	for s < len(buf) {
		switch  {
		case buf[s] > 32 && buf[s] != 127:
			return s
		default:
			s++
		}
	}
	return s
}


func parseValue(buf []byte, start int) (Value, int, error) {

	var err error
	var value Value
	start = skipInvalidChars(buf, start)
	switch buf[start] {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		value = &jsonValueInt{}
	case ValueStringStartChar:
		value = &jsonValueString{}
	case ValueObjectStartChar:
		value = &jsonValueObject{}
	case ValueArrayStartChar:
		value = &jsonValueArray{}
	default:
		return nil, -1, fmt.Errorf("unknown type: %s when parse value", string(buf[:start+1]))
	}

	if start, err = value.parse(buf, start); err == nil {
		return value, start, nil
	}
	return nil, -1, err
}
