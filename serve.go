package main

import (
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/TrueBlocks/trueblocks-core/sdk/v3"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/colors"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/logger"
)

func (a *App) serve() {
	ready := make(chan bool)
	go sdk.NewDaemon(getApiUrl()).Start(ready)
	<-ready
	logger.Info(colors.Yellow + "API server started..." + colors.Off)
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
