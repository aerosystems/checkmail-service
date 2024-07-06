package models

import (
	"time"
)

type Domain struct {
	Id        int
	Name      string
	Type      DomainType
	Coverage  DomainCoverage
	CreatedAt time.Time
	UpdatedAt time.Time
}

type DomainType struct {
	slug string
}

func (d DomainType) String() string {
	return d.slug
}

func DomainTypeFromString(s string) DomainType {
	switch s {
	case "whitelist":
		return WhitelistType
	case "blacklist":
		return BlacklistType
	default:
		return UndefinedType
	}
}

var (
	UndefinedType = DomainType{"undefined"}
	WhitelistType = DomainType{"whitelist"}
	BlacklistType = DomainType{"blacklist"}
)

type DomainCoverage struct {
	slug string
}

func (d DomainCoverage) String() string {
	return d.slug
}

var (
	EqualsCoverage   = DomainCoverage{"equals"}
	ContainsCoverage = DomainCoverage{"contains"}
	LeftCoverage     = DomainCoverage{"left"}
	RightCoverage    = DomainCoverage{"right"}
)
