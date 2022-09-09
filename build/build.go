package build

import (
	"github.com/ipfs/go-log"
)

var (
	Version  = "unknown" // Client software version. Set with `-ldflags -X` during the build.
	Revision = "unknown" // Client software revision. Set with `-ldflags -X` during the build.
	logger   = log.Logger("keep-build")
)

func LogVersion() {
	logger.Info("version=", Version)
	logger.Info("revision=", Revision)
}
