package homedir

import (
	"os"
	"sync"

	"github.com/mitchellh/go-homedir"
)

const defaultPath = "/usr/bin:/bin:/usr/sbin:/sbin"

var pathLock sync.Mutex

// Dir returns the current user's home directory using a method that works for
// cross-compiled binaries with CGO disabled.
func Dir() (string, error) {
	// Call homedir.Dir() first to see if we can resolve the home directory
	// without any special workarounds. This also allows us to avoid applying
	// workarounds after the home directory is cached by homedir.Dir().
	d, err := homedir.Dir()
	if err == nil {
		return d, err
	}

	// Our first attempt failed, let's check if $PATH is empty and if so try to
	// work around the issue where the fallback to lookup the user's home
	// directory by shelling out fails when $PATH is empty and try again. See
	// ADA-13055 for a customer reported instance of this.
	//
	// We take the pathLock to ensure multiple invocations aren't changing $PATH
	// out from underneath each other (this does not protect us against
	// modifications to $PATH elsewhere).
	pathLock.Lock()
	defer pathLock.Unlock()
	if path := os.Getenv("PATH"); path == "" {
		defer os.Setenv("PATH", path)
		os.Setenv("PATH", defaultPath)
		return homedir.Dir()
	}
	return d, err
}
