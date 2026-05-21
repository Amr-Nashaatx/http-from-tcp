package request

import (
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type chunkReader struct {
	data            string
	numBytesPerRead int
	pos             int
}

func (cr *chunkReader) Read(buf []byte) (n int, err error) {
	// return EOF of pos reached end of file
	if cr.pos >= len(cr.data) {
		return 0, io.EOF
	}
	// read until endIndex
	endIndex := min(cr.pos+cr.numBytesPerRead, len(cr.data))

	nRead := copy(buf, []byte(cr.data)[cr.pos:endIndex])

	cr.pos += nRead
	return nRead, nil
}
func TestPushChunkToSlice(t *testing.T) {
	slc := make([]byte, 8)
	chk := []byte("hello")

	slc = PushChunkToSlice(slc, chk, 5)
	assert.Contains(t, string(slc), "hello")

	slc = []byte("hello")
	newSpace := make([]byte, 8)
	slc = append(slc, newSpace...)
	chk = []byte("world")

	slc = PushChunkToSlice(slc, chk, 5)
	assert.Contains(t, string(slc), "helloworld")

	slc = []byte("hello")
	chk = []byte("world")
	slc = PushChunkToSlice(slc, chk, 5)
	fmt.Println(slc)
	assert.Contains(t, string(slc), "helloworld")

}

func TestRequestFromReader(t *testing.T) {
	reader1 := &chunkReader{
		data:            "GET / HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
		numBytesPerRead: 3,
	}

	r, err := RequestFromReader(reader1)

	// Test: Good GET Request line with path
	require.NoError(t, err, "An error occured")
	require.NotNil(t, r)
	assert.Equal(t, "GET", r.RequestLine.Method)
	assert.Equal(t, "/", r.RequestLine.RequestTarget)
	assert.Equal(t, "1.1", r.RequestLine.HttpVersion)

	// Test: Good GET Request line with path
	reader2 := &chunkReader{
		data:            "GET /coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
		numBytesPerRead: 1,
	}
	r1, err1 := RequestFromReader(reader2)
	require.NoError(t, err1)
	require.NotNil(t, r1)
	assert.Equal(t, "GET", r1.RequestLine.Method)
	assert.Equal(t, "/coffee", r1.RequestLine.RequestTarget)
	assert.Equal(t, "1.1", r1.RequestLine.HttpVersion)

	// Test: Invalid method
	reader3 := &chunkReader{
		data:            "/cofee POST HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
		numBytesPerRead: 4,
	}
	_, err2 := RequestFromReader(reader3)
	require.Error(t, err2)
	require.ErrorIs(t, err2, InvalidMethodError)

	// Test Invalid target
	reader4 := &chunkReader{
		data:            "POST HTTP/1.1 /coffee\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
		numBytesPerRead: 2,
	}
	_, err3 := RequestFromReader(reader4)
	require.Error(t, err3)
	require.ErrorIs(t, err3, InvalidTargetError)

	// Test: Invalid number of parts in request line
	_, err = RequestFromReader(strings.NewReader("/coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"))
	require.Error(t, err)
}
