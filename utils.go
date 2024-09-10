package main

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
)

// EstablishFolder creates folders given a list of folders
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

func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func pingRpc(providerUrl string) error {
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
