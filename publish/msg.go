package publish

import (
	"fmt"
	"log"
	"net"

	"github.com/ecadlabs/tezos-bot/models"
)

const (
	proposalZone = "proposal.tezz.ie"
	addressZone  = "tz.tezz.ie"
)

// GetStatusString composes a status string based on available vanity data
func GetStatusString(ballot *models.Ballot) string {
	templateBasic := `Tezos address %s voted "%s" %son #Tezos proposal "%s"%s`
	templateVanity := `Tezos baker "%s" /%s voted "%s" %son #Tezos proposal "%s"%s`

	var proposalVanityName string
	protocolName, err := LookupTZName(ballot.ProposalHash, proposalZone)
	if err != nil {
		proposalVanityName = ballot.ProposalHash
	} else {
		proposalVanityName = protocolName
	}

	templateRolls := ""
	if ballot.Rolls != 0 {
		templateRolls = fmt.Sprintf("with %d rolls ", ballot.Rolls)
	}

	templateQuorum := "and quorum has been reached"
	percentTowardQuorum := ballot.PercentTowardQuorum()
	if percentTowardQuorum > 0 {
		templateQuorum = fmt.Sprintf("with %.2f%% remaining to reach %.2f%% quorum", percentTowardQuorum, ballot.Quorum)
	}

	templatePhase := "for the promotion phase."

	if ballot.IsTesting {
		templatePhase = "for the exploration phase."
	}

	templateStatus := fmt.Sprintf("\n\nVote status is %.2f%% yay/%.2f%% nay, %s %s", ballot.CountingPercentYay(), ballot.CountingPercentNay(), templateQuorum, templatePhase)

	// tz.tezz.ie is an experimental DNS zone to resolve vanity names from tz
	// addresses
	address, err := LookupTZName(ballot.PKH, addressZone)

	if err != nil {
		log.Printf("No address found for %s, err: %s", ballot.PKH, err)
		return fmt.Sprintf(templateBasic, ballot.PKH, ballot.Ballot, templateRolls, proposalVanityName, templateStatus)
	}
	log.Printf("Address %s found for %s, ", address, ballot.PKH)
	return fmt.Sprintf(templateVanity, address, ballot.PKH, ballot.Ballot, templateRolls, proposalVanityName, templateStatus)

}

// GetProposalSummaryString get status message for daily proposal summary
func GetProposalSummaryString(summary *models.ProposalSummary) string {
	proposalName, err := LookupTZName(summary.ProposalHash, proposalZone)

	if err != nil {
		log.Printf("No protocol found for %s, err: %s", summary.ProposalHash, err)
		proposalName = summary.ProposalHash
	} else {
		proposalName = fmt.Sprintf("%s (%s)", proposalName, summary.ProposalHash)
	}

	return fmt.Sprintf("Proposal upvotes: #Tezos proposal %s received %d upvotes in cycle %d, and now has %d votes.", proposalName, summary.NewSupporters, summary.Cycle, summary.SupporterCount)
}

// GetWinningProposalString get status message for proposal that moved to exploration phase
func GetWinningProposalString(summary *models.ProposalSummary) string {
	proposalName, err := LookupTZName(summary.ProposalHash, proposalZone)

	if err != nil {
		log.Printf("No protocol found for %s, err: %s", summary.ProposalHash, err)
		proposalName = summary.ProposalHash
	} else {
		proposalName = fmt.Sprintf("%s (%s)", proposalName, summary.ProposalHash)
	}

	return fmt.Sprintf("Proposal period complete: proposal %s received the most upvotes (%d) and is advancing to the exploration vote period.", proposalName, summary.SupporterCount)
}

// GetProposalInjectString retrieve the template for proposal injection status
func GetProposalInjectString(proposal *models.Proposal) string {
	address, err := LookupTZName(proposal.PKH, proposalZone)

	if err != nil {
		log.Printf("No address found for %s, err: %s", proposal.PKH, err)
	}

	templateAddress := proposal.PKH

	if address != "" {
		templateAddress = fmt.Sprintf("%s /%s", address, proposal.PKH)
	}

	return fmt.Sprintf("New #Tezos proposal injected! %s injected proposal %s in voting period %d.", templateAddress, proposal.ProposalHash, proposal.Period)
}

// GetProposalUpvoteString retrieve the template for proposal upvote status
func GetProposalUpvoteString(proposal *models.Proposal) string {
	address, err := LookupTZName(proposal.PKH, proposalZone)

	if err != nil {
		log.Printf("No address found for %s, err: %s", proposal.PKH, err)
	}

	templateRolls := ""
	if proposal.Rolls != 0 {
		templateRolls = fmt.Sprintf("with %d rolls ", proposal.Rolls)
	}

	templateAddress := proposal.PKH

	if address != "" {
		templateAddress = fmt.Sprintf("%s /%s", address, proposal.PKH)
	}

	return fmt.Sprintf("Address %s upvoted proposal %s %sin voting period %d.", templateAddress, proposal.ProposalHash, templateRolls, proposal.Period)
}

// GetProtocolString retrieve the template for protocol change status
func GetProtocolString(proto string) string {
	lookupKey := fmt.Sprintf("%s", proto)

	protocolName, err := LookupTZName(lookupKey, proposalZone)

	if err != nil {
		log.Printf("No protocol found for %s, err: %s", lookupKey, err)
		protocolName = proto
	}

	return fmt.Sprintf("Protocol %s is now live on mainnet! #Tezos", protocolName)
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
