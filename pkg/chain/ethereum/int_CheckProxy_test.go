package ethereum

/*
// export KEEP_SLOW_TEST="run"
	var envFn = os.Getenv("KEEP_TEST_CFG")
KeepRandomBeacon/KeepRandomBeaconImplV1:
    1. Call KeepRandomBeaconImplV1.minimumStake() to verify that you get back a 0 and you are talking directly to the contract.
    2. Call KeepRandomBeacon(proxy).minimumStake() to verify that you are talking to the proxy. (Get back a non-zero value)
*/

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

var WaitForTx chan string

func init() {
	WaitForTx = make(chan string, 100)
}

func SendSuccess(s string) {
	WaitForTx <- "Success: " + s
}

func SendError(s string) {
	WaitForTx <- "Error: " + s
}

// Should add a timeout capability
func WaitForUpTo60Sec(s string) {
	// TODO See: https://gobyexample.com/timers
}

func WaitForMsg() {
	result := <-WaitForTx
	fmt.Printf("Results [%s]\n", result)
}

// This test will check that the proxy is set up and that initial configuration
// values have been created.
// KeepGroup/KeepGroupImplV1:
// 	1. Call to KeepGroupImplV1 directly to verity that group size is 0 - it should be if the proxy setup is correct.
// 	2. Call to KeepGroupImplV1 via the proxy address (KeepGroup) verify that group size is not zero
func TestKeepGroupProxy(t *testing.T) {

	SetAddressToCallImpl(TestConfig, EthConn)

	// ebc, err := EthConn.BlockCounter() // (chain.BlockCounter, error) {
	// func (kg *keepGroup) GroupSize() (int, error) {
	size, err := EthConn.keepGroupContract.GroupSize()
	if err != nil {
		if fmt.Sprintf("%s", err) == "abi: unmarshalling empty output" {
			t.Errorf("error calling GroupSize/May be incorrect contract address: [%v]", err)
		} else {
			t.Errorf("error calling GroupSize: [%v]", err)
		}
		SetAddressToCallProxy(TestConfig, EthConn)
		return
	}

	if size != 0 {
		t.Errorf(
			"\nexpected: [%v]\nactual:   [%v]\n",
			0,
			size,
		)
	}

	SetAddressToCallProxy(TestConfig, EthConn)

	// ebc, err := EthConn.BlockCounter() // (chain.BlockCounter, error) {
	// func (kg *keepGroup) GroupSize() (int, error) {
	size, err = EthConn.keepGroupContract.GroupSize()
	if err != nil {
		t.Errorf("error calling GroupSize via proxy: [%v]", err)
	}

	if size == 0 {
		t.Errorf(
			"\nexpected: [> 0 value]\nactual:   [%v]\n",
			size,
		)
	}

}

