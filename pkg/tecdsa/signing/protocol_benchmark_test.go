package signing

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/bnb-chain/tss-lib/ecdsa/keygen"
	"github.com/keep-network/keep-core/pkg/crypto/ephemeral"
	"github.com/keep-network/keep-core/pkg/internal/testutils"
	"github.com/keep-network/keep-core/pkg/protocol/group"
	"github.com/keep-network/keep-core/pkg/tecdsa"
	"io/ioutil"
	"math/big"
	"path/filepath"
	"runtime"
	"testing"
	"time"
)

const (
	benchmarkGroupSize                         = 51
	benchmarkDishonestThreshold                = 50
	benchmarkTestDataDirFormat                 = "%s/benchmarkdata"
	benchmarkPrivateKeyShareTestDataFileFormat = "private_key_share_data_%d.json"
)

func BenchmarkTssRoundOne(b *testing.B) {
	for i := 0; i < b.N; i++ {
		members, err := initializeTssRoundOneMembersGroupForBenchmark(
			benchmarkDishonestThreshold,
			benchmarkGroupSize,
		)
		if err != nil {
			b.Fatal(err)
		}

		ctx, cancelCtx := context.WithCancel(
			context.Background(),
		)

		b.ResetTimer()

		_, err = members[0].tssRoundOne(ctx)
		if err != nil {
			b.Fatal(err)
		}

		cancelCtx()
	}
}

func BenchmarkTssRoundTwo(b *testing.B) {
	for i := 0; i < b.N; i++ {
		members, tssRoundOneMessages, err := initializeTssRoundTwoMembersGroupForBenchmark(
			benchmarkDishonestThreshold,
			benchmarkGroupSize,
		)
		if err != nil {
			b.Fatal(err)
		}

		member := members[0]

		var receivedTssRoundOneMessages []*tssRoundOneMessage
		for _, tssRoundOneMessage := range tssRoundOneMessages {
			if tssRoundOneMessage.senderID != member.id {
				receivedTssRoundOneMessages = append(
					receivedTssRoundOneMessages,
					tssRoundOneMessage,
				)
			}
		}

		ctx, cancelCtx := context.WithCancel(context.Background())

		b.ResetTimer()

		_, err = member.tssRoundTwo(
			ctx,
			receivedTssRoundOneMessages,
		)
		if err != nil {
			cancelCtx()
			b.Fatal(err)
		}

		cancelCtx()
	}
}

func BenchmarkTssRoundThree(b *testing.B) {
	for i := 0; i < b.N; i++ {
		members, tssRoundTwoMessages, err := initializeTssRoundThreeMembersGroupForBenchmark(
			benchmarkDishonestThreshold,
			benchmarkGroupSize,
		)
		if err != nil {
			b.Fatal(err)
		}

		member := members[0]

		var receivedTssRoundTwoMessages []*tssRoundTwoMessage
		for _, tssRoundTwoMessage := range tssRoundTwoMessages {
			if tssRoundTwoMessage.senderID != member.id {
				receivedTssRoundTwoMessages = append(
					receivedTssRoundTwoMessages,
					tssRoundTwoMessage,
				)
			}
		}

		ctx, cancelCtx := context.WithCancel(context.Background())

		b.ResetTimer()

		_, err = member.tssRoundThree(
			ctx,
			receivedTssRoundTwoMessages,
		)
		if err != nil {
			cancelCtx()
			b.Fatal(err)
		}

		cancelCtx()
	}
}

func BenchmarkTssRoundFour(b *testing.B) {
	for i := 0; i < b.N; i++ {
		members, tssRoundThreeMessages, err := initializeTssRoundFourMembersGroupForBenchmark(
			benchmarkDishonestThreshold,
			benchmarkGroupSize,
		)
		if err != nil {
			b.Fatal(err)
		}

		member := members[0]

		var receivedTssRoundThreeMessages []*tssRoundThreeMessage
		for _, tssRoundThreeMessage := range tssRoundThreeMessages {
			if tssRoundThreeMessage.senderID != member.id {
				receivedTssRoundThreeMessages = append(
					receivedTssRoundThreeMessages,
					tssRoundThreeMessage,
				)
			}
		}

		ctx, cancelCtx := context.WithCancel(context.Background())

		b.ResetTimer()

		_, err = member.tssRoundFour(
			ctx,
			receivedTssRoundThreeMessages,
		)
		if err != nil {
			cancelCtx()
			b.Fatal(err)
		}

		cancelCtx()
	}
}

