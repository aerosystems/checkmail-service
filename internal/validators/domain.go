package validators

import (
	CustomError "github.com/aerosystems/checkmail-service/pkg/custom_error"
	"os"
	"regexp"
	"strings"
)

const (
	trustDomainTypes     = "blacklist,whitelist,undefined"
	trustDomainCoverages = "begins,ends,equals,contains"
)

func ValidateDomainTypes(tpe string) *CustomError.Error {
	trustTypes := strings.Split(trustDomainTypes, ",")
	if !Contains(trustTypes, tpe) {
		return CustomError.New(400003, "domain type does not exist in trusted types")
	}
	return nil
}

func ValidateDomainCoverage(coverage string) *CustomError.Error {
	trustCoverages := strings.Split(os.Getenv(trustDomainCoverages), ",")
	if !Contains(trustCoverages, coverage) {
		return CustomError.New(400004, "domain coverage does not exist in trusted coverages")
	}
	return nil
}

func ValidateDomainName(domainName string) *CustomError.Error {
	domainRegex := regexp.MustCompile(`^(?:[_a-z0-9](?:[_a-z0-9-]{0,61}[a-z0-9])?\.)+(?:[a-z](?:[a-z0-9-]{0,61}[a-z0-9])?)?$`)
	if !domainRegex.MatchString(domainName) {
		return CustomError.New(400002, "domain name does not valid")
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
