package build

var (
	Version  = "unknown" // Client software version. Set with `-ldflags -X` during the build.
	Revision = "unknown" // Client software revision. Set with `-ldflags -X` during the build.
)
