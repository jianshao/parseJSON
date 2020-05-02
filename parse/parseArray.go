package parse

import (
	"fmt"
	"strconv"
)

type jsonValueArray struct {
	elementType int
	elements    []Value
}

func (a *jsonValueArray)parse(buf []byte, start int) (int, error) {

	start = skipInvalidChars(buf, start)
	if buf[start] != ValueArrayStartChar {
		return -1, fmt.Errorf("")
	}

	a.elementType = InvalidType
	if a.elements == nil {
		a.elements = make([]Value, 0)
	}

	pos := start + 1
	var value Value
	var err error
	for pos < len(buf) {
		if value, pos, err = parseValue(buf, pos); err != nil {
			return -1, err
		}

		if a.elementType == InvalidType {
			a.elementType = value.getType()
		} else if a.elementType != value.getType() {
			return -1, fmt.Errorf("element in array type is %d, but fount %d", a.elementType, value.getType())
		}
		a.elements = append(a.elements, value)
		fmt.Println("array element: ", value)

		pos = skipInvalidChars(buf, pos)
		if buf[pos] == ValueArrayEndChar {
			start = pos + 1
			break
		} else if buf[pos] == ValueEndSep {
			pos = pos + 1
		} else {
			return -1, fmt.Errorf("expect \",\" at end, but not found")
		}
	}
	return start, nil
}

func (a *jsonValueArray)getType() int {
	return Array
}

func (a *jsonValueArray)get(key string) (Value, error) {
	if index, err := strconv.Atoi(key); err == nil {
		if index == 0 {
			return nil, fmt.Errorf("index in array begin with 1")
		}
		if len(a.elements) >= index {
			return a.elements[index - 1], nil
		} else {
			return nil, fmt.Errorf("the pos should be less than %d but %s", len(a.elements) + 1, key)
		}
	}
	return nil, fmt.Errorf("expect a number but %s", key)
}