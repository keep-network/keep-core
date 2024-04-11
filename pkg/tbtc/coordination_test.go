package tbtc

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
	"reflect"
	"testing"
	"time"

	"github.com/go-test/deep"
	"github.com/keep-network/keep-core/pkg/bitcoin"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/chain/local_v1"
	"github.com/keep-network/keep-core/pkg/generator"
	"github.com/keep-network/keep-core/pkg/net"
	netlocal "github.com/keep-network/keep-core/pkg/net/local"
	"github.com/keep-network/keep-core/pkg/operator"
	"github.com/keep-network/keep-core/pkg/protocol/group"
	"github.com/keep-network/keep-core/pkg/tecdsa"
	"golang.org/x/exp/slices"

	"github.com/keep-network/keep-core/internal/testutils"
)

func TestCoordinationWindow_ActivePhaseEndBlock(t *testing.T) {
	window := newCoordinationWindow(900)

	testutils.AssertIntsEqual(
		t,
		"active phase end block",
		980,
		int(window.activePhaseEndBlock()),
	)
}

func TestCoordinationWindow_EndBlock(t *testing.T) {
	window := newCoordinationWindow(900)

	testutils.AssertIntsEqual(
		t,
		"end block",
		1000,
		int(window.endBlock()),
	)
}

func TestCoordinationWindow_IsAfter(t *testing.T) {
	window := newCoordinationWindow(1800)

	previousWindow := newCoordinationWindow(900)
	sameWindow := newCoordinationWindow(1800)
	nextWindow := newCoordinationWindow(2700)

	testutils.AssertBoolsEqual(
		t,
		"result for nil",
		true,
		window.isAfter(nil),
	)
	testutils.AssertBoolsEqual(
		t,
		"result for previous window",
		true,
		window.isAfter(previousWindow),
	)
	testutils.AssertBoolsEqual(
		t,
		"result for same window",
		false,
		window.isAfter(sameWindow),
	)
	testutils.AssertBoolsEqual(
		t,
		"result for next window",
		false,
		window.isAfter(nextWindow),
	)
}

func TestCoordinationWindow_Index(t *testing.T) {
	tests := map[string]struct {
		coordinationBlock uint64
		expectedIndex     uint64
	}{
		"block 0": {
			coordinationBlock: 0,
			expectedIndex:     0,
		},
		"block 900": {
			coordinationBlock: 900,
			expectedIndex:     1,
		},
		"block 1800": {
			coordinationBlock: 1800,
			expectedIndex:     2,
		},
		"block 9000": {
			coordinationBlock: 9000,
			expectedIndex:     10,
		},
		"block 9001": {
			coordinationBlock: 9001,
			expectedIndex:     0,
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			window := newCoordinationWindow(test.coordinationBlock)

			testutils.AssertIntsEqual(
				t,
				"index",
				int(test.expectedIndex),
				int(window.index()),
			)
		})
	}
}

func TestWatchCoordinationWindows(t *testing.T) {
	watchBlocksFn := func(ctx context.Context) <-chan uint64 {
		blocksChan := make(chan uint64)

		go func() {
			ticker := time.NewTicker(1 * time.Millisecond)
			defer ticker.Stop()

			block := uint64(0)

			for {
				select {
				case <-ticker.C:
					block++
					blocksChan <- block
				case <-ctx.Done():
					return
				}
			}
		}()

		return blocksChan
	}

	receivedWindows := make([]*coordinationWindow, 0)
	onWindowFn := func(window *coordinationWindow) {
		receivedWindows = append(receivedWindows, window)
	}

	ctx, cancelCtx := context.WithTimeout(
		context.Background(),
		2000*time.Millisecond,
	)
	defer cancelCtx()

	go watchCoordinationWindows(ctx, watchBlocksFn, onWindowFn)

	<-ctx.Done()

	testutils.AssertIntsEqual(t, "received windows", 2, len(receivedWindows))
	testutils.AssertIntsEqual(
		t,
		"first window",
		900,
		int(receivedWindows[0].coordinationBlock),
	)
	testutils.AssertIntsEqual(
		t,
		"second window",
		1800,
		int(receivedWindows[1].coordinationBlock),
	)
}

