package models

import (
	"time"
)

type Domain struct {
	Id        int    `gorm:"primaryKey;unique;autoIncrement"`
	Name      string `gorm:"uniqueIndex:idx_name_coverage" example:"gmail.com"`
	Type      string `example:"whitelist"`
	Coverage  string `gorm:"uniqueIndex:idx_name_coverage" example:"equals"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type KindDomain string

const (
	WhitelistDomain KindDomain = "whitelist"
	BlacklistDomain KindDomain = "blacklist"
	UndefinedDomain KindDomain = "undefined"
)

type KindCoverage string

const (
	EqualsCoverage   KindCoverage = "equals"
	ContainsCoverage KindCoverage = "contains"
	LeftCoverage     KindCoverage = "left"
	RightCoverage    KindCoverage = "right"
)
