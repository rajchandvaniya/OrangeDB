package core_test

import (
	"fmt"
	"testing"

	"github.com/rajchandvaniya/orangedb/core"
)

func TestSimpleStrings(t *testing.T) {
	cases := map[string]string{
		"+OK\r\n": "OK",
	}
	for input, expectedOutput := range cases {
		output, _ := core.Decode([]byte(input))
		if output != expectedOutput {
			t.Fail()
		}
	}
}

func TestBulkStrings(t *testing.T) {
	cases := map[string]string{
		"$5\r\nhello\r\n": "hello",
		"$0\r\n\r\n":      "",
	}
	for input, expectedOutput := range cases {
		output, _ := core.Decode([]byte(input))
		if output != expectedOutput {
			t.Fail()
		}
	}
}

func TestInt64(t *testing.T) {
	cases := map[string]int64{
		":100\r\n":  100,
		":0\r\n":    0,
		":-369\r\n": -369,
	}
	negativeCases := map[string]int64{
		":12rf\r\n": 0,
	}
	for input, expectedOutput := range cases {
		output, _ := core.Decode([]byte(input))
		if output != expectedOutput {
			t.Fail()
		}
	}
	for input, _ := range negativeCases {
		_, err := core.Decode([]byte(input))
		if (err == nil) || (err.Error() != "non digit character found") {
			t.Fail()
		}
	}
}

func TestError(t *testing.T) {
	cases := map[string]string{
		"-Error Message\r\n": "Error Message",
	}
	for input, expectedOutput := range cases {
		output, _ := core.Decode([]byte(input))
		if output != expectedOutput {
			t.Fail()
		}
	}
}

func TestErrorDecode(t *testing.T) {
	cases := map[string][]interface{}{
		"*0\r\n":                                                   {},
		"*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n":                     {"hello", "world"},
		"*3\r\n:1\r\n:2\r\n:3\r\n":                                 {1, 2, 3},
		"*5\r\n:1\r\n:-2\r\n:3\r\n:4\r\n$5\r\nhello\r\n":           {1, -2, 3, 4, "hello"},
		"*2\r\n*3\r\n:1\r\n:2\r\n:3\r\n*2\r\n+Hello\r\n-World\r\n": {[]int{1, 2, 3}, []interface{}{"Hello", "World"}},
		"*-1\r\n": nil,
	}
	for input, expectedOutput := range cases {
		output, _ := core.Decode([]byte(input))
		arrayOutput := output.([]interface{})
		if len(arrayOutput) != len(expectedOutput) {
			t.Fail()
		}
		for i := range expectedOutput {
			fmt.Println("expected", expectedOutput[i], "actual", arrayOutput[i])
			if fmt.Sprintf("%v", arrayOutput[i]) != fmt.Sprintf("%v", expectedOutput[i]) {
				t.Fail()
			}
		}
	}
}
