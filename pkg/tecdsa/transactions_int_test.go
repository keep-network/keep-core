package tecdsa

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"testing"

	"github.com/btcsuite/btcutil"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
)

type Transaction struct {
	TxID               string
	DestinationAddress string
	TxMessage          *wire.MsgTx
}

type Node struct {
	wif     string
	address string
}

type SourceTransaction struct {
	hash     string
	pointer  uint32
	pkScript []byte
}

func TestTransactionSigning(t *testing.T) {
	node1 := Node{
		wif:     "91xmXfzdXXqgf2DLJguwTG4eMmbYBeEngfATwCNaBxqcTMoM7RC",
		address: "miBLnDq73eZRSZ92qzAQYqbpGg95YBxPFF",
	}
	node2 := Node{
		wif:     "923CjseKgQf7Xz185dmYUJer9i8jsb9Cd18Rtec4DFKeiBZg3wi",
		address: "n1EaXu1KH6QqutB5pwZjacs1t4yxx3m3Ex",
	}
	sourceWif := node1.wif
	sourceAddress := node1.address
	destinationAddress := node2.address
	amount := int64(912)

	sourcePKScript, err := hex.DecodeString("76a9141d32d40498522678af8491e995f21ed4caa9116f88ac")
	if err != nil {
		t.Fatal(err)
	}

	sourceTransaction := SourceTransaction{
		hash:     "048b45b0ca5d4f00b223ac1ae56b5c2168203b853aab3e4972382ec2e5a5f820",
		pointer:  uint32(0),
		pkScript: sourcePKScript,
	}

	networkParams := &chaincfg.TestNet3Params

	if result, err := validatePubKeyScriptAddress(sourceTransaction.pkScript, sourceAddress, networkParams); !result || err != nil {
		t.Fatal(err)
	}

	transaction, err := CreateTransaction(sourceWif, destinationAddress, amount, sourceTransaction, networkParams)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("Generated transaction ID: %s\n", transaction.TxID)

	// SendTransactionToPeer(transaction)
	// SendTransactionToQBitNinjaAPI(transaction.TxMessage)
	SendTransactionToRPC(transaction)
}

func validatePubKeyScriptAddress(pubKeyScript []byte, expectedAddress string, networkParams *chaincfg.Params) (bool, error) {
	_, address, _, err := txscript.ExtractPkScriptAddrs(pubKeyScript, networkParams)
	if err != nil {
		return false, err
	}

	if address[0].String() != expectedAddress {
		return false, fmt.Errorf("Address from Public Key Script doesn't match source address.\nActual: %s\nExpected: %s", address[0], expectedAddress)
	}
	return true, nil
}

func CreateTransaction(wifString string,
	destination string,
	amount int64,
	sourceTx SourceTransaction,
	networkParams *chaincfg.Params,
) (*Transaction, error) {
	wif, err := btcutil.DecodeWIF(wifString)
	if err != nil {
		return nil, err
	}
	// Initialize transaction message
	msgTx := wire.NewMsgTx(wire.TxVersion)

	// Source transaction hash
	sourceTxHash, err := chainhash.NewHashFromStr(sourceTx.hash)
	if err != nil {
		return nil, err
	}
	// Input
	previousTxOut := wire.NewOutPoint(sourceTxHash, sourceTx.pointer)
	txIn := wire.NewTxIn(previousTxOut, nil, nil)
	msgTx.AddTxIn(txIn)

	// Destination Address
	destinationAddress, err := btcutil.DecodeAddress(destination, networkParams)
	if err != nil {
		return nil, err
	}
	// Destination Pub Key Script
	outPKScript, _ := txscript.PayToAddrScript(destinationAddress)
	// Transaction Output
	txOut := wire.NewTxOut(amount, outPKScript)
	msgTx.AddTxOut(txOut)

	// Sign Transaction
	sigScript, err := txscript.SignatureScript(msgTx, 0, sourceTx.pkScript, txscript.SigHashAll, wif.PrivKey, true)
	if err != nil {
		return nil, err
	}
	msgTx.TxIn[0].SignatureScript = sigScript

	// Validate Transaction
	validationEngine, err := txscript.NewEngine(sourceTx.pkScript, msgTx, 0, txscript.StandardVerifyFlags, nil, nil, amount)
	if err != nil {
		return nil, err
	}
	if err := validationEngine.Execute(); err != nil {
		return nil, err
	}

	return &Transaction{
		TxID:               msgTx.TxHash().String(),
		DestinationAddress: destinationAddress.EncodeAddress(),
		TxMessage:          msgTx,
	}, nil
}

