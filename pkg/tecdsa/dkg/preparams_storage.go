package dkg

import (
	"encoding/hex"
	"fmt"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/crypto"

	"github.com/ipfs/go-log"

	"github.com/keep-network/keep-common/pkg/persistence"
)

const (
	dirName = "preparams"
)

type preParamsStorage struct {
	// mutex is a single struct-wide lock that ensures all functions
	// of the storage are thread-safe.
	mutex sync.Mutex

	persistence persistence.Handle
	logger      log.StandardLogger
}

func newPreParamsStorage(
	persistence persistence.Handle,
	logger log.StandardLogger,
) preParamsStorage {
	return preParamsStorage{
		persistence: persistence,
		logger:      logger,
	}
}

// Saves provided PreParams to the storage.
func (p *preParamsStorage) Save(pp *PreParams) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	ppBytes, err := pp.Marshal()
	if err != nil {
		return fmt.Errorf("marshalling of the preparams failed: [%v]", err)
	}

	fileName := getFileName(ppBytes, pp.creationTimestamp)

	if err := p.persistence.Save(
		ppBytes,
		dirName,
		fileName,
	); err != nil {
		return fmt.Errorf("saving preparams failed: [%w]", err)
	}

	return nil
}

// Deletes provided PreParams from the storage.
func (p *preParamsStorage) Delete(pp *PreParams) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	ppBytes, err := pp.Marshal()
	if err != nil {
		return fmt.Errorf("marshalling of the preparams failed: [%v]", err)
	}
	fileName := getFileName(ppBytes, pp.creationTimestamp)

	p.logger.Debugf("deleting preparams [%s]...", fileName)

	return p.persistence.Delete(dirName, fileName)
}

// ReadAll reads all the PreParams stored in the storage and returns them as a
// slice.
func (p *preParamsStorage) ReadAll() ([]*PreParams, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	allPreParams := make([]*PreParams, 0)

	descriptorsChan, errorsChan := p.persistence.ReadAll()

	// Two goroutines read from descriptors and errors channels and either
	// add the PreParams to the result slice or outputs a log error.
	// The reason for using two goroutines at the same time - one for
	// descriptors and one for errors - is that channels do not have to be
	// buffered, and we do not know in what order the information is written to
	// channels.
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		for descriptor := range descriptorsChan {
			// Read only the files located in the `dirName` subdirectory.
			if descriptor.Directory() != dirName {
				continue
			}

			content, err := descriptor.Content()
			if err != nil {
				p.logger.Errorf(
					"could not read PreParams from file [%s] in directory [%s]: [%v]",
					descriptor.Name(),
					descriptor.Directory(),
					err,
				)
				continue
			}

			preParams := &PreParams{}
			if err = preParams.Unmarshal(content); err != nil {
				p.logger.Errorf(
					"could not unmarshal PreParams from file [%s] in directory [%s]: [%v]",
					descriptor.Name(),
					descriptor.Directory(),
					err,
				)
				continue
			}

			allPreParams = append(allPreParams, preParams)
		}

		wg.Done()
	}()

	go func() {
		for err := range errorsChan {
			p.logger.Errorf(
				"could not load preparams from disk: [%v]",
				err,
			)
		}

		wg.Done()
	}()

	wg.Wait()

	return allPreParams, nil
}

func getFileName(ppBytes []byte, ppCreationTimestamp time.Time) string {
	return fmt.Sprintf(
		"pp_%d_%s",
		// Use timestamp in the filename so that when the data are read from the
		// disk the First In, First Out algorithm applies.
		ppCreationTimestamp.UnixMilli(),
		// Add part of the hash for an ultra unlikely scenario that two saves
		// happen in the exactly the same millisecond.
		hex.EncodeToString(crypto.Keccak256(ppBytes))[:7],
	)
}
