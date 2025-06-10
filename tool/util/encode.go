package util

import "encoding/base64"

func Base64Encode(input string) string {
	return base64.StdEncoding.EncodeToString([]byte(input))
}

func Base64Decode(input string) string {
	res, _ := base64.StdEncoding.DecodeString(input)
	return string(res)
}
