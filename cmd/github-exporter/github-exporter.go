package main

import (
	"github.com/clambin/github-exporter/internal/collector"
	"github.com/clambin/github-exporter/internal/github"
	"github.com/clambin/go-common/httpclient"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/exp/slog"
	"net/http"
	"os"
	"time"
)

var (
	configFilename string
	BuildVersion   = "change-me"
	cmd            = &cobra.Command{
		Use:     "mediamon",
		Short:   "Prometheus exporter for various media applications. Currently supports Transmission, OpenVPN Client, Sonarr, Radarr and Plex.",
		Run:     Main,
		Version: BuildVersion,
	}
)

func main() {
	if viper.GetBool("debug") {
		slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug})))
	}

	if err := cmd.Execute(); err != nil {
		slog.Error("failed to start", "err", err)
		os.Exit(1)
	}
}

func Main(_ *cobra.Command, _ []string) {
	c := collector.Collector{
		Users: viper.GetStringSlice("repos.user"),
		Repos: viper.GetStringSlice("repos.repo"),
		Client: github.Client{
			HTTPClient: &http.Client{
				Transport: httpclient.NewRoundTripper(httpclient.WithCache(httpclient.DefaultCacheTable, time.Hour, 24*time.Hour)),
				Timeout:   10 * time.Second,
			},
			Token: viper.GetString("git.token"),
		},
		IncludeArchived: viper.GetBool("repos.archived"),
	}

	prometheus.MustRegister(c)

	http.Handle("/metrics", promhttp.Handler())
	_ = http.ListenAndServe(":9090", nil)
}

func init() {
	cobra.OnInitialize(initConfig)
	cmd.Flags().StringVar(&configFilename, "config", "", "Configuration file")
	cmd.Flags().Bool("debug", false, "Log debug messages")
	_ = viper.BindPFlag("debug", cmd.Flags().Lookup("debug"))
}

func initConfig() {
	if configFilename != "" {
		viper.SetConfigFile(configFilename)
	} else {
		viper.AddConfigPath("/etc/github-exporter/")
		viper.AddConfigPath("$HOME/.github-exporter")
		viper.AddConfigPath(".")
		viper.SetConfigName("config")
	}

	viper.SetDefault("debug", false)
	viper.SetDefault("repos.user", []string{})
	viper.SetDefault("repos.repo", []string{})
	viper.SetDefault("repos.archived", false)
	viper.SetDefault("git.token", "")

	viper.SetEnvPrefix("GITHUB_EXPORTER")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		slog.Error("failed to read config file", "err", err)
		os.Exit(1)
	}
}
