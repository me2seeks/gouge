package settings

import (
	"regexp"
	"strings"
)

var UserAllowAll = regexp.MustCompile("")

func ParseAuth(auth string) (string, string) {
	if strings.Contains(auth, ":") {
		pair := strings.SplitN(auth, ":", 2)
		return pair[0], pair[1]
	}
	return "", ""
}
