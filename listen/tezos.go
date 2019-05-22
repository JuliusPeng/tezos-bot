package listen

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ecadlabs/go-tezos"
	"github.com/ecadlabs/tezos-bot/models"
)

// TezosConfig interface with method necessary to obtain tezos listener configurable parameter
type TezosConfig interface {
	GetRPCURL() string
	GetChainID() string
	GetRetryCount() int
	IsHistory() bool
}

// TezosListener is a struct containing information necessary to monitor the tezos chain
type TezosListener struct {
	service    *tezos.Service
	votesChan  chan *models.Ballot
	cache      *cache
	signals    chan bool
	config     TezosConfig
	bStreaming BlockStreamingFunc
}

// NewTezosListener create a new TezosListener
func NewTezosListener(config TezosConfig) (*TezosListener, error) {
	client, err := tezos.NewRPCClient(http.DefaultClient, config.GetRPCURL())
	if err != nil {
		return nil, err
	}

	var bStreamingFunc BlockStreamingFunc = MonitorBlockStreamingFunc

	if config.IsHistory() {
		bStreamingFunc = HistoryBlockStreamingFunc
	}

	return &TezosListener{
		service:    &tezos.Service{Client: client},
		cache:      newCache(),
		votesChan:  make(chan *models.Ballot),
		signals:    make(chan bool),
		config:     config,
		bStreaming: bStreamingFunc,
	}, nil
}

// Start start monitoring the chain and push new Ballot in the ballot channel
func (t *TezosListener) Start() {
	ctx := context.Background()
	cBlockHash := make(chan string)
	defer close(cBlockHash)
	go func() {
		for {
			select {
			case _ = <-t.signals:
				return
			default:
				err := t.bStreaming(ctx, t.config, t.service, cBlockHash)
				if err != nil {
					panic(err.Error())
				}
				return
			}
		}
	}()

	for hash := range cBlockHash {
		// cBlockHash channel can emit the same has multiple time
		// In order to avoid duplicate we check if it has already been processed
		if !t.cache.Has(hash) {
			t.cache.Add(hash)
			block, err := t.service.GetBlock(ctx, t.config.GetChainID(), hash)

			if err != nil {
				fmt.Printf("Block: %s skipped because of error: %s\n", hash, err.Error())
				continue
			}

			fmt.Printf("TezosListener: Inspecting block %s for new ballot operations.\n", block.Hash)

			ballotOps := []*tezos.BallotOperationElem{}
			for _, group := range block.Operations {
				for _, op := range group {
					ballotOps = append(ballotOps, tezos.FilterBallotOps(op.Contents)...)
				}
			}

			ballots, err := t.service.GetBallots(ctx, t.config.GetChainID(), hash)

			if err != nil {
				fmt.Printf("Block: %s skipped because of error: %s\n", hash, err.Error())
				continue
			}

			listings, err := t.service.GetBallotListings(ctx, t.config.GetChainID(), hash)

			if err != nil {
				fmt.Printf("Block: %s skipped because of error: %s\n", hash, err.Error())
				continue
			}

			totalRolls := int64(0)
			for _, entry := range listings {
				totalRolls += entry.Rolls
			}

			if totalRolls == 0 {
				// Unlikely to occurs
				fmt.Printf("Block: %s skipped because no rolls", hash)
				continue
			}

			quorum, err := t.getQuorum(ctx, hash)

			if err != nil {
				fmt.Printf("Block: %s skipped because of error: %s\n", hash, err.Error())
				continue
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
					TotalRolls:   float64(totalRolls),
				}
				t.votesChan <- ballot
			}
		}
	}
}

func (t *TezosListener) getQuorum(ctx context.Context, block string) (float64, error) {
	quorum, err := t.service.GetCurrentQuorum(ctx, t.config.GetChainID(), block)

	if err != nil {
		return 0, err
	}

	return float64(quorum) / 100, nil
}

// Stop stop the tezos listener
func (t *TezosListener) Stop() {
	t.signals <- true
}

// GetNewVotes returns a Ballot channel
func (t *TezosListener) GetNewVotes() chan *models.Ballot {
	return t.votesChan
}