//	3. KeepGroup(proxy).setGroupSizeThreshold - to new values and verity values changed.
// TestSetGroupSizeThreshold will set group size threshold and then re-set it back to the original.
func TestSetGroupSizeThreshold(t *testing.T) {

	var runSlow = os.Getenv("KEEP_SLOW_TEST")
	if runSlow != "run" {
		fmt.Printf("Skipping TestSetGroupSizeThreshold - slow test, set KEEP_SLOW_TEST='run' to run this test\n")
		return
	}

	fmt.Printf("\n=============================\nTest: Set Group Size and Threshold\n=============================\n\n")

	// Pull out the original sizes.
	size, err := EthConn.keepGroupContract.GroupSize()
	if err != nil {
		t.Errorf("error calling GroupSize: [%v]", err)
		return
	}

	threshold, err := EthConn.keepGroupContract.GroupThreshold()
	if err != nil {
		t.Errorf("error calling GroupThreshold: [%v]", err)
		return
	}

	tx, err := EthConn.SetGroupSizeThreshold(256, 129)
	if err != nil {
		t.Errorf("error creating transaction: tx error: %s\n", err)
		return
	}
	if tx == nil {
		t.Errorf("error creating transaction: tx should not be nil - may have used an account with no funds\n")
		return
	}

	// ---------------------------------------
	fmt.Printf("tx=%s\n", PrintAsJson(tx))
	// ---------------------------------------

	fmt.Printf("Sleeping 60 seconds - waiting for blocks to occur on chain\n")
	// fmt.Printf("Typeof(tx)=%T     typeof(tx.Hash)=%T\n", tx, tx.Hash)
	txHashStr := fmt.Sprintf("0x%x", tx.Hash())
	GetTheReceipt(txHashStr, ChkForNoEvent)
	WaitForMsg()

	// Pull out the modified data.
	newSize, err := EthConn.keepGroupContract.GroupSize()
	if err != nil {
		t.Errorf("error calling GroupSize: [%v]", err)
		return
	}

	newThreshold, err := EthConn.keepGroupContract.GroupThreshold()
	if err != nil {
		t.Errorf("error calling GroupThreshold: [%v]", err)
		return
	}

	// ---------------------------------------
	fmt.Printf("newSize=%d newThreshold=%d\n", newSize, newThreshold)
	// ---------------------------------------

	if newSize != 256 {
		t.Errorf("\nexpected: [%v]\nactual:   [%v]\n", 256, newSize)
	}
	if newThreshold != 129 {
		t.Errorf("\nexpected: [%v]\nactual:   [%v]\n", 129, newThreshold)
	}

	// Reset to original values
	tx, err = EthConn.SetGroupSizeThreshold(size, threshold)
	if tx == nil {
		t.Errorf("error creating transaction: tx should not be nil\n")
		return
	}

	fmt.Printf("Sleeping 60 seconds - waiting for blocks to occur on chain\n")
	time.Sleep(60 * time.Second)

	// Pull out the original sizes.
	resetSize, err := EthConn.keepGroupContract.GroupSize()
	if err != nil {
		t.Errorf("error calling GroupSize: [%v]", err)
		return
	}

	fmt.Printf("Sleeping 60 seconds - waiting for blocks to occur on chain\n")
	time.Sleep(60 * time.Second)

	resetThreshold, err := EthConn.keepGroupContract.GroupThreshold()
	if err != nil {
		t.Errorf("error calling GroupThreshold: [%v]", err)
		return
	}

	if size != resetSize {
		t.Errorf("\nexpected: [%v]\nactual:   [%v]\n", size, resetSize)
	}
	if threshold != resetThreshold {
		t.Errorf("\nexpected: [%v]\nactual:   [%v]\n", threshold, resetThreshold)
	}

}

func ChkForNoEvent(s string, rp GetTransactionReceiptType) {
}

func GetTheReceipt(txHash string, chkEvent func(s string, rp GetTransactionReceiptType)) {

	fmt.Printf("TxHash: %s\n", txHash)

	go func(URLToCall, txHash string) {
		var s string
		var status int
		var err error
		var rp GetTransactionReceiptType
		for i := 0; i < 100; i++ {
			s, rp, status, err = FetchReceipt(URLToCall, txHash)
			if err == nil {
				break
			}
			fmt.Printf(".")
			time.Sleep(5 * time.Second)
			fmt.Printf(".")
		}
		if err != nil {
			fmt.Printf("error getting transaction: %s\n", err)
			SendError("Failed to get transaction")
		} else if status == 0 {
			fmt.Printf("failed s=->%s<-\n", s)
			SendError("")
		} else {
			fmt.Printf("Transaction succeeded success Receipt:%s\n", s)
			SendSuccess("Transaction worked.")
			// xyzzy - TODO - check for event with correct signature.
			chkEvent(s, rp)
		}
	}(TestConfig.URLRPC, txHash)
}

// Pulled from go-ethereum source and fixed
type TxdataConverted struct {
	AccountNonce uint64          `json:"nonce"`
	Price        *big.Int        `json:"gasPrice"`
	GasLimit     uint64          `json:"gas"`
	Recipient    *common.Address `json:"to"`
	Amount       *big.Int        `json:"value"`
	Payload      []byte          `json:"input"`
	V            *big.Int        `json:"v"`
	R            *big.Int        `json:"r"`
	S            *big.Int        `json:"s"`
	Hash         *common.Hash    `json:"hash"`
}

type Txdata struct {
	AccountNonce string `json:"nonce"`
	Price        string `json:"gasPrice"`
	GasLimit     string `json:"gas"`
	Recipient    string `json:"to"`
	Amount       string `json:"value"`
	Payload      string `json:"input"`
	V            string `json:"v"`
	R            string `json:"r"`
	S            string `json:"s"`
	Hash         string `json:"hash"`
}

