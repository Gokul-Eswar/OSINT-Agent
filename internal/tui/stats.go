package tui

import (
	"fmt"
	"runtime"

	"github.com/spectre/spectre/internal/storage"
)

func GetSystemStats() string {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	// Fetch some DB stats
	cases, _ := storage.ListCases()
	
	return fmt.Sprintf(
		"SPECTRE System Status\n\n"+
		"OS:      %s\n"+
		"Arch:    %s\n"+
		"Memory:  %v MB\n\n"+
		"Database Stats:\n"+
		"- Active Cases: %d\n"+
		"- Storage:      Local SQLite\n\n"+
		"(esc: back)",
		runtime.GOOS,
		runtime.GOARCH,
		m.Alloc/1024/1024,
		len(cases),
	)
}

