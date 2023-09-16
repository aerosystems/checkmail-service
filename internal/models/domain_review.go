package models

import "time"

type DomainReview struct {
	ID        uint      `json:"-" gorm:"primaryKey;unique;autoIncrement"`
	Name      string    `json:"name" gorm:"index:idx_name" example:"gmail.com"`
	Type      string    `json:"type" example:"whitelist"`
	CreatedAt time.Time `json:"createdAt" example:"2021-01-01T00:00:00Z"`
	UpdatedAt time.Time `json:"updatedAt" example:"2021-01-01T00:00:00Z"`
}

type DomainReviewRepository interface {
	FindByName(name string) (*DomainReview, error)
	Create(domain *DomainReview) error
	Delete(domain *DomainReview) error
}
