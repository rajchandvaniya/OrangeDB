package core

import "errors"

func Decode(data []byte) (interface{}, error) {
	if len(data) == 0 {
		return nil, errors.New("empty data")
	}
	value, _, err := decodeOne(data)
	return value, err
}

func decodeOne(data []byte) (interface{}, int, error) {
	if len(data) == 0 {
		return nil, 0, errors.New("empty data")
	}
	switch data[0] {
	case '+':
		return decodeSimpleString(data)
	case '$':
		return decodeBulkString(data)
	case ':':
		return decodeInt64(data)
	case '-':
		return decodeError(data)
	case '*':
		return decodeArray(data)
	default:
		return nil, 0, errors.New("corrupted data")
	}
}

// reads a RESP encoded simple string from data and returns
// the string, the delta, and the error
func decodeSimpleString(data []byte) (string, int, error) {
	// first character: +
	pos := 1
	for ; data[pos] != '\r'; pos++ {
	}
	return string(data[1:pos]), pos + 2, nil
}

// reads a RESP encoded string from data and returns
// the string, the delta, and the error
func decodeBulkString(data []byte) (string, int, error) {
	// first character: $
	pos := 1

	// reading the length
	len, delta := readLength(data[pos:])
	pos += delta

	// reading `len` bytes as string
	return string(data[pos:(pos + len)]), pos + len + 2, nil
}

// reads a RESP encoded integer from data and returns
// the intger value, the delta, and the error
func decodeInt64(data []byte) (int64, int, error) {
	// first character: :
	pos := 1
	var num int64
	for ; data[pos] != '\r'; pos++ {
		num = num*10 + int64(data[pos]-'0')
	}
	return num, pos + 2, nil
}

// reads a RESP encoded error from data and returns
// the error string, the delta, and the error
func decodeError(data []byte) (string, int, error) {
	// first character: -
	return decodeSimpleString(data)
}

// reads a RESP encoded array from data and returns
// the array, the delta, and the error
func decodeArray(data []byte) ([]interface{}, int, error) {
	// first character: *
	pos := 1

	// reading the length
	len, delta := readLength(data[pos:])
	pos += delta

	elems := make([]interface{}, len)

	for i := range elems {
		elem, delta, err := decodeOne(data[pos:])
		if err != nil {
			return nil, 0, err
		}
		elems[i] = elem
		pos += delta
	}
	return elems, pos, nil
}

// reads the length typically the first integer of the string
// until hit by an non-digit byte and returns
// the integer and the delta = length + 2 (CRLF)
func readLength(data []byte) (int, int) {
	len, pos := 0, 0
	for ; data[pos] != '\r'; pos++ {
		len = len*10 + int(data[pos]-'0')
	}
	return len, pos + 2
}
