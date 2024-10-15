package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/k3env/frm-exporter/config"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func recordMetrics() {
	hc := http.Client{}
	go func() {
		for {
			slog.Info("Start scraping")
			res, err := hc.Get(fmt.Sprintf("%s/getPower", cfg.FRM.Url))
			if err != nil {
				continue
			}
			bs, err := io.ReadAll(res.Body)
			if err != nil {
				continue
			}
			var pwinfo []PowerInfo
			err = json.Unmarshal(bs, &pwinfo)
			if err != nil {
				continue
			}

			for _, info := range pwinfo {
				powerProduced.WithLabelValues(fmt.Sprintf("Circuit #%d", info.CircuitID)).Set(info.PowerProduction)
				powerConsumed.WithLabelValues(fmt.Sprintf("Circuit #%d", info.CircuitID)).Set(info.PowerConsumed)
				powerMaxConsumed.WithLabelValues(fmt.Sprintf("Circuit #%d", info.CircuitID)).Set(info.PowerMaxConsumed)
				powerCapacity.WithLabelValues(fmt.Sprintf("Circuit #%d", info.CircuitID)).Set(info.PowerCapacity)
			}
			slog.Info("Done")
			time.Sleep(cfg.FRM.ScrapeInterval)
		}
	}()
}

var (
	powerProduced    = promauto.NewGaugeVec(prometheus.GaugeOpts{Name: "frm_power_produced"}, []string{"circuit"})
	powerConsumed    = promauto.NewGaugeVec(prometheus.GaugeOpts{Name: "frm_power_consumed"}, []string{"circuit"})
	powerMaxConsumed = promauto.NewGaugeVec(prometheus.GaugeOpts{Name: "frm_power_max_consumed"}, []string{"circuit"})
	powerCapacity    = promauto.NewGaugeVec(prometheus.GaugeOpts{Name: "frm_power_capacity"}, []string{"circuit"})
	cfg              *config.Config
)

func main() {
	var err error
	var path string

	flag.StringVar(&path, "config", "config.yaml", "Path to config file")
	flag.Parse()

	cfg, err = config.LoadConfig(path)
	if err != nil {
		slog.Error("Error loading config:", slog.String("error", err.Error()))
		os.Exit(1)
	}
	slog.Info("Scraping FRM from server", slog.String("frm-url", cfg.FRM.Url))
	recordMetrics()

	slog.Info("Starting server", slog.String("addr", cfg.Web.Address))
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(cfg.Web.Address, nil)
}