func BenchmarkTssRoundFive(b *testing.B) {
	for i := 0; i < b.N; i++ {
		members, tssRoundFourMessages, err := initializeTssRoundFiveMembersGroupForBenchmark(
			benchmarkDishonestThreshold,
			benchmarkGroupSize,
		)
		if err != nil {
			b.Fatal(err)
		}

		member := members[0]

		var receivedTssRoundFourMessages []*tssRoundFourMessage
		for _, tssRoundFourMessage := range tssRoundFourMessages {
			if tssRoundFourMessage.senderID != member.id {
				receivedTssRoundFourMessages = append(
					receivedTssRoundFourMessages,
					tssRoundFourMessage,
				)
			}
		}

		ctx, cancelCtx := context.WithCancel(context.Background())

		b.ResetTimer()

		_, err = member.tssRoundFive(
			ctx,
			receivedTssRoundFourMessages,
		)
		if err != nil {
			cancelCtx()
			b.Fatal(err)
		}

		cancelCtx()
	}
}

func BenchmarkTssRoundSix(b *testing.B) {
	for i := 0; i < b.N; i++ {
		members, tssRoundFiveMessages, err := initializeTssRoundSixMembersGroupForBenchmark(
			benchmarkDishonestThreshold,
			benchmarkGroupSize,
		)
		if err != nil {
			b.Fatal(err)
		}

		member := members[0]

		var receivedTssRoundFiveMessages []*tssRoundFiveMessage
		for _, tssRoundFiveMessage := range tssRoundFiveMessages {
			if tssRoundFiveMessage.senderID != member.id {
				receivedTssRoundFiveMessages = append(
					receivedTssRoundFiveMessages,
					tssRoundFiveMessage,
				)
			}
		}

		ctx, cancelCtx := context.WithCancel(context.Background())

		b.ResetTimer()

		_, err = member.tssRoundSix(
			ctx,
			receivedTssRoundFiveMessages,
		)
		if err != nil {
			cancelCtx()
			b.Fatal(err)
		}

		cancelCtx()
	}
}

func BenchmarkTssRoundSeven(b *testing.B) {
	for i := 0; i < b.N; i++ {
		members, tssRoundSixMessages, err := initializeTssRoundSevenMembersGroupForBenchmark(
			benchmarkDishonestThreshold,
			benchmarkGroupSize,
		)
		if err != nil {
			b.Fatal(err)
		}

		member := members[0]

		var receivedTssRoundSixMessages []*tssRoundSixMessage
		for _, tssRoundSixMessage := range tssRoundSixMessages {
			if tssRoundSixMessage.senderID != member.id {
				receivedTssRoundSixMessages = append(
					receivedTssRoundSixMessages,
					tssRoundSixMessage,
				)
			}
		}

		ctx, cancelCtx := context.WithCancel(context.Background())

		b.ResetTimer()

		_, err = member.tssRoundSeven(
			ctx,
			receivedTssRoundSixMessages,
		)
		if err != nil {
			cancelCtx()
			b.Fatal(err)
		}

		cancelCtx()
	}
}

