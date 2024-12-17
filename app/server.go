package app

import (
	"fmt"
	"net"
	"os"
	"strings"

	sdk "github.com/TrueBlocks/trueblocks-sdk/v4"
)

// RunServer runs the API server in a goroutine, waits for it to be ready and then
// returns. All available trueblocks-core endpoints are available.
func (a *App) RunServer() {
	ready := make(chan bool)
	go sdk.NewDaemon(getApiUrl()).Start(ready)
	<-ready
}

var apiPort = ""

func getApiUrl() string {
	if apiPort != "" {
		return "localhost:" + apiPort
	}

	apiPort = strings.ReplaceAll(os.Getenv("TB_TEST_API_SERVER"), ":", "")
	if apiPort == "" {
		apiPort = fmt.Sprintf("%d", findAvailablePort())
	}
	return "localhost:" + apiPort
}

func findAvailablePort() int {
	preferredPorts := []int{8080, 8088, 9090, 9099}
	for _, port := range preferredPorts {
		address := fmt.Sprintf(":%d", port)
		listener, err := net.Listen("tcp", address)
		if err == nil {
			defer listener.Close()
			return port
		}
	}
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		return 0
	}
	defer listener.Close()
	addr := listener.Addr().(*net.TCPAddr)
	return addr.Port
}
