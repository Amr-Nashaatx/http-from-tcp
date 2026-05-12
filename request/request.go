package request

import (
	"fmt"
	"io"
	"slices"
	"strings"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	Method        string
	RequestTarget string
	HttpVersion   string
}

var InvalidMethodError error = fmt.Errorf("Invalid request method")
var InvalidReqLineFormat error = fmt.Errorf("invalid request line format")
var InvalidTargetError error = fmt.Errorf("Target in request line has invalid format")
var InvalidHttpVersionError error = fmt.Errorf("Http version in request line has invalid format")

func parseRequestLine(line string) (*RequestLine, error) {
	methods := []string{"GET", "POST", "DELETE", "PUT", "PATCH"}

	parts := strings.Split(line, " ")
	if len(parts) != 3 {
		return nil, InvalidReqLineFormat
	}
	reqMethod := parts[0]
	reqTarget := parts[1]
	httpVer := parts[2]

	if ok := slices.Contains(methods, reqMethod); !ok {
		return nil, InvalidMethodError
	}

	if ok := strings.HasPrefix(reqTarget, "/"); !ok {
		return nil, InvalidTargetError
	}

	if ok := strings.HasPrefix(httpVer, "HTTP/"); !ok {
		return nil, InvalidHttpVersionError
	}

	httpVer = strings.TrimPrefix(httpVer, "HTTP/")
	return &RequestLine{reqMethod, reqTarget, httpVer}, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	reqStr, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(reqStr), "\r\n")

	reqLine, parseErr := parseRequestLine(lines[0])
	if parseErr != nil {
		return nil, parseErr
	}
	return &Request{*reqLine}, nil
}
