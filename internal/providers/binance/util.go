package binance

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"net/url"
)

func computeHmac256(strMessage string, strSecret string) string {
	key := []byte(strSecret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(strMessage))

	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func createSign(mapParams url.Values, secretKey string) string {
	strParams := mapParams.Encode()
	return url.QueryEscape(computeHmac256(strParams, secretKey))
}
