package parse

import "fmt"

type jsonValueInt struct {
	v int
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

func (i *jsonValueInt)parse(buf []byte, start int) (int, error) {
	start = skipInvalidChars(buf, start)
	base := 10
	if buf[start] == '0' {
		if len(buf) - 1 == start {
			/* 0是最后一个字符，这是异常情况 */
			return -1, fmt.Errorf("the last char can not be \"0\"")
		} else {
			next := buf[start + 1]
			if next == 'x' {
				/* 0x 表示16进制 */
				start = start + 2
				base = 16
			}
		}
	}

	value := 0
	for ; start < len(buf); start++ {
		v := getIntFromByte(buf[start], base)
		if v != -1 {
			value = value * base + v
		} else {
			if buf[start] != ValueEndSep &&
				buf[start] != ValueObjectEndChar &&
				buf[start] != ValueArrayEndChar &&
				buf[start] != '\n' &&
				buf[start] != ' ' &&
				buf[start] != '\t' {
				return -1, fmt.Errorf("invalid number")
			} else {
				break
			}
		}
	}

	i.v = value
	/* 不管是遇到逗号还是最终结束，都是可以接受的结果 */
	return start, nil
}

func (i *jsonValueInt)getType() int {
	return Int
}

func (i *jsonValueInt)get(key string) (Value, error) {
	return i, nil
}