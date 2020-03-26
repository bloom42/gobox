// Get application directories such as config and cache.
package appdir

// Dirs requests application directories paths.
type Dirs interface {
	// Get the user-specific config directory.
	UserConfig() (string, error)
	// Get the user-specific cache directory.
	UserCache() (string, error)
	// Get the user-specific logs directory.
	UserLogs() (string, error)
	// Get the user-specific data directory.
	UserData() (string, error)
}

// New creates a new App with the provided name.
func New(name string) Dirs {
	return &dirs{name: name}
}
