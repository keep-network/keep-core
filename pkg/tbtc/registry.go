package tbtc

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/hex"
	"fmt"
	"sync"

	"github.com/keep-network/keep-core/pkg/bitcoin"

	"github.com/keep-network/keep-common/pkg/persistence"
)

// walletRegistry is the component that holds the data of the wallets managed
// by the given node. All functions of the registry are safe for concurrent use.
type walletRegistry struct {
	// mutex is a single struct-wide lock that ensures all functions
	// of the registry are thread-safe.
	mutex sync.Mutex

	// walletCache is a cache of maintained wallets. The cache's key is the
	// uncompressed public key of the given wallet.
	walletCache map[string]*walletCacheValue

	// walletStorage is the handle to the wallet storage responsible for
	// wallet persistence.
	walletStorage *walletStorage
}

type walletCacheValue struct {
	// SHA-256+RIPEMD-160 hash computed over the compressed ECDSA public key of
	// the wallet.
	walletPublicKeyHash [20]byte
	// Array of wallet signers controlled by this node.
	signers []*signer
}

// newWalletRegistry creates a new instance of the walletRegistry.
func newWalletRegistry(persistence persistence.ProtectedHandle) *walletRegistry {
	walletStorage := newWalletStorage(persistence)

	// Pre-populate the wallet cache using the wallet storage.
	walletCache := make(map[string]*walletCacheValue)
	walletSigners := walletStorage.loadSigners()
	if len(walletSigners) > 0 {
		for walletStorageKey, signers := range walletSigners {
			// We need to extract the wallet from the signers array. The
			// walletStorage.loadSigners function guarantees there is always
			// at least one signer for the given walletStorageKey so, we
			// don't need to check len(signers). Then, we can just take the
			// wallet from the first signer as the wallet is same for all of
			// them.
			wallet := signers[0].wallet
			walletPublicKeyHash := bitcoin.PublicKeyHash(wallet.publicKey)

			walletCache[walletStorageKey] = &walletCacheValue{
				walletPublicKeyHash: walletPublicKeyHash,
				signers:             signers,
			}

			logger.Infof(
				"wallet signing group [0x%v] loaded from storage "+
					"with [%v] members and wallet public key hash [0x%x]",
				walletStorageKey,
				len(signers),
				walletPublicKeyHash,
			)
		}
	} else {
		logger.Infof("no wallet signing groups found in the storage")
	}

	return &walletRegistry{
		walletCache:   walletCache,
		walletStorage: walletStorage,
	}
}

// registerSigner registers the given signer using in the walletRegistry.
func (wr *walletRegistry) registerSigner(signer *signer) error {
	wr.mutex.Lock()
	defer wr.mutex.Unlock()

	err := wr.walletStorage.saveSigner(signer)
	if err != nil {
		return fmt.Errorf("cannot save signer in the storage: [%w]", err)
	}

	walletStorageKey := getWalletStorageKey(signer.wallet.publicKey)

	// If the wallet cache does not have the given entry yet, initialize
	// the value and compute the wallet public key hash. This way, the hash
	// is computed only once. No need to initialize signers slice as
	// appending works with nil values.
	if _, ok := wr.walletCache[walletStorageKey]; !ok {
		wr.walletCache[walletStorageKey] = &walletCacheValue{
			walletPublicKeyHash: bitcoin.PublicKeyHash(signer.wallet.publicKey),
		}
	}

	wr.walletCache[walletStorageKey].signers = append(
		wr.walletCache[walletStorageKey].signers,
		signer,
	)

	return nil
}

// getSigners gets all signers for the given wallet held by the walletRegistry.
func (wr *walletRegistry) getSigners(
	walletPublicKey *ecdsa.PublicKey,
) []*signer {
	wr.mutex.Lock()
	defer wr.mutex.Unlock()

	if value, ok := wr.walletCache[getWalletStorageKey(walletPublicKey)]; ok {
		return value.signers
	}

	return nil
}

// getWalletByPublicKeyHash gets the given wallet by its 20-byte wallet
// public key hash. Second boolean return value denotes whether the wallet
// was found in the registry or not.
func (wr *walletRegistry) getWalletByPublicKeyHash(
	walletPublicKeyHash [20]byte,
) (wallet, bool) {
	wr.mutex.Lock()
	defer wr.mutex.Unlock()

	for _, value := range wr.walletCache {
		if value.walletPublicKeyHash == walletPublicKeyHash {
			// All signers belong to one wallet. Take that wallet from the
			// first signer.
			return value.signers[0].wallet, true
		}
	}

	return wallet{}, false
}

// walletStorage is the component that persists data of the wallets managed
// by the given node using the underlying persistence layer. It should be
// used directly only by the walletRegistry.
type walletStorage struct {
	// persistence is the handle to the underlying persistence layer.
	persistence persistence.ProtectedHandle
}

// newWalletStorage creates a new instance of the walletStorage.
func newWalletStorage(persistence persistence.ProtectedHandle) *walletStorage {
	return &walletStorage{persistence}
}

// saveSigner saves the given signer using the underlying persistence layer
// of the walletStorage. It does not add the signer to any in-memory cache
// and should not be called from any other place than walletRegistry.
func (ws *walletStorage) saveSigner(signer *signer) error {
	signerBytes, err := signer.Marshal()
	if err != nil {
		return fmt.Errorf("could not marshal signer: [%w]", err)
	}

	err = ws.persistence.Save(
		signerBytes,
		getWalletStorageKey(signer.wallet.publicKey),
		fmt.Sprintf("/membership_%v", signer.signingGroupMemberIndex),
	)
	if err != nil {
		return fmt.Errorf(
			"could not save membership using the "+
				"underlying persistence layer: [%w]",
			err,
		)
	}

	return nil
}

// loadSigners loads all signers stored using the underlying persistence layer.
// This function should not be called from any other place than walletRegistry.
func (ws *walletStorage) loadSigners() map[string][]*signer {
	signersByWallet := make(map[string][]*signer)

	descriptorsChan, errorsChan := ws.persistence.ReadAll()

	// Two goroutines read from descriptors and errors channels and either
	// add the signer to the result map or outputs a log error.
	// The reason for using two goroutines at the same time - one for
	// descriptors and one for errors - is that channels do not have to be
	// buffered, and we do not know in what order the information is written to
	// channels.
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

			signersByWallet[walletStorageKey] = append(
				signersByWallet[walletStorageKey],
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

	return signersByWallet
}

// getWalletStorageKey compute the wallet storage key that is used to identify
// the given wallet for caching and storage purposes.
func getWalletStorageKey(walletPublicKey *ecdsa.PublicKey) string {
	walletPublicKeyBytes := elliptic.Marshal(
		walletPublicKey.Curve,
		walletPublicKey.X,
		walletPublicKey.Y,
	)

	// Strip the 04 prefix to limit the key length to 128 characters in order
	// to make it usable as a directory name.
	return hex.EncodeToString(walletPublicKeyBytes)[2:]
}
