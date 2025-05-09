package resp

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

type Aof_method struct {
	file       os.File
	buffer     bufio.Reader
	mutex_lock sync.Mutex
}

func Create_aof(file_name string) Aof_method {
	file, err := os.OpenFile(file_name, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		fmt.Println("Open file error: ", err)
		return Aof_method{}
	}

	aof_method := Aof_method{
		file:   *file,
		buffer: *bufio.NewReader(file),
	}

	go func() {
		for {
			aof_method.mutex_lock.Lock()
			aof_method.file.Sync()
			aof_method.mutex_lock.Unlock()
			time.Sleep(time.Second)
		}
	}()

	return aof_method
}

func Aof_close(aof_method Aof_method) error {
	aof_method.mutex_lock.Lock()
	defer aof_method.mutex_lock.Unlock()

	return aof_method.file.Close()
}

func Write_aof(aof_method Aof_method, ci Client_input) error {
	aof_method.mutex_lock.Lock()

	_, err := aof_method.file.Write(writerr(ci))
	if err != nil {
		fmt.Println("Write file error: ", err)
		return err
	}

	aof_method.mutex_lock.Unlock()

	return nil
}

func Read_aof(aof_method Aof_method, call_main func(ci Client_input)) error {
	aof_method.mutex_lock.Lock()
	defer aof_method.mutex_lock.Unlock()

	input := Resp_input_buffer(&aof_method.file)

	for {
		client_input, err := Read_buffer(input)
		if err == nil {
			call_main(client_input)
		}
		if err == io.EOF {
			break
		}
	}

	return nil
}
