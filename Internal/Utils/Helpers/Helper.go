package Helper

import (
	"encoding/base64"
	"strings"
)

func Base64ToBytes(base64str string) ([]byte, error) {
	if strings.HasPrefix(base64str, "data:") {
		base64str = strings.Split(base64str, ",")[1]
	}
	return base64.StdEncoding.DecodeString(base64str)
}
