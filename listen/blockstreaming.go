package listen

import (
	"context"
	"fmt"
	"time"

	tezos "github.com/ecadlabs/go-tezos"
)

// BlockStreamingFunc function that emit a of block
type BlockStreamingFunc func(ctx context.Context, config TezosConfig, service *tezos.Service, results chan<- string) error

// MonitorBlockStreamingFunc emit the hash of each new head
func MonitorBlockStreamingFunc(ctx context.Context, config TezosConfig, service *tezos.Service, results chan<- string) error {
	cMonitorBlock := make(chan *tezos.MonitorBlock)
	defer close(cMonitorBlock)
	go func() {
		for block := range cMonitorBlock {
			results <- block.Hash
		}
	}()

	errCount := 0
	for {
		err := service.GetMonitorHeads(ctx, config.GetChainID(), cMonitorBlock)
		if err != nil {
			if errCount > config.GetRetryCount() {
				return fmt.Errorf("Unable to connect to rpc node after %d tries", config.GetRetryCount())
			}

			errCount++
			time.Sleep(time.Duration(errCount) * time.Second)
			fmt.Printf("Error encountered while trying to connect to rpc node: %s\n", err.Error())
		} else {
			errCount = 0
		}
	}
}

// HistoryBlockStreamingFunc start from genesis and emit the hash of each block until current head
func HistoryBlockStreamingFunc(ctx context.Context, config TezosConfig, service *tezos.Service, results chan<- string) error {
	head, err := service.GetBlock(ctx, config.GetChainID(), "head")
	if err != nil {
		return err
	}

	level := head.Header.Level
	i := 0
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
