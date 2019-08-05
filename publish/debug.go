package publish

import (
	"fmt"

	"github.com/ecadlabs/tezos-bot/models"
)

// DebugPublisher is a simple publish that logs ballot directly to stdout
type DebugPublisher struct{}

// Publish logs ballot directly to stdout
func (d *DebugPublisher) Publish(ballot *models.Ballot) error {
	status := GetStatusString(ballot)
	fmt.Printf("(%d) %s\n", len(status), status)
	return nil
}

// PublishProtoChange a new protocol change message to stdout
func (d *DebugPublisher) PublishProtoChange(proto string) error {
	status := GetProtocolString(proto)
	fmt.Printf("(%d) %s\n", len(status), status)
	return nil
}

// PublishProposalInjection a new proposal injection message to stdout
func (d *DebugPublisher) PublishProposalInjection(proposal *models.Proposal) error {
	status := GetProposalInjectString(proposal)
	fmt.Printf("(%d) %s\n", len(status), status)
	return nil
}

// PublishProposalSummary a new proposal summary message to stdout
func (d *DebugPublisher) PublishProposalSummary(proposal *models.ProposalSummary) error {
	status := GetProposalSummaryString(proposal)
	fmt.Printf("(%d) %s\n", len(status), status)
	return nil
}

// PublishWinningProposalSummary a new winning proposal summary message to stdout
func (d *DebugPublisher) PublishWinningProposalSummary(proposal *models.ProposalSummary) error {
	status := GetWinningProposalString(proposal)
	fmt.Printf("(%d) %s\n", len(status), status)
	return nil
}

// PublishProposalUpvote a new proposal upvote message to stdout
func (d *DebugPublisher) PublishProposalUpvote(proposal *models.Proposal) error {
	status := GetProposalUpvoteString(proposal)
	fmt.Printf("(%d) %s\n", len(status), status)
	return nil
}
