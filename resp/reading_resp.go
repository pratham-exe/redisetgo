package resp

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

const (
	ARRAY   = '*'
	BULK    = '$'
	STRING  = '+'
	INTEGER = ':'
	ERROR   = '-'
)

type Client_input struct {
	Tipe  string
	Bulk  string
	Str   string
	Array []Client_input
}

type Input_buffer struct {
	reader *bufio.Reader
}

func Resp_input_buffer(con io.Reader) *Input_buffer {
	return &Input_buffer{reader: bufio.NewReader(con)}
}

func Read_buffer(rb *Input_buffer) Client_input {
	bite, err := rb.reader.ReadByte()
	if err != nil {
		fmt.Println("Read Buffer Error: ", err)
		return Client_input{}
	}

	switch bite {
	case ARRAY:
		return read_array(rb)
	case BULK:
		return read_bulk(rb)
	default:
		fmt.Printf("IDK: %v", string(bite))
		return Client_input{}
	}
}

func read_array(rb *Input_buffer) Client_input {
	ci := Client_input{}
	ci.Tipe = "Array"

	size := read_size(rb)
	ci.Array = make([]Client_input, size)

	for i := 0; i < size; i++ {
		ci_read := Read_buffer(rb)
		ci.Array = append(ci.Array, ci_read)
	}

	return ci
}

func read_bulk(rb *Input_buffer) Client_input {
	ci := Client_input{}
	ci.Tipe = "Bulk"

	size := read_size(rb)
	Bulk_str := make([]byte, size)
	rb.reader.Read(Bulk_str)
	ci.Bulk = string(Bulk_str)

	_ = read_line(rb)

	return ci
}

func read_size(rb *Input_buffer) int {
	line_crlf := read_line(rb)

	i64, err := strconv.ParseInt(string(line_crlf), 10, 64)
	if err != nil {
		fmt.Println("Parse Int Error: ", err)
		return 0
	}

	return int(i64)
}

func read_line(rb *Input_buffer) (line []byte) {
	for {
		bite, err := rb.reader.ReadByte()
		if err != nil {
			fmt.Println("Read Line Error: ", err)
			return nil
		}

		line = append(line, bite)
		if len(line) >= 2 && line[len(line)-2] == '\r' {
			break
		}
	}
	return line[:len(line)-2]
}
