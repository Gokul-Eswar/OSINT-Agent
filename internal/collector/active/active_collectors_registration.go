package active

import "github.com/Gokul-Eswar/Spectre/internal/collector"

func RegisterBuiltins() {
	collector.Register(&HTTPCollector{})
	collector.Register(&PortCollector{})
	collector.Register(&ScreenshotCollector{})
	collector.Register(NewSocialCollector())
}
