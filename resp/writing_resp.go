package resp

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

type Output_buffer struct {
	writer *bufio.Writer
}

func Resp_output_buffer(con io.Writer) *Output_buffer {
	return &Output_buffer{writer: bufio.NewWriter(con)}
}

func Write_buffer(wb *Output_buffer, rsp Client_input) int {
	bites := writerr(rsp)
	size, err := wb.writer.Write(bites)
	if err != nil {
		fmt.Println("Write Buffer Error: ", err)
		return 0
	}

	err = wb.writer.Flush()
	if err != nil {
		fmt.Println("Buffer Flush Error: ", err)
		return 0
	}

	return size
}

func writerr(rsp Client_input) []byte {
	switch rsp.Tipe {
	case "array":
		return writerr_array(rsp)
	case "bulk":
		return writerr_bulk(rsp)
	case "string":
		return writerr_string(rsp)
	case "error":
		return writerr_error(rsp)
	case "nill":
		return writerr_nill()
	case "integer":
		return writerr_integer(rsp)
	default:
		return []byte{}
	}
}

func writerr_string(rsp Client_input) []byte {
	var bites []byte
	bites = append(bites, STRING)
	bites = append(bites, rsp.Str...)
	bites = append(bites, '\r', '\n')

	return bites
}

func writerr_bulk(rsp Client_input) []byte {
	var bites []byte
	bites = append(bites, BULK)
	bites = append(bites, strconv.Itoa(len(rsp.Bulk))...)
	bites = append(bites, '\r', '\n')
	bites = append(bites, rsp.Bulk...)
	bites = append(bites, '\r', '\n')

	return bites
}

func writerr_array(rsp Client_input) []byte {
	var bites []byte
	bites = append(bites, ARRAY)
	bites = append(bites, strconv.Itoa(len(rsp.Array))...)
	bites = append(bites, '\r', '\n')

	for i := 0; i < len(rsp.Array); i++ {
		bites = append(bites, writerr(rsp.Array[i])...)
	}

	return bites
}

func writerr_error(rsp Client_input) []byte {
	var bites []byte
	bites = append(bites, ERROR)
	bites = append(bites, rsp.Str...)
	bites = append(bites, '\r', '\n')

	return bites
}

func writerr_nill() []byte {
	return []byte("$-1\r\n")
}

func writerr_integer(rsp Client_input) []byte {
	var bites []byte
	bites = append(bites, INTEGER)
	num := strconv.Itoa(rsp.Num)
	bites = append(bites, num...)
	bites = append(bites, '\r', '\n')

	return bites
}
