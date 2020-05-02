package parse

import "fmt"

type jsonValueString struct {
	v string
}

func (s *jsonValueString)parse(buf []byte, start int) (int, error) {

	start = skipInvalidChars(buf, start)
	if buf[start] != ValueStringStartChar {
		return -1, fmt.Errorf("expect \" but %s", string(buf[start]))
	}

	for e := start + 1; e < len(buf); e++ {
		if buf[e] == ValueStringEndChar {
			s.v = string(buf[start+1:e])
			return e+1, nil
		}
	}
	return -1, fmt.Errorf("expect \" at end")
}

func (s *jsonValueString)getType() int {
	return String
}

func (s *jsonValueString)get(key string) (Value, error) {
	return s, nil
}