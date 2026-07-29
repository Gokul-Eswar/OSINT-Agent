package report

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"html"
	"os"
	"strings"

	"github.com/Gokul-Eswar/Spectre/internal/storage"
)

// GenerateJSONReport exports case metadata, entities, and relationships into a structured JSON report.
func GenerateJSONReport(caseID string, outputPath string) error {
	c, err := storage.GetCase(caseID)
	if err != nil || c == nil {
		return fmt.Errorf("case not found: %w", err)
	}

	entities, _ := storage.ListEntitiesByCase(caseID)
	rels, _ := storage.ListRelationshipsByCase(caseID)
	evidence, _ := storage.ListEvidenceByCase(caseID)
	timeline, _ := storage.GetCaseTimeline(caseID)
	analysis, _ := storage.GetLatestAnalysis(caseID)

	reportData := map[string]interface{}{
		"case":          c,
		"entities":      entities,
		"relationships": rels,
		"evidence":      evidence,
		"timeline":      timeline,
		"analysis":      analysis,
	}

	bytes, err := json.MarshalIndent(reportData, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(outputPath, bytes, 0644)
}

// GenerateCSVReport exports discovered case entities into a CSV file.
func GenerateCSVReport(caseID string, outputPath string) error {
	entities, err := storage.ListEntitiesByCase(caseID)
	if err != nil {
		return fmt.Errorf("failed to list entities: %w", err)
	}

	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Header
	if err := writer.Write([]string{"ID", "Type", "Value", "Source", "DiscoveredAt"}); err != nil {
		return err
	}

	for _, e := range entities {
		row := []string{e.ID, e.Type, e.Value, e.Source, e.DiscoveredAt.Format("2006-01-02 15:04:05")}
		if err := writer.Write(row); err != nil {
			return err
		}
	}

	return nil
}

// GenerateHTMLReport creates a self-contained HTML report document for a case.
func GenerateHTMLReport(caseID string, outputPath string) error {
	c, err := storage.GetCase(caseID)
	if err != nil || c == nil {
		return fmt.Errorf("case not found")
	}

	entities, _ := storage.ListEntitiesByCase(caseID)
	rels, _ := storage.ListRelationshipsByCase(caseID)
	analysis, _ := storage.GetLatestAnalysis(caseID)

	var sb strings.Builder
	sb.WriteString("<!DOCTYPE html><html><head><meta charset='UTF-8'><title>SPECTRE Intelligence Report</title>")
	sb.WriteString("<style>body{font-family:sans-serif;margin:40px;background:#121212;color:#e0e0e0;}")
	sb.WriteString("h1,h2{color:#00e5ff;}table{width:100%;border-collapse:collapse;margin-top:10px;}")
	sb.WriteString("th,td{border:1px solid #333;padding:8px;text-align:left;}th{background:#1e1e1e;color:#00e5ff;}")
	sb.WriteString(".card{background:#1e1e1e;padding:15px;border-radius:6px;margin-bottom:20px;}</style></head><body>")

	sb.WriteString(fmt.Sprintf("<h1>SPECTRE Intelligence Report</h1><div class='card'><h2>Case: %s</h2><p>ID: %s</p></div>", html.EscapeString(c.Name), html.EscapeString(c.ID)))

	if analysis != nil {
		sb.WriteString("<div class='card'><h2>AI Insights</h2><h3>Findings</h3><ul>")
		for _, f := range analysis.Findings {
			sb.WriteString(fmt.Sprintf("<li>%s</li>", html.EscapeString(f)))
		}
		sb.WriteString("</ul><h3>Risks</h3><ul>")
		for _, r := range analysis.Risks {
			sb.WriteString(fmt.Sprintf("<li>%s</li>", html.EscapeString(r)))
		}
		sb.WriteString("</ul></div>")
	}

	sb.WriteString("<h2>Discovered Entities</h2><table><tr><th>Type</th><th>Value</th><th>Source</th></tr>")
	for _, e := range entities {
		sb.WriteString(fmt.Sprintf("<tr><td>%s</td><td>%s</td><td>%s</td></tr>", html.EscapeString(e.Type), html.EscapeString(e.Value), html.EscapeString(e.Source)))
	}
	sb.WriteString("</table>")

	if len(rels) > 0 {
		sb.WriteString("<h2>Relationships</h2><table><tr><th>From</th><th>Type</th><th>To</th></tr>")
		for _, r := range rels {
			sb.WriteString(fmt.Sprintf("<tr><td>%s</td><td>%s</td><td>%s</td></tr>", html.EscapeString(r.FromEntityID), html.EscapeString(r.Type), html.EscapeString(r.ToEntityID)))
		}
		sb.WriteString("</table>")
	}

	sb.WriteString("</body></html>")

	return os.WriteFile(outputPath, []byte(sb.String()), 0644)
}
