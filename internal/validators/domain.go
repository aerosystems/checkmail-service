package validators

import (
	CustomErrors "github.com/aerosystems/checkmail-service/internal/common/custom_errors"
	"regexp"
	"strings"
)

const (
	trustDomainTypes     = "blacklist,whitelist,undefined"
	trustDomainCoverages = "begins,ends,equals,contains"
)

func ValidateDomainTypes(tpe string) error {
	trustTypes := strings.Split(trustDomainTypes, ",")
	if !Contains(trustTypes, tpe) {
		return CustomErrors.ErrDomainTrustedTypes
	}
	return nil
}

func ValidateDomainCoverage(coverage string) error {
	trustCoverages := strings.Split(trustDomainCoverages, ",")
	if !Contains(trustCoverages, coverage) {
		return CustomErrors.ErrDomainCoverage
	}
	return nil
}

func ValidateDomainName(domainName string) error {
	domainRegex := regexp.MustCompile(`^(?:[_a-z0-9](?:[_a-z0-9-]{0,61}[a-z0-9])?\.)+(?:[a-z](?:[a-z0-9-]{0,61}[a-z0-9])?)?$`)
	if !domainRegex.MatchString(domainName) {
		return CustomErrors.ErrDomainNotValid
	}
	return nil
}

func Contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}
