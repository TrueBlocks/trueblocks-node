package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"text/template"
	"time"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/colors"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/logger"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/walk"
	sdk "github.com/TrueBlocks/trueblocks-sdk/v3"
)

func (a *App) scrape(wg *sync.WaitGroup) {
	defer wg.Done()

	a.Busy = true
	go func() {
		for a.Busy {
			cnt := 0
			walk.ForEveryFileInFolder(a.Config.IndexPath(), func(path string, vP any) (bool, error) {
				cnt++
				return true, nil
			}, nil)
			fmt.Fprintf(os.Stderr, "Synced % 5d files\r", cnt)
			time.Sleep(time.Millisecond * 3000)
		}
	}()

	opts := sdk.InitOptions{}
	if _, _, err := opts.Init(); err != nil { // blooms only, if that fails
		if _, _, err := opts.InitAll(); err != nil { // try --all
			a.Logger.Error("", "error", err)
			return
		}
	}
	a.Busy = false

	dataFilename := filepath.Join(a.Config.ConfigPath, "scraper.report")
	a.Logger.Info("Scraping...", "fn", dataFilename, "config", a.Config.String())

	for {
		a.Busy = true
		fmt.Print(colors.Green, "Scraping...", colors.Off)
		a.Busy = false
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
		go a.scrapeOnce(dataFilename, &wwg)
		wwg.Wait()
		quit = true
		fmt.Println(colors.Green, "Done.", colors.Off)
		time.Sleep(time.Millisecond * 1000)
		fmt.Print("\r \r")
		time.Sleep(time.Millisecond * 4000)
	}
}

func (a *App) scrapeOnce(dataFilename string, wwg *sync.WaitGroup) {
	defer wwg.Done()

	opts := sdk.ScrapeOptions{
		BlockCnt: 100,
		Globals: sdk.Globals{
			Chain: a.Config.DefaultChain,
		},
	}

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
