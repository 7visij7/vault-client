package base64

import (
        "encoding/base64"
)

func Base64Encode(str string) string {
    return base64.StdEncoding.EncodeToString([]byte(str))
}


func base64Decode(str string) (string, bool) {
    data, err := base64.StdEncoding.DecodeString(str)
    if err != nil {
        return "", false
    }
    return string(data), true
}

func CheckBase64(value string) string {
	decodeValue, err := base64Decode(value)
	if err {
		return decodeValue
	} 
	return value
}