package adapters

import (
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	"github.com/aerosystems/checkmail-service/internal/entities"
	"github.com/aerosystems/common-service/pkg/gcpclient"
)

const (
	inspectTopicName = "inspect"
	timestampLayout  = "2006-01-02 15:04:05"
)

type StatAdapter struct {
	pubsub gcpclient.PubSubClient
}

type InspectEvent struct {
	ProjectToken string `json:"project_token"`
	Data         string `json:"data"`
	Domain       string `json:"domain"`
	DomainType   string `json:"domain_type"`
	CreatedAt    string `json:"created_at"`
}

func modelToEvent(inspect entities.InspectEvent) *InspectEvent {
	return &InspectEvent{
		ProjectToken: inspect.ProjectToken,
		Data:         inspect.Data,
		Domain:       inspect.Domain,
		DomainType:   inspect.DomainType.String(),
		CreatedAt:    inspect.CreatedAt.Format(timestampLayout),
	}
}

func NewStatAdapter(client gcpclient.PubSubClient) *StatAdapter {
	return &StatAdapter{
		pubsub: client,
	}
}

func (s StatAdapter) PublishStat(ctx context.Context, inspectEvent entities.InspectEvent) error {
	topic := s.pubsub.Client.Topic(inspectTopicName)
	defer topic.Stop()

	event := modelToEvent(inspectEvent)
	eventData, err := json.Marshal(event)
	if err != nil {
		return err
	}

	res := topic.Publish(ctx, &pubsub.Message{
		Data: eventData,
	})

	_, err = res.Get(s.pubsub.Ctx)
	if err != nil {
		return err
	}

	return nil
}
