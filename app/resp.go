package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"
)

// Errors
var (
	ErrInvalidEncoding  = errors.New("invalid encoding")
	ErrTypeNotSupported = errors.New("RESP is type not supported")
)

const (
	BULK_STR_TYPE = '$'
	ARRAY_TYPE    = '*'
)

const CR = '\r'
const LF = '\n'

var CRLF = []byte{CR, LF}

type RESPDecoder struct {
	Reader io.Reader
}

// RESP Types
type BulkString = []byte
type Array = []interface{}

func (d *RESPDecoder) Decode() (interface{}, error) {
	t := make([]byte, 1)
	_, err := d.Reader.Read(t)
	if err != nil {
		fmt.Printf("Error while reading the type: %v\n", err)
		return nil, err
	}
	switch t[0] {
	case BULK_STR_TYPE:
		return d.DecodeBulkString()
	case ARRAY_TYPE:
		return d.DecodeArray()
	default:
		return nil, ErrTypeNotSupported
	}
}

func (d *RESPDecoder) DecodeArray() (Array, error) {
	s, err := getLength(d.Reader)
	if err != nil {
		fmt.Printf("Error in extracting size: %v\n", err)
		return nil, ErrInvalidEncoding
	}
	fmt.Printf("Size of the Array: %d\n", s)

	arr := make(Array, 0)
	for _ = range s {
		elem, err := d.Decode()

		if err != nil {
			if err == io.EOF {
				break
			} else {
				fmt.Printf("Error in decoding: %v\n", err)
				return nil, err
			}
		}
		arr = append(arr, elem)
	}

	return arr, nil
}

func (d *RESPDecoder) DecodeBulkString() (BulkString, error) {
	l, err := getLength(d.Reader)
	if err != nil {
		fmt.Printf("Error in extracting length: %v\n", err)
		return nil, err
	}
	fmt.Printf("Length of the string: %d\n", l)

	// Read Bytes
	str := make([]byte, l)
	_, err = d.Reader.Read(str)
	if err != nil {
		return nil, err
	}

	// Consume CRLF bytes
	crlf := make([]byte, 2)
	_, err = d.Reader.Read(crlf)
	if err != nil {
		return nil, err
	} else if !bytes.Equal(crlf, CRLF) {
		fmt.Printf("CRLF not found,  %v \n", crlf)
		return nil, ErrInvalidEncoding
	}

	return str, nil
}

func getLength(r io.Reader) (int, error) {
	buff := make([]byte, 1)
	l := ""
	for {
		_, err := r.Read(buff)
		if err != nil {
			return 0, err
		}

		if buff[0] == CR {
			break
		} else {
			l += string(buff)
		}
	}

	_, err := r.Read(buff)
	if err != nil {
		return 0, err
	}

	return strconv.Atoi(l)
}

type RESPEncoder struct {
	Writer io.Writer
}

func (e *RESPEncoder) Encode(val interface{}) error {
	switch v := val.(type) {
	case []byte:
		// ByteString
		bs := make([]byte, 0)
		bs = append(bs, '$')
		bs = append(bs, []byte(strconv.Itoa(len(v)))...)
		bs = append(bs, CRLF...)
		bs = append(bs, v...)
		bs = append(bs, CRLF...)

		_, err := e.Writer.Write(bs)
		return err
	case string:
		s := "+"
		s += v
		s += string(CRLF)
		_, err := e.Writer.Write([]byte(s))
		return err
	default:
		return fmt.Errorf("%w: %v", ErrTypeNotSupported, v)

	}
}