func BenchmarkTssRoundEight(b *testing.B) {
	for i := 0; i < b.N; i++ {
		members, tssRoundSevenMessages, err := initializeTssRoundEightMembersGroupForBenchmark(
			benchmarkDishonestThreshold,
			benchmarkGroupSize,
		)
		if err != nil {
			b.Fatal(err)
		}

		member := members[0]

		var receivedTssRoundSevenMessages []*tssRoundSevenMessage
		for _, tssRoundSevenMessage := range tssRoundSevenMessages {
			if tssRoundSevenMessage.senderID != member.id {
				receivedTssRoundSevenMessages = append(
					receivedTssRoundSevenMessages,
					tssRoundSevenMessage,
				)
			}
		}

		ctx, cancelCtx := context.WithCancel(context.Background())

		b.ResetTimer()

		_, err = member.tssRoundEight(
			ctx,
			receivedTssRoundSevenMessages,
		)
		if err != nil {
			cancelCtx()
			b.Fatal(err)
		}

		cancelCtx()
	}
}

func BenchmarkTssRoundNine(b *testing.B) {
	for i := 0; i < b.N; i++ {
		members, tssRoundEightMessages, err := initializeTssRoundNineMembersGroupForBenchmark(
			benchmarkDishonestThreshold,
			benchmarkGroupSize,
		)
		if err != nil {
			b.Fatal(err)
		}

		member := members[0]

		var receivedTssRoundEightMessages []*tssRoundEightMessage
		for _, tssRoundEightMessage := range tssRoundEightMessages {
			if tssRoundEightMessage.senderID != member.id {
				receivedTssRoundEightMessages = append(
					receivedTssRoundEightMessages,
					tssRoundEightMessage,
				)
			}
		}

		ctx, cancelCtx := context.WithCancel(context.Background())

		b.ResetTimer()

		_, err = member.tssRoundNine(
			ctx,
			receivedTssRoundEightMessages,
		)
		if err != nil {
			cancelCtx()
			b.Fatal(err)
		}

		cancelCtx()
	}
}

func BenchmarkTssFinalize(b *testing.B) {
	for i := 0; i < b.N; i++ {
		members, tssRoundNineMessages, err := initializeFinalizingMembersGroupForBenchmark(
			benchmarkDishonestThreshold,
			benchmarkGroupSize,
		)
		if err != nil {
			b.Fatal(err)
		}

		member := members[0]

		var receivedTssRoundNineMessages []*tssRoundNineMessage
		for _, tssRoundNineMessage := range tssRoundNineMessages {
			if tssRoundNineMessage.senderID != member.id {
				receivedTssRoundNineMessages = append(
					receivedTssRoundNineMessages,
					tssRoundNineMessage,
				)
			}
		}

		ctx, cancelCtx := context.WithCancel(context.Background())

		b.ResetTimer()

		err = member.tssFinalize(
			ctx,
			receivedTssRoundNineMessages,
		)
		if err != nil {
			cancelCtx()
			b.Fatal(err)
		}

		cancelCtx()
	}
}

func loadPrivateKeyShareTestFixturesForBenchmark(count int) (
	[]keygen.LocalPartySaveData,
	error,
) {
	makeTestFixtureFilePath := func(partyIndex int) string {
		_, callerFileName, _, _ := runtime.Caller(0)
		srcDirName := filepath.Dir(callerFileName)
		fixtureDirName := fmt.Sprintf(benchmarkTestDataDirFormat, srcDirName)
		return fmt.Sprintf(
			"%s/"+benchmarkPrivateKeyShareTestDataFileFormat,
			fixtureDirName,
			partyIndex,
		)
	}

	shares := make([]keygen.LocalPartySaveData, 0, count)

	for j := 0; j < count; j++ {
		fixtureFilePath := makeTestFixtureFilePath(j)

		// #nosec G304 (file path provided as taint input)
		// This line is used to read a test fixture file.
		// There is no user input.
		bz, err := ioutil.ReadFile(fixtureFilePath)
		if err != nil {
			return nil, fmt.Errorf(
				"could not open the test fixture for party [%d] "+
					"in the expected location [%s]: [%w]",
				j,
				fixtureFilePath,
				err,
			)
		}
		var share keygen.LocalPartySaveData
		if err = json.Unmarshal(bz, &share); err != nil {
			return nil, fmt.Errorf(
				"could not unmarshal fixture data for party [%d] "+
					"located at [%s]: [%w]",
				j,
				fixtureFilePath,
				err,
			)
		}
		shares = append(shares, share)
	}
	return shares, nil
}

