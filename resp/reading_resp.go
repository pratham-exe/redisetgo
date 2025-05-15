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
	Num   int
	Array []Client_input
}

type Input_buffer struct {
	reader *bufio.Reader
}

func Resp_input_buffer(con io.Reader) *Input_buffer {
	return &Input_buffer{reader: bufio.NewReader(con)}
}

func Read_buffer(rb *Input_buffer) (Client_input, error) {
	bite, err := rb.reader.ReadByte()
	if err != nil {
		fmt.Println("Read Buffer Error: ", err)
		return Client_input{}, err
	}

	switch bite {
	case ARRAY:
		return read_array(rb)
	case BULK:
		return read_bulk(rb)
	default:
		fmt.Printf("IDK: %v", string(bite))
		return Client_input{}, nil
	}
}

func read_array(rb *Input_buffer) (Client_input, error) {
	ci := Client_input{}
	ci.Tipe = "array"

	size, err := read_size(rb)
	if err != nil {
		return ci, err
	}

	ci.Array = make([]Client_input, size)

	for i := 0; i < size; i++ {
		ci_read, err := Read_buffer(rb)
		if err != nil {
			return ci, err
		}

		ci.Array[i] = ci_read
	}

	return ci, nil
}

func read_bulk(rb *Input_buffer) (Client_input, error) {
	ci := Client_input{}
	ci.Tipe = "bulk"

	size, err := read_size(rb)
	if err != nil {
		return ci, err
	}

	Bulk_str := make([]byte, size)
	rb.reader.Read(Bulk_str)
	ci.Bulk = string(Bulk_str)

	_, _ = read_line(rb)

	return ci, nil
}

func read_size(rb *Input_buffer) (int, error) {
	line_crlf, err := read_line(rb)
	if err != nil {
		return 0, err
	}

	i64, err := strconv.ParseInt(string(line_crlf), 10, 64)
	if err != nil {
		fmt.Println("Parse Int Error: ", err)
		return 0, err
	}

	return int(i64), nil
}

func read_line(rb *Input_buffer) (line []byte, err error) {
	for {
		bite, err := rb.reader.ReadByte()
		if err != nil {
			fmt.Println("Read Line Error: ", err)
			return nil, err
		}

		line = append(line, bite)
		if len(line) >= 2 && line[len(line)-2] == '\r' {
			break
		}
	}
	return line[:len(line)-2], nil
}
