package listen

import (
	"context"
	"fmt"
	"log"
	"time"

	tezos "github.com/ecadlabs/go-tezos"
)

const (
	HEAD_BLOCK = "head"
)

// BlockStreamingFunc function that emit a of block
type BlockStreamingFunc func(ctx context.Context, config TezosConfig, service *tezos.Service, results chan<- string) error

// MonitorBlockStreamingFunc emit the hash of each new head
func MonitorBlockStreamingFunc(ctx context.Context, config TezosConfig, service *tezos.Service, results chan<- string) error {
	cMonitorBlock := make(chan *tezos.MonitorBlock)
	errCount := 0
	defer close(cMonitorBlock)
	go func() {
		for block := range cMonitorBlock {
			// Reset the error count on new block
			errCount = 0
			results <- block.Hash
		}
	}()

	for {
		err := service.GetMonitorHeads(ctx, config.GetChainID(), cMonitorBlock)
		if err != nil {
			// Retry connection until retry count is reached
			if errCount > config.GetRetryCount() {
				return fmt.Errorf("Unable to connect to rpc node after %d tries", config.GetRetryCount())
			}

			errCount++
			log.Printf("Error encountered while trying to connect to rpc node (err count: %d): %s\n", errCount, err.Error())
			time.Sleep(time.Duration(errCount) * time.Second)
		}
	}
}

// HistoryBlockStreamingFunc start from genesis and emit the hash of each block until current head
func HistoryBlockStreamingFunc(ctx context.Context, config TezosConfig, service *tezos.Service, results chan<- string) error {
	head, err := service.GetBlock(ctx, config.GetChainID(), HEAD_BLOCK)
	if err != nil {
		return err
	}

	level := head.Header.Level
	i := config.GetHistoryStartingBlock()
	for i <= level {
		block, err := service.GetBlock(ctx, config.GetChainID(), fmt.Sprintf("%d", i))
		if err != nil {
			return err
		}
		results <- block.Hash
		time.Sleep(time.Second * 2)
		i++
	}
	return nil
}
