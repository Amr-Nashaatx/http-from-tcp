package request

import (
	"bytes"
	"slices"
	"strings"
)

func parseRequestLine(data []byte) (*RequestLine, int, error) {
	methods := []string{"GET", "POST", "DELETE", "PUT", "PATCH"}

	delimPos := bytes.Index(data, []byte("\r\n"))
	if delimPos == -1 {
		return nil, 0, nil
	}
	parts := strings.Split(string(data[:delimPos]), " ")
	if len(parts) != 3 {
		return nil, 0, InvalidReqLineFormat
	}
	reqMethod := parts[0]
	reqTarget := parts[1]
	httpVer := parts[2]

	if !slices.Contains(methods, reqMethod) {
		return nil, 0, InvalidMethodError
	}

	if !strings.HasPrefix(reqTarget, "/") {
		return nil, 0, InvalidTargetError
	}

	if !strings.HasPrefix(httpVer, "HTTP/") {
		return nil, 0, InvalidHttpVersionError
	}

	httpVer = strings.TrimPrefix(httpVer, "HTTP/")
	return &RequestLine{reqMethod, reqTarget, httpVer}, delimPos + 2, nil
}

func PushChunkToSlice(s []byte, chunk []byte, chunkLength int) []byte {
	lastElemIdx := slices.Index(s, byte(0))
	if lastElemIdx == -1 {
		lastElemIdx = len(s)
	}
	if !(len(s)-lastElemIdx >= chunkLength) {
		addSpace := make([]byte, chunkLength)
		s = append(s, addSpace...)
	}
	for i := range len(chunk) {
		s[i+lastElemIdx] = chunk[i]
	}

	return s
}
