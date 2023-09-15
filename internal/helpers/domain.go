package helpers

import (
	"errors"
	"strings"
)

func GetRootDomain(domain string) (string, error) {
	arrDomain := strings.Split(domain, ".")
	if len(arrDomain) < 2 {
		return "", errors.New("domain is not valid")
	}
	return arrDomain[len(arrDomain)-1], nil
}
