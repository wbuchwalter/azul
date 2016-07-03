package auth

import "encoding/base64"

func GetAuthInfo(username string, password string) []byte {
	var token []byte
	base64.StdEncoding.Encode(token, []byte(username+":"+password))
	return token
}
