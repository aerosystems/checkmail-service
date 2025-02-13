package entities

import "time"

type InspectEvent struct {
	ProjectToken string
	Data         string
	Domain       string
	DomainType   Type
	CreatedAt    time.Time
}
