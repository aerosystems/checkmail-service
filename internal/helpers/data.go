package helpers

import (
	"regexp"
)

func ValidateDomain(domainName string) bool {
	domainRegex := regexp.MustCompile(`^(?:[_a-z0-9](?:[_a-z0-9-]{0,61}[a-z0-9])?\.)+(?:[a-z](?:[a-z0-9-]{0,61}[a-z0-9])?)?$`)
	return domainRegex.MatchString(domainName)
}
