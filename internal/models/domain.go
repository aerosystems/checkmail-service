package models

import "time"

type Domain struct {
	ID        uint      `json:"-" gorm:"primaryKey;unique;autoIncrement"`
	Name      string    `json:"name" gorm:"unique" example:"gmail.com"`
	Type      string    `json:"type" example:"whitelist"`
	Coverage  string    `json:"coverage" example:"full"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

type DomainRepository interface {
	FindByID(id int) (*Domain, error)
	FindByName(name string) (*Domain, error)
	Create(domain *Domain) error
	Update(domain *Domain) error
	Delete(domain *Domain) error
}
