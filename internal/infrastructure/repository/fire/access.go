package fire

import (
	"cloud.google.com/go/firestore"
	"context"
	"errors"
	CustomErrors "github.com/aerosystems/checkmail-service/internal/common/custom_errors"
	"github.com/aerosystems/checkmail-service/internal/models"
	"google.golang.org/api/iterator"
	"time"
)

const (
	collectionName = "access"
)

type ApiAccessRepo struct {
	client *firestore.Client
}

func NewApiAccessRepo(client *firestore.Client) *ApiAccessRepo {
	return &ApiAccessRepo{
		client: client,
	}
}

type ApiAccessFire struct {
	Token            string    `firestore:"token"`
	SubscriptionType string    `firestore:"subscription_type"`
	AccessTime       time.Time `firestore:"access_time"`
}

func (a *ApiAccessFire) ToModel() (models.Access, error) {
	subscriptionType := models.SubscriptionTypeFromString(a.SubscriptionType)
	if subscriptionType == models.UnknownSubscriptionType {
		return models.Access{}, errors.New("unknown subscription type")
	}
	return models.Access{
		Token:            a.Token,
		SubscriptionType: models.SubscriptionTypeFromString(a.SubscriptionType),
		AccessTime:       a.AccessTime,
	}, nil
}

func (r *ApiAccessRepo) Get(ctx context.Context, token string) (*models.Access, error) {
	var accessFire ApiAccessFire
	iter := r.client.Collection(collectionName).Where("token", "==", token).Documents(ctx)
	for {
		doc, err := iter.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			return nil, err
		}
		if err := doc.DataTo(&accessFire); err == nil {
			break
		}
	}
	if accessFire == (ApiAccessFire{}) {
		return nil, CustomErrors.ErrApiKeyNotFound
	}
	access, err := accessFire.ToModel()
	if err != nil {
		return nil, err
	}
	return &access, nil
}
