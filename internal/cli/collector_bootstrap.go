package cli

import (
	"fmt"
	"sync"

	"github.com/spectre/spectre/internal/collector"
	activecollector "github.com/spectre/spectre/internal/collector/active"
	dnscollector "github.com/spectre/spectre/internal/collector/dns"
	geocollector "github.com/spectre/spectre/internal/collector/geo"
	githubcollector "github.com/spectre/spectre/internal/collector/github"
	whoiscollector "github.com/spectre/spectre/internal/collector/whois"
)

var collectorBootstrapOnce sync.Once
var collectorBootstrapErr error

func ensureCollectorBootstrap() error {
	collectorBootstrapOnce.Do(func() {
		activecollector.RegisterBuiltins()
		dnscollector.Register()
		geocollector.Register()
		githubcollector.Register()
		whoiscollector.Register()

		plugins, err := collector.DiscoverPlugins()
		if err != nil {
			collectorBootstrapErr = fmt.Errorf("failed to discover external plugins: %w", err)
			return
		}

		for _, p := range plugins {
			collector.Register(p)
		}
	})

	return collectorBootstrapErr
}
