package resp

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

const (
	ARRAY = '*'
	BULK  = '$'
)

type Client_input struct {
	tipe  string
	bulk  string
	array []Client_input
}

type Buffer struct {
	reader *bufio.Reader
}

func Resp_buffer(con io.Reader) *Buffer {
	return &Buffer{reader: bufio.NewReader(con)}
}

func Read_buffer(rb *Buffer) Client_input {
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
		fmt.Printf("Unknown: %v", string(bite))
		return Client_input{}
	}
}

func read_array(rb *Buffer) Client_input {
	ci := Client_input{}
	ci.tipe = "Array"

	size := read_size(rb)
	ci.array = make([]Client_input, size)

	for i := 0; i < size; i++ {
		ci_read := Read_buffer(rb)
		ci.array = append(ci.array, ci_read)
	}

	return ci
}

func read_bulk(rb *Buffer) Client_input {
	ci := Client_input{}
	ci.tipe = "Bulk"

	size := read_size(rb)
	bulk_str := make([]byte, size)
	rb.reader.Read(bulk_str)
	ci.bulk = string(bulk_str)

	_ = read_line(rb)

	return ci
}

func read_size(rb *Buffer) int {
	line_crlf := read_line(rb)

	i64, err := strconv.ParseInt(string(line_crlf), 10, 64)
	if err != nil {
		fmt.Println("Parse Int Error: ", err)
		return 0
	}

	return int(i64)
}

func read_line(rb *Buffer) (line []byte) {
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
