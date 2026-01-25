package tui

import (
	"fmt"
	"strings"

	"github.com/spectre/spectre/internal/storage"
)

func RenderTimeline(caseID string) string {
	if caseID == "" {
		return "No case selected."
	}

	events, err := storage.GetCaseTimeline(caseID)
	if err != nil {
		return "Error fetching timeline: " + err.Error()
	}

	if len(events) == 0 {
		return "Timeline is empty. Collect some evidence first."
	}

	var s strings.Builder
	s.WriteString(fmt.Sprintf("TIMELINE — %s\n", caseID))
	s.WriteString("────────────────────────────────────\n\n")

	for _, e := range events {
		timeStr := e.Timestamp.Format("15:04:05")
		typeStr := StyleMuted.Render(fmt.Sprintf("[%s]", e.Type))
		s.WriteString(fmt.Sprintf("%s  %s  %s\n", timeStr, typeStr, e.Description))
	}

	return s.String()
}
