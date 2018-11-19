package gjkr

import (
	"fmt"
	"math/big"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	"github.com/keep-network/keep-core/pkg/beacon/relay/event"
	"github.com/keep-network/keep-core/pkg/beacon/relay/result"
	"github.com/keep-network/keep-core/pkg/chain/local"
	"github.com/pschlump/MiscLib"
	"github.com/pschlump/godebug"
)

func TestChallengeResult(t *testing.T) {

	type op struct {
		cmd        string
		opt0       int
		opt1       int64
		resultHash []byte
	}

	var tests = map[string]struct {
		runTest       bool
		thresholdSize int
		groupSize     int
		resultValue   int64
		cmds          []op
	}{
		"001 success with 8 of 15 group size": {
			runTest:       false,
			thresholdSize: 8,
			groupSize:     15,
			resultValue:   2,
			cmds: []op{
				{cmd: "init"},
				{cmd: "publish", opt1: /*groupPublicKey*/ 1,
					resultHash: []byte{129, 169, 11, 49, 6, 91, 33, 192, 47, 59, 124, 52, 156, 242, 148,
						51, 37, 195, 18, 222, 25, 74, 74, 245, 109, 140, 243, 253, 185, 54, 214, 108},
				},
				{cmd: "challenge", opt0: /*correctResult*/ 0, opt1: /*groupPublicKey*/ 1,
					resultHash: []byte{129, 169, 11, 49, 6, 91, 33, 192, 47, 59, 124, 52, 156, 242, 148,
						51, 37, 195, 18, 222, 25, 74, 74, 245, 109, 140, 243, 253, 185, 54, 214, 108},
				},
				{cmd: "send", opt1: /*groupPublicKey*/ 1},
				{cmd: "block"},
				{cmd: "vote", opt0: /*result-subscript*/ 0, opt1: /*groupPublicKey*/ 1,
					resultHash: []byte{129, 169, 11, 49, 6, 91, 33, 192, 47, 59, 124, 52, 156, 242, 148,
						51, 37, 195, 18, 222, 25, 74, 74, 245, 109, 140, 243, 253, 185, 54, 214, 108},
				},
				{cmd: "block"},
				{cmd: "vote", opt0: /*result-subscript*/ 0, opt1: /*groupPublicKey*/ 1,
					resultHash: []byte{129, 169, 11, 49, 6, 91, 33, 192, 47, 59, 124, 52, 156, 242, 148,
						51, 37, 195, 18, 222, 25, 74, 74, 245, 109, 140, 243, 253, 185, 54, 214, 108},
				},
				{cmd: "block"},
				{cmd: "vote", opt0: /*result-subscript*/ 0, opt1: /*groupPublicKey*/ 1,
					resultHash: []byte{129, 169, 11, 49, 6, 91, 33, 192, 47, 59, 124, 52, 156, 242, 148,
						51, 37, 195, 18, 222, 25, 74, 74, 245, 109, 140, 243, 253, 185, 54, 214, 108},
				},
				{cmd: "block"},
				{cmd: "dump-data"},
				{cmd: "validate-data", opt0: 1},
			},
		},
		"002 success with 8 of 15 group size": {
			runTest:       false,
			thresholdSize: 8,
			groupSize:     15,
			resultValue:   2,
			cmds: []op{
				{cmd: "init"},
				{cmd: "publish", opt1: /*groupPublicKey*/ 1,
					resultHash: []byte{129, 169, 11, 49, 6, 91, 33, 192, 47, 59, 124, 52, 156, 242, 148,
						51, 37, 195, 18, 222, 25, 74, 74, 245, 109, 140, 243, 253, 185, 54, 214, 108},
				},
				//	{cmd: "publish", opt1: /*groupPublicKey*/ 1,
				//		resultHash: []byte{229, 169, 11, 49, 6, 91, 33, 192, 47, 59, 124, 52, 156, 242, 148,
				//			51, 37, 195, 18, 222, 25, 74, 74, 245, 109, 140, 243, 253, 185, 54, 214, 108},
				//	},
				{cmd: "challenge", opt0: /*correctResult*/ 0, opt1: /*groupPublicKey*/ 1,
					resultHash: []byte{129, 169, 11, 49, 6, 91, 33, 192, 47, 59, 124, 52, 156, 242, 148,
						51, 37, 195, 18, 222, 25, 74, 74, 245, 109, 140, 243, 253, 185, 54, 214, 108},
				},
				//				{cmd: "challenge", opt0: /*correctResult*/ 0, opt1: /*groupPublicKey*/ 1,
				//					resultHash: []byte{229, 169, 11, 49, 6, 91, 33, 192, 47, 59, 124, 52, 156, 242, 148,
				//						51, 37, 195, 18, 222, 25, 74, 74, 245, 109, 140, 243, 253, 185, 54, 214, 108},
				//				},
				{cmd: "send", opt1: /*groupPublicKey*/ 1},
				{cmd: "block"},
				{cmd: "vote", opt0: /*result-subscript*/ 0, opt1: /*groupPublicKey*/ 1,
					resultHash: []byte{129, 169, 11, 49, 6, 91, 33, 192, 47, 59, 124, 52, 156, 242, 148,
						51, 37, 195, 18, 222, 25, 74, 74, 245, 109, 140, 243, 253, 185, 54, 214, 108},
				},
				{cmd: "block"},
				{cmd: "vote", opt0: /*result-subscript*/ 0, opt1: /*groupPublicKey*/ 1,
					resultHash: []byte{129, 169, 11, 49, 6, 91, 33, 192, 47, 59, 124, 52, 156, 242, 148,
						51, 37, 195, 18, 222, 25, 74, 74, 245, 109, 140, 243, 253, 185, 54, 214, 108},
				},
				{cmd: "block"},
				{cmd: "vote", opt0: /*result-subscript*/ 0, opt1: /*groupPublicKey*/ 1,
					resultHash: []byte{129, 169, 11, 49, 6, 91, 33, 192, 47, 59, 124, 52, 156, 242, 148,
						51, 37, 195, 18, 222, 25, 74, 74, 245, 109, 140, 243, 253, 185, 54, 214, 108},
				},
				{cmd: "block"},
				{cmd: "dump-data"},
				{cmd: "validate-data", opt0: 1},
			},
		},
	}

	for testName, test := range tests {

		one_time := true

		x := t.Run(testName[3:], func(t *testing.T) {

			if test.runTest {

				fmt.Printf("Test: %s AT: %s\n", testName, godebug.LF()) // DEL-XXX

				var members []*PublishingMember
				var member *PublishingMember
				var chainRelay chain.Interface
				var err error
				var vc *DKG

				// Implemente a Von Neumann fetch execute cycle based tester.  `pc` is the
				// subscript of the instruciton.   `cmd` is the instruction. `cmd.cmd` is
				// the opcode for the instruction.
				for pc, cmd := range test.cmds {
					switch cmd.cmd {

					// Step 0 - initialize
					case "init":
						members, err = initializePublishingMembersGroup2(test.thresholdSize, test.groupSize)
						if err != nil {
							t.Errorf("%s", err)
						}
						member = members[0] // WHY? Fix This FIXME xyzzy

						chainRelay = member.protocolConfig.chain.ThresholdRelay()

						// func (pm *PublishingMember) setupCurrentBlockHeight() error {
						err = setup2CurrentBlockHeight(member.protocolConfig.chain)

					// Step 1 - Publish result
					case "publish":
						groupPublicKey := cmd.opt1 // take 1st opt and use that for publishing - it is a int64
						result1 := &result.Result{GroupPublicKey: big.NewInt(groupPublicKey)}
						expectedEvent := &event.PublishedResult{
							PublisherID: member.ID,
							Hash:        cmd.resultHash,
						}
						if one_time && chainRelay.IsResultPublished(result1) {
							one_time = false
							t.Errorf("Result is already published on chain")
						}

						eventPublish, err := member.PublishResult(result1, 5)
						if err != nil {
							t.Errorf("\nexpected: %s\nactual:   %s\n", "", err)
						}
						if !reflect.DeepEqual(expectedEvent, eventPublish) {
							t.Errorf("\nexpected: %v\nactual:   %v\n", expectedEvent, eventPublish)
						}

						if !chainRelay.IsResultPublished(result1) {
							t.Errorf("Result should be published on chain")
						}

					case "print-block-height":
						fmt.Printf("Block Height: ?????\n")

					case "challenge":
						//func (vc *DKG) Challenge(resultHash string, correctResult int) {
						vc = &DKG{
							P:         big.NewInt(100), //
							Q:         big.NewInt(9),   //
							TMax:      8,               // M_Max
							TConflict: 4,               // T_conflict
							TNow:      1,               // T_Now - Current block height at the current time
							TFirst:    1,               // T_First -
						}
						vc.PerGroup = make(map[string]ValidationGroupState)
						group := fmt.Sprintf("%s", big.NewInt(cmd.opt1))
						vc.PerGroup[group] = ValidationGroupState{
							TFirst:           12,                // T_first - Block height for the group - when the first result event occurred block height
							AllResults:       []result.Result{}, // Set of all results
							AllVotes:         []ResultVotes{},   // Set of all results
							LeadResult:       0,                 // Position of lead result
							AlreadySubmitted: false,             //
						}

						resultHash := cmd.resultHash
						if resultHash == nil {
							fmt.Printf("%sOOPS - bad command, pc=%d\n%s ", MiscLib.ColorRed, pc, MiscLib.ColorReset)
						}
						correctResult := cmd.opt0
						go vc.Challenge(resultHash, correctResult)

					case "block":
						go vc.NextBlock()
						time.Sleep(1 * time.Second)

					case "dump-data":
						vc.DumpData("dump data from main")

					case "validate-data":
						go vc.ValidateResults(t)
						time.Sleep(1 * time.Second)

					case "vote":
						resultHash := cmd.resultHash
						go vc.Vote(resultHash, big.NewInt(cmd.opt1))
						time.Sleep(1 * time.Second)

					case "send":
						Submission <- result.Result{
							Type:           result.SuccessIA,
							GroupPublicKey: big.NewInt(cmd.opt1),
							Disqualified:   []int{},
							Inactive:       []int{},
						}
						time.Sleep(1 * time.Second)

					default:
						fmt.Printf("Fatal Error: bad command at pc=%d of %+v\n", pc, cmd)
						t.Fatalf("Bad Command")
					}
				}
			}

		})

		_ = x

	}

}

