package request

import (
	"errors"
	"io"
)

const parserDone = 2
const parserInit = 1

type Request struct {
	RequestLine RequestLine
	state       int // 1 initialized, 2 done
}
type RequestLine struct {
	Method        string
	RequestTarget string
	HttpVersion   string
}

func (r *Request) parse(data []byte) (int, error) {
	var nRead int
	var err error
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

	return nRead, err
}

func read(reader io.Reader, acc *[]byte) error {
	readBuffer := make([]byte, 8)
	nRead, err := reader.Read(readBuffer)
	*acc = append(*acc, readBuffer[:nRead]...)
	if errors.Is(err, io.EOF) {
		return err
	}
	return err
}
func RequestFromReader(reader io.Reader) (*Request, error) {
	req := &Request{state: parserInit}
	data := make([]byte, 0)
	read(reader, &data)
	for req.state != parserDone {
		n, parseErr := req.parse(data)
		if parseErr != nil {
			return nil, parseErr
		}
		if n == 0 {
			read(reader, &data)
		}
	}
	return req, nil
}
