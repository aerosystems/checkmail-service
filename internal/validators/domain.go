package validators

import (
	"github.com/aerosystems/checkmail-service/internal/models"
	"os"
	"regexp"
	"strings"
)

const (
	trustDomainTypes     = "blacklist,whitelist,undefined"
	trustDomainCoverages = "begins,ends,equals,contains"
)

func ValidateDomainTypes(tpe string) *models.Error {
	trustTypes := strings.Split(trustDomainTypes, ",")
	if !Contains(trustTypes, tpe) {
		return &models.Error{
			Code:    400003,
			Message: "domain type does not exist in trusted types",
		}
	}
	return nil
}

func ValidateDomainCoverage(coverage string) *models.Error {
	trustCoverages := strings.Split(os.Getenv(trustDomainCoverages), ",")
	if !Contains(trustCoverages, coverage) {
		return &models.Error{
			Code:    400004,
			Message: "domain coverage does not exist in trusted coverages",
		}
	}
	return nil
}

func ValidateDomainName(domainName string) *models.Error {
	domainRegex := regexp.MustCompile(`^(?:[_a-z0-9](?:[_a-z0-9-]{0,61}[a-z0-9])?\.)+(?:[a-z](?:[a-z0-9-]{0,61}[a-z0-9])?)?$`)
	if !domainRegex.MatchString(domainName) {
		return &models.Error{
			Code:    400002,
			Message: "domain name does not valid",
		}
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
