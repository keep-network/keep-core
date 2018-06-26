package cmd

import (
	"fmt"
	"testing"

	"github.com/keep-network/keep-core/util"
	ma "github.com/multiformats/go-multiaddr"
)

func stringsToMultiAddr(t *testing.T, s []string) []ma.Multiaddr {
	multiAddrs := make([]ma.Multiaddr, len(s))
	for _, addr := range s {
		ma, err := ma.NewMultiaddr(addr)
		util.Ok(t, err)
		multiAddrs = append(multiAddrs, ma)
	}
	return multiAddrs
}

func TestMultiAddrIPs(t *testing.T) {
	myIPAddress, err := GetPreferredOutboundIP()
	util.Ok(t, err)
	myIPv4MultiAddr := fmt.Sprintf("/ip4/%s/tcp/27001", myIPAddress)
	typical := []string{
		"/ip4/127.0.0.1/tcp/27001",
		"/ip6/::1/tcp/27001",
		"/ip4/192.168.10.103/tcp/27001",
	}
	duplicate := []string{
		"/ip4/127.0.0.1/tcp/27001",
		"/ip6/::1/tcp/27001",
		"/ip4/192.168.10.103/tcp/27001",
		"/ip4/192.168.10.103/tcp/27001",
		myIPv4MultiAddr,
		myIPv4MultiAddr,
	}
	only_home := []string{
		"/ip4/127.0.0.1/tcp/27001",
	}
	only_ipv6 := []string{
		"/ip6/::1/tcp/27001",
	}
	empty := []string{}
	tests := map[string]struct {
		maIPs          []ma.Multiaddr
		expectedString string
	}{
		"typical multi address set": {
			maIPs:          stringsToMultiAddr(t, typical),
			expectedString: "192.168.10.103",
		},
		"duplicate address set": {
			maIPs:          stringsToMultiAddr(t, duplicate),
			expectedString: myIPAddress,
		},
		"only home address": {
			maIPs:          stringsToMultiAddr(t, only_home),
			expectedString: "127.0.0.1",
		},
		"only ipv6 address": {
			maIPs:          stringsToMultiAddr(t, only_ipv6),
			expectedString: "127.0.0.1",
		},
		"empty address": {
			maIPs:          stringsToMultiAddr(t, empty),
			expectedString: "127.0.0.1",
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			acutal := GetIPv4Address(test.maIPs)
			util.Equals(t, test.expectedString, acutal)
		})
	}
}
