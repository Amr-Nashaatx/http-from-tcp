package request

import (
	"errors"
	"fmt"
	"io"
)

type Request struct {
	RequestLine RequestLine
	state       int // 1 initialized, 2 done
}
type RequestLine struct {
	Method        string
	RequestTarget string
	HttpVersion   string
}

func NewRequest() *Request {
	return &Request{state: 1}
}
func (r *Request) parse(data []byte) (int, error) {
	var nRead int
	var err error
	switch r.state {
	case 1:
		reqLine, n, parseErr := parseRequestLine(data)
		if parseErr != nil {
			return 0, parseErr
		}
		if n == 0 {
			return 0, nil
		}
		nRead = n
		err = parseErr
		r.RequestLine = *reqLine
		r.state = 2
	case 2:
		return 0, fmt.Errorf("error: trying to read data in a done state")
	default:
		return 0, fmt.Errorf("error: unknown state")

	}

	return nRead, err
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	readBuffer := make([]byte, 8)
	acc := make([]byte, 0)
	req := NewRequest()
	for {
		nRead, err := reader.Read(readBuffer)
		acc = append(acc, readBuffer[:nRead]...)
		if errors.Is(err, io.EOF) {
			fmt.Println("reached eof")
			break
		}
		n, parseErr := req.parse(acc)
		if parseErr != nil {
			fmt.Println(parseErr)
			return nil, parseErr
		}
		if n == 0 {
			continue
		} else {
			break
		}
	}
	return req, nil
}
