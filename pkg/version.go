package pkg

//nolint:gochecknoglobals //reason:app build properties
var (
	// Release returns the release version
	Release = "UNKNOWN"
	// CompileTime returns the git repository URL
	CompileTime = "UNKNOWN"
	// Commit returns the short sha from git
	Commit = "UNKNOWN"
)
