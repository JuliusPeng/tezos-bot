package main

import (
	"fmt"

	"github.com/ecadlabs/tezos-bot/config"
	"github.com/ecadlabs/tezos-bot/listen"
	"github.com/ecadlabs/tezos-bot/publish"
	"github.com/ecadlabs/tezos-bot/service"
)

func main() {

	c := config.Config{
		RPCURL:          "https://mainnet-node.tzscan.io",
		ChainID:         "main",
		RetryCount:      100,
		History:         false,
		MonitorVote:     true,
		MonitorProtocol: true,
	}

	c.Load("./config.yaml")

	l, err := listen.NewTezosListener(c)

	if err != nil {
		fmt.Printf(err.Error())
		return
	}

	var p service.VotePublisher

	if c.GetTwitterAccessToken() == "" {
		fmt.Printf("Twitter access token not configured posting vote to stdout\n")
		p = &publish.DebugPublisher{}
	} else {
		p, err = publish.NewTwitterPublisher(c)

		if err != nil {
			fmt.Printf(err.Error())
			return
		}
	}

	s := service.New(l, p)
	s.Start()
}
