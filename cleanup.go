package graceful

var cleanupFuncs = []func(){}

// AddCleanup function to be executed when shutting down
func AddCleanup(f func()) {
	cleanupFuncs = append(cleanupFuncs, f)
}

// Cleanup executes all registered cleanup functions
func Cleanup() {
	for _, f := range cleanupFuncs {
		f()
	}
}