func TestCoordinationExecutor_Coordinate(t *testing.T) {
	// Uncompressed public key corresponding to the 20-byte public key hash:
	// aa768412ceed10bd423c025542ca90071f9fb62d.
	publicKeyHex, err := hex.DecodeString(
		"0471e30bca60f6548d7b42582a478ea37ada63b402af7b3ddd57f0c95bb6843175" +
			"aa0d2053a91a050a6797d85c38f2909cb7027f2344a01986aa2f9f8ca7a0c289",
	)
	if err != nil {
		t.Fatal(err)
	}

	// 20-byte public key hash corresponding to the public key above.
	buffer, err := hex.DecodeString("aa768412ceed10bd423c025542ca90071f9fb62d")
	if err != nil {
		t.Fatal(err)
	}
	var publicKeyHash [20]byte
	copy(publicKeyHash[:], buffer)

	parseScript := func(script string) bitcoin.Script {
		parsed, err := hex.DecodeString(script)
		if err != nil {
			t.Fatal(err)
		}

		return parsed
	}

	coordinationBlock := uint64(900)

	type operatorFixture struct {
		chain              Chain
		address            chain.Address
		channel            net.BroadcastChannel
		waitForBlockHeight func(ctx context.Context, blockHeight uint64) error
	}

	generateOperator := func(privateKey int64) *operatorFixture {
		// Generate operators with deterministic addresses that don't change
		// between test runs. This is required to assert the leader selection.
		privateKeyBigInt := big.NewInt(privateKey)
		x, y := local_v1.DefaultCurve.ScalarBaseMult(privateKeyBigInt.Bytes())

		localChain := ConnectWithKey(
			&operator.PrivateKey{
				PublicKey: operator.PublicKey{
					Curve: operator.Secp256k1,
					X:     x,
					Y:     y,
				},
				D: privateKeyBigInt,
			},
			100*time.Millisecond,
		)

		localChain.setBlockHashByNumber(
			coordinationBlock-32,
			"1422996cbcbc38fc924a46f4df5f9064279d3ab43396e58386dac9b87440d64f",
		)

		operatorAddress, err := localChain.operatorAddress()
		if err != nil {
			t.Fatal(err)
		}

		_, operatorPublicKey, err := localChain.OperatorKeyPair()
		if err != nil {
			t.Fatal(err)
		}

		broadcastChannel, err := netlocal.ConnectWithKey(operatorPublicKey).
			BroadcastChannelFor("test")
		if err != nil {
			t.Fatal(err)
		}

		broadcastChannel.SetUnmarshaler(func() net.TaggedUnmarshaler {
			return &coordinationMessage{}
		})

		waitForBlockHeight := func(ctx context.Context, blockHeight uint64) error {
			blockCounter, err := localChain.BlockCounter()
			if err != nil {
				return err
			}

			wait, err := blockCounter.BlockHeightWaiter(blockHeight)
			if err != nil {
				return err
			}

			select {
			case <-wait:
			case <-ctx.Done():
			}

			return nil
		}

		return &operatorFixture{
			chain:              localChain,
			address:            operatorAddress,
			channel:            broadcastChannel,
			waitForBlockHeight: waitForBlockHeight,
		}
	}

	operator1 := generateOperator(1)
	operator2 := generateOperator(2)
	operator3 := generateOperator(3)

	coordinatedWallet := wallet{
		publicKey: unmarshalPublicKey(publicKeyHex),
		signingGroupOperators: []chain.Address{
			operator2.address,
			operator3.address,
			operator1.address,
			operator1.address,
			operator3.address,
			operator2.address,
			operator2.address,
			operator3.address,
			operator1.address,
			operator1.address,
		},
	}

	proposalGenerator := newMockCoordinationProposalGenerator(
		func(
			walletPublicKeyHash [20]byte,
			actionsChecklist []WalletActionType,
		) (CoordinationProposal, error) {
			for _, action := range actionsChecklist {
				if walletPublicKeyHash == publicKeyHash && action == ActionRedemption {
					return &RedemptionProposal{
						RedeemersOutputScripts: []bitcoin.Script{
							parseScript("00148db50eb52063ea9d98b3eac91489a90f738986f6"),
							parseScript("76a9148db50eb52063ea9d98b3eac91489a90f738986f688ac"),
						},
						RedemptionTxFee: big.NewInt(10000),
					}, nil
				}
			}

			return &NoopProposal{}, nil
		},
	)

	membershipValidator := group.NewMembershipValidator(
		&testutils.MockLogger{},
		coordinatedWallet.signingGroupOperators,
		Connect().Signing(),
	)

	protocolLatch := generator.NewProtocolLatch()

	generateExecutor := func(operator *operatorFixture) *coordinationExecutor {
		return newCoordinationExecutor(
			operator.chain,
			coordinatedWallet,
			coordinatedWallet.membersByOperator(operator.address),
			operator.address,
			proposalGenerator,
			operator.channel,
			membershipValidator,
			protocolLatch,
			operator.waitForBlockHeight,
		)
	}

	window := newCoordinationWindow(coordinationBlock)

	type report struct {
		operatorIndex int
		result        *coordinationResult
		err           error
	}

	reportChan := make(chan *report, 3)

	for i, currentOperator := range []*operatorFixture{
		operator1,
		operator2,
		operator3,
	} {
		go func(operatorIndex int, operator *operatorFixture) {
			result, err := generateExecutor(operator).coordinate(window)

			reportChan <- &report{
				operatorIndex: operatorIndex,
				result:        result,
				err:           err,
			}
		}(i+1, currentOperator)
	}

	reports := make([]*report, 0)
loop:
	//lint:ignore S1000 for-select is used as the channel is not closed by senders.
	for {
		select {
		case r := <-reportChan:
			reports = append(reports, r)

			if len(reports) == 3 {
				break loop
			}
		}
	}

	slices.SortFunc(reports, func(i, j *report) int {
		return i.operatorIndex - j.operatorIndex
	})

	testutils.AssertIntsEqual(t, "reports count", 3, len(reports))

	expectedResult := &coordinationResult{
		wallet: coordinatedWallet,
		window: window,
		leader: operator2.address,
		proposal: &RedemptionProposal{
			RedeemersOutputScripts: []bitcoin.Script{
				parseScript("00148db50eb52063ea9d98b3eac91489a90f738986f6"),
				parseScript("76a9148db50eb52063ea9d98b3eac91489a90f738986f688ac"),
			},
			RedemptionTxFee: big.NewInt(10000),
		},
		faults: nil,
	}

	expectedReports := []*report{
		{
			operatorIndex: 1,
			result:        expectedResult,
			err:           nil,
		},
		{
			operatorIndex: 2,
			result:        expectedResult,
			err:           nil,
		},
		{
			operatorIndex: 3,
			result:        expectedResult,
			err:           nil,
		},
	}
	if !reflect.DeepEqual(expectedReports, reports) {
		t.Errorf(
			"unexpected reports:\n"+
				"expected: %v\n"+
				"actual:   %v",
			expectedReports,
			reports,
		)

	}

	testutils.AssertBoolsEqual(
		t,
		"protocol latch state",
		false,
		protocolLatch.IsExecuting(),
	)
}

