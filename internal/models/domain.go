package models

import (
	"time"
)

type Domain struct {
	ID        int       `json:"-" gorm:"primaryKey;unique;autoIncrement"`
	Name      string    `json:"name" gorm:"uniqueIndex:idx_name_coverage" example:"gmail.com"`
	Type      string    `json:"type" example:"whitelist"`
	Coverage  string    `json:"coverage" gorm:"uniqueIndex:idx_name_coverage" example:"equals"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

type DomainRepository interface {
	FindAll() (*[]Domain, error)
	FindByID(id int) (*Domain, error)
	FindByName(name string) (*Domain, error)
	Create(domain *Domain) error
	Update(domain *Domain) error
	Delete(domain *Domain) error
	MatchEquals(name string) (*Domain, error)
	MatchContains(name string) (*Domain, error)
	MatchBegins(name string) (*Domain, error)
	MatchEnds(name string) (*Domain, error)
	Count() (map[string]int, error)
}
