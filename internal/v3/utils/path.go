package utils

import "regexp"

// CheckModuleSlug validates if a module slug follows the required format:
// - Must start with alphanumeric characters
// - Must contain a hyphen or slash separator
// - After separator, must start with a lowercase letter
// - Can contain lowercase letters, numbers, and underscores after first letter
// Example valid slugs: "myOrg/module", "company-mymodule"
func CheckModuleSlug(slug string) bool {
	r, _ := regexp.Compile(`^[a-zA-Z0-9]+[-\/][a-z][a-z0-9_]*$`)
	return r.MatchString(slug)
}

// CheckReleaseSlug validates if a release slug follows the required format:
// - Must start with a valid module slug (see above)
// - Must contain another hyphen or slash separator
// - Must end with a semantic version (X.Y.Z)
// - May optionally include a pre-release or build metadata after version
// Example valid slugs: "myOrg/module/1.2.3", "company-mymodule/2.0.0-beta.1"
func CheckReleaseSlug(slug string) bool {
	r, _ := regexp.Compile(`^[a-zA-Z0-9]+[-\/][a-z][a-z0-9_]*[-\/][0-9]+\.[0-9]+\.[0-9]+(?:[\-+].+)?$`)
	return r.MatchString(slug)
}
