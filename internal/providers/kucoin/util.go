package kucoin

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
)

func createSign(msg, timestamp, secretKey string) string {
	hm := hmac.New(sha256.New, []byte(secretKey))
	if len(timestamp) > 0 {
		hm.Write([]byte(timestamp + msg))
	} else {
		hm.Write([]byte(msg))
	}
	data := hm.Sum(nil)
	return base64.StdEncoding.EncodeToString(data)
}
