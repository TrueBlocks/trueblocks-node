package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/logger"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
	sdk "github.com/TrueBlocks/trueblocks-sdk/v3"
)

func (a *App) scrape(wg *sync.WaitGroup) {
	defer wg.Done()

	if a.InitMode == None {
		return
	}

	for _, chain := range a.Config.Targets {
		opts := sdk.InitOptions{
			Globals: sdk.Globals{
				Chain: chain,
			},
		}

		// a.Busy = true
		// go func() {
		// 	for a.Busy {
		// 		cnt := 0
		// 		path := filepath.Join(a.Config.IndexPath(), chain)
		// 		walk.ForEveryFileInFolder(path, func(path string, vP any) (bool, error) {
		// 			cnt++
		// 			return true, nil
		// 		}, nil)
		// 		fmt.Fprintf(os.Stderr, "Initializing %s: downloaded % 5d files\r", chain, cnt)
		// 		time.Sleep(2 * time.Second) // reporting speed
		// 	}
		// }()

		a.Busy = true
		logger.SetLoggerWriter(os.Stderr)
		var err error
		if a.InitMode == All {
			_, _, err = opts.InitAll()
		} else if a.InitMode == Blooms {
			_, _, err = opts.Init()
		}
		logger.SetLoggerWriter(io.Discard)
		a.Busy = false

		if err != nil {
			a.Logger.Error("", "error", err)
			if !strings.HasPrefix(err.Error(), "no record found in the Unchained Index") {
				return
			} else {
				a.Logger.Warn("No record found in the Unchained Index for chain", "chain", chain)
			}
		}

		// a.Busy = false
	}

	for {
		for _, chain := range a.Config.Targets {
			if report, err := a.scrapeOnce(chain); err != nil {
				a.Logger.Error("ScrapeRunOnce failed", "error", err)
				time.Sleep(time.Duration(a.Sleep) * time.Second)

			} else {
				caughtUp := report.Staged < 30
				msg := fmt.Sprintf("Behind (%s)...", report.Chain)
				sMsg := fmt.Sprintf("%d secs", 0)
				if caughtUp {
					msg = fmt.Sprintf("Caught up (%s)...", report.Chain)
					sMsg = fmt.Sprintf("%d secs", a.Sleep)
				}
				a.Logger.Info(msg,
					"head", fmt.Sprintf("% 9d", report.Head),
					"unripe", -report.Unripe,
					"staged", -report.Staged,
					"index", -report.Finalized,
					"blockCnt", report.BlockCnt,
					"sleep", sMsg,
				)
				if caughtUp {
					time.Sleep(time.Duration(a.Sleep) * time.Second)
				}
			}
		}
	}
}

func (a *App) scrapeOnce(chain string) (*Report, error) {
	// TODO: Allow user to specify block_cnt
	blockCnt := 100
	opts := sdk.ScrapeOptions{
		BlockCnt: uint64(blockCnt),
		Globals: sdk.Globals{
			Chain: chain,
		},
	}

	fmt.Fprintf(os.Stderr, "Scraping pass %s (%d blocks)...                \r", chain, blockCnt)
	if _, meta, err := opts.ScrapeRunOnce(); err != nil {
		return nil, err
	} else {
		return NewReportFromMeta(meta, chain, blockCnt), nil
	}
}

type Report struct {
	Chain     string `json:"chain"`
	BlockCnt  int    `json:"blockCnt"`
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

func NewReportFromMeta(meta *types.MetaData, chain string, blockCnt int) *Report {
	return &Report{
		Chain:     chain,
		BlockCnt:  blockCnt,
		Head:      int(meta.Latest),
		Unripe:    int(meta.Latest) - int(meta.Unripe),
		Staged:    int(meta.Latest) - int(meta.Staging),
		Finalized: int(meta.Latest) - int(meta.Finalized),
		Time:      time.Now().Format("01-02 15:04:05"),
	}
}
