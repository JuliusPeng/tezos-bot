package publish

import (
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/ecadlabs/tezos-bot/models"
)

// TwitterConfig interface with method necessary to obtain twitter publisher configurable parameter
type TwitterConfig interface {
	GetTwitterConsummerID() string
	GetTwitterConsummerKey() string
	GetTwitterAccessTokenSecret() string
	GetTwitterAccessToken() string
}

// TwitterPublisher publisher that post new ballot on twitter
type TwitterPublisher struct {
	client *twitter.Client
}

// NewTwitterPublisher create a new TwitterPublisher
func NewTwitterPublisher(config TwitterConfig) (*TwitterPublisher, error) {
	c := oauth1.NewConfig(config.GetTwitterConsummerID(), config.GetTwitterConsummerKey())
	token := oauth1.NewToken(config.GetTwitterAccessToken(), config.GetTwitterAccessTokenSecret())

	httpClient := c.Client(oauth1.NoContext, token)

	client := twitter.NewClient(httpClient)

	// Just making a simple call to verify that the token is ok
	_, _, err := client.Timelines.HomeTimeline(&twitter.HomeTimelineParams{
		Count: 1,
	})

	if err != nil {
		return nil, err
	}

	return &TwitterPublisher{
		client: client,
	}, nil
}

// Publish a new ballot as a tweet
func (t *TwitterPublisher) Publish(ballot *models.Ballot) error {
	status := GetStatusString(ballot)
	_, _, err := t.client.Statuses.Update(status, nil)
	return err
}

// PublishProtoChange a new protocol change message as a tweet
func (t *TwitterPublisher) PublishProtoChange(proto string) error {
	status := GetProtocolString(proto)
	_, _, err := t.client.Statuses.Update(status, nil)
	return err
}
