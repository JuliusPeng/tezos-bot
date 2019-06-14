package listen

import (
	"context"
	"fmt"
	"log"

	tezos "github.com/ecadlabs/go-tezos"
	"github.com/ecadlabs/tezos-bot/models"
)

func (t *TezosListener) lookForBallot(ctx context.Context, block *tezos.Block, periodKind tezos.PeriodKind) error {
	hash := block.Hash
	log.Printf("TezosListener: Inspecting block %s for new ballot operations.\n", block.Hash)

	ballotOps := []*tezos.BallotOperationElem{}
	for _, group := range block.Operations {
		for _, op := range group {
			ballotOps = append(ballotOps, tezos.FilterBallotOps(op.Contents)...)
		}
	}

	ballots, err := t.service.GetBallots(ctx, t.config.GetChainID(), hash)

	if err != nil {
		return err
	}

	listings, err := t.service.GetBallotListings(ctx, t.config.GetChainID(), hash)

	if err != nil {
		return err
	}

	totalRolls := int64(0)
	for _, entry := range listings {
		totalRolls += entry.Rolls
	}

	if totalRolls == 0 {
		// Unlikely to occurs
		return fmt.Errorf("No rolls found in this block")
	}

	quorum, err := t.getQuorum(ctx, hash)

	if err != nil {
		return err
	}

	for _, ballotOp := range ballotOps {
		rolls := int64(0)
		for _, entry := range listings {
			if entry.PKH == ballotOp.Source {
				rolls = entry.Rolls
			}
		}
		ballot := &models.Ballot{
			PKH:          ballotOp.Source,
			Ballot:       ballotOp.Ballot,
			ProposalHash: ballotOp.Proposal,
			Rolls:        rolls,
			Yay:          ballots.Yay,
			Nay:          ballots.Nay,
			Pass:         ballots.Pass,
			Quorum:       quorum,
			IsTesting:    periodKind.IsTestingVote(),
			TotalRolls:   float64(totalRolls),
		}
		t.votesChan <- ballot
	}
	return nil
}

func (t *TezosListener) getQuorum(ctx context.Context, block string) (float64, error) {
	quorum, err := t.service.GetCurrentQuorum(ctx, t.config.GetChainID(), block)

	if err != nil {
		return 0, err
	}

	return float64(quorum) / 100, nil
}
