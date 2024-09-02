package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"path/filepath"
	"sync"
	"text/template"
	"time"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/colors"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/logger"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
	sdk "github.com/TrueBlocks/trueblocks-sdk/v3"
)

func scraper(config *Config, wg *sync.WaitGroup) {
	defer wg.Done()

	opts := sdk.InitOptions{}
	if _, _, err := opts.Init(); err != nil { // blooms only, if that fails
		// if _, _, err := opts.InitAll(); err != nil { // try --all
		logger.Error(err)
		return
		// }
	}

	dataFilename := filepath.Join(config.OutputPath, "meta.json")
	for {
		screenMutex.Lock()
		fmt.Print(colors.Green, "Scraping...", colors.Off)
		quit := false
		go func() {
			for {
				if quit {
					break
				}
				time.Sleep(time.Millisecond * 1000)
				fmt.Print(".")
			}
		}()
		wwg := sync.WaitGroup{}
		wwg.Add(1)
		go scrapeOnce(dataFilename, &wwg)
		wwg.Wait()
		quit = true
		fmt.Println(colors.Green, "Done.", colors.Off)
		time.Sleep(time.Millisecond * 1000)
		fmt.Print("\r \r")
		screenMutex.Unlock()
		time.Sleep(time.Millisecond * 4000)
	}
}

func scrapeOnce(dataFilename string, wwg *sync.WaitGroup) {
	defer wwg.Done()

	opts := sdk.ScrapeOptions{
		BlockCnt: 500,
		Globals: sdk.Globals{
			Chain: "mainnet",
		},
	}

	w := logger.GetLoggerWriter()
	logger.SetLoggerWriter(io.Discard)

	if _, meta, err := opts.ScrapeRunOnce(); err != nil {
		logger.Error(err)
	} else {
		tmpl := `Head (H): {{.Head}}
Unripe:    H - {{.Unripe}}
Staged:    H - {{.Staged}}
Finalized: H - {{.Finalized}}
{{.Time}}
`

		t, err := template.New("myTemplate").Parse(tmpl)
		if err != nil {
			panic(err)
		}
		var buf bytes.Buffer
		report := NewReportFromMeta(meta)
		err = t.Execute(&buf, report)
		if err != nil {
			panic(err)
		}
		// fmt.Println("\n" + buf.String())
		file.StringToAsciiFile(dataFilename, buf.String())
	}

	logger.SetLoggerWriter(w)
}

type Report struct {
	Head      int    `json:"head"`
	Unripe    int    `json:"unripe"`
	Staged    int    `json:"staged"`
	Finalized int    `json:"finalized"`
	Time      string `json:"time"`
}

func (r *Report) String() string {
	bytes, _ := json.Marshal(r)
	return string(bytes)
}

func NewReportFromMeta(meta *types.MetaData) *Report {
	return &Report{
		Head:      int(meta.Latest),
		Unripe:    int(meta.Latest) - int(meta.Unripe),
		Staged:    int(meta.Latest) - int(meta.Staging),
		Finalized: int(meta.Latest) - int(meta.Finalized),
		Time:      time.Now().Format("01-02 15:04:05"),
	}
}
