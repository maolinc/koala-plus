package utilx

import "strings"

func GetNewPerms(perms, prefix string) string {
	if strings.HasPrefix(perms, prefix) {
		return perms
	}
	return prefix + perms
}
