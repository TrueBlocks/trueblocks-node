package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/logger"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
	sdk "github.com/TrueBlocks/trueblocks-sdk/v3"
)

// RunScraper runs the scraper in a goroutine. It will scrape the chains in the configuration
// file and sleep for the duration specified in the configuration file.
func (a *App) RunScraper(wg *sync.WaitGroup) {
	defer wg.Done()

	if a.InitMode != None {
		a.Logger.Debug("Entering init mode", "mode", a.InitMode)
		for _, chain := range a.Config.Targets {
			a.Logger.Debug("For chain", "chain", chain)
			opts := sdk.InitOptions{
				Globals: sdk.Globals{
					Chain: chain,
				},
			}

			originalHandler := a.Logger.Handler()
			a.Logger = slog.New(slog.NewTextHandler(nil, nil))
			logger.SetLoggerWriter(os.Stderr)
			var err error
			if a.InitMode == All {
				_, _, err = opts.InitAll()
			} else if a.InitMode == Blooms {
				_, _, err = opts.Init()
			}
			logger.SetLoggerWriter(io.Discard)
			a.Logger = slog.New(originalHandler)

			if err != nil {
				a.Logger.Error("", "error", err)
				if !strings.HasPrefix(err.Error(), "no record found in the Unchained Index") {
					return
				} else {
					a.Logger.Warn("No record found in the Unchained Index for chain", "chain", chain)
				}
			}
		}
	}

	a.Logger.Info("Entering scrape loop: ", "sleep", a.Sleep)
	time.Sleep(2 * time.Second)

	for {
		a.Logger.Debug("Entering scraper loop")
		for _, chain := range a.Config.Targets {
			a.Logger.Debug("For chain", "chain", chain)
			time.Sleep(1 * time.Second)
			if report, err := a.scrapeOnce(chain); err != nil {
				a.Logger.Error("ScrapeRunOnce failed", "error", err)
				time.Sleep(time.Duration(a.Sleep) * time.Second)

			} else {
				a.Logger.Info("In the loop", "sleep", a.Sleep)
				time.Sleep(20 * time.Second)
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
				} else {
					time.Sleep(1 * time.Second)
				}
			}
		}
	}
}

func (a *App) scrapeOnce(chain string) (*Report, error) {
	// TODO: Allow user to specify block_cnt
	blockCnt := 30
	opts := sdk.ScrapeOptions{
		BlockCnt: uint64(blockCnt),
		Globals: sdk.Globals{
			Chain: chain,
		},
	}

	fmt.Fprintf(os.Stderr, "Scraping pass %s (%d blocks)...                \r", chain, blockCnt)
	if msg, meta, err := opts.ScrapeRunOnce(); err != nil {
		return nil, err
	} else {
		a.Logger.Info(msg[0].String())
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