func initializeEphemeralKeyPairGeneratingMembersGroupForBenchmark(
	dishonestThreshold int,
	groupSize int,
) ([]*ephemeralKeyPairGeneratingMember, error) {
	signingGroup := group.NewGroup(dishonestThreshold, groupSize)

	testData, err := loadPrivateKeyShareTestFixturesForBenchmark(groupSize)
	if err != nil {
		return nil, fmt.Errorf("failed to load test data: [%v]", err)
	}

	var members []*ephemeralKeyPairGeneratingMember
	for i := 1; i <= groupSize; i++ {
		id := group.MemberIndex(i)

		members = append(members, &ephemeralKeyPairGeneratingMember{
			member: &member{
				logger:            &testutils.MockLogger{},
				id:                id,
				group:             signingGroup,
				sessionID:         sessionID,
				message:           big.NewInt(100),
				privateKeyShare:   tecdsa.NewPrivateKeyShare(testData[i-1]),
				identityConverter: &identityConverter{keys: testData[i-1].Ks},
			},
			ephemeralKeyPairs: make(map[group.MemberIndex]*ephemeral.KeyPair),
		})
	}

	return members, nil
}
func initializeSymmetricKeyGeneratingMembersGroupForBenchmark(
	dishonestThreshold int,
	groupSize int,
) (
	[]*symmetricKeyGeneratingMember,
	[]*ephemeralPublicKeyMessage,
	error,
) {
	var symmetricKeyGeneratingMembers []*symmetricKeyGeneratingMember
	var ephemeralPublicKeyMessages []*ephemeralPublicKeyMessage

	ephemeralKeyPairGeneratingMembers, err :=
		initializeEphemeralKeyPairGeneratingMembersGroupForBenchmark(
			dishonestThreshold,
			groupSize,
		)
	if err != nil {
		return nil, nil, fmt.Errorf(
			"cannot generate ephemeral key pair generating "+
				"members group: [%v]",
			err,
		)
	}

	for _, member := range ephemeralKeyPairGeneratingMembers {
		message, err := member.generateEphemeralKeyPair()
		if err != nil {
			return nil, nil, fmt.Errorf(
				"cannot generate ephemeral key pair for member [%v]: [%v]",
				member.id,
				err,
			)
		}

		symmetricKeyGeneratingMembers = append(
			symmetricKeyGeneratingMembers,
			member.initializeSymmetricKeyGeneration(),
		)
		ephemeralPublicKeyMessages = append(ephemeralPublicKeyMessages, message)
	}

	return symmetricKeyGeneratingMembers, ephemeralPublicKeyMessages, nil
}

func initializeTssRoundOneMembersGroupForBenchmark(
	dishonestThreshold int,
	groupSize int,
) (
	[]*tssRoundOneMember,
	error,
) {
	var tssRoundOneMembers []*tssRoundOneMember

	symmetricKeyGeneratingMembers, ephemeralPublicKeyMessages, err :=
		initializeSymmetricKeyGeneratingMembersGroupForBenchmark(
			dishonestThreshold,
			groupSize,
		)
	if err != nil {
		return nil, fmt.Errorf(
			"cannot generate symmetric key generating members group: [%v]",
			err,
		)
	}

	for _, member := range symmetricKeyGeneratingMembers {
		var receivedMessages []*ephemeralPublicKeyMessage
		for _, message := range ephemeralPublicKeyMessages {
			if message.senderID != member.id {
				receivedMessages = append(receivedMessages, message)
			}
		}

		err := member.generateSymmetricKeys(receivedMessages)
		if err != nil {
			return nil, fmt.Errorf(
				"cannot generate symmetric keys for member [%v]: [%v]",
				member.id,
				err,
			)
		}

		tssRoundOneMembers = append(
			tssRoundOneMembers,
			member.initializeTssRoundOne(),
		)
	}

	return tssRoundOneMembers, nil
}