// Requires local RPC server to be runing
// Command:
// btcd --testnet --notls --rpcuser testNetRPCuser --rpcpass testNetRPCpass
func SendTransactionToRPC(transaction *Transaction) {
	client, err := rpcclient.New(&rpcclient.ConnConfig{
		HTTPPostMode: true,
		DisableTLS:   true,
		Host:         "127.0.0.1:18334",
		User:         "testNetRPCuser",
		Pass:         "testNetRPCpass",
	}, nil)
	if err != nil {
		log.Fatalf("error creating new btc client: %v", err)
	}

	txID, err := client.SendRawTransaction(transaction.TxMessage, false)
	if err != nil {
		log.Fatalf("error sendMany: %v", err)
	}
	log.Printf("SendRawTransaction completed! tx sha is: %s", txID.String())
}

// List of DNS Seeds
// https://github.com/btcsuite/btcd/blob/9a2f9524024889e129a5422aca2cff73cb3eabf6/chaincfg/params.go#L405
func SendTransactionToPeer(transaction *Transaction) error {
	// Create version message data.
	lastBlock := int32(234234)
	tcpAddrMe := &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 8333}
	me := wire.NewNetAddress(tcpAddrMe, wire.SFNodeNetwork)
	tcpAddrYou := &net.TCPAddr{IP: net.ParseIP("192.168.0.1"), Port: 8333}
	you := wire.NewNetAddress(tcpAddrYou, wire.SFNodeNetwork)
	nonce, err := wire.RandomUint64()
	if err != nil {
		return err
	}

	versionMsg := wire.NewMsgVersion(me, you, nonce, lastBlock)

	// connect to this socket
	peerIP := "136.144.129.222"
	peerPort := "18333"
	conn1, _ := net.Dial("tcp", fmt.Sprintf("%s:%s", peerIP, peerPort))

	// Use the most recent protocol version supported by the package and the
	// main bitcoin network.
	pver := wire.ProtocolVersion
	btcnet := wire.TestNet3

	// Writes a bitcoin message msg to conn using the protocol version
	// pver, and the bitcoin network btcnet.  The return is a possible
	// error.
	err = wire.WriteMessage(conn1, versionMsg, pver, btcnet)
	if err != nil {
		return fmt.Errorf("Error writing %s", err)
	}

	msg1, _, err := wire.ReadMessage(conn1, pver, btcnet)
	if err != nil {
		return fmt.Errorf("Error writing %s", err)
	}
	fmt.Printf("MSG1: %v\n", msg1)
	msg2, _, err := wire.ReadMessage(conn1, pver, btcnet)
	if err != nil {
		return fmt.Errorf("Error writing %s", err)
	}
	fmt.Printf("MSG2: %v\n", msg2)

	err = wire.WriteMessage(conn1, wire.NewMsgInv(), pver, btcnet)
	if err != nil {
		return fmt.Errorf("Error writing %s", err)
	}

	// TRANSACTION
	err = wire.WriteMessage(conn1, transaction.TxMessage, pver, btcnet)
	if err != nil {
		return fmt.Errorf("Error writing %s", err)
	}

	// // connect to this socket
	// peerIp2 := "94.130.216.246"
	// conn2, _ := net.Dial("tcp", fmt.Sprintf("%s:%s", peerIp2, peerPort))
	// fmt.Printf("MSG1: %v\n", msg1)
	// msgAfterTx, _, err := wire.ReadMessage(conn2, pver, btcnet)
	// if err != nil {
	// 	fmt.Errorf("Error writing %s", err)
	// }
	// fmt.Printf("msgAfterTx: %v\n", msgAfterTx)

	return nil
}

func SendTransactionToQBitNinjaAPI(msgTx *wire.MsgTx) error {
	apiURL := "https://tapi.qbit.ninja/transactions"

	fmt.Printf("Send transaction: %v\nto: %s\n", encodeMsgTx(msgTx), apiURL)

	body := []byte(fmt.Sprintf("\"%s\"", encodeMsgTx(msgTx)))
	request, _ := http.NewRequest("POST", apiURL, bytes.NewBuffer(body))
	request.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		return fmt.Errorf("Error when sending request to the server [%s]", err)
	}

	respBody, _ := ioutil.ReadAll(response.Body)
	defer response.Body.Close()

	fmt.Printf("Response status: %s\n", response.Status)
	fmt.Printf("Response body:\n%s\n", respBody)

	return nil
}

func encodeMsgTx(msgTx *wire.MsgTx) string {
	var signedTx bytes.Buffer
	msgTx.Serialize(&signedTx)
	return hex.EncodeToString(signedTx.Bytes())
}
