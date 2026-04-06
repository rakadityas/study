package signing

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

func SignV1(secret string, timestamp string, body []byte) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(timestamp))
	mac.Write([]byte("."))
	mac.Write(body)
	return hex.EncodeToString(mac.Sum(nil))
}

func VerifyV1(secret string, timestamp string, body []byte, hexSig string) bool {
	expected, err := hex.DecodeString(hexSig)
	if err != nil {
		return false
	}
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(timestamp))
	mac.Write([]byte("."))
	mac.Write(body)
	return hmac.Equal(mac.Sum(nil), expected)
}

