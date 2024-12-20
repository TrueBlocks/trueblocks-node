package app

import (
	"fmt"
	"os"
	"strings"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/rpc"
	sdk "github.com/TrueBlocks/trueblocks-sdk/v4"
)

// RunServer runs the API server in a goroutine, waits for it to be ready and then
// returns. All available trueblocks-core endpoints are available.
func (a *App) RunServer() (apiUrl string, err error) {
	apiUrl = getApiUrl()
	ready := make(chan bool)
	go sdk.NewDaemon(apiUrl).Start(ready)
	<-ready
	return
}

// apiPort is a global variable to store the port of the API server
// of which there can be only one.
var apiPort = ""

func getApiUrl() string {
	// use the same port if we've already started the server...
	if apiPort != "" {
		return "localhost:" + apiPort
	}

	apiPort = strings.ReplaceAll(os.Getenv("TB_TEST_API_SERVER"), ":", "")
	if apiPort == "" {
		apiPort = fmt.Sprintf("%d", rpc.FindAvailablePort())
	}

	return "localhost:" + apiPort
}
