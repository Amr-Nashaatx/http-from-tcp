package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHeaders(t *testing.T) {
	// Test: Valid single header
	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	assert.NoError(t, err)
	assert.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["Host"])
	assert.Equal(t, 23, n)
	assert.False(t, done)

	// Test: Invalid spacing header
	headers = NewHeaders()
	data = []byte("       Host: localhost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	assert.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Valid single header with extra whitespace
	headers = NewHeaders()
	data = []byte("Host: localhost:42069        \r\n\r\n")
	n, done, err = headers.Parse(data)
	assert.NoError(t, err)
	assert.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["Host"])
	assert.Equal(t, len(data)-len([]byte("\r\n")), n)
	assert.False(t, done)

	// Test: Valid single header with invalid character
	headers = NewHeaders()
	data = []byte("H©st: localhost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "©")

	// Valid 2 headers with existing headers
	headers = NewHeaders()
	data = []byte("Host: localhost:42069        \r\n\r\n")
	n, done, err = headers.Parse(data)
	assert.NoError(t, err)
	assert.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["Host"])
	assert.Equal(t, len(data)-len([]byte("\r\n")), n)
	assert.False(t, done)

	n, done, err = headers.Parse([]byte("Content-Type: application/json \r\n\r\n"))
	assert.Equal(t, "application/json", headers["Content-type"])
	assert.Equal(t, "localhost:42069", headers["Host"])
}
