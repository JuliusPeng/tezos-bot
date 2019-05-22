package publish

import (
	"fmt"
	"log"
	"net"

	"github.com/ecadlabs/tezos-bot/models"
)

// GetStatusString composes a status string based on available vanity data
func GetStatusString(ballot *models.Ballot) string {
	templateBasic := `Tezos address %s voted "%s" %son #Tezos proposal "%s"%s`
	templateVanity := `Tezos baker "%s" /%s voted "%s" %son #Tezos proposal "%s"%s`
	// TODO(jev) update to query Proposal vanity name for DNS
	proposalVanityName := "Athens A"

	templateRolls := ""
	if ballot.Rolls != 0 {
		templateRolls = fmt.Sprintf("with %d rolls ", ballot.Rolls)
	}

	templateQuorum := "and quorum has been reached."
	percentTowardQuorum := ballot.PercentTowardQuorum()
	if percentTowardQuorum > 0 {
		templateQuorum = fmt.Sprintf("with %.2f%% remaining to reach %.2f%% quorum.", percentTowardQuorum, ballot.Quorum)
	}
	templateStatus := fmt.Sprintf("\n\nVote status is %.2f%%/yay %.2f%%/nay, %s", ballot.CountingPercentYay(), ballot.CountingPercentNay(), templateQuorum)

	// tz.tezz.ie is an experimental DNS zone to resolve vanity names from tz
	// addresses
	address, err := LookupTZName(ballot.PKH, "tz.tezz.ie")

	if err != nil {
		log.Printf("No address found for %s, err: %s", ballot.PKH, err)
		return fmt.Sprintf(templateBasic, ballot.PKH, ballot.Ballot, templateRolls, proposalVanityName, templateStatus)
	}
	log.Printf("Address %s found for %s, ", address, ballot.PKH)
	return fmt.Sprintf(templateVanity, address, ballot.PKH, ballot.Ballot, templateRolls, proposalVanityName, templateStatus)

}

// LookupTZName queries DNS for a txt record corresponding to a TZ address.
func LookupTZName(address, zone string) (string, error) {
	query := fmt.Sprintf("%s.%s", address, zone)
	rrs, err := net.LookupTXT(query)
	if err != nil {
		return "", err
	}
	return string(rrs[0]), nil
}
