package validators

import (
	"fmt"
	CustomError "github.com/aerosystems/checkmail-service/pkg/custom_error"
	"os"
	"regexp"
	"strings"
)

func ValidateDomainTypes(tpe string) *CustomError.Error {
	trustTypes := strings.Split(os.Getenv("TRUST_DOMAIN_TYPES"), ",")
	if !Contains(trustTypes, tpe) {
		return CustomError.New(400003, fmt.Sprintf("domain type %s does not exist in trusted types", tpe))
	}
	return nil
}

func ValidateDomainCoverage(coverage string) *CustomError.Error {
	trustCoverages := strings.Split(os.Getenv("TRUST_DOMAIN_COVERAGES"), ",")
	if !Contains(trustCoverages, coverage) {
		return CustomError.New(400004, fmt.Sprintf("domain coverage %s does not exist in trusted coverages", coverage))
	}
	return nil
}

func ValidateDomainName(domainName string) *CustomError.Error {
	domainRegex := regexp.MustCompile(`^(?:[_a-z0-9](?:[_a-z0-9-]{0,61}[a-z0-9])?\.)+(?:[a-z](?:[a-z0-9-]{0,61}[a-z0-9])?)?$`)
	if !domainRegex.MatchString(domainName) {
		return CustomError.New(400005, fmt.Sprintf("domain name %s does not valid", domainName))
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
