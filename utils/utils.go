package utils

import (
	"bytes"
	"fmt"
	"net/http"
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

// PingRpc sends a ping request to the RPC provider, returns an error or nil on success.
func PingRpc(providerUrl string) error {
	jsonData := []byte(`{ "jsonrpc": "2.0", "method": "web3_clientVersion", "id": 6 }`)
	req, err := http.NewRequest("POST", providerUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	clientHTTP := &http.Client{}
	resp, err := clientHTTP.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