func TestCoordinationExecutor_GetSeed(t *testing.T) {
	coordinationBlock := uint64(900)

	localChain := Connect()

	localChain.setBlockHashByNumber(
		coordinationBlock-32,
		"1322996cbcbc38fc924a46f4df5f9064279d3ab43396e58386dac9b87440d64f",
	)

	// Uncompressed public key corresponding to the 20-byte public key hash:
	// aa768412ceed10bd423c025542ca90071f9fb62d.
	publicKeyHex, err := hex.DecodeString(
		"0471e30bca60f6548d7b42582a478ea37ada63b402af7b3ddd57f0c95bb6843175" +
			"aa0d2053a91a050a6797d85c38f2909cb7027f2344a01986aa2f9f8ca7a0c289",
	)
	if err != nil {
		t.Fatal(err)
	}

	coordinatedWallet := wallet{
		// Set only relevant fields.
		publicKey: unmarshalPublicKey(publicKeyHex),
	}

	executor := &coordinationExecutor{
		// Set only relevant fields.
		chain:             localChain,
		coordinatedWallet: coordinatedWallet,
	}

	seed, err := executor.getSeed(coordinationBlock)
	if err != nil {
		t.Fatal(err)
	}

	// Expected seed is sha256(wallet_public_key_hash | safe_block_hash).
	expectedSeed := "e55c779d6d83183409ddc90c6cd5130567f0593349a9c82494b402048ec2d03d"

	testutils.AssertStringsEqual(
		t,
		"coordination seed",
		expectedSeed,
		hex.EncodeToString(seed[:]),
	)
}

func TestCoordinationExecutor_GetLeader(t *testing.T) {
	seedBytes, err := hex.DecodeString(
		"9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08",
	)
	if err != nil {
		t.Fatal(err)
	}

	var seed [32]byte
	copy(seed[:], seedBytes)

	coordinatedWallet := wallet{
		// Set only relevant fields.
		signingGroupOperators: []chain.Address{
			"957ECF59507a6A74b8d98747f07a74De270D3CC3", // member 1
			"5E14c0f27612fbfB7A6FE40b5A6Ec997fA62fc04", // member 2
			"D2662604f8b4540336fBd3c1F48d7e9cdFbD079c", // member 3
			"7CBD87ABC182216A7Aa0E8d19aA21abFA2511383", // member 4
			"FAc73b03884d94a08a5c6c7BB12Ac0b20571F162", // member 5
			"705C76445651530fe0D25eeE287b6164cE2c7216", // member 6
			"7CBD87ABC182216A7Aa0E8d19aA21abFA2511383", // member 7  (same operator as member 4)
			"405ad1f632b49A0617fbdc1fD427aF54BA9Bb3dd", // member 8
			"7CBD87ABC182216A7Aa0E8d19aA21abFA2511383", // member 9  (same operator as member 4)
			"5E14c0f27612fbfB7A6FE40b5A6Ec997fA62fc04", // member 10 (same operator as member 2)
		},
	}

	executor := &coordinationExecutor{
		// Set only relevant fields.
		coordinatedWallet: coordinatedWallet,
	}

	leader := executor.getLeader(seed)

	testutils.AssertStringsEqual(
		t,
		"coordination leader",
		"D2662604f8b4540336fBd3c1F48d7e9cdFbD079c",
		leader.String(),
	)
}

func TestCoordinationExecutor_GetActionsChecklist(t *testing.T) {
	tests := map[string]struct {
		coordinationBlock uint64
		expectedChecklist []WalletActionType
	}{
		// Incorrect coordination window.
		"block 0": {
			coordinationBlock: 0,
			expectedChecklist: nil,
		},
		"block 900": {
			coordinationBlock: 900,
			expectedChecklist: []WalletActionType{ActionRedemption},
		},
		// Incorrect coordination window.
		"block 901": {
			coordinationBlock: 901,
			expectedChecklist: nil,
		},
		"block 1800": {
			coordinationBlock: 1800,
			expectedChecklist: []WalletActionType{ActionRedemption},
		},
		"block 2700": {
			coordinationBlock: 2700,
			expectedChecklist: []WalletActionType{ActionRedemption},
		},
		// Heartbeat randomly selected for the 4th coordination window.
		"block 3600": {
			coordinationBlock: 3600,
			expectedChecklist: []WalletActionType{
				ActionRedemption,
				ActionDepositSweep,
				ActionHeartbeat,
			},
		},
		"block 4500": {
			coordinationBlock: 4500,
			expectedChecklist: []WalletActionType{ActionRedemption},
		},
		"block 5400": {
			coordinationBlock: 5400,
			expectedChecklist: []WalletActionType{
				ActionRedemption,
			},
		},
		"block 6300": {
			coordinationBlock: 6300,
			expectedChecklist: []WalletActionType{ActionRedemption},
		},
		"block 7200": {
			coordinationBlock: 7200,
			expectedChecklist: []WalletActionType{
				ActionRedemption,
				ActionDepositSweep,
			},
		},
		"block 8100": {
			coordinationBlock: 8100,
			expectedChecklist: []WalletActionType{ActionRedemption},
		},
		"block 9000": {
			coordinationBlock: 9000,
			expectedChecklist: []WalletActionType{ActionRedemption},
		},
		"block 9900": {
			coordinationBlock: 9900,
			expectedChecklist: []WalletActionType{ActionRedemption},
		},
		"block 10800": {
			coordinationBlock: 10800,
			expectedChecklist: []WalletActionType{
				ActionRedemption,
				ActionDepositSweep,
			},
		},
		"block 11700": {
			coordinationBlock: 11700,
			expectedChecklist: []WalletActionType{ActionRedemption},
		},
		"block 12600": {
			coordinationBlock: 12600,
			expectedChecklist: []WalletActionType{
				ActionRedemption,
			},
		},
		"block 13500": {
			coordinationBlock: 13500,
			expectedChecklist: []WalletActionType{ActionRedemption},
		},
		"block 14400": {
			coordinationBlock: 14400,
			expectedChecklist: []WalletActionType{
				ActionRedemption,
				ActionDepositSweep,
			},
		},
	}

	executor := &coordinationExecutor{}

	for testName, test := range tests {
		t.Run(
			testName, func(t *testing.T) {
				window := newCoordinationWindow(test.coordinationBlock)

				// Build an arbitrary seed based on the coordination block number.
				seed := sha256.Sum256(
					big.NewInt(int64(window.coordinationBlock) + 2).Bytes(),
				)

				checklist := executor.getActionsChecklist(window.index(), seed)

				if diff := deep.Equal(
					checklist,
					test.expectedChecklist,
				); diff != nil {
					t.Errorf(
						"compare failed: %v\nactual: %s\nexpected: %s",
						diff,
						checklist,
						test.expectedChecklist,
					)
				}
			},
		)
	}
}

func TestCoordinationExecutor_ExecuteLeaderRoutine(t *testing.T) {
	// Uncompressed public key corresponding to the 20-byte public key hash:
	// aa768412ceed10bd423c025542ca90071f9fb62d.
	publicKeyHex, err := hex.DecodeString(
		"0471e30bca60f6548d7b42582a478ea37ada63b402af7b3ddd57f0c95bb6843175" +
			"aa0d2053a91a050a6797d85c38f2909cb7027f2344a01986aa2f9f8ca7a0c289",
	)
	if err != nil {
		t.Fatal(err)
	}

	// 20-byte public key hash corresponding to the public key above.
	buffer, err := hex.DecodeString("aa768412ceed10bd423c025542ca90071f9fb62d")
	if err != nil {
		t.Fatal(err)
	}
	var publicKeyHash [20]byte
	copy(publicKeyHash[:], buffer)

	coordinatedWallet := wallet{
		// Set only relevant fields.
		publicKey: unmarshalPublicKey(publicKeyHex),
	}

	// Deliberately use an unsorted list of members indexes to make sure the
	// leader routine sorts them before determining the coordination message
	// sender.
	membersIndexes := []group.MemberIndex{77, 5, 10}

	proposalGenerator := newMockCoordinationProposalGenerator(
		func(
			walletPublicKeyHash [20]byte,
			actionsChecklist []WalletActionType,
		) (
			CoordinationProposal,
			error,
		) {
			for _, action := range actionsChecklist {
				if walletPublicKeyHash == publicKeyHash && action == ActionHeartbeat {
					return &HeartbeatProposal{
						Message: [16]byte{0x01, 0x02},
					}, nil
				}
			}

			return &NoopProposal{}, nil
		},
	)

	provider := netlocal.Connect()

	broadcastChannel, err := provider.BroadcastChannelFor("test")
	if err != nil {
		t.Fatal(err)
	}

	broadcastChannel.SetUnmarshaler(func() net.TaggedUnmarshaler {
		return &coordinationMessage{}
	})

	executor := &coordinationExecutor{
		// Set only relevant fields.
		coordinatedWallet: coordinatedWallet,
		membersIndexes:    membersIndexes,
		proposalGenerator: proposalGenerator,
		broadcastChannel:  broadcastChannel,
	}

	actionsChecklist := []WalletActionType{
		ActionRedemption,
		ActionDepositSweep,
		ActionMovedFundsSweep,
		ActionMovingFunds,
		ActionHeartbeat,
	}

	ctx, cancelCtx := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelCtx()

	var message *coordinationMessage
	broadcastChannel.Recv(ctx, func(m net.Message) {
		cm, ok := m.Payload().(*coordinationMessage)
		if !ok {
			t.Fatal("unexpected message type")
		}

		// Capture the message for later assertions.
		message = cm

		// Cancel the context to proceed with the test quicker.
		cancelCtx()
	})

	proposal, err := executor.executeLeaderRoutine(ctx, 900, actionsChecklist)
	if err != nil {
		t.Fatal(err)
	}

	<-ctx.Done()

	expectedProposal := &HeartbeatProposal{
		Message: [16]byte{0x01, 0x02},
	}

	if !reflect.DeepEqual(expectedProposal, proposal) {
		t.Errorf(
			"unexpected proposal: \n"+
				"expected: %v\n"+
				"actual:   %v",
			expectedProposal,
			proposal,
		)
	}

	expectedMessage := &coordinationMessage{
		senderID:            5,
		coordinationBlock:   900,
		walletPublicKeyHash: publicKeyHash,
		proposal:            expectedProposal,
	}

	if !reflect.DeepEqual(expectedMessage, message) {
		t.Errorf(
			"unexpected message: \n"+
				"expected: %v\n"+
				"actual:   %v",
			expectedMessage,
			message,
		)
	}
}

