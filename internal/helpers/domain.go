package helpers

import (
	"fmt"
	"os"
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

func Contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}
