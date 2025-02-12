package entities

import (
	"time"
)

type Access struct {
	Token            string
	SubscriptionType SubscriptionType
	AccessCount      int
	AccessTime       time.Time
}

type SubscriptionType struct {
	slug string
}

var (
	UnknownSubscriptionType  = SubscriptionType{"unknown"}
	TrialSubscriptionType    = SubscriptionType{"trial"}
	StartupSubscriptionType  = SubscriptionType{"startup"}
	BusinessSubscriptionType = SubscriptionType{"business"}
)

func (k SubscriptionType) String() string {
	return k.slug
}

func SubscriptionTypeFromString(kind string) SubscriptionType {
	switch kind {
	case TrialSubscriptionType.String():
		return TrialSubscriptionType
	case StartupSubscriptionType.String():
		return StartupSubscriptionType
	case BusinessSubscriptionType.String():
		return BusinessSubscriptionType
	default:
		return UnknownSubscriptionType
	}
}