/*
	{
				"address": "0x18b8562f6013356a6787013ada7ed168b62208c5",
				"blockHash": "0xeb8853c4422758dad770efbca50f82e185dba52fe3afd4e4e0a6046674228726",
				"blockNumber": "0xbdf91",
				"data": "0x00010203040506070809101112131415161718192ec12f07dd2e695d9f223031",
				"logIndex": "0x0",
				"removed": false,
				"topics": [
					"0x37ff73cebd127696e16cf29468c47e7573b3e5f42dd5ecdb1f4af06101050fef"
				],
				"transactionHash": "0x3050fb72ac12867db2f18a5795ee9981f2106122280149c062b6e75d375b6bb9",
				"transactionIndex": "0x0"
			}
*/
type GetLogsReceiptBody struct {
	Address          string
	BlockHash        string   `json:"blockHash"`
	BlockNumber      string   `json:"blockNumber"`
	Data             string   `json:"data"`
	LogIndex         string   `json:"logIndex"`
	Removed          bool     `json:"removed"`
	Topics           []string `json:"topics"`
	TransactionHash  string   `json:"transactionHash"`
	TransactionIndex string   `json:"transactionIndex"`
}

type GetTransactionReceiptBodyType struct {
	BlockHash         string               // hex decode
	BlockNumber       string               // hex decode
	ContractAddress   string               // hex decode, only set when loading a contract
	CumulativeGasUsed string               //
	From              string               //
	GasUsed           string               //
	Logs              []GetLogsReceiptBody //
	LogsBloom         string               // hex decoe - long
	Status            string               // hex decode, 0 indicates failure, 1 indicates success
	To                string               // address
	TransactionHash   string               // Hash
	TransactionIndex  string               // hex decode
	StatusInt         int                  //
	HasLog            bool                 //
}

type GetTransactionReceiptType struct {
	JsonRPC string
	Id      int
	Result  GetTransactionReceiptBodyType
}

var n_id = 1

// Search in topics for this item
func EventSignature(event string) string {
	eventSignature := []byte(event)
	hash := crypto.Keccak256Hash(eventSignature)
	return fmt.Sprintf("%s", hash)
}

func FetchReceipt(URLToCall, txHash string) (rv string, rp GetTransactionReceiptType, status int, err error) {

	values := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "eth_getTransactionReceipt",
		"id":      n_id,
		"params":  []string{txHash},
	}
	n_id++

	jsonValue, err := json.Marshal(values)
	if err != nil {
		// xyzzy messge
		return
	}

	// URLToCall := "http://10.51.245.213:8545/"

	resp, err := http.Post(URLToCall, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		if db0101 {
			fmt.Printf("Error: %s\n", err)
		}
		return
	}

	if db0101 {
		fmt.Printf("resp: %s err: %s\n", PrintAsJson(resp), err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		if db0101 {
			fmt.Printf("Failed to read body: %s\n", err)
		}
		return
	}
	respstatus := resp.StatusCode
	if respstatus == 200 {
		if db0101 {
			fmt.Printf("Status 200/success - Body is: ->%s<-\n", string(body))
		}
	}

	var bodyDecode GetTransactionReceiptType

	err = json.Unmarshal(body, &bodyDecode)
	if err != nil {
		if db0101 {
			fmt.Printf("Error: %s failed to decode body\n", err)
		}
		return
	}
	rp = bodyDecode

	bodyDecode.Result.HasLog = len(bodyDecode.Result.Logs) > 0
	StatusInt, err := strconv.ParseInt(bodyDecode.Result.Status, 0, 64)
	if err != nil {
		if db0101 {
			fmt.Printf("Error: %s failed to parse status\n", err)
		}
		return
	}

	bodyDecode.Result.StatusInt = int(StatusInt)
	status = bodyDecode.Result.StatusInt
	rv = fmt.Sprintf("\n%s\n", PrintAsJson(bodyDecode))

	return
}