func initializeTssRoundTwoMembersGroupForBenchmark(
	dishonestThreshold int,
	groupSize int,
) (
	[]*tssRoundTwoMember,
	[]*tssRoundOneMessage,
	error,
) {
	var tssRoundTwoMembers []*tssRoundTwoMember
	var tssRoundOneMessages []*tssRoundOneMessage

	tssRoundOneMembers, err :=
		initializeTssRoundOneMembersGroupForBenchmark(
			dishonestThreshold,
			groupSize,
		)
	if err != nil {
		return nil, nil, fmt.Errorf(
			"cannot generate TSS round one members group: [%v]",
			err,
		)
	}

	for _, member := range tssRoundOneMembers {
		ctx, cancelCtx := context.WithTimeout(
			context.Background(),
			10*time.Second,
		)

		message, err := member.tssRoundOne(ctx)
		if err != nil {
			cancelCtx()
			return nil, nil, fmt.Errorf(
				"cannot do TSS round one for member [%v]: [%v]",
				member.id,
				err,
			)
		}

		tssRoundTwoMembers = append(
			tssRoundTwoMembers,
			member.initializeTssRoundTwo(),
		)
		tssRoundOneMessages = append(tssRoundOneMessages, message)

		cancelCtx()
	}

	return tssRoundTwoMembers, tssRoundOneMessages, nil
}

func initializeTssRoundThreeMembersGroupForBenchmark(
	dishonestThreshold int,
	groupSize int,
) (
	[]*tssRoundThreeMember,
	[]*tssRoundTwoMessage,
	error,
) {
	var tssRoundThreeMembers []*tssRoundThreeMember
	var tssRoundTwoMessages []*tssRoundTwoMessage

	tssRoundTwoMembers, tssRoundOneMessages, err :=
		initializeTssRoundTwoMembersGroupForBenchmark(
			dishonestThreshold,
			groupSize,
		)
	if err != nil {
		return nil, nil, fmt.Errorf(
			"cannot generate TSS round two members group: [%v]",
			err,
		)
	}

	for _, member := range tssRoundTwoMembers {
		var receivedTssRoundOneMessages []*tssRoundOneMessage
		for _, tssRoundOneMessage := range tssRoundOneMessages {
			if tssRoundOneMessage.senderID != member.id {
				receivedTssRoundOneMessages = append(
					receivedTssRoundOneMessages,
					tssRoundOneMessage,
				)
			}
		}

		ctx, cancelCtx := context.WithTimeout(
			context.Background(),
			10*time.Second,
		)

		tssRoundTwoMessage, err := member.tssRoundTwo(
			ctx,
			receivedTssRoundOneMessages,
		)
		if err != nil {
			cancelCtx()
			return nil, nil, fmt.Errorf(
				"cannot do TSS round two for member [%v]: [%v]",
				member.id,
				err,
			)
		}

		tssRoundThreeMembers = append(
			tssRoundThreeMembers,
			member.initializeTssRoundThree(),
		)
		tssRoundTwoMessages = append(tssRoundTwoMessages, tssRoundTwoMessage)

		cancelCtx()
	}

	return tssRoundThreeMembers, tssRoundTwoMessages, nil
}

