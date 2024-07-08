package models

import "time"

type Review struct {
	Id        int    `gorm:"primaryKey;unique;autoIncrement"`
	Name      string `gorm:"index:idx_name" example:"gmail.com"`
	Type      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
