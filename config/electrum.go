package config

import (
	"embed"
	"fmt"
	"math/rand"
	"strings"

	"github.com/keep-network/keep-core/pkg/bitcoin"
)

//go:embed _electrum_urls/*
var electrumURLs embed.FS

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
func (c *Config) resolveElectrum() error {
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
	selectedURL := urls[rand.Intn(len(urls))]

	logger.Infof("Auto-selecting Electrum server: [%v]", selectedURL)

	// Set only the URL in the original config. Other fields may be already set,
	// and we don't want to override them.
	c.Bitcoin.Electrum.URL = selectedURL

	return nil
}
