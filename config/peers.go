package config

import (
	"embed"
	"fmt"
	"strings"

	"github.com/keep-network/keep-core/config/network"
)

//go:embed _peers/*
var peersData embed.FS

// readPeers reads peers from an embedded file for the given `network`.
func readPeers(clientNetwork network.Type) ([]string, error) {
	peers, err := peersData.ReadFile(fmt.Sprintf("_peers/%s", clientNetwork))
	if err != nil {
		return nil, err
	}

	return cleanStrings(strings.Split(string(peers), "\n")), nil
}

// cleanStrings iterates over entires in a slice and trims spaces from the beginning
// and the end of a string. It also removes empty entries or entries commented with `#`.
func cleanStrings(s []string) []string {
	var peers []string
	for _, str := range s {
		str = strings.TrimSpace(str)
		if str == "" || strings.HasPrefix(str, "#") {
			continue
		}

		peers = append(peers, str)
	}
	return peers
}

// resolvePeers checks if peers are already set. If the peers list is empty it
// reads the peers from the embedded peers list for the given network.
func (c *Config) resolvePeers(clientNetwork network.Type) error {
	// Return if peers are already set.
	if len(c.LibP2P.Peers) > 0 {
		return nil
	}

	// For unknown and developer networks we don't expect the default peers to be
	// embedded in the client. The user should configure them in the config file.
	if clientNetwork == network.Developer || clientNetwork == network.Unknown {
		logger.Warnf(
			"peers were not configured for [%s] network; "+
				"see network section in configuration",
			clientNetwork,
		)
		return nil
	}

	logger.Debugf(
		"peers were not configured for [%s] network; reading defaults",
		clientNetwork,
	)

	peers, err := readPeers(clientNetwork)
	if err != nil {
		return fmt.Errorf("failed to read default peers: [%v]", err)
	}

	c.LibP2P.Peers = peers

	return nil
}
