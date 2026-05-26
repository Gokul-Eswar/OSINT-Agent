package active

import "github.com/spectre/spectre/internal/collector"

func RegisterBuiltins() {
	collector.Register(&HTTPCollector{})
	collector.Register(&PortCollector{})
	collector.Register(&ScreenshotCollector{})
	collector.Register(NewSocialCollector())
}