func TestCoordinationExecutor_ExecuteFollowerRoutine(t *testing.T) {
	// Uncompressed public key corresponding to the 20-byte public key hash:
	// aa768412ceed10bd423c025542ca90071f9fb62d.
	publicKeyHex, err := hex.DecodeString(
		"0471e30bca60f6548d7b42582a478ea37ada63b402af7b3ddd57f0c95bb6843175" +
			"aa0d2053a91a050a6797d85c38f2909cb7027f2344a01986aa2f9f8ca7a0c289",
	)
	if err != nil {
		t.Fatal(err)
	}

	parseScript := func(script string) bitcoin.Script {
		parsed, err := hex.DecodeString(script)
		if err != nil {
			t.Fatal(err)
		}

		return parsed
	}

	generateOperator := func() struct {
		address chain.Address
		channel net.BroadcastChannel
	} {
		localChain := Connect()

		operatorAddress, err := localChain.operatorAddress()
		if err != nil {
			t.Fatal(err)
		}

		_, operatorPublicKey, err := localChain.OperatorKeyPair()
		if err != nil {
			t.Fatal(err)
		}

		broadcastChannel, err := netlocal.ConnectWithKey(operatorPublicKey).
			BroadcastChannelFor("test")
		if err != nil {
			t.Fatal(err)
		}

		broadcastChannel.SetUnmarshaler(func() net.TaggedUnmarshaler {
			return &coordinationMessage{}
		})
		// Register an unmarshaler for the signingDoneMessage that will
		// be uses to test the case with the wrong message type.
		broadcastChannel.SetUnmarshaler(func() net.TaggedUnmarshaler {
			return &signingDoneMessage{}
		})

		return struct {
			address chain.Address
			channel net.BroadcastChannel
		}{
			address: operatorAddress,
			channel: broadcastChannel,
		}
	}

	leader := generateOperator()
	follower1 := generateOperator()
	follower2 := generateOperator()

	coordinatedWallet := wallet{
		publicKey: unmarshalPublicKey(publicKeyHex),
		signingGroupOperators: []chain.Address{
			follower1.address,
			follower2.address,
			leader.address,
			leader.address,
			follower2.address,
			follower1.address,
			follower1.address,
			follower2.address,
			leader.address,
			leader.address,
		},
	}

	leaderID := coordinatedWallet.membersByOperator(leader.address)[0]

	localChain := Connect()

	membershipValidator := group.NewMembershipValidator(
		&testutils.MockLogger{},
		coordinatedWallet.signingGroupOperators,
		localChain.Signing(),
	)

	// Set up the executor for follower 1.
	executor := &coordinationExecutor{
		// Set only relevant fields.
		chain:               localChain,
		coordinatedWallet:   coordinatedWallet,
		membersIndexes:      coordinatedWallet.membersByOperator(follower1.address),
		operatorAddress:     follower1.address,
		broadcastChannel:    follower1.channel,
		membershipValidator: membershipValidator,
	}

	ctx, cancelCtx := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelCtx()

	go func() {
		// Give the follower routine some time to start and set up the
		// broadcast channel handler.
		time.Sleep(1 * time.Second)

		// Send message of wrong type.
		err := leader.channel.Send(ctx, &signingDoneMessage{
			senderID:      leaderID,
			message:       big.NewInt(100),
			attemptNumber: 2,
			signature: &tecdsa.Signature{
				R:          big.NewInt(200),
				S:          big.NewInt(300),
				RecoveryID: 3,
			},
			endBlock: 4500,
		})
		if err != nil {
			t.Error(err)
			return
		}

		// Send message from self.
		err = follower1.channel.Send(ctx, &coordinationMessage{
			senderID:            coordinatedWallet.membersByOperator(follower1.address)[0],
			coordinationBlock:   900,
			walletPublicKeyHash: executor.walletPublicKeyHash(),
			proposal:            &NoopProposal{},
		})
		if err != nil {
			t.Error(err)
			return
		}

		// Send message with invalid membership.
		err = leader.channel.Send(ctx, &coordinationMessage{
			// Leader operator uses senderID controlled by follower 2.
			senderID:            coordinatedWallet.membersByOperator(follower2.address)[0],
			coordinationBlock:   900,
			walletPublicKeyHash: executor.walletPublicKeyHash(),
			proposal:            &NoopProposal{},
		})
		if err != nil {
			t.Error(err)
			return
		}

		// Send message with wrong coordination block.
		err = leader.channel.Send(ctx, &coordinationMessage{
			// Proper block is 900.
			senderID:            leaderID,
			coordinationBlock:   901,
			walletPublicKeyHash: executor.walletPublicKeyHash(),
			proposal:            &NoopProposal{},
		})
		if err != nil {
			t.Error(err)
			return
		}

		// Send message with wrong wallet.
		err = leader.channel.Send(ctx, &coordinationMessage{
			senderID:            leaderID,
			coordinationBlock:   900,
			walletPublicKeyHash: [20]byte{0x01},
			proposal:            &NoopProposal{},
		})
		if err != nil {
			t.Error(err)
			return
		}

		// Send message that impersonates the leader.
		err = follower2.channel.Send(ctx, &coordinationMessage{
			senderID:            coordinatedWallet.membersByOperator(follower2.address)[0],
			coordinationBlock:   900,
			walletPublicKeyHash: executor.walletPublicKeyHash(),
			proposal:            &NoopProposal{},
		})
		if err != nil {
			t.Error(err)
			return
		}

		// Send message with not allowed action proposal.
		err = leader.channel.Send(ctx, &coordinationMessage{
			// Heartbeat proposal is not allowed for this window.
			senderID:            leaderID,
			coordinationBlock:   900,
			walletPublicKeyHash: executor.walletPublicKeyHash(),
			proposal: &HeartbeatProposal{
				Message: [16]byte{0x01, 0x02},
			},
		})
		if err != nil {
			t.Error(err)
			return
		}

		// Send a proper message.
		err = leader.channel.Send(ctx, &coordinationMessage{
			senderID:            leaderID,
			coordinationBlock:   900,
			walletPublicKeyHash: executor.walletPublicKeyHash(),
			proposal: &RedemptionProposal{
				RedeemersOutputScripts: []bitcoin.Script{
					parseScript("00148db50eb52063ea9d98b3eac91489a90f738986f6"),
					parseScript("76a9148db50eb52063ea9d98b3eac91489a90f738986f688ac"),
				},
				RedemptionTxFee: big.NewInt(10000),
			},
		})
		if err != nil {
			t.Error(err)
			return
		}
	}()

	proposal, faults, err := executor.executeFollowerRoutine(
		ctx,
		leader.address,
		900,
		[]WalletActionType{ActionRedemption, ActionNoop},
	)
	if err != nil {
		t.Fatal(err)
	}

	expectedProposal := &RedemptionProposal{
		RedeemersOutputScripts: []bitcoin.Script{
			parseScript("00148db50eb52063ea9d98b3eac91489a90f738986f6"),
			parseScript("76a9148db50eb52063ea9d98b3eac91489a90f738986f688ac"),
		},
		RedemptionTxFee: big.NewInt(10000),
	}

	if !reflect.DeepEqual(expectedProposal, proposal) {
		t.Errorf(
			"unexpected proposal: \n"+
				"expected: %v\n"+
				"actual:   %v",
			expectedProposal,
			proposal,
		)
	}

	expectedFaults := []*coordinationFault{
		{
			culprit:   follower2.address,
			faultType: FaultLeaderImpersonation,
		},
		{
			culprit:   leader.address,
			faultType: FaultLeaderMistake,
		},
	}
	if !reflect.DeepEqual(expectedFaults, faults) {
		t.Errorf(
			"unexpected faults: \n"+
				"expected: %v\n"+
				"actual:   %v",
			expectedProposal,
			proposal,
		)
	}
}

