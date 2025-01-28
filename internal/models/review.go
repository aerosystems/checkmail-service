package models

import "time"

type Review struct {
	Id        int
	Name      string
	Type      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
