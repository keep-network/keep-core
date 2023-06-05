package config

import (
	"embed"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/keep-network/keep-core/pkg/bitcoin"
)

//go:embed _electrum_urls/*
var electrumURLs embed.FS

// Keep Alive Interval value used for Blockstream's electrum connections.
// This value is used only if a Blockstream's server is randomly selected from
// the list of embedded Electrum servers. It does not apply if a Blockstream's
// server connection is explicitly set in the client's configuration.
var blockstreamKeepAliveInterval = 55 * time.Second

// readElectrumUrls reads Electrum URLs from an embedded file for the
// given Bitcoin network.
func readElectrumUrls(network bitcoin.Network) (
	[]string,
	error,
) {
	file, err := electrumURLs.ReadFile(fmt.Sprintf("_electrum_urls/%s", network))
	if err != nil {
		return nil, fmt.Errorf("cannot read URLs file: [%v]", err)
	}

	urlsStrings := cleanStrings(strings.Split(string(file), "\n"))

	return urlsStrings, nil
}

// resolveElectrum checks if Electrum is already configured. If the Electrum URL
// is empty it reads the Electrum configs from the embedded list for the given
// network and picks up one randomly.
func (c *Config) resolveElectrum(rng *rand.Rand) error {
	network := c.Bitcoin.Network

	// Return if Electrum is already set.
	if len(c.Bitcoin.Electrum.URL) > 0 {
		return nil
	}

	// For unknown and regtest networks we don't expect the Electrum configs to be
	// embedded in the client. The user should configure it in the config file.
	if network == bitcoin.Regtest || network == bitcoin.Unknown {
		logger.Warnf(
			"Electrum configs were not configured for [%s] network; "+
				"see bitcoin section in configuration",
			network,
		)
		return nil
	}

	logger.Debugf(
		"Electrum was not configured for [%s] bitcoin network; "+
			"reading defaults",
		network,
	)

	urls, err := readElectrumUrls(network)
	if err != nil {
		return fmt.Errorf("failed to read default Electrum URLs: [%v]", err)
	}

	// #nosec G404 (insecure random number source (rand))
	// Picking up an Electrum server does not require secure randomness.
	selectedURL := urls[rng.Intn(len(urls))]

	logger.Infof("Auto-selecting Electrum server: [%v]", selectedURL)

	// Set only the URL in the original config. Other fields may be already set,
	// and we don't want to override them.
	c.Bitcoin.Electrum.URL = selectedURL

	// Blockstream's servers timeout session after 60 seconds of inactivity which
	// is much shorter than expected 600 seconds. To workaround connection drops
	// and logs pollution with warning we reduce the KeepAliveInterval for the
	// Blockstream's servers to less than 60 seconds.
	if c.Bitcoin.Electrum.KeepAliveInterval == 0 &&
		strings.Contains(selectedURL, "electrum.blockstream.info") {
		c.Bitcoin.Electrum.KeepAliveInterval = blockstreamKeepAliveInterval
	}

	return nil
}
