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
