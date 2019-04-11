package publish

import (
	"fmt"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/ecadlabs/tezos-bot/models"
	"golang.org/x/oauth2"
)

// TwitterConfig interface with method necessary to obtain twitter publisher configurable parameter
type TwitterConfig interface {
	GetTwitterAccessToken() string
}

// TwitterPublisher publisher that post new ballot on twitter
type TwitterPublisher struct {
	client *twitter.Client
}

// NewTwitterPublisher create a new TwitterPublisher
func NewTwitterPublisher(config TwitterConfig) (*TwitterPublisher, error) {
	twitterConf := &oauth2.Config{}
	token := &oauth2.Token{AccessToken: config.GetTwitterAccessToken()}
	httpClient := twitterConf.Client(oauth2.NoContext, token)

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
	status := fmt.Sprintf("%s voted %s for proposal %s", ballot.PKH, ballot.Ballot, ballot.ProposalHash)

	_, _, err := t.client.Statuses.Update(status, nil)
	return err
}
