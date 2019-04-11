package publish

import (
	"fmt"

	"github.com/ecadlabs/tezos-bot/models"
)

// DebugPublisher is a simple publish that logs ballot directly to stdout
type DebugPublisher struct{}

// Publish logs ballot directly to stdout
func (d *DebugPublisher) Publish(ballot *models.Ballot) error {
	fmt.Printf("%s voted %s for proposal %s", ballot.PKH, ballot.Ballot, ballot.ProposalHash)
	return nil
}
