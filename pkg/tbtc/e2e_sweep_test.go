package tbtc

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/ipfs/go-log"
	"github.com/keep-network/keep-common/pkg/persistence"
	"github.com/keep-network/keep-core/pkg/bitcoin"
	"github.com/keep-network/keep-core/pkg/bitcoin/electrum"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/chain/local_v1"
	"github.com/keep-network/keep-core/pkg/generator"
	"github.com/keep-network/keep-core/pkg/net/local"
	"github.com/keep-network/keep-core/pkg/operator"
	"github.com/keep-network/keep-core/pkg/tecdsa"
	"os"
	"os/exec"
	"strings"
	"sync"
	"testing"
	"time"
)

// Input parameters.
const (
	clientsCount          = 10
	walletPublicKeyString = "f863124437b72fe0ec0359592053b281c6ca625194073418c35fb28435cf3a8c74d45b17cf60541eeef4f52786670da9710ed804ccb45e93a8d4cc7e15a5fee0"
	fundingTxHash         = "237ff2157b4d4bd5a5dad182e07e8fa270bc762e0f2f16c8dcf060b61652d551"
	fundingTxOutputIndex  = 0
	fundingTxOutputValue  = 1800000
	depositor             = "68ad60CC5e8f3B7cC53beaB321cf0e6036962dBc"
	blindingFactor        = "7b8ec9ad2b95a36f"
	walletPublicKeyHash   = "621dee7e4e7a9273fcca01855cda53243b9820e2"
	refundPublicKeyHash   = "7ac2d9378a1c47e589dfb8095ca95ed2140d2726"
	refundLocktime        = "42fd2c65"
)