func initializeTssRoundFourMembersGroupForBenchmark(
	dishonestThreshold int,
	groupSize int,
) (
	[]*tssRoundFourMember,
	[]*tssRoundThreeMessage,
	error,
) {
	var tssRoundFourMembers []*tssRoundFourMember
	var tssRoundThreeMessages []*tssRoundThreeMessage

	tssRoundThreeMembers, tssRoundTwoMessages, err :=
		initializeTssRoundThreeMembersGroupForBenchmark(
			dishonestThreshold,
			groupSize,
		)
	if err != nil {
		return nil, nil, fmt.Errorf(
			"cannot generate TSS round three members group: [%v]",
			err,
		)
	}

	for _, member := range tssRoundThreeMembers {
		var receivedTssRoundTwoMessages []*tssRoundTwoMessage
		for _, tssRoundTwoMessage := range tssRoundTwoMessages {
			if tssRoundTwoMessage.senderID != member.id {
				receivedTssRoundTwoMessages = append(
					receivedTssRoundTwoMessages,
					tssRoundTwoMessage,
				)
			}
		}

		ctx, cancelCtx := context.WithTimeout(
			context.Background(),
			10*time.Second,
		)

		tssRoundThreeMessage, err := member.tssRoundThree(
			ctx,
			receivedTssRoundTwoMessages,
		)
		if err != nil {
			cancelCtx()
			return nil, nil, fmt.Errorf(
				"cannot do TSS round three for member [%v]: [%v]",
				member.id,
				err,
			)
		}

		tssRoundFourMembers = append(
			tssRoundFourMembers,
			member.initializeTssRoundFour(),
		)
		tssRoundThreeMessages = append(tssRoundThreeMessages, tssRoundThreeMessage)

		cancelCtx()
	}

	return tssRoundFourMembers, tssRoundThreeMessages, nil
}

func initializeTssRoundFiveMembersGroupForBenchmark(
	dishonestThreshold int,
	groupSize int,
) (
	[]*tssRoundFiveMember,
	[]*tssRoundFourMessage,
	error,
) {
	var tssRoundFiveMembers []*tssRoundFiveMember
	var tssRoundFourMessages []*tssRoundFourMessage

	tssRoundFourMembers, tssRoundThreeMessages, err :=
		initializeTssRoundFourMembersGroupForBenchmark(
			dishonestThreshold,
			groupSize,
		)
	if err != nil {
		return nil, nil, fmt.Errorf(
			"cannot generate TSS round four members group: [%v]",
			err,
		)
	}

	for _, member := range tssRoundFourMembers {
		var receivedTssRoundThreeMessages []*tssRoundThreeMessage
		for _, tssRoundThreeMessage := range tssRoundThreeMessages {
			if tssRoundThreeMessage.senderID != member.id {
				receivedTssRoundThreeMessages = append(
					receivedTssRoundThreeMessages,
					tssRoundThreeMessage,
				)
			}
		}

		ctx, cancelCtx := context.WithTimeout(
			context.Background(),
			10*time.Second,
		)

		tssRoundFourMessage, err := member.tssRoundFour(
			ctx,
			receivedTssRoundThreeMessages,
		)
		if err != nil {
			cancelCtx()
			return nil, nil, fmt.Errorf(
				"cannot do TSS round four for member [%v]: [%v]",
				member.id,
				err,
			)
		}

		tssRoundFiveMembers = append(
			tssRoundFiveMembers,
			member.initializeTssRoundFive(),
		)
		tssRoundFourMessages = append(tssRoundFourMessages, tssRoundFourMessage)

		cancelCtx()
	}

	return tssRoundFiveMembers, tssRoundFourMessages, nil
}

func initializeTssRoundSixMembersGroupForBenchmark(
	dishonestThreshold int,
	groupSize int,
) (
	[]*tssRoundSixMember,
	[]*tssRoundFiveMessage,
	error,
) {
	var tssRoundSixMembers []*tssRoundSixMember
	var tssRoundFiveMessages []*tssRoundFiveMessage

	tssRoundFiveMembers, tssRoundFourMessages, err :=
		initializeTssRoundFiveMembersGroupForBenchmark(
			dishonestThreshold,
			groupSize,
		)
	if err != nil {
		return nil, nil, fmt.Errorf(
			"cannot generate TSS round five members group: [%v]",
			err,
		)
	}

	for _, member := range tssRoundFiveMembers {
		var receivedTssRoundFourMessages []*tssRoundFourMessage
		for _, tssRoundFourMessage := range tssRoundFourMessages {
			if tssRoundFourMessage.senderID != member.id {
				receivedTssRoundFourMessages = append(
					receivedTssRoundFourMessages,
					tssRoundFourMessage,
				)
			}
		}

		ctx, cancelCtx := context.WithTimeout(
			context.Background(),
			10*time.Second,
		)

		tssRoundFiveMessage, err := member.tssRoundFive(
			ctx,
			receivedTssRoundFourMessages,
		)
		if err != nil {
			cancelCtx()
			return nil, nil, fmt.Errorf(
				"cannot do TSS round five for member [%v]: [%v]",
				member.id,
				err,
			)
		}

		tssRoundSixMembers = append(
			tssRoundSixMembers,
			member.initializeTssRoundSix(),
		)
		tssRoundFiveMessages = append(tssRoundFiveMessages, tssRoundFiveMessage)

		cancelCtx()
	}

	return tssRoundSixMembers, tssRoundFiveMessages, nil
}

