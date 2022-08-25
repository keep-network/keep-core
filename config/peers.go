package config

import (
	"embed"
	"fmt"
	"strings"

	"github.com/keep-network/keep-common/pkg/chain/ethereum"
)

//go:embed _peers/*
var peersData embed.FS

// readPeers reads peers from an embedded file for the given ethereum `network`.
func readPeers(network ethereum.Network) ([]string, error) {
	peers, err := peersData.ReadFile(fmt.Sprintf("_peers/%s", network))
	if err != nil {
		return nil, err
	}

	return cleanStrings(strings.Split(string(peers), "\n")), nil
}

// cleanStrings iterates over entires in a slice and trims spaces from the beginning
// and the end of a string. It also removes empty entries or entries commented with `#`.
func cleanStrings(s []string) []string {
	var r []string
	for _, str := range s {
		str = strings.TrimSpace(str)
		if str == "" || strings.HasPrefix(str, "#") {
			continue
		}

		r = append(r, str)
	}
	return r
}

// resolvePeers checks if peers are already set. If the peers list is empty it
// reads the peers from the embedded peers list for the given network.
func (c *Config) resolvePeers() error {
	network := c.Ethereum.Network

	// Return if peers are already set.
	if len(c.LibP2P.Peers) > 0 {
		return nil
	}

	// For developer network we don't expect the default peers to be embedded in
	// the client. The user should configure them in the config file.
	if network == ethereum.Developer {
		logger.Warnf(
			"peers were not configured for [%s] network; "+
				"see network section in configuration",
			network,
		)
		return nil
	}

	logger.Debugf(
		"peers were not configured for [%s] ethereum network; reading defaults",
		network,
	)

	peers, err := readPeers(network)
	if err != nil {
		return fmt.Errorf("failed to read default peers: [%v]", err)
	}

	c.LibP2P.Peers = peers

	return nil
}
