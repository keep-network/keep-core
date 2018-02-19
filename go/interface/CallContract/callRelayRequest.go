package main

import (
	"bufio"
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi/bind" // "www.2c-why.com/job/ethrpc/ethABI/bind"
	"github.com/pschlump/MiscLib"
	// "github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/keep.network/keep-core/go/interface/lib/KStart" // "github.com/keep.network/keep-core/go/interface/m4/KStart"
)

func main() {

	// Create RPC connection to a remote node
	// addr := "http://10.51.245.75:8545"
	// addr := "http://10.51.245.75:8546" // WebSocket connect
	addr := "/Users/corwin/Projects/eth/data/geth.ipc"
	conn, err := ethclient.Dial(addr)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v at address: %s", err, addr)
	}

	caddr := "0xe705ab560794cf4912960e5069d23ad6420acde7"
	fmt.Printf("Connected - v12 %s(Call requestRelay to create event)%s - test on testnet(%s)\n", MiscLib.ColorCyan, MiscLib.ColorReset, addr)

	// Instantiate the contract and call a method
	kstart, err := KStart.NewKStartTransactor(common.HexToAddress(caddr), conn)
	if err != nil {
		// log.Fatalf("Failed to instantiate a Token contract: %v", err)
		log.Fatalf("Failed to instantiate a KStartTranactor contract: %v", err)
	}

	fmt.Printf("----------------------------------------------------------------------------------\n")
	fmt.Printf("Connect to KStart contract successful\n")
	fmt.Printf("----------------------------------------------------------------------------------\n")

	keyFile := "./UTC--2018-02-15T19-57-35.216297214Z--6ffba2d0f4c8fd7961f516af43c55fe2d56f6044"
	keyFilePassword := "password"

	file, err := os.Open(keyFile)
	if err != nil {
		log.Fatalf("Failed to open keyfile: %v, %s", err, keyFile)
	}

	opts, err := bind.NewTransactor(bufio.NewReader(file), keyFilePassword)
	if err != nil {
		log.Fatalf("Failed to read keyfile: %v, %s", err, keyFile)
	}

	seed := []byte("abcdef")
	_payment := big.NewInt(12)
	_blockReward := big.NewInt(12)
	_seed := big.NewInt(0).SetBytes(seed)

	// function requestRelay(uint256 _payment, uint256 _blockReward, uint256 _seed) public returns ( uint256 RequestID ) {
	tx, err := kstart.RequestRelay(opts, _payment, _blockReward, _seed)
	if err != nil {
		log.Fatalf("Failed to call KStart.requestRelay: %s", err)
	}

	fmt.Printf("----------------------------------------------------------------------------------\n")
	fmt.Printf("%sKStart.requestRelay called%s tx=%s\n", MiscLib.ColorGreen, MiscLib.ColorReset, tx)
	fmt.Printf("----------------------------------------------------------------------------------\n")

}
