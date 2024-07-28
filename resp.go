package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

const (
	STRING  = '+'
	ERROR   = '-'
	INTEGER = ':'
	BULK    = '$'
	ARRAY   = '*'
)

type Value struct {
	typ   string
	str   string
	num   int
	bulk  string
	array []Value
}

type Resp struct {
	reader *bufio.Reader
}

func NewResp(rd io.Reader) *Resp {
	return &Resp{reader: bufio.NewReader(rd)}
}

func (r *Resp) readLine() (line []byte, n int, err error) {
	for {
		b, err := r.reader.ReadByte()
		if err != nil {
			return nil, 0, err
		}
		n += 1
		line = append(line, b)
		if len(line) >= 2 && line[len(line)-2] == '\r' && line[len(line)-1] == '\n' {
			break
		}
	}
	return line[:len(line)-2], n, nil
}

func (r *Resp) readInteger() (x int, n int, err error) {
	line, n, err := r.readLine()
	if err != nil {
		return 0, n, err
	}

	i64, err := strconv.ParseInt(string(line), 10, 64)
	if err != nil {
		return 0, n, err
	}

	return int(i64), n, nil
}

func (r *Resp) Read() (Value, error) {
	_type, err := r.reader.ReadByte()
	if err != nil {
		return Value{}, err
	}

	switch _type {
	case ARRAY:
		return r.readArray()
	case BULK:
		return r.readBulk()
	default:
		fmt.Println("Unknown type: %v", _type)
		return Value{}, nil
	}
}

func (r *Resp) readArray() (Value, error) {
	v := Value{typ: "array"}

	n, _, err := r.readInteger()
	if err != nil {
		return v, err
	}
	v.array = make([]Value, n)
	for i := 0; i < n; i++ {
		v.array[i], err = r.Read()
		if err != nil {
			return v, err
		}
	}

	return v, nil
}

func (r *Resp) readBulk() (Value, error) {
	v := Value{typ: "bulk"}

	n, _, err := r.readInteger()
	if err != nil {
		return v, err
	}

	bulk := make([]byte, n)
	r.reader.Read(bulk)
	v.bulk = string(bulk)
	r.readLine()
	return v, nil
}

func (v Value) Marshal() []byte {
	switch v.typ {
	case "array":
		return v.marshalArray()
	case "bulk":
		return v.marshalBulk()
	case "string":
		return v.marshalString()
	case "null":
		return v.marshalNull()
	case "error":
		return v.marshalError()
	default:
		return []byte{}
	}
}

func (v Value) marshalString() []byte {
	return []byte("+" + v.str + "\r\n")
}

func (v Value) marshalBulk() []byte {
	return []byte("$" + strconv.Itoa(len(v.bulk)) + "\r\n" + v.bulk + "\r\n")
}

func (v Value) marshalArray() []byte {
	var res []byte
	res = append(res, '*')
	res = append(res, strconv.Itoa(len(v.array))...)
	res = append(res, '\r', '\n')
	for _, val := range v.array {
		res = append(res, val.Marshal()...)
	}
	return res
}

func (v Value) marshalNull() []byte {
	return []byte("$-1\r\n")
}

func (v Value) marshalError() []byte {
	return []byte("-" + v.str + "\r\n")
}

type Writer struct {
	writer io.Writer
}

func NewWriter(wr io.Writer) *Writer {
	return &Writer{writer: wr}
}

func (w *Writer) Write(v Value) error {
	var bytes = v.Marshal()
	_, err := w.writer.Write(bytes)
	if err != nil {
		return err
	}
	return nil
}
