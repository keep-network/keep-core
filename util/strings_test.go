package util

import (
	"fmt"
	"testing"
)

func TestMultiAddrIPs(t *testing.T) {

	tests := map[string]struct {
		list          []string
		hasDuplicates bool
	}{
		"no entries": {
			list:          []string{},
			hasDuplicates: false,
		},
		"one entry": {
			list: []string{
				"/ip4/127.0.0.1/tcp/27001/ipfs/12D3KooWKRyzVWW6ChFjQjK4miCty85Niy49tpPV95XdKu1BcvMA",
			},
			hasDuplicates: false,
		},
		"multiple unique entries": {
			list: []string{
				"/ip4/127.0.0.1/tcp/27001/ipfs/12D3KooWKRyzVWW6ChFjQjK4miCty85Niy49tpPV95XdKu1BcvMA",
				"/ip4/127.0.0.1/tcp/27002/ipfs/12D3KooWKRyzVWW6ChFjQjK4miCty85Niy49tpPV95XdKu1BcvMA",
			},
			hasDuplicates: false,
		},
		"duplicate entries": {
			list: []string{
				"/ip4/127.0.0.1/tcp/27001/ipfs/12D3KooWKRyzVWW6ChFjQjK4miCty85Niy49tpPV95XdKu1BcvMA",
				"/ip4/127.0.0.1/tcp/27001/ipfs/12D3KooWKRyzVWW6ChFjQjK4miCty85Niy49tpPV95XdKu1BcvMA",
			},
			hasDuplicates: true,
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			actual := DuplicatesExist(test.list)
			Equals(t, test.hasDuplicates, actual)
		})
	}
}

func TestMultiAddrIPsWithJoin(t *testing.T) {
	var emptyErrMsg = ""
	const multiAddr = "/ip4/127.0.0.1/tcp/27001/ipfs/12D3KooWKRyzVWW6ChFjQjK4miCty85Niy49tpPV95XdKu1BcvMA"
	const multiAddr2 = "/ip4/127.0.0.1/tcp/27002/ipfs/12D3KooWKRyzVWW6ChFjQjK4miCty85Niy49tpPV95XdKu1BcvMA"
	tests := map[string]struct {
		list   []string
		errMsg string
	}{
		"no entries": {
			list:   []string{},
			errMsg: emptyErrMsg,
		},
		"one entry": {
			list: []string{
				multiAddr,
			},
			errMsg: emptyErrMsg,
		},
		"multiple unique entries": {
			list: []string{
				multiAddr,
				multiAddr2,
			},
			errMsg: emptyErrMsg,
		},
		"duplicate entries": {
			list: []string{
				multiAddr,
				multiAddr,
				multiAddr2,
				multiAddr2,
				multiAddr2,
			},
			errMsg: fmt.Sprintf("Node.Peers invalid; duplicates found: %s %s", multiAddr, multiAddr2),
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {

			var got string
			if DuplicatesExist(test.list) {
				got = fmt.Sprintf("Node.Peers invalid; duplicates found: %s",
					Join(Duplicates(test.list), " "))
			}
			Equals(t, test.errMsg, got)
		})
	}
}
