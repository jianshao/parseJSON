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
)

type jsonValueInt int64
type jsonValueString string
type jsonValueObject struct {
	fields map[jsonKey]*jsonValue
}

type jsonValueArray struct {
	array []*jsonValue
}

type jsonValue struct {
	Type int
	Value interface{}
}

type jsonKey struct {
	Type int
	Value interface{}
}



func buildJsonValue(Type int, value interface{}) *jsonValue {
	switch Type {
	case Int:
		newValue := jsonValueInt(value.(int))
		return &jsonValue{Type:Type, Value:&newValue}
	case String:
		newValue := jsonValueString(value.([]byte))
		return &jsonValue{Type:Type, Value:&newValue}
	case Object:
		return &jsonValue{Type:Type, Value:value}
	default:
		return nil
	}
}


func getIntFromByte(b byte, base int) int {
	result := -1
	switch  {
	case b <= '9' && b >= '0':
		result = int(b - '0')
	case b >= 'a' && b <= 'f':
		result = int(b - 'a' + 10)
	case b >= 'A' && b <= 'F':
		result = int(b - 'A' + 10)
	default:
		return -1
	}

	if result >= base {
		result = -1
	}

	return result
}

func parseValueInt(buf []byte, s int) (*jsonValue, int, error) {
	base := 10
	if buf[s] == '0' {
		if len(buf) - 1 == s {
			return buildJsonValue(Int, 0), s+1, nil
		} else if buf[s+1] == 'x' {
			s = s + 2
			base = 16
		}
	}

	value := 0
	for ; s < len(buf); s++ {
		v := getIntFromByte(buf[s], base)
		if v != -1 {
			value = value * base + v
		} else {
			break
		}
	}

	/* 不管是遇到逗号还是最终结束，都是可以接受的结果 */
	return buildJsonValue(Int, value), s, nil
}

func parseValueString(buf []byte, s int) (*jsonValue, int, error) {

	e := s + 1
	for ; e < len(buf); e++ {
		if buf[e] == ValueStringEndChar {
			break
		}
	}
	if e == len(buf) {
		return nil, -1, fmt.Errorf("expect \" at end")
	}

	return buildJsonValue(String, buf[s+1:e]), e+1, nil
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


func parseValueObject(buf []byte, s int) (*jsonValue, int, error) {

	object := &jsonValueObject{
		fields:make(map[jsonKey]*jsonValue),
	}
	s = s + 1
	for s < len(buf) {
		key, pos, err := parseKey(buf, s)
		if err != nil {
			return nil, -1, nil
		}
		fmt.Println("")

		s = pos
		pos = skipInvalidChars(buf, s)
		if pos >= len(buf) || buf[pos] != ':' {
			return nil, -1, fmt.Errorf("expect ':' at %s", string(buf[s:pos+1]))
		}

		value, pos, err := parseValue(buf, pos+1)
		if err != nil {
			return nil, -1, err
		}
		object.fields[*key] = value

		pos = skipInvalidChars(buf, pos)
		if buf[pos] == ValueEndSep {
			s = pos + 1
		} else if buf[pos] == ValueObjectEndChar {
			s = pos
			break
		} else {
			return nil, -1, fmt.Errorf("expect } at end")
		}
	}
	if buf[s] != ValueObjectEndChar {
		return nil, -1, fmt.Errorf("expect } at end")
	}
	return buildJsonValue(Object, object), s+1, nil
}

func parseValueArray(buf []byte, s int) (*jsonValue, int, error) {	var commonType = InvalidType
	result := make([]*jsonValue, 0)
	arrayValue := &jsonValueArray{
		array:result,
	}
	s++

	for {
		s = skipInvalidChars(buf, s)
		value, s, err := parseValue(buf, s)
		if err != nil {
			return nil, -1, err
		}
		if commonType == InvalidType {
			commonType = value.Type
		} else if commonType != value.Type {
			return nil, -1, fmt.Errorf("invalid type")
		}
		result = append(result, value)
		fmt.Println("array value: ", value)

		s = skipInvalidChars(buf, s)
		if buf[s] == ValueArrayEndChar {
			break
		} else if buf[s] != ValueEndSep {
			return nil, -1, fmt.Errorf("expect \",\" at end, but not found")
		} else {
			s++
		}
	}
	return buildJsonValue(Array, arrayValue), s+1, nil
}

func parseValue(buf []byte, s int) (*jsonValue, int, error) {

	s = skipInvalidChars(buf, s)
	switch buf[s] {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return parseValueInt(buf, s)
	case ValueStringStartChar:
		return parseValueString(buf, s)
	case ValueObjectStartChar:
		return parseValueObject(buf, s)
	case ValueArrayStartChar:
		return parseValueArray(buf, s)
	default:
		return nil, -1, fmt.Errorf("unknown type: %d when parse value", buf[s])
	}
}

func buildJsonKey(value []byte) *jsonKey {
	return &jsonKey{Type:String, Value:jsonValueString(value)}
}

func parseKey(buf []byte, s int) (*jsonKey, int, error) {
	var resultKey *jsonKey = nil
	var result int = -1

	s = skipInvalidChars(buf, s)
	if buf[s] == ValueStringStartChar {
		for e := s + 1 ; e < len(buf); e++ {
			if buf[e] == ValueStringEndChar {
				resultKey = buildJsonKey(buf[s+1:e])
				result = e + 1
				break
			}
		}
		if -1 == result {
			return nil, -1, fmt.Errorf("expect \" at end")
		}
	} else {
		return nil, -1, fmt.Errorf("expect \" after %s when parse key", string(buf[0:s+1]))
	}

	//fmt.Println("parse key:", string(buf[s+1: result-1]))
	return resultKey, result, nil
}

func parse(buf []byte) (*jsonValue, error) {

	root, pos, err := parseValue(buf, 0)
	if -1 == pos || nil == root {
		return nil, err
	}

	pos = skipInvalidChars(buf, pos)
	if pos != len(buf) {
		return nil, fmt.Errorf("extra data found")
	}

	return root, nil
}

func (object *jsonValueObject)Get(name string) *jsonValue {
	key := jsonKey{Type:String, Value:jsonValueString(name)}
	return object.fields[key]
}

func (json *jsonValue)Get(s []string) (*jsonValue, error) {
	root := json
	i := 0
	for ; i < len(s); i++ {
		if nil == root || root.Type != Object {
			break
		}
		root = root.Value.(*jsonValueObject).Get(s[i])
	}
	if root == nil || i != len(s) {
		return nil, fmt.Errorf("not exist")
	}

	return root, nil
}