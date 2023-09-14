package models

import (
	"strings"
	"time"
)

type Domain struct {
	ID        uint      `json:"-" gorm:"primaryKey;unique;autoIncrement"`
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

func (d *Domain) Match(domainName string) bool {
	switch d.Coverage {
	case "equals":
		return strings.Contains(domainName, d.Name)
	case "begins":
		return strings.HasPrefix(domainName, d.Name)
	case "ends":
		return strings.HasSuffix(domainName, d.Name)
	default:
		return false
	}
}
