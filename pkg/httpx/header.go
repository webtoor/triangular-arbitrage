package httpx

import (
	"net/textproto"
)

// Normalize formats the input header to the formation of "Xxx-Xxx".
func Normalize(header string) string {
	return textproto.CanonicalMIMEHeaderKey(header)
}
