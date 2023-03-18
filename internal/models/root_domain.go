package models

import "time"

type RootDomain struct {
	ID        uint      `json:"id" gorm:"primaryKey;unique;autoIncrement"`
	Name      string    `json:"name" gorm:"unique"`
	Type      string    `json:"type"`
	TLD       string    `json:"tld"`
	Source    string    `json:"source"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
