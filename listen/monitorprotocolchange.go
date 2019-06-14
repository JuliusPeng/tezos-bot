package listen

import (
	"context"
	"log"

	tezos "github.com/ecadlabs/go-tezos"
)

func (t *TezosListener) lookForProtocolChange(ctx context.Context, block *tezos.Block) error {
	log.Printf("TezosListener: Inspecting block %s for protocol changes.\n", block.Hash)

	pred := lastBlock
	if lastBlock == nil {
		predHash := block.Header.Predecessor
		b, err := t.service.GetBlock(ctx, t.config.GetChainID(), predHash)
		if err != nil {
			return err
		}
		pred = b
	}

	if block.Protocol != pred.Protocol {
		t.protoChan <- block.Protocol
	}

	lastBlock = block

	return nil
}