func initializePublishingMembersGroup2(threshold, groupSize int) ([]*PublishingMember, error) {
	chain := local.Connect(10, 4)
	blockCounter, err := chain.BlockCounter()
	if err != nil {
		return nil, err
	}
	err = blockCounter.WaitForBlocks(1)
	if err != nil {
		return nil, err
	}

	initialBlockHeight, err := blockCounter.CurrentBlock()
	if err != nil {
		return nil, err
	}

	dkg := &DKG{
		chain:              chain,
		expectedDuration:   4,
		blockStep:          1,
		initialBlockHeight: initialBlockHeight,
	}

	group := &Group{
		groupSize:          groupSize,
		dishonestThreshold: threshold,
	}

	var members []*PublishingMember

	for i := 1; i <= groupSize; i++ {
		id := i
		members = append(members,
			&PublishingMember{
				SharingMember: &SharingMember{
					QualifiedMember: &QualifiedMember{
						SharesJustifyingMember: &SharesJustifyingMember{
							CommittingMember: &CommittingMember{
								memberCore: &memberCore{
									ID:             id,
									group:          group,
									protocolConfig: dkg,
								},
							},
						},
					},
				},
			})
		group.RegisterMemberID(id)
	}
	return members, nil
}

func (vc *DKG) ValidateResults(t *testing.T) {
	found := false
	for x := range Complete {
		if strings.HasPrefix(x.FinishState, "success-") {
			found = true
		}
	}
	if !found {
		t.Errorf("Uable to find the stucces state\n")
	}
}

func (vc *DKG) DumpData(msg string) {
}

func (vc *DKG) NextBlock() {
	// Increment to next block - call processing function
	TNow := getCurrentBlockHeight() // has to have geth connect info
	fmt.Printf("....... Sleep started block %d ......\n", TNow)
	time.Sleep(501 * time.Millisecond)
	TNow = getCurrentBlockHeight() // has to have geth connect info
	fmt.Printf("....... wakie wakie block %d ......\n", TNow)
	Control <- ControlType{
		Cmd: "block-no",
	}
}
