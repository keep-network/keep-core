package dkg

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"sync"

	"github.com/ipfs/go-log/v2"

	"github.com/keep-network/keep-common/pkg/persistence"
	"github.com/keep-network/keep-core/pkg/generator"
)

const (
	dirName = "preparams"
)

// PersistedPreParams is an alias for Persisted PreParams used in generator.Persistence
// interface implementation.
type PersistedPreParams = generator.Persisted[PreParams]

type preParamsStorage struct {
	// mutex is a single struct-wide lock that ensures all functions
	// of the storage are thread-safe.
	mutex sync.Mutex

	persistence persistence.BasicHandle
	logger      log.StandardLogger
}

func newPreParamsStorage(
	persistence persistence.BasicHandle,
	logger log.StandardLogger,
) preParamsStorage {
	return preParamsStorage{
		persistence: persistence,
		logger:      logger,
	}
}

// Saves provided PreParams to the storage.
func (p *preParamsStorage) Save(pp *PreParams) (*PersistedPreParams, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	ppBytes, err := pp.Marshal()
	if err != nil {
		return nil, fmt.Errorf("marshalling of the preparams failed: [%v]", err)
	}
	ppHash := sha256.Sum256(ppBytes)

	fileName := fmt.Sprintf(
		"pp_%d_%s",
		// Use timestamp in the filename so that when the data are read from the
		// disk the First In, First Out algorithm applies.
		pp.creationTimestamp.UnixMilli(),
		// Add part of the hash for an ultra unlikely scenario that two saves
		// happen in the exactly the same millisecond.
		hex.EncodeToString(ppHash[:7]),
	)

	if err := p.persistence.Save(
		ppBytes,
		dirName,
		fileName,
	); err != nil {
		return nil, fmt.Errorf("saving preparams failed: [%w]", err)
	}

	return &PersistedPreParams{Data: *pp, ID: fileName}, nil
}

// Deletes provided PreParams from the storage.
func (p *preParamsStorage) Delete(pp *PersistedPreParams) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.logger.Debugf("deleting preparams [%s]...", pp.ID)

	return p.persistence.Delete(dirName, pp.ID)
}

// ReadAll reads all the PreParams stored in the storage and returns them as a
// slice.
func (p *preParamsStorage) ReadAll() ([]*PersistedPreParams, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	allPreParams := make([]*PersistedPreParams, 0)

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

			persistedPreParams := &PersistedPreParams{}
			if err = persistedPreParams.Data.Unmarshal(content); err != nil {
				p.logger.Errorf(
					"could not unmarshal PreParams from file [%s] in directory [%s]: [%v]",
					descriptor.Name(),
					descriptor.Directory(),
					err,
				)
				continue
			}
			persistedPreParams.ID = descriptor.Name()

			allPreParams = append(allPreParams, persistedPreParams)
		}

		sort.Slice(allPreParams, func(i, j int) bool {
			return allPreParams[i].Data.creationTimestamp.
				Before(allPreParams[j].Data.creationTimestamp)
		})

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
