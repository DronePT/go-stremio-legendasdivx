package utils

import (
	"encoding/base64"
	"strings"
)

func ParseUserData(config string) (string, string) {
	decodedCredentials, _ := base64.RawStdEncoding.DecodeString(config)
	credentials := strings.Split(string(decodedCredentials), ":")

	return credentials[0], credentials[1]
}
