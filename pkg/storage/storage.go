package storage

import (
	"fmt"
	"path"
	"path/filepath"

	"github.com/keep-network/keep-common/pkg/persistence"
)

// Config stores meta-info about keeping data on disk
type Config struct {
	// Path to the persistent storage directory on disk.
	Dir string
}

const (
	// The key store persistent storage directory is a place where the sensitive
	// data are stored. The greatest care should be taken to backup the directory
	// and never loose any data stored there.
	// Losing data stored in this directory is a serious protocol violation.
	keyStoreDirName = "keystore"
	// The work persistent storage directory keeps data that should persist
	// the client restart. Losing data stored in this directory may
	// lead to losing rewards as a result of inactivity but is not
	// a protocol violation.
	workDirName = "work"
)

// Storage is a disk persistent storage for the client.
type Storage struct {
	keystoreDir        string
	workDir            string
	encryptionPassword string
}

// Initialize initializes a disk storage with `keystore` and `work` directories.
// The provided `encryptionPassword` will be used to encrypt the work persisted
// to the storage.
func Initialize(config Config, encryptionPassword string) (Storage, error) {
	storage := Storage{}

	storageRootDir := filepath.Clean(config.Dir)

	if err := persistence.EnsureDirectoryExists(
		storageRootDir,
		keyStoreDirName,
	); err != nil {
		return storage, fmt.Errorf(
			"cannot create storage directory for keystore: [%w]",
			err,
		)
	}
	storage.keystoreDir = filepath.Join(storageRootDir, keyStoreDirName)

	if err := persistence.EnsureDirectoryExists(
		storageRootDir,
		workDirName,
	); err != nil {
		return storage, fmt.Errorf(
			"cannot create storage directory for work: [%w]",
			err,
		)
	}
	storage.workDir = filepath.Join(storageRootDir, workDirName)

	storage.encryptionPassword = encryptionPassword

	return storage, nil
}

// InitializeKeyStorePersistence initializes a disk persistence under keystore parent.
func (s *Storage) InitializeKeyStorePersistence(dir string) (
	persistence.Handle,
	error,
) {
	return s.initializePersistence(s.keystoreDir, dir)
}

// InitializeWorkPersistence initializes a disk persistence under work parent.
func (s *Storage) InitializeWorkPersistence(dir string) (
	persistence.Handle,
	error,
) {
	return s.initializePersistence(s.workDir, dir)
}

// initializePersistence creates a persistent directory under a parent directory.
// It returns an error is the parent directory doesn't exist.
func (s *Storage) initializePersistence(parentDir string, dir string) (
	persistence.Handle,
	error,
) {
	if err := persistence.EnsureDirectoryExists(parentDir, dir); err != nil {
		return nil, fmt.Errorf(
			"cannot create storage directory [%s] in [%s]: [%w]",
			dir,
			parentDir,
			err,
		)
	}

	path := path.Join(parentDir, dir)

	diskHandle, err := persistence.NewDiskHandle(path)
	if err != nil {
		return nil, fmt.Errorf("cannot create [%s] disk handle: [%w]", path, err)
	}

	return persistence.NewEncryptedPersistence(
		diskHandle,
		s.encryptionPassword,
	), nil
}
