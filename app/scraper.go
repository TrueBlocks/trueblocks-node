package app

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
// file and sleep between each run for the duration specified with --sleep.
func (a *App) RunScraper(wg *sync.WaitGroup) {
	defer wg.Done()

	if a.InitMode != None {
		a.Logger.Info("Entering init mode", "mode", a.InitMode)

		reports := make([]*scraperReport, 0, len(a.Config.Targets))
		for _, chain := range a.Config.Targets {
			if rep, err := a.initOneChain(chain); err != nil {
				if !strings.HasPrefix(err.Error(), "no record found in the Unchained Index") {
					a.Logger.Error("Error", "err", err)
				} else {
					a.Logger.Warn("No record found in the Unchained Index for chain", "chain", chain)
				}
			} else {
				reports = append(reports, rep)
			}
		}

		for _, report := range reports {
			a.ReportOneScrape(report)
		}
	}

	a.Logger.Info("Entering scraper loop", "sleep", a.Sleep, "targets", a.Config.Targets)
	time.Sleep(2 * time.Second)

	runCount := 0
	for {
		caughtUp := true
		msg := []any{"sleep", a.Sleep}
		for _, chain := range a.Config.Targets {
			if report, err := a.scrapeOneChain(chain); err != nil {
				a.Logger.Error("ScrapeRunOnce failed", "error", err)
				time.Sleep(1 * time.Second)

			} else {
				msg = append(msg, report.Chain, -report.Staged)
				// TODO: This should be per-chain from the config file
				if report.Staged > (28 + 4) {
					caughtUp = false
				}
				// a.ReportOneScrape(report)
			}
		}

		if caughtUp {
			if runCount%5 == 0 || a.Sleep > 10 {
				fmt.Fprintf(os.Stderr, "%s\r", strings.Repeat(" ", 120))
				a.Logger.Info("caught up", msg...)
			}
			runCount++
			time.Sleep(time.Duration(a.Sleep) * time.Second)
		} else {
			time.Sleep(1 * time.Second)
		}
	}
}

type scraperReport struct {
	Chain     string `json:"chain"`
	BlockCnt  int    `json:"blockCnt"`
	Head      int    `json:"head"`
	Unripe    int    `json:"unripe"`
	Staged    int    `json:"staged"`
	Finalized int    `json:"finalized"`
	Time      string `json:"time"`
}

func (r *scraperReport) String() string {
	bytes, _ := json.Marshal(r)
	return string(bytes)
}

func newScraperReportFromMeta(meta *types.MetaData, chain string, blockCnt int) *scraperReport {
	return &scraperReport{
		Chain:     chain,
		BlockCnt:  blockCnt,
		Head:      int(meta.Latest),
		Unripe:    int(meta.Latest) - int(meta.Unripe),
		Staged:    int(meta.Latest) - int(meta.Staging),
		Finalized: int(meta.Latest) - int(meta.Finalized),
		Time:      time.Now().Format("01-02 15:04:05"),
	}
}

func (a *App) ReportOneScrape(report *scraperReport) {
	msg := fmt.Sprintf("behind (% 10.10s)...", report.Chain)
	if report.Staged < 30 {
		msg = fmt.Sprintf("atHead (% 10.10s)...", report.Chain)
	}
	a.Logger.Info(msg,
		"head", report.Head,
		"unripe", -report.Unripe,
		"staged", -report.Staged,
		"finalized", -report.Finalized,
		"blockCnt", report.BlockCnt,
	)
}

func (a *App) initOneChain(chain string) (*scraperReport, error) {
	a.Logger.Info("For chain", "chain", chain)

	originalHandler := a.Logger.Handler()
	defer func() {
		logger.SetLoggerWriter(io.Discard)
		a.Logger = slog.New(originalHandler)
		os.Setenv("TB_NODE_HEADLESS", "")
	}()
	a.Logger = slog.New(slog.NewTextHandler(nil, nil))
	logger.SetLoggerWriter(os.Stderr)
	os.Setenv("TB_NODE_HEADLESS", "true")

	opts := sdk.InitOptions{
		Globals: sdk.Globals{
			Chain: chain,
		},
	}

	var err error
	var meta *types.MetaData
	if a.InitMode == All {
		_, meta, err = opts.InitAll()
	} else if a.InitMode == Blooms {
		_, meta, err = opts.Init()
	}
	if err != nil {
		return nil, err
	}

	return newScraperReportFromMeta(meta, chain, a.BlockCnt), nil
}

// ----------------------------------------------------------------------------------
func (a *App) scrapeOneChain(chain string) (*scraperReport, error) {
	originalHandler := a.Logger.Handler()
	defer func() {
		logger.SetLoggerWriter(io.Discard)
		a.Logger = slog.New(originalHandler)
		os.Setenv("TB_NODE_HEADLESS", "")
	}()
	a.Logger = slog.New(slog.NewTextHandler(nil, nil))
	logger.SetLoggerWriter(os.Stderr)
	os.Setenv("TB_NODE_HEADLESS", "true")

	opts := sdk.ScrapeOptions{
		BlockCnt: uint64(a.BlockCnt),
		Globals: sdk.Globals{
			Chain: chain,
		},
	}

	if msg, meta, err := opts.ScrapeRunOnce(); err != nil {
		return nil, err
	} else {
		if len(msg) > 0 {
			a.Logger.Info(msg[0].String())
		}
		return newScraperReportFromMeta(meta, chain, a.BlockCnt), nil
	}
}