func initializeTssRoundSevenMembersGroupForBenchmark(
	dishonestThreshold int,
	groupSize int,
) (
	[]*tssRoundSevenMember,
	[]*tssRoundSixMessage,
	error,
) {
	var tssRoundSevenMembers []*tssRoundSevenMember
	var tssRoundSixMessages []*tssRoundSixMessage

	tssRoundSixMembers, tssRoundFiveMessages, err :=
		initializeTssRoundSixMembersGroupForBenchmark(
			dishonestThreshold,
			groupSize,
		)
	if err != nil {
		return nil, nil, fmt.Errorf(
			"cannot generate TSS round six members group: [%v]",
			err,
		)
	}

	for _, member := range tssRoundSixMembers {
		var receivedTssRoundFiveMessages []*tssRoundFiveMessage
		for _, tssRoundFiveMessage := range tssRoundFiveMessages {
			if tssRoundFiveMessage.senderID != member.id {
				receivedTssRoundFiveMessages = append(
					receivedTssRoundFiveMessages,
					tssRoundFiveMessage,
				)
			}
		}

		ctx, cancelCtx := context.WithTimeout(
			context.Background(),
			10*time.Second,
		)

		tssRoundSixMessage, err := member.tssRoundSix(
			ctx,
			receivedTssRoundFiveMessages,
		)
		if err != nil {
			cancelCtx()
			return nil, nil, fmt.Errorf(
				"cannot do TSS round six for member [%v]: [%v]",
				member.id,
				err,
			)
		}

		tssRoundSevenMembers = append(
			tssRoundSevenMembers,
			member.initializeTssRoundSeven(),
		)
		tssRoundSixMessages = append(tssRoundSixMessages, tssRoundSixMessage)

		cancelCtx()
	}

	return tssRoundSevenMembers, tssRoundSixMessages, nil
}

func initializeTssRoundEightMembersGroupForBenchmark(
	dishonestThreshold int,
	groupSize int,
) (
	[]*tssRoundEightMember,
	[]*tssRoundSevenMessage,
	error,
) {
	var tssRoundEightMembers []*tssRoundEightMember
	var tssRoundSevenMessages []*tssRoundSevenMessage

	tssRoundSevenMembers, tssRoundSixMessages, err :=
		initializeTssRoundSevenMembersGroupForBenchmark(
			dishonestThreshold,
			groupSize,
		)
	if err != nil {
		return nil, nil, fmt.Errorf(
			"cannot generate TSS round seven members group: [%v]",
			err,
		)
	}

	for _, member := range tssRoundSevenMembers {
		var receivedTssRoundSixMessages []*tssRoundSixMessage
		for _, tssRoundSixMessage := range tssRoundSixMessages {
			if tssRoundSixMessage.senderID != member.id {
				receivedTssRoundSixMessages = append(
					receivedTssRoundSixMessages,
					tssRoundSixMessage,
				)
			}
		}

		ctx, cancelCtx := context.WithTimeout(
			context.Background(),
			10*time.Second,
		)

		tssRoundSevenMessage, err := member.tssRoundSeven(
			ctx,
			receivedTssRoundSixMessages,
		)
		if err != nil {
			cancelCtx()
			return nil, nil, fmt.Errorf(
				"cannot do TSS round seven for member [%v]: [%v]",
				member.id,
				err,
			)
		}

		tssRoundEightMembers = append(
			tssRoundEightMembers,
			member.initializeTssRoundEight(),
		)
		tssRoundSevenMessages = append(tssRoundSevenMessages, tssRoundSevenMessage)

		cancelCtx()
	}

	return tssRoundEightMembers, tssRoundSevenMessages, nil
}