func TestE2ESweep(t *testing.T) {
	_ = log.SetLogLevel("*", "INFO")

	persistenceDir := "persistence"

	walletPublicKey := unmarshalEcdsaPublicKey(walletPublicKeyString)

	depositToSweep := unmarshalDeposit(
		fundingTxHash,
		fundingTxOutputIndex,
		fundingTxOutputValue,
		depositor,
		blindingFactor,
		walletPublicKeyHash,
		refundPublicKeyHash,
		refundLocktime,
	)

	err := os.RemoveAll(persistenceDir)
	if err != nil {
		t.Fatal(err)
	}

	err = os.Mkdir(persistenceDir, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}

	clientsPersistences := make([]persistence.ProtectedHandle, clientsCount)

	for i := 0; i < clientsCount; i++ {
		client := fmt.Sprintf("keep-client-%v", i)
		pod := fmt.Sprintf("%v-0", client)
		clientPersistence := fmt.Sprintf("%v/%v", persistenceDir, client)

		err = os.MkdirAll(
			fmt.Sprintf(
				"%v/current/%v",
				clientPersistence,
				walletPublicKeyString,
			), os.ModePerm,
		)
		if err != nil {
			t.Fatal(err)
		}

		err = os.MkdirAll(
			fmt.Sprintf("%v/archive", clientPersistence),
			os.ModePerm,
		)
		if err != nil {
			t.Fatal(err)
		}

		err = os.MkdirAll(
			fmt.Sprintf("%v/snapshot", clientPersistence),
			os.ModePerm,
		)
		if err != nil {
			t.Fatal(err)
		}

		err = exec.Command(
			"kubectl",
			"cp",
			fmt.Sprintf(
				"%v:mnt/keep-client/data/keystore/tbtc/current/%v",
				pod,
				walletPublicKeyString,
			),
			fmt.Sprintf(
				"%v/current/%v",
				clientPersistence,
				walletPublicKeyString,
			),
		).Run()
		if err != nil {
			t.Fatal(err)
		}
		fmt.Printf("fetched key shares from [%v]\n", client)

		output, err := exec.Command(
			"kubectl",
			"get",
			"secret",
			"eth-account-passphrases",
			"-o",
			fmt.Sprintf("jsonpath='{.data.account-%v}'", i),
		).Output()
		if err != nil {
			t.Fatal(err)
		}

		decodedOutput, err := base64.StdEncoding.DecodeString(
			strings.Trim(string(output), "'"),
		)
		if err != nil {
			t.Fatal(err)
		}

		password := string(decodedOutput)

		fmt.Printf("fetched encryption password for [%v]\n", client)

		protectedPersistence, err := persistence.NewProtectedDiskHandle(
			clientPersistence,
		)
		if err != nil {
			t.Fatal(err)
		}

		clientsPersistences[i] = persistence.NewEncryptedProtectedPersistence(
			protectedPersistence,
			password,
		)
	}

	fmt.Printf("fetched everything from kubernetes cluster\n")

	groupParameters := &GroupParameters{
		GroupSize:       100,
		GroupQuorum:     90,
		HonestThreshold: 51,
	}

	operatorPrivateKey, operatorPublicKey, err := operator.GenerateKeyPair(
		local_v1.DefaultCurve,
	)
	if err != nil {
		t.Fatal(err)
	}

	localChain := ConnectWithKey(operatorPrivateKey)

	operatorAddress, err := localChain.Signing().PublicKeyToAddress(
		operatorPublicKey,
	)
	if err != nil {
		t.Fatal(err)
	}

	var operators []chain.Address
	for i := 0; i < groupParameters.GroupSize; i++ {
		operators = append(operators, operatorAddress)
	}

	localProvider := local.ConnectWithKey(operatorPublicKey)

	node, err := newNode(
		groupParameters,
		localChain,
		localProvider,
		&compositeProtectedPersistence{clientsPersistences},
		&mockPersistenceHandle{},
		generator.StartScheduler(),
		Config{},
	)
	if err != nil {
		t.Fatal(err)
	}

	// Override signing group operators with local operator address in order to
	// pass the validation of the network layer.
	signers := node.walletRegistry.getSigners(walletPublicKey)
	for _, signer := range signers {
		signer.wallet.signingGroupOperators = operators
	}

	fmt.Printf("started local tbtc node\n")

	electrumClient, err := electrum.Connect(
		context.Background(),
		electrum.Config{
			URL:                 "electrum.blockstream.info:60002",
			Protocol:            electrum.SSL,
			RequestTimeout:      2 * time.Second,
			RequestRetryTimeout: 4 * time.Second,
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("connected to electrum server\n")

	sweepTxBuilder, err := assembleDepositSweepTransaction(
		electrumClient,
		walletPublicKey,
		nil,
		[]*deposit{depositToSweep},
		1000,
	)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("assembled unsigned sweep tx\n")

	sighashes, err := sweepTxBuilder.ComputeSignatureHashes()
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("computed signature hashes\n")

	executor, ok, err := node.getSigningExecutor(walletPublicKey)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("no signing executor for wallet")
	}

	fmt.Printf("computing transaction signatures...\n")

	bc, err := localChain.BlockCounter()
	if err != nil {
		t.Fatal(err)
	}

	currentBlock, err := bc.CurrentBlock()
	if err != nil {
		t.Fatal(err)
	}

	signatures, err := executor.signBatch(
		context.Background(),
		sighashes,
		currentBlock + 1,
	)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("computed transaction signatures\n")

	signatureContainers := make([]*bitcoin.SignatureContainer, len(signatures))
	for i, s := range signatures {
		signatureContainers[i] = &bitcoin.SignatureContainer{
			R:         s.R,
			S:         s.S,
			PublicKey: walletPublicKey,
		}
	}

	sweepTx, err := sweepTxBuilder.AddSignatures(signatureContainers)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("signed sweep tx: %x\n", sweepTx.Serialize())
}

func unmarshalEcdsaPublicKey(publicKey string) *ecdsa.PublicKey {
	publicKeyBytes, err := hex.DecodeString("04" + publicKey)
	if err != nil {
		panic(err)
	}

	x, y := elliptic.Unmarshal(tecdsa.Curve, publicKeyBytes)

	return &ecdsa.PublicKey{
		Curve: tecdsa.Curve,
		X:     x,
		Y:     y,
	}
}

func unmarshalDeposit(
	fundingTxHash string,
	fundingTxOutputIndex uint32,
	fundingTxOutputValue int64,
	depositor string,
	blindingFactor string,
	walletPublicKeyHash string,
	refundPublicKeyHash string,
	refundLocktime string,
	) *deposit {
	hexToSlice := func(hexString string) []byte {
		bytes, err := hex.DecodeString(hexString)
		if err != nil {
			panic(err)
		}
		return bytes
	}

	hash, err := bitcoin.NewHashFromString(
		fundingTxHash,
		bitcoin.InternalByteOrder,
	)
	if err != nil {
		panic(err)
	}

	d := new(deposit)

	d.utxo = &bitcoin.UnspentTransactionOutput{
		Outpoint: &bitcoin.TransactionOutpoint{
			TransactionHash: hash,
			OutputIndex:     fundingTxOutputIndex,
		},
		Value: fundingTxOutputValue,
	}

	copy(d.depositor[:], hexToSlice(depositor))
	copy(d.blindingFactor[:], hexToSlice(blindingFactor))
	copy(d.walletPublicKeyHash[:], hexToSlice(walletPublicKeyHash))
	copy(d.refundPublicKeyHash[:], hexToSlice(refundPublicKeyHash))
	copy(d.refundLocktime[:], hexToSlice(refundLocktime))

	return d
}

type compositeProtectedPersistence struct {
	persistences []persistence.ProtectedHandle
}

func (cpp *compositeProtectedPersistence) ReadAll() (
	<-chan persistence.DataDescriptor,
	<-chan error,
) {
	compositeDataChan := make(chan persistence.DataDescriptor)
	compositeErrChan := make(chan error)

	go func() {
		defer close(compositeDataChan)
		defer close(compositeErrChan)

		for _, p := range cpp.persistences {
			dataChan, errChan := p.ReadAll()

			var wg sync.WaitGroup
			wg.Add(2)

			go func() {
				for descriptor := range dataChan {
					compositeDataChan <- descriptor
				}

				wg.Done()
			}()

			go func() {
				for err := range errChan {
					compositeErrChan <- err
				}

				wg.Done()
			}()

			wg.Wait()
		}
	}()

	return compositeDataChan, compositeErrChan
}

func (cpp *compositeProtectedPersistence) Save(
	data []byte,
	directory string,
	name string,
) error {
	panic("not supported")
}

func (cpp *compositeProtectedPersistence) Archive(directory string) error {
	panic("not supported")
}

func (cpp *compositeProtectedPersistence) Snapshot(
	data []byte,
	directory string,
	name string,
) error {
	panic("not supported")
}