//	4. KeepGroup(proxy).createGroup - create a group
//	5. KeepGroup(proxy).numberOfGroups - should return 1, should get a groupStartedEvent
//	6. KeepGroup(proxy).groupIndex - should return 0
//	6. KeepGroup(proxy).getGroupPubkey(0) - should return pub key used to create group
//	7. KeepGroup(proxy).emitEventGroupExists with correct key should emit a True event
//	8. KeepGroup(proxy).emitEventGroupExists with in-correct key should emit a False event
func TestGroupCreation(t *testing.T) {

	var runSlow = os.Getenv("KEEP_SLOW_TEST")
	if runSlow != "run" {
		fmt.Printf("Skipping TestGroupCreation - slow test, set KEEP_SLOW_TEST='run' to run this test\n")
		return
	}

	// createGroup
	/*
	   func (kg *keepGroup) CreateGroup(
	   	groupPubKey []byte,
	   ) (*types.Transaction, error) {
	*/
	// func (ec *ethereumChain) CreateGroup(groupPubKey []byte) (*types.Transaction, error) {

	fmt.Printf("\n=============================\nTest: Group creation\n=============================\n\n")

	// TODO - get initial number of groups
	//	5. KeepGroup(proxy).numberOfGroups - should return 1, should get a groupStartedEvent
	initialNumberOfGroups, err := EthConn.NumberOfGroups()
	if err != nil {
		t.Errorf("error calling NumberOfGroup: [%v]", err)
		return
	}

	gPubKeyArray := [32]byte{
		0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09,
		0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19,
		0x20, 0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29,
		0x30, 0x31,
	}
	b, err := GenerateRandomBytes(10)
	for i := 0; i < 10; i++ {
		gPubKeyArray[i+20] = b[i]
	}
	gPubKey := gPubKeyArray[:]
	fmt.Printf("Group Public Key: %x\n", gPubKey)

	tx, err := EthConn.CreateGroup(gPubKey)
	if err != nil {
		t.Errorf("error creating transaction: tx error: %s\n", err)
		return
	}
	if tx == nil {
		t.Errorf("error creating transaction: tx should not be nil - may have used an account with no funds\n")
		return
	}

	if db0102 {
		fmt.Printf("tx=%s\n", PrintAsJson(tx))
	}

	fmt.Printf("Sleeping 60 seconds - waiting for blocks to occur on chain\n")
	txHashStr := fmt.Sprintf("0x%x", tx.Hash())
	GetTheReceipt(txHashStr,
		func(s string, rp GetTransactionReceiptType) {
			ev := "GroupStartedEvent(bytes32)"
			es := EventSignature(ev)
			esStr := fmt.Sprintf("0x%x", es)
			for _, log := range rp.Result.Logs {
				for _, topic := range log.Topics {
					if esStr == topic {
						fmt.Printf("Results: [Event %s found]\n", ev)
						return // success just return
					}
				}
			}
			t.Errorf("failed to find topic 0x%x - missing event %s\n", es, ev)
		},
	)
	WaitForMsg()

	//	5. KeepGroup(proxy).numberOfGroups - should return 1, should get a groupStartedEvent
	numberOfGroups, err := EthConn.NumberOfGroups()
	if err != nil {
		t.Errorf("error calling NumberOfGroup: [%v]", err)
		return
	}

	if db0000 {
		fmt.Printf("Original Number of Groups: %d\nNew Number of Groups %d\n",
			initialNumberOfGroups,
			numberOfGroups,
		)
	}

	if numberOfGroups <= initialNumberOfGroups {
		t.Errorf(
			"\nexpected: [%v] to be smaller than...\nactual:   [%v]\n",
			initialNumberOfGroups,
			numberOfGroups,
		)
	}

	//	6. KeepGroup(proxy).getGroupIndex - should return 0
	idx, err := EthConn.GetGroupIndex(gPubKey)
	if err != nil {
		t.Errorf("error GetGroupIndex: [%v]\n", err)
		return
	}

	if db0000 {
		fmt.Printf("Index for this group: %d\n", idx)
	}

	//	6. KeepGroup(proxy).getGroupPubkey(0) - should return pub key used to create group
	groupPubKeyFromGeth, err := EthConn.GetGroupPubKey(idx)
	if err != nil {
		t.Errorf("error GetGroupPubKey: [%v]\n", err)
		return
	}
	if !bytes.Equal(gPubKey, groupPubKeyFromGeth) {
		t.Errorf(
			"\nexpected: [%x]\nactual:   [%x]\n",
			gPubKey,
			groupPubKeyFromGeth,
		)
	}

	// TODO
	//	7. KeepGroup(proxy).emitEventGroupExists with correct key should emit a True event
	//	8. KeepGroup(proxy).emitEventGroupExists with in-correct key should emit a False event
	// err = doWatch("KeepGroup", "GroupErrorCode", true)
	//     emit GroupErrorCode(20);
	// err = doWatch("KeepGroup", "GroupStartedEvent", true)
	// emit GroupStartedEvent(groupPubKey);

	// abi: improperly formatted output

}

const db0101 = false
const db0102 = true
const db0000 = true
