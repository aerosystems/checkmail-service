package models

import (
	"regexp"
	"time"
)

type Domain struct {
	Name      string
	Type      Type
	Match     Match
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Type struct {
	slug string
}

var (
	UndefinedType = Type{"undefined"}
	BlacklistType = Type{"blacklist"}
	WhitelistType = Type{"whitelist"}
)

func (d Type) String() string {
	return d.slug
}

func DomainTypeFromString(s string) Type {
	switch s {
	case BlacklistType.String():
		return BlacklistType
	case WhitelistType.String():
		return WhitelistType
	default:
		return UndefinedType
	}
}

type Match struct {
	slug string
}

var (
	UndefinedMatch = Match{"undefined"}
	PrefixMatch    = Match{"prefix"}
	SuffixMatch    = Match{"suffix"}
	EqualsMatch    = Match{"equals"}
	ContainsMatch  = Match{"contains"}
)

func (d Match) String() string {
	return d.slug
}

func DomainMatchFromString(s string) Match {
	switch s {
	case PrefixMatch.String():
		return PrefixMatch
	case SuffixMatch.String():
		return SuffixMatch
	case EqualsMatch.String():
		return EqualsMatch
	case ContainsMatch.String():
		return ContainsMatch
	default:
		return UndefinedMatch
	}
}

func ValidateDomainName(domainName string) error {
	domainRegex := regexp.MustCompile(`^(?:[_a-z0-9](?:[_a-z0-9-]{0,61}[a-z0-9])?\.)+(?:[a-z](?:[a-z0-9-]{0,61}[a-z0-9])?)?$`)
	if !domainRegex.MatchString(domainName) {
		return ErrDomainNotValid
	}
	return nil
}
