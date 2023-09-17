package models

import (
	"time"
)

type Filter struct {
	ID           uint      `json:"Id" gorm:"primaryKey;unique;autoIncrement"`
	Name         string    `json:"name" gorm:"uniqueIndex:idx_name_coverage_project_token" example:"gmail.com"`
	Type         string    `json:"type" example:"whitelist"`
	Coverage     string    `json:"coverage" gorm:"uniqueIndex:idx_name_coverage_project_token" example:"equals"`
	ProjectToken string    `json:"-" gorm:"uniqueIndex:idx_name_coverage_project_token" example:"38fa45ebb919g5d966122bf9g42a38ceb1e4f6eddf1da70ef00afbdc38197d8f"`
	CreatedAt    time.Time `json:"CreatedAt" example:"2024-01-01T00:00:00Z"`
	UpdatedAt    time.Time `json:"UpdatedAt" example:"2024-01-01T00:00:00Z"`
}

type FilterRepository interface {
	FindAll() (*[]Filter, error)
	FindByID(id int) (*Filter, error)
	FindByProjectToken(projectToken string) (*Filter, error)
	Create(domain *Filter) error
	Update(domain *Filter) error
	Delete(domain *Filter) error
	MatchEquals(domainName, projectToken string) (*Filter, error)
	MatchContains(domainName, projectToken string) (*Filter, error)
	MatchBegins(domainName, projectToken string) (*Filter, error)
	MatchEnds(domainName, projectToken string) (*Filter, error)
}
