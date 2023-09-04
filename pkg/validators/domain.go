package validators

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
)

func ValidateDomainTypes(tpe string) error {
	trustTypes := strings.Split(os.Getenv("TRUST_DOMAIN_TYPES"), ",")
	if !Contains(trustTypes, tpe) {
		return fmt.Errorf("domain Type %s exists in trusted types", tpe)
	}
	return nil
}

func ValidateDomainCoverages(coverage string) error {
	trustCoverages := strings.Split(os.Getenv("TRUST_DOMAIN_COVERAGES"), ",")
	if !Contains(trustCoverages, coverage) {
		return fmt.Errorf("domain Coverage %s exists in trusted coverages", coverage)
	}
	return nil
}

func ValidateDomain(domainName string) error {
	domainRegex := regexp.MustCompile(`^(?:[_a-z0-9](?:[_a-z0-9-]{0,61}[a-z0-9])?\.)+(?:[a-z](?:[a-z0-9-]{0,61}[a-z0-9])?)?$`)
	if !domainRegex.MatchString(domainName) {
		return errors.New("domain name is not valid")
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
