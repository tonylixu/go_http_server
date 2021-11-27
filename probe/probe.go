package probe

import "os"

const liveFile = "/tmp/live"

// Create will create a file for the liveness check
func Create() error {
	_, err := os.Create(liveFile)
	return err
}

// Remove will remove the file created for liveness probe
func Remove() error {
	return os.Remove(liveFile)
}

// Exists checks if the file created for liveness probe exists
func Exists() bool {
	if _, err := os.Stat(liveFile); err == nil {
		return true
	}
	return false
}
