package main

//go:generate make

//
// This package demonstrates the go:generate process tied to a Makefile.
//
// This package verifies that all the generated code has been built and comiles.
//
// To run go generate just:
//
//	$ go generate
//
// Then:
//
//	$ go run gen.go
//

import (
	"fmt"

	"github.com/keep-network/keep-core/pkg/chain/gen/KeepRandomBeacon"
	"github.com/keep-network/keep-core/pkg/chain/gen/KeepRandomBeaconImplV1"
	"github.com/keep-network/keep-core/pkg/chain/gen/KeepToken"
	"github.com/keep-network/keep-core/pkg/chain/gen/StakingProxy"
	"github.com/keep-network/keep-core/pkg/chain/gen/TokenGrant"
	"github.com/keep-network/keep-core/pkg/chain/gen/TokenStaking"
)

func main() {
	_ = KeepRandomBeacon.KeepRandomBeacon{}
	_ = KeepRandomBeaconImplV1.KeepRandomBeaconImplV1{}
	_ = KeepToken.KeepToken{}
	_ = StakingProxy.StakingProxy{}
	_ = TokenGrant.TokenGrant{}
	_ = TokenStaking.TokenStaking{}
	fmt.Printf("PASS\n")
}