func initializeTssRoundNineMembersGroupForBenchmark(
	dishonestThreshold int,
	groupSize int,
) (
	[]*tssRoundNineMember,
	[]*tssRoundEightMessage,
	error,
) {
	var tssRoundNineMembers []*tssRoundNineMember
	var tssRoundEightMessages []*tssRoundEightMessage

	tssRoundEightMembers, tssRoundSevenMessages, err :=
		initializeTssRoundEightMembersGroupForBenchmark(
			dishonestThreshold,
			groupSize,
		)
	if err != nil {
		return nil, nil, fmt.Errorf(
			"cannot generate TSS round eight members group: [%v]",
			err,
		)
	}

	for _, member := range tssRoundEightMembers {
		var receivedTssRoundSevenMessages []*tssRoundSevenMessage
		for _, tssRoundSevenMessage := range tssRoundSevenMessages {
			if tssRoundSevenMessage.senderID != member.id {
				receivedTssRoundSevenMessages = append(
					receivedTssRoundSevenMessages,
					tssRoundSevenMessage,
				)
			}
		}

		ctx, cancelCtx := context.WithTimeout(
			context.Background(),
			10*time.Second,
		)

		tssRoundEightMessage, err := member.tssRoundEight(
			ctx,
			receivedTssRoundSevenMessages,
		)
		if err != nil {
			cancelCtx()
			return nil, nil, fmt.Errorf(
				"cannot do TSS round eight for member [%v]: [%v]",
				member.id,
				err,
			)
		}

		tssRoundNineMembers = append(
			tssRoundNineMembers,
			member.initializeTssRoundNine(),
		)
		tssRoundEightMessages = append(tssRoundEightMessages, tssRoundEightMessage)

		cancelCtx()
	}

	return tssRoundNineMembers, tssRoundEightMessages, nil
}

func initializeFinalizingMembersGroupForBenchmark(
	dishonestThreshold int,
	groupSize int,
) (
	[]*finalizingMember,
	[]*tssRoundNineMessage,
	error,
) {
	var finalizingMembers []*finalizingMember
	var tssRoundNineMessages []*tssRoundNineMessage

	tssRoundNineMembers, tssRoundEightMessages, err :=
		initializeTssRoundNineMembersGroupForBenchmark(
			dishonestThreshold,
			groupSize,
		)
	if err != nil {
		return nil, nil, fmt.Errorf(
			"cannot generate TSS round nine members group: [%v]",
			err,
		)
	}

	for _, member := range tssRoundNineMembers {
		var receivedTssRoundEightMessages []*tssRoundEightMessage
		for _, tssRoundEightMessage := range tssRoundEightMessages {
			if tssRoundEightMessage.senderID != member.id {
				receivedTssRoundEightMessages = append(
					receivedTssRoundEightMessages,
					tssRoundEightMessage,
				)
			}
		}

		ctx, cancelCtx := context.WithTimeout(
			context.Background(),
			10*time.Second,
		)

		tssRoundNineMessage, err := member.tssRoundNine(
			ctx,
			receivedTssRoundEightMessages,
		)
		if err != nil {
			cancelCtx()
			return nil, nil, fmt.Errorf(
				"cannot do TSS round nine for member [%v]: [%v]",
				member.id,
				err,
			)
		}

		finalizingMembers = append(
			finalizingMembers,
			member.initializeFinalization(),
		)
		tssRoundNineMessages = append(tssRoundNineMessages, tssRoundNineMessage)

		cancelCtx()
	}

	return finalizingMembers, tssRoundNineMessages, nil
}
