package dkg2

import (
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"sync"
	"testing"
	"time"

	relayChain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	"github.com/keep-network/keep-core/pkg/beacon/relay/gjkr"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/chain/local"
)

func TestPhase14_pt1(t *testing.T) {
	threshold := 2
	groupSize := 5
	blockStep := 2 // T_step

	chainHandle, _ /*initialBlock*/, err := initChainHandle2(threshold, groupSize)
	if err != nil {
		t.Fatal(err)
	}

	type runCode struct {
		op              string
		requestID       *big.Int
		groupPubKey     *big.Int
		dkgResult       *relayChain.DKGResult
		intVal          int
		intVal2         int
		resultToPublish *relayChain.DKGResult
	}

	var tests = map[string]struct {
		runIt           bool
		correctResult   *relayChain.DKGResult
		publishingIndex int
		steps           []runCode
	}{
		"vote test with no data, nothing submitted": {
			runIt: false,
			correctResult: &relayChain.DKGResult{
				GroupPublicKey: []byte{40, 00},
			},
			publishingIndex: 0,
			steps: []runCode{
				{
					op: "go-phase14",
				},
			},
		},
		"send a Vote - 1 vote after start": {
			runIt: false,
			correctResult: &relayChain.DKGResult{
				GroupPublicKey: []byte{40, 01},
			},
			publishingIndex: 0,
			steps: []runCode{
				{op: "setup", intVal: 102},
				{op: "submit-result",
					requestID: big.NewInt(102),
					resultToPublish: &relayChain.DKGResult{
						Success:        true,
						GroupPublicKey: []byte{40, 01},
					},
				},
				{op: "go-phase14"}, // Process result
				{op: "send-vote",
					requestID: big.NewInt(102),
					dkgResult: &relayChain.DKGResult{
						Success:        true,
						GroupPublicKey: []byte{40, 01},
					},
				},
				{op: "dump-submissions", requestID: big.NewInt(102)},
				{op: "validate-votes",
					requestID: big.NewInt(102), // request to check
					intVal:    2,               // # of votes to expect
					intVal2:   0,               // positon in submission set
				},
			},
		},
		"multiple votes for a singe successful result waiting for timeout": {
			runIt: false,
			correctResult: &relayChain.DKGResult{
				GroupPublicKey: []byte{40, 01},
			},
			publishingIndex: 0,
			steps: []runCode{
				{op: "setup", intVal: 102},
				{op: "submit-result",
					requestID: big.NewInt(102),
					resultToPublish: &relayChain.DKGResult{
						Success:        true,
						GroupPublicKey: []byte{40, 01},
					},
				},
				{op: "sleep", intVal: 500},
				{op: "go-phase14"}, // Process result
				{op: "send-vote",
					requestID: big.NewInt(102),
					dkgResult: &relayChain.DKGResult{
						Success:        true,
						GroupPublicKey: []byte{40, 01},
					},
				},
				{op: "dump-submissions", requestID: big.NewInt(102)},
				{op: "validate-votes",
					requestID: big.NewInt(102), // request to check
					intVal:    2,               // # of votes to expect
					intVal2:   0,               // positon in submission set
				},
				{op: "sleep", intVal: 1500},
				{op: "send-vote",
					requestID: big.NewInt(102),
					dkgResult: &relayChain.DKGResult{
						Success:        true,
						GroupPublicKey: []byte{40, 01},
					},
				},
				{op: "validate-votes",
					requestID: big.NewInt(102), // request to check
					intVal:    3,               // # of votes to expect
					intVal2:   0,               // positon in submission set
				},
				{op: "sleep", intVal: 1500},
			},
		},
		"pass the voting threshold": {
			runIt: true,
			correctResult: &relayChain.DKGResult{
				GroupPublicKey: []byte{40, 01},
			},
			publishingIndex: 0,
			steps: []runCode{
				{op: "setup", intVal: 102},
				{op: "submit-result",
					requestID: big.NewInt(102),
					resultToPublish: &relayChain.DKGResult{
						Success:        true,
						GroupPublicKey: []byte{40, 01},
					},
				},
				{op: "sleep", intVal: 500},
				{op: "go-phase14"}, // Process result
				{op: "send-vote",
					requestID: big.NewInt(102),
					dkgResult: &relayChain.DKGResult{
						Success:        true,
						GroupPublicKey: []byte{40, 01},
					},
				},
				{op: "dump-submissions", requestID: big.NewInt(102)},
				{op: "validate-votes",
					requestID: big.NewInt(102), // request to check
					intVal:    2,               // # of votes to expect
					intVal2:   0,               // positon in submission set
				},
				{op: "sleep", intVal: 100},
				{op: "send-vote",
					requestID: big.NewInt(102),
					dkgResult: &relayChain.DKGResult{
						Success:        true,
						GroupPublicKey: []byte{40, 01},
					},
				},
				{op: "validate-votes",
					requestID: big.NewInt(102), // request to check
					intVal:    3,               // # of votes to expect
					intVal2:   0,               // positon in submission set
				},
				{op: "send-vote",
					requestID: big.NewInt(102),
					dkgResult: &relayChain.DKGResult{
						Success:        true,
						GroupPublicKey: []byte{40, 01},
					},
				},
				{op: "sleep", intVal: 100},
			},
		},
		"pass the voting threshold, 2nd time with new requestID": {
			runIt: true,
			correctResult: &relayChain.DKGResult{
				GroupPublicKey: []byte{40, 01},
			},
			publishingIndex: 0,
			steps: []runCode{
				{op: "setup", intVal: 103},
				{op: "submit-result",
					requestID: big.NewInt(103),
					resultToPublish: &relayChain.DKGResult{
						Success:        true,
						GroupPublicKey: []byte{40, 01},
					},
				},
				{op: "sleep", intVal: 500},
				{op: "go-phase14"}, // Process result
				{op: "send-vote",
					requestID: big.NewInt(103),
					dkgResult: &relayChain.DKGResult{
						Success:        true,
						GroupPublicKey: []byte{40, 01},
					},
				},
				{op: "dump-submissions", requestID: big.NewInt(103)},
				{op: "validate-votes",
					requestID: big.NewInt(103), // request to check
					intVal:    2,               // # of votes to expect
					intVal2:   0,               // positon in submission set
				},
				{op: "sleep", intVal: 100},
				{op: "send-vote",
					requestID: big.NewInt(103),
					dkgResult: &relayChain.DKGResult{
						Success:        true,
						GroupPublicKey: []byte{40, 01},
					},
				},
				{op: "validate-votes",
					requestID: big.NewInt(103), // request to check
					intVal:    3,               // # of votes to expect
					intVal2:   0,               // positon in submission set
				},
				{op: "send-vote",
					requestID: big.NewInt(103),
					dkgResult: &relayChain.DKGResult{
						Success:        true,
						GroupPublicKey: []byte{40, 01},
					},
				},
				{op: "sleep", intVal: 100},
			},
		},
	}

	thresholdRelayChain := chainHandle.ThresholdRelay()
	var wg sync.WaitGroup

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			if !test.runIt {
				return
			}
			publisher := &Publisher{}

			for pc, ex := range test.steps {
				// fmt.Printf("------ Running %s ------\n", ex.op)
				switch ex.op {
				case "setup":
					publisher = &Publisher{
						ID:               gjkr.MemberID(test.publishingIndex + 1),
						RequestID:        big.NewInt(int64(ex.intVal)),
						publishingIndex:  test.publishingIndex,
						chainHandle:      chainHandle,
						blockStep:        blockStep,
						conflictDuration: 8, // T_conflict
						votingThreshold:  3, // T_max
					}

				case "sleep-1-sec":
					// fmt.Printf("*** Sleep for 1 sec ***\n")
					time.Sleep(1 * time.Second)
				case "sleep":
					// fmt.Printf("*** Sleep for %d millisecond ***\n", ex.intVal)
					time.Sleep(time.Duration(ex.intVal) * time.Millisecond)
				case "call-phase14": // blocking call to publisher.Phase14!
					publisher.Phase14(test.correctResult)
				case "submit-result":
					promise := thresholdRelayChain.SubmitDKGResult(ex.requestID, ex.resultToPublish)
					_ = promise // local test will immediately fulfil so can be ignored.
				case "dump-submissions":
					submissions := thresholdRelayChain.GetDKGSubmissions(ex.requestID)
					if os.Getenv("publish_verbose") == "yes" {
						fmt.Printf("submissions ->%s<-\n", convertToJSON(submissions))
					}
				case "validate-votes": // intVal: 2},
					submissions := thresholdRelayChain.GetDKGSubmissions(ex.requestID)
					votes := submissions.DKGSubmissions[ex.intVal2].Votes
					if votes != ex.intVal {
						t.Errorf("Invalid number of votes\n")
					}
				case "send-vote":
					dkgResultHash := ex.dkgResult.Hash()
					thresholdRelayChain.DKGResultVote(ex.requestID, dkgResultHash)
				case "go-phase14":
					wg.Add(1)
					go func() {
						defer wg.Done()
						err := publisher.Phase14(test.correctResult)
						if err != nil {
							fmt.Printf("**** Error: [%v]\n", err)
							t.Errorf("Returned an error from Phase14 - [%v]\n", err)
						}
					}()
					time.Sleep(1 * time.Second)

				default:
					fmt.Printf("In test [%s] invalid op [%s] at %d\n", testName, ex.op, pc)
				}
			}
		})
	}

	wg.Wait()
}

func initChainHandle2(threshold, groupSize int) (chainHandle chain.Handle, initialBlock int, err error) {
	chainHandle = local.Connect(groupSize, threshold)
	blockCounter, err := chainHandle.BlockCounter() // PJS - save blockCounter?
	if err != nil {
		return nil, -1, err
	}
	err = blockCounter.WaitForBlocks(1)
	if err != nil {
		return nil, -1, err
	}

	initialBlock, err = blockCounter.CurrentBlock() // PJS - need CurrentBlock to make this work
	if err != nil {
		return nil, -1, err
	}
	return
}

// convertToJSON return the JSON encoded version of the data with tab indentation.
func convertToJSON(v interface{}) string {
	// s, err := json.Marshal ( v )
	s, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		return fmt.Sprintf("Error:%s", err)
	} else {
		return string(s)
	}
}
