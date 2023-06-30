package main

import (
	"github.com/clambin/github-exporter/internal/collector"
	"github.com/clambin/github-exporter/internal/collector/client"
	"github.com/clambin/github-exporter/internal/collector/limiter"
	"github.com/clambin/go-common/httpclient"
	"github.com/google/go-github/v53/github"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/exp/slog"
	"golang.org/x/oauth2"
	"net/http"
	"os"
	"time"
)

var (
	configFilename string
	BuildVersion   = "change-me"
	cmd            = &cobra.Command{
		Use:     "github-exporter",
		Short:   "Prometheus exporter for GitHub repositories",
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

func Main(cmd *cobra.Command, _ []string) {
	slog.Info(cmd.Name()+" started", "version", cmd.Version, "cache", viper.GetDuration("git.cache"))

	ctx := cmd.Context()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: viper.GetString("git.token")},
	)
	tc := oauth2.NewClient(ctx, ts)

	tp := httpclient.NewRoundTripper(
		httpclient.WithMetrics("github", "monitor", "github-exporter"),
		httpclient.WithRoundTripper(tc.Transport),
	)

	l := limiter.New(25, tp, "github", "monitor", "github-exporter")

	c := collector.Collector{
		GitHubCache: collector.GitHubCache{
			Client: &client.Client{
				Client: github.NewClient(&http.Client{
					Transport: l,
					Timeout:   10 * time.Second},
				),
			},
			Users:           viper.GetStringSlice("repos.user"),
			Repos:           viper.GetStringSlice("repos.repo"),
			IncludeArchived: viper.GetBool("repos.archived"),
			Lifetime:        viper.GetDuration("git.cache"),
		},
	}

	prometheus.MustRegister(tp)
	prometheus.MustRegister(&c)
	prometheus.MustRegister(l)

	http.Handle("/metrics", promhttp.Handler())
	_ = http.ListenAndServe(viper.GetString("addr"), nil)
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
	viper.SetDefault("addr", ":9090")
	viper.SetDefault("repos.user", []string{})
	viper.SetDefault("repos.repo", []string{})
	viper.SetDefault("repos.archived", false)
	viper.SetDefault("git.token", "")
	viper.SetDefault("git.cache", time.Hour)

	viper.SetEnvPrefix("GITHUB_EXPORTER")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		slog.Error("failed to read config file", "err", err)
		os.Exit(1)
	}
}
