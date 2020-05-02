package parse

import "fmt"

type jsonValueObject struct {
	fields map[jsonValueString]Value
}

func (o *jsonValueObject)parse(buf []byte, start int) (int, error) {
	if o.fields == nil {
		o.fields = make(map[jsonValueString]Value)
	}

	var pos int
	var err error
	var value Value
	start = start + 1
	for start < len(buf) {

		key := jsonValueString{}
		if pos, err = key.parse(buf, start); err != nil {
			return -1, err
		}
		//fmt.Println("key: ", key.v, " pos: ", pos)

		pos = skipInvalidChars(buf, pos)
		if -1 == pos || pos >= len(buf) || buf[pos] != KeyValueDelimiter {
			return -1, fmt.Errorf("expect \":\" after %s", buf[:pos])
		}

		if value, pos, err = parseValue(buf, pos+1); err != nil {
			return -1, err
		}
		o.fields[key] = value
		//fmt.Println("value: ", value)

		pos = skipInvalidChars(buf, pos)
		if buf[pos] == ValueEndSep {
			start = pos + 1
		} else if buf[pos] == ValueObjectEndChar {
			start = pos + 1
			break
		} else {
			return -1, fmt.Errorf("expect } at end")
		}
	}
	return start, nil
}

func (o *jsonValueObject)getType() int {
	return Object
}

func (o *jsonValueObject)get(key string) (Value, error) {
	k := jsonValueString{v:key}
	if v, ok := o.fields[k]; ok {
		return v, nil
	}

	return nil, fmt.Errorf("key: %s not existed", key)
}