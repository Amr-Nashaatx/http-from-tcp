package headers

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/Amr-Nashaatx/http-from-tcp/utils"
)

type Headers map[string]string

func NewHeaders() Headers {
	return make(Headers)
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	clrf := []byte("\r\n")
	// does the data starts with clrf? if so we know headers sections is done
	if bytes.HasPrefix(data, clrf) {
		return 0, true, nil
	}
	// does the data contain clrf?, if not we still need more
	clrfPos := bytes.Index(data, clrf)
	if clrfPos == -1 {
		return 0, false, nil
	}

	// extract field name and field value
	parts := strings.SplitN(string(data[:clrfPos]), ":", 2)
	if len(parts) < 2 {
		return 0, false, fmt.Errorf("Invalid header format: missing colon separator")
	}
	fieldName := strings.ToLower(parts[0])
	fieldValue := strings.Trim(parts[1], " ")

	// field name should not contain any spaces
	if strings.Contains(fieldName, " ") {
		return 0, false, fmt.Errorf("Invalid field name")
	}

	// check for invalid characters in field name
	for token := range utils.GetTokensFromText(fieldName) {
		if !IsValidToken(token) {
			return 0, false, fmt.Errorf("Invalid token %s, in field name", token)
		}
	}

	// check if the field name exists on headers map, if so we append the field value to it.
	if h[fieldName] != "" {
		h[fieldName] = h[fieldName] + ", " + fieldValue
	} else {
		h[fieldName] = fieldValue
	}

	return clrfPos + 2, false, nil
}
