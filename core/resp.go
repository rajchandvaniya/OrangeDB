package core

import (
	"errors"
	"fmt"
)

const RESP_NIL = "$-1\r\n"

func Encode(value interface{}, isSimple bool) []byte {
	switch typ := value.(type) {
	case string:
		if isSimple {
			return []byte(fmt.Sprintf("+%s\r\n", typ))
		} else {
			return []byte(fmt.Sprintf("$%d\r\n%s\r\n", len(typ), typ))
		}
	case int64:
		return []byte(fmt.Sprintf(":%v\r\n", typ))
	default:
		return []byte(RESP_NIL)
	}
}

func EncodeError(err error) []byte {
	return []byte(fmt.Sprintf("-%s\r\n", err))
}

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
	case '-':
		return decodeError(data)
	case '$':
		return decodeBulkString(data)
	case ':':
		return decodeInt64(data)
	case '*':
		return decodeArray(data)
	default:
		return nil, 0, errors.New("invalid data type received, RESP supports (+,-,:,$,*)")
	}
}

func decodeSimpleString(data []byte) (string, int, error) {
	// first character: +
	pos := 1
	for ; data[pos] != '\r'; pos++ {
	}
	return string(data[1:pos]), pos + 2, nil
}

func decodeError(data []byte) (string, int, error) {
	// first character: -
	return decodeSimpleString(data)
}

func decodeBulkString(data []byte) (string, int, error) {
	// reading the length
	len, pos, err := decodeInt64(data)
	if err != nil {
		return "", 0, err
	}

	// reading `len` bytes as string
	return string(data[pos:(pos + int(len))]), pos + int(len) + 2, nil
}

func decodeInt64(data []byte) (int64, int, error) {
	// first character: :
	pos := 1
	var num int64 = 0
	var isNegative = false
	if data[pos] == '-' {
		isNegative = true
		pos++
	}

	for ; data[pos] != '\r'; pos++ {
		if !(data[pos] >= '0' && data[pos] <= '9') {
			return 0, 0, errors.New("non digit character found")
		}
		num = num*10 + int64(data[pos]-'0')
	}

	if isNegative {
		num *= -1
	}
	return num, pos + 2, nil
}

func decodeArray(data []byte) ([]interface{}, int, error) {
	numElems, pos, err := decodeInt64(data)
	if err != nil {
		return nil, 0, err
	}

	// null array
	if numElems == -1 {
		return nil, pos + 2, nil
	}

	elems := make([]interface{}, numElems)
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
