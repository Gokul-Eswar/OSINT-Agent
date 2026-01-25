package tui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spectre/spectre/internal/analysis"
	"github.com/spectre/spectre/internal/core"
)

type AnalysisStepMsg int
type AnalysisFinishedMsg struct {
	Result *core.Analysis
}
type AnalysisErrorMsg string
type TickMsg time.Time

func tickCmd() tea.Cmd {
	return tea.Tick(800*time.Millisecond, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

func StartAnalysis(caseID string, modelName string) tea.Cmd {
	return tickCmd()
}

func PerformActualAnalysis(caseID string, modelName string) tea.Cmd {
	return func() tea.Msg {
		res, err := analysis.AnalyzeCase(caseID, modelName)
		if err != nil {
			return AnalysisErrorMsg(err.Error())
		}
		return AnalysisFinishedMsg{Result: res}
	}
}

func FormatAnalysis(res *core.Analysis) string {
	if res == nil {
		return ""
	}

	var s strings.Builder
	s.WriteString(StyleHeader.Render("PRELIMINARY FINDINGS") + "\n")
	for _, f := range res.Findings {
		s.WriteString(fmt.Sprintf(" • %s\n", f))
	}

	s.WriteString("\n" + StyleHeader.Render("IDENTIFIED RISKS") + "\n")
	for _, r := range res.Risks {
		s.WriteString(fmt.Sprintf(" ⚠ %s\n", r))
	}

	s.WriteString("\n" + StyleHeader.Render("RECOMMENDED NEXT STEPS") + "\n")
	for _, n := range res.NextSteps {
		s.WriteString(fmt.Sprintf(" → %s\n", n))
	}

	s.WriteString(fmt.Sprintf("\nConfidence Level: %.2f", res.Confidence))
	return s.String()
}