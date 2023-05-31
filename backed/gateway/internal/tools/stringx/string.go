package stringx

import (
	"net/http"
	"strings"
)

func HasPrefixAndJoin(str, pre string) string {
	if !strings.HasPrefix(str, pre) {
		return str + pre
	}
	return str
}

func FixHttpMethod(method string) (string, bool) {
	m := strings.ToUpper(method)
	if m == http.MethodGet || m == http.MethodPost || m == http.MethodDelete || m == http.MethodPut ||
		m == http.MethodConnect || m == http.MethodHead || m == http.MethodOptions || m == http.MethodPatch ||
		m == http.MethodTrace {
		return m, true
	}
	return "", false
}
