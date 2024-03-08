package utils

import "regexp"

func CheckModuleSlug(slug string) bool {
	r, _ := regexp.Compile(`^[a-zA-Z0-9]+[-\/][a-z][a-z0-9_]*$`)
	return r.MatchString(slug)
}

func CheckReleaseSlug(slug string) bool {
	r, _ := regexp.Compile(`^[a-zA-Z0-9]+[-\/][a-z][a-z0-9_]*[-\/][0-9]+\.[0-9]+\.[0-9]+(?:[\-+].+)?$`)
	return r.MatchString(slug)
}
