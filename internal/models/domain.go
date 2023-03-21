package models

import (
	"strings"
	"time"
)

type Domain struct {
	ID        uint      `json:"-" gorm:"primaryKey;unique;autoIncrement"`
	Name      string    `json:"name" gorm:"unique" example:"gmail.com"`
	Type      string    `json:"type" example:"whitelist"`
	Coverage  string    `json:"coverage" example:"equals"`
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
}

func (d *Domain) GetType(domainName string) *Domain {
	switch d.Coverage {
	case "equals":
		if domainName == d.Name {
			return d
		}
	case "contains":
		if strings.Contains(domainName, d.Name) {
			return d
		}
	case "begins":
		if strings.HasPrefix(domainName, d.Name) {
			return d
		}
	case "ends":
		if strings.HasSuffix(domainName, d.Name) {
			return d
		}
	}
	return nil
}
