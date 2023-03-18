package models

import "time"

type Domain struct {
	ID        uint      `json:"id" gorm:"primaryKey;unique;autoIncrement"`
	Name      string    `json:"name" gorm:"unique"`
	Type      string    `json:"type"`
	Coverage  string    `json:"coverage"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

type DomainRepository interface {
	FindByID(ID int) (*Domain, error)
	FindByName(token string) (*Domain, error)
	Create(project *Domain) error
	Update(project *Domain) error
	Delete(project *Domain) error
}
