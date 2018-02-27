package EthBlockCounter

import (
	"fmt"
	"testing"
	"time"
)

func (eth *EthBlockCounter) testStop() {
	eth.ticker.Stop()
	eth.ticker2.Stop()
	for height, waiters := range eth.waiters {
		for _, waiter := range waiters {
			go func(w chan int) { w <- height }(waiter)
		}
	}
}

func Test_EthBlockCounter(t *testing.T) {

	fmt.Printf("This test takes about 1.5 min to run if it is failing and betwen 20 and 40 seconds if it is working.\n")

	GethServer := "ws://192.168.0.157:8546" // TODO - input data
	GethServer = "http://192.168.0.157:8545"

	bc := NewEthBlockCounter(GethServer, 100)

	// To turn on debuging print statements
	// eth := (bc).(*EthBlockCounter)
	// eth.debugFlag = true

	go func() {
		time.Sleep(time.Millisecond * 92500)
		eth := (bc).(*EthBlockCounter)
		eth.testStop()
		t.Errorf("Failed to catch the event\n")
		// Unblock all the waiters
	}()

	bc.WaitForBlocks(2)

}
