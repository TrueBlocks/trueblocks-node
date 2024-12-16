package utils

import (
	"os"
)

// EstablishFolder ensures that the given folder exists. If any folders in the path must
// be created, they will be created with 0755 permissions. Returns an error or nil on sucess.
func EstablishFolder(rootPath string) error {
	_, err := os.Stat(rootPath)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(rootPath, 0755)
			if err != nil {
				return err
			}
		} else {
			// If there's an error other than not exist...we fail
			return err
		}
	}
	return nil
}

// FileExists returns true if the file exists and is not a folder, false otherwise.
func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
