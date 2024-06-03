package main

import (
	"github.com/clambin/github-exporter/internal/collector"
	"github.com/clambin/github-exporter/internal/stats"
	ghc "github.com/clambin/github-exporter/internal/stats/github"
	"github.com/clambin/go-common/http/metrics"
	"github.com/clambin/go-common/http/roundtripper"
	"github.com/google/go-github/v62/github"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"log/slog"
	"net/http"
	"os"
	"time"
)

var (
	configFilename string
	version        = "change-me"
	cmd            = &cobra.Command{
		Use:     "github-exporter",
		Short:   "Prometheus exporter for GitHub repositories",
		Run:     Main,
		Version: version,
	}
)

func main() {
	if err := cmd.Execute(); err != nil {
		slog.Error("failed to start", "err", err)
		os.Exit(1)
	}
}

func Main(cmd *cobra.Command, _ []string) {
	var opts slog.HandlerOptions

	if viper.GetBool("debug") {
		opts.Level = slog.LevelDebug
	}
	logger := slog.New(slog.NewJSONHandler(os.Stderr, &opts))

	logger.Info(cmd.Name()+" started", "version", cmd.Version, "cache", viper.GetDuration("git.cache"))

	ctx := cmd.Context()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: viper.GetString("git.token")},
	)
	tc := oauth2.NewClient(ctx, ts)

	rm := metrics.NewRequestMetrics(metrics.Options{Namespace: "github", Subsystem: "exporter"})
	im1 := metrics.NewInflightMetrics("github", "exporter", map[string]string{"stage": "pre"})
	im2 := metrics.NewInflightMetrics("github", "exporter", map[string]string{"stage": "post"})
	prometheus.MustRegister(rm, im1, im2)

	tp := roundtripper.New(
		roundtripper.WithInflightMetrics(im1),
		roundtripper.WithLimiter(25),
		roundtripper.WithInflightMetrics(im2),
		roundtripper.WithRequestMetrics(rm),
		roundtripper.WithRoundTripper(tc.Transport),
	)

	c := collector.Collector{
		Client: stats.Client{
			GitHubClient: (*ghc.Client)(github.NewClient(&http.Client{Transport: tp, Timeout: 10 * time.Second})),
			Logger:       logger.With("component", "github"),
		},
		Users:           viper.GetStringSlice("repos.user"),
		Repos:           viper.GetStringSlice("repos.repo"),
		IncludeArchived: viper.GetBool("repos.archived"),
		Lifetime:        viper.GetDuration("git.cache"),
		Logger:          logger.With("component", "collector"),
	}
	prometheus.MustRegister(&c)

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
