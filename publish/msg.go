package publish

import (
	"bytes"
	"fmt"
	"log"
	"math"
	"net"
	"strconv"
	"strings"
	"text/template"

	"github.com/ecadlabs/tezos-bot/models"
)

const (
	proposalZone = "proposal.tezz.ie"
	addressZone  = "tz.tezz.ie"
)

type statusTmplData struct {
	AccountName         string
	ProposalName        string
	Rolls               int64
	Phase               string
	Voted               string
	Period              int
	Ballot              string
	PercentYay          float64
	PercentNay          float64
	PercentTowardQuorum float64
	Quorum              float64
	QuorumReached       bool
}

func Percent(s float64) string {
	return fmt.Sprintf("%s%%", strconv.FormatFloat(math.Round(s*100)/100, 'f', -1, 64))
}

var (
	funcMap = template.FuncMap{
		"Title":   strings.Title,
		"Percent": Percent,
	}
	statusTmpl = template.Must(template.New("status.tmpl").Funcs(funcMap).ParseFiles("./templates/status.tmpl"))
)

// GetStatusString composes a status string based on available vanity data
func GetStatusString(ballot *models.Ballot) (string, error) {
	var tpl bytes.Buffer
	if err := statusTmpl.Execute(&tpl, statusTmplData{
		AccountName:         lookupOrDefault(ballot.PKH, addressZone),
		Rolls:               ballot.Rolls,
		ProposalName:        lookupOrDefault(ballot.ProposalHash, proposalZone),
		Phase:               ballot.Phase(),
		Period:              ballot.Period,
		Ballot:              ballot.Ballot,
		PercentYay:          ballot.CountingPercentYay(),
		PercentNay:          ballot.CountingPercentNay(),
		PercentTowardQuorum: ballot.PercentTowardQuorum(),
		Quorum:              ballot.Quorum,
		QuorumReached:       ballot.PercentTowardQuorum() <= 0,
	}); err != nil {
		return "", err
	}
	return tpl.String(), nil
}

func lookupOrDefault(hash string, zone string) string {
	address, err := LookupTZName(hash, zone)

	if err != nil {
		log.Printf("No value found for %s, err: %s", hash, err)
		return hash
	}
	return address
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
	address, err := LookupTZName(proposal.PKH, addressZone)

	if err != nil {
		log.Printf("No address found for %s, err: %s", proposal.PKH, err)
	}

	templateRolls := ""
	if proposal.Rolls != 0 {
		templateRolls = fmt.Sprintf("with %d rolls ", proposal.Rolls)
	}

	templateAddress := proposal.PKH

	if address != "" {
		templateAddress = address
	}

	proposalName := lookupOrDefault(proposal.ProposalHash, proposalZone)

	return fmt.Sprintf("Address %s upvoted proposal %s %sin voting period %d.\n\nhttps://tezblock.io/account/%s?tab=Votes", templateAddress, proposalName, templateRolls, proposal.Period, proposal.PKH)
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
