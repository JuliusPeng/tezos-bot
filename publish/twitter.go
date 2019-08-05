package publish

import (
	"log"

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

	log.Println("Verifying twitter credentials...")
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
	log.Printf("(twitter) Published status: %s\n", status)
	return err
}

// PublishProtoChange a new protocol change message as a tweet
func (t *TwitterPublisher) PublishProtoChange(proto string) error {
	status := GetProtocolString(proto)
	_, _, err := t.client.Statuses.Update(status, nil)
	log.Printf("(twitter) Published status: %s\n", status)
	return err
}

// PublishProposalInjection a new proposal injection message as a tweet
func (t *TwitterPublisher) PublishProposalInjection(proposal *models.Proposal) error {
	status := GetProposalInjectString(proposal)
	_, _, err := t.client.Statuses.Update(status, nil)
	log.Printf("(twitter) Published status: %s\n", status)
	return err
}

// PublishProposalSummary a new proposal summary message as a tweet
func (t *TwitterPublisher) PublishProposalSummary(proposal *models.ProposalSummary) error {
	status := GetProposalSummaryString(proposal)
	_, _, err := t.client.Statuses.Update(status, nil)
	log.Printf("(twitter) Published status: %s\n", status)
	return err
}

// PublishWinningProposalSummary a new winning proposal summary message as a tweet
func (t *TwitterPublisher) PublishWinningProposalSummary(proposal *models.ProposalSummary) error {
	status := GetWinningProposalString(proposal)
	_, _, err := t.client.Statuses.Update(status, nil)
	log.Printf("(twitter) Published status: %s\n", status)
	return err
}

// PublishProposalUpvote a new proposal upvote message to twitter
func (t *TwitterPublisher) PublishProposalUpvote(proposal *models.Proposal) error {
	status := GetProposalUpvoteString(proposal)
	_, _, err := t.client.Statuses.Update(status, nil)
	if err != nil {
		log.Printf("(twitter) Published status: %s\n", status)
	}
	return err
}
