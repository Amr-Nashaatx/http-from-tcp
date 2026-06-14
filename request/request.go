package request

import (
	"errors"
	"fmt"
	"io"
	"slices"
	"strconv"

	"github.com/Amr-Nashaatx/http-from-tcp/headers"
)

const parserInit = 1
const parseHeaders = 2
const parseBody = 3
const parserDone = 4
const maxHeadersBufferSize = 8192 // 8KB

// The request struct wich contains all parsed info about the incoming HTTP request
type Request struct {
	RequestLine RequestLine
	Headers     headers.Headers
	Body        []byte
	state       int
}

func (r *Request) parse(data []byte) (int, error) {
	switch r.state {
	case parserInit:
		// Parse request-line
		reqLine, n, parseErr := parseRequestLine(data)
		if parseErr != nil {
			return 0, parseErr
		}
		/* 	n = 0 is a signal that means parseRequestLine didn't find a CLRF,
		which means the bytes passed don't constitute a full herader line
		in this case it returns n = 0 to signal (need more data)
		*/
		if n == 0 {
			return 0, nil
		}
		// Update the request object (the reciever of this function)
		r.RequestLine = *reqLine
		r.state = parseHeaders

		return n, parseErr
	case parseHeaders:
		// Parse headers
		hn, hDone, hParseErr := r.Headers.Parse(data)

		if hParseErr != nil {
			return 0, hParseErr
		}
		// check for n == 0, for the same reason as above
		if hn == 0 && hDone == false {
			return 0, nil
		}

		if hDone == true {
			r.state = parseBody
		}
		return hn, hParseErr
	case parseBody:
		contentLengthHeader := r.Headers.Get("content-length")
		// if no content-length field then request has no body
		if contentLengthHeader == "" {
			r.state = parserDone
			return 0, nil
		}
		contentLength, convErr := strconv.Atoi(contentLengthHeader)
		if convErr != nil {
			return 0, fmt.Errorf("Invalid content-length header")
		}
		r.Body = slices.Clone(data)
		if len(r.Body) < contentLength {
			return 0, nil
		} else if len(r.Body) > contentLength {
			return len(r.Body), fmt.Errorf("body excceeds declared length")
		} else {
			r.state = parserDone
			return len(data), nil
		}
	default:
		return 0, fmt.Errorf("Invalid parser state")
	}
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
	req := &Request{state: parserInit, Headers: headers.NewHeaders()}
	data := make([]byte, 0)
	read(reader, &data)

	// parse request-line
	for req.state != parserDone {
		n, parseErr := req.parse(data)
		if parseErr != nil {
			return nil, parseErr
		}
		if n == 0 {
			if req.state == parserDone {
				break
			}
			readErr := read(reader, &data)
			/* if buffer exceeds max size, reject the request.
			this handles the case when sender keeps sending header data
			without any line terminator (\r\n) in which case memory will accumulate indefinitely
			*/
			if len(data) > maxHeadersBufferSize {
				return nil, fmt.Errorf("headers buffer exceeded maximum size of %d bytes", maxHeadersBufferSize)
			}
			// if we reached EOF while still parsing headers, that's an error
			if errors.Is(readErr, io.EOF) && req.state == parseHeaders {
				return nil, fmt.Errorf("unexpected EOF: incomplete headers")
			}
			if errors.Is(readErr, io.EOF) && req.state == parseBody {
				return nil, fmt.Errorf("unexpected EOF: imcomplete body")
			}
			if readErr != nil && !errors.Is(readErr, io.EOF) {
				return nil, readErr
			}

			// when n > 0, it mean we finished a stage in parsing, so we flush out the data buffer
		} else if n > 0 {
			data = slices.Clone(data[n:])
		}
	}
	return req, nil
}
