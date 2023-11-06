package models

import "time"

type RootDomain struct {
	ID        int       `json:"id" gorm:"primaryKey;unique;autoIncrement"`
	Name      string    `json:"name" gorm:"uniqueIndex:idx_name"`
	Type      string    `json:"type"`
	TLD       string    `json:"tld"`
	Source    string    `json:"source"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type RootDomainRepository interface {
	FindByID(id int) (*RootDomain, error)
	FindByName(name string) (*RootDomain, error)
	Create(rootDomain *RootDomain) error
	Update(rootDomain *RootDomain) error
	Delete(rootDomain *RootDomain) error
}
