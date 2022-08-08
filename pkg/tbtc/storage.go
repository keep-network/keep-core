package tbtc

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/hex"
	"fmt"
	"github.com/keep-network/keep-common/pkg/persistence"
	"sync"
)

// walletStorage is the component that persists data of the wallets managed
// by the given node. All functions of the storage are safe for concurrent use.
type walletStorage struct {
	// mutex is a single struct-wide lock that ensures all functions
	// of the storage are thread-safe.
	mutex sync.Mutex

	// persistence is the handle to the underlying persistence layer.
	persistence persistence.Handle

	// wallets is a cache of maintained wallets. The cache's key is the
	// uncompressed public key of the given wallet. The cache's value is
	// a slice of the wallet signers controlled by this node.
	wallets map[string][]*signer
}

// connectWalletStorage connects the underlying persistence layer and
// returns the walletStorage handle.
func connectWalletStorage(persistence persistence.Handle) *walletStorage {
	walletStorage := &walletStorage{
		persistence: persistence,
		wallets:     make(map[string][]*signer),
	}

	walletStorage.prepopulateCache()

	return walletStorage
}

// saveSigner saves the given signer using the underlying persistence layer.
func (ws *walletStorage) saveSigner(signer *signer) error {
	ws.mutex.Lock()
	defer ws.mutex.Unlock()

	walletStorageKey := getWalletStorageKey(signer.wallet.publicKey)

	signerBytes, err := signer.Marshal()
	if err != nil {
		return fmt.Errorf("could not marshal signer: [%w]", err)
	}

	err = ws.persistence.Save(
		signerBytes,
		walletStorageKey,
		fmt.Sprintf("/membership_%v", signer.signingGroupMemberIndex),
	)
	if err != nil {
		return fmt.Errorf(
			"could not save membership using the "+
				"underlying persistence layer: [%w]",
			err,
		)
	}

	ws.wallets[walletStorageKey] = append(ws.wallets[walletStorageKey], signer)

	return nil
}

// prepopulateCache loads all signers using the underlying persistence layer
// and uses them to pre-populate the in-memory cache.
func (ws *walletStorage) prepopulateCache() {
	ws.mutex.Lock()
	defer ws.mutex.Unlock()

	ws.wallets = make(map[string][]*signer)

	descriptorsChan, errorsChan := ws.persistence.ReadAll()

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		for descriptor := range descriptorsChan {
			content, err := descriptor.Content()
			if err != nil {
				logger.Errorf(
					"could not get content from file [%v] "+
						"in directory [%v]: [%v]",
					descriptor.Name(),
					descriptor.Directory(),
					err,
				)
				continue
			}

			signer := &signer{}
			if err := signer.Unmarshal(content); err != nil {
				logger.Errorf(
					"could not unmarshal signer from file [%v] "+
						"in directory [%v]: [%v]",
					descriptor.Name(),
					descriptor.Directory(),
					err,
				)
				continue
			}

			walletStorageKey := getWalletStorageKey(signer.wallet.publicKey)

			ws.wallets[walletStorageKey] = append(
				ws.wallets[walletStorageKey],
				signer,
			)
		}

		wg.Done()
	}()

	go func() {
		for err := range errorsChan {
			logger.Errorf(
				"could not load signer from disk: [%v]",
				err,
			)
		}

		wg.Done()
	}()

	wg.Wait()

	for walletStorageKey, signers := range ws.wallets {
		logger.Infof(
			"wallet's signing group [0x%v] loaded with [%v] members",
			walletStorageKey,
			len(signers),
		)
	}
}

func (ws *walletStorage) getSigners(
	walletPublicKey *ecdsa.PublicKey,
) []*signer {
	ws.mutex.Lock()
	defer ws.mutex.Unlock()

	return ws.wallets[getWalletStorageKey(walletPublicKey)]
}

// getWalletStorageKey compute the wallet storage key that is used to identify
// the given wallet for storage purposes.
func getWalletStorageKey(walletPublicKey *ecdsa.PublicKey) string {
	walletPublicKeyBytes := elliptic.Marshal(
		walletPublicKey.Curve,
		walletPublicKey.X,
		walletPublicKey.Y,
	)

	return hex.EncodeToString(walletPublicKeyBytes)
}