func TestCoordinationExecutor_ExecuteFollowerRoutine_WithIdleLeader(t *testing.T) {
	// Uncompressed public key corresponding to the 20-byte public key hash:
	// aa768412ceed10bd423c025542ca90071f9fb62d.
	publicKeyHex, err := hex.DecodeString(
		"0471e30bca60f6548d7b42582a478ea37ada63b402af7b3ddd57f0c95bb6843175" +
			"aa0d2053a91a050a6797d85c38f2909cb7027f2344a01986aa2f9f8ca7a0c289",
	)
	if err != nil {
		t.Fatal(err)
	}

	generateOperator := func() chain.Address {
		operatorPrivateKey, operatorPublicKey, err := operator.GenerateKeyPair(
			local_v1.DefaultCurve,
		)
		if err != nil {
			t.Fatal(err)
		}

		operatorAddress, err := ConnectWithKey(operatorPrivateKey).
			Signing().
			PublicKeyToAddress(operatorPublicKey)
		if err != nil {
			t.Fatal(err)
		}

		return operatorAddress
	}

	leader := generateOperator()
	follower1 := generateOperator()
	follower2 := generateOperator()

	coordinatedWallet := wallet{
		publicKey: unmarshalPublicKey(publicKeyHex),
		signingGroupOperators: []chain.Address{
			follower1,
			follower2,
			leader,
			leader,
			follower2,
			follower1,
			follower1,
			follower2,
			leader,
			leader,
		},
	}

	provider := netlocal.Connect()

	broadcastChannel, err := provider.BroadcastChannelFor("test-idle")
	if err != nil {
		t.Fatal(err)
	}

	executor := &coordinationExecutor{
		// Set only relevant fields.
		coordinatedWallet: coordinatedWallet,
		broadcastChannel:  broadcastChannel,
	}

	ctx, cancelCtx := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancelCtx()

	_, faults, err := executor.executeFollowerRoutine(
		ctx,
		leader,
		900,
		[]WalletActionType{ActionRedemption, ActionNoop},
	)

	expectedErr := fmt.Errorf("coordination message not received on time")
	if !reflect.DeepEqual(expectedErr, err) {
		t.Errorf(
			"unexpected error: \n"+
				"expected: %v\n"+
				"actual:   %v",
			expectedErr,
			err,
		)
	}

	expectedFaults := []*coordinationFault{
		{
			culprit:   leader,
			faultType: FaultLeaderIdleness,
		},
	}
	if !reflect.DeepEqual(expectedFaults, faults) {
		t.Errorf(
			"unexpected faults: \n"+
				"expected: %v\n"+
				"actual:   %v",
			expectedFaults,
			faults,
		)
	}
}

type mockCoordinationProposalGenerator struct {
	delegate func(
		walletPublicKeyHash [20]byte,
		actionsChecklist []WalletActionType,
	) (CoordinationProposal, error)
}

func newMockCoordinationProposalGenerator(
	delegate func(
		walletPublicKeyHash [20]byte,
		actionsChecklist []WalletActionType,
	) (CoordinationProposal, error),
) *mockCoordinationProposalGenerator {
	return &mockCoordinationProposalGenerator{
		delegate: delegate,
	}
}

func (mcpg *mockCoordinationProposalGenerator) Generate(
	request *CoordinationProposalRequest,
) (CoordinationProposal, error) {
	return mcpg.delegate(request.WalletPublicKeyHash, request.ActionsChecklist)
}
