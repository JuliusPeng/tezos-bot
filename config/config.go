package config

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

// Config struct containing all configurable parameter for the tezos bot
type Config struct {
	RPCURL                   string `yaml:"rpc_url"`
	TwitterAccessToken       string `yaml:"twitter_access_token"`
	TwitterAccessTokenSecret string `yaml:"twitter_access_token_secret"`
	TwitterConsummerID       string `yaml:"twitter_consummer_id"`
	TwitterConsummerKey      string `yaml:"twitter_consummer_key"`
	ChainID                  string `yaml:"chain"`
	RetryCount               int    `yaml:"retry_count"`
	History                  bool   `yaml:"history"`
	HistoryStartingBlock     int    `yaml:"history_starting_block"`
	MonitorVote              bool   `yaml:"monitor_vote"`
	MonitorProtocol          bool   `yaml:"monitor_protocol"`
}

// GetHistoryStartingBlock return the starting block from which the bot should start monitring
func (c Config) GetHistoryStartingBlock() int {
	return c.HistoryStartingBlock
}

// GetRPCURL returns the tezos rpc endpoint
func (c Config) GetRPCURL() string {
	return c.RPCURL
}

// IsHistory returns true if listener should read history
func (c Config) IsHistory() bool {
	return c.History
}

// GetRetryCount returns the max retry attempt to connect to tezos node
func (c Config) GetRetryCount() int {
	return c.RetryCount
}

// GetChainID returns the chain ID
func (c Config) GetChainID() string {
	return c.ChainID
}

// IsMonitorVote returns true if should monitor vote
func (c Config) IsMonitorVote() bool {
	return c.MonitorVote
}

// IsMonitorProtocol returns true if should monitor protocol change
func (c Config) IsMonitorProtocol() bool {
	return c.MonitorProtocol
}

// GetTwitterAccessToken returns the twitter access token
func (c Config) GetTwitterAccessToken() string {
	return c.TwitterAccessToken
}

// GetTwitterAccessTokenSecret returns the twitter access token secret
func (c Config) GetTwitterAccessTokenSecret() string {
	return c.TwitterAccessTokenSecret
}

// GetTwitterConsummerID returns the twitter consummer id
func (c Config) GetTwitterConsummerID() string {
	return c.TwitterConsummerID
}

// GetTwitterConsummerKey returns the twitter consummer key
func (c Config) GetTwitterConsummerKey() string {
	return c.TwitterConsummerKey
}

// Load read a config file and unmarshal it into the config struct
func (c *Config) Load(name string) error {
	buf, err := ioutil.ReadFile(name)
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(buf, c); err != nil {
		return err
	}

	return nil
}
