package report

import (
	"fmt"
	"time"

	"github.com/jung-kurt/gofpdf"
	"github.com/spectre/spectre/internal/core"
	"github.com/spectre/spectre/internal/storage"
)

// GeneratePDFReport creates a professional PDF report for the case.
func GeneratePDFReport(caseID string) (string, error) {
	c, err := storage.GetCase(caseID)
	if err != nil {
		return "", err
	}
	if c == nil {
		return "", fmt.Errorf("case not found")
	}

	entities, _ := storage.ListEntitiesByCase(caseID)
	// rels, _ := storage.ListRelationshipsByCase(caseID)
	timeline, _ := storage.GetCaseTimeline(caseID)
	analysis, _ := storage.GetLatestAnalysis(caseID)
	evidence, _ := storage.ListEvidenceByCase(caseID)

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetHeaderFunc(func() {
		pdf.SetFont("Arial", "I", 8)
		pdf.CellFormat(0, 10, fmt.Sprintf("SPECTRE Intelligence Report - Case %s", caseID), "", 0, "R", false, 0, "")
		pdf.Ln(10)
	})
	pdf.SetFooterFunc(func() {
		pdf.SetY(-15)
		pdf.SetFont("Arial", "I", 8)
		pdf.CellFormat(0, 10, fmt.Sprintf("Page %d", pdf.PageNo()), "", 0, "C", false, 0, "")
	})

	// --- Cover Page ---
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 24)
	pdf.CellFormat(0, 60, "", "", 1, "", false, 0, "") // Spacer
	pdf.CellFormat(0, 10, "INTELLIGENCE REPORT", "", 1, "C", false, 0, "")
	
	pdf.SetFont("Arial", "", 16)
	pdf.CellFormat(0, 10, fmt.Sprintf("Case: %s", c.Name), "", 1, "C", false, 0, "")
	
	pdf.SetFont("Arial", "", 12)
	pdf.CellFormat(0, 10, fmt.Sprintf("ID: %s", c.ID), "", 1, "C", false, 0, "")
	pdf.CellFormat(0, 10, fmt.Sprintf("Date: %s", time.Now().Format("2006-01-02")), "", 1, "C", false, 0, "")
	
	pdf.Ln(50)
	pdf.SetFont("Courier", "", 10)
	pdf.MultiCell(0, 5, "CONFIDENTIAL // INTERNAL USE ONLY", "", "C", false)

	// --- Executive Summary ---
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(0, 10, "Executive Summary")
	pdf.Ln(15)

	pdf.SetFont("Arial", "", 11)
	if analysis != nil {
		pdf.MultiCell(0, 6, "AI Analysis provided the following insights based on collected evidence:", "", "L", false)
		pdf.Ln(5)
		
		pdf.SetFont("Arial", "B", 12)
		pdf.Cell(0, 10, "Key Findings")
		pdf.Ln(8)
		pdf.SetFont("Arial", "", 11)
		for _, f := range analysis.Findings {
			pdf.MultiCell(0, 6, fmt.Sprintf("- %s", f), "", "L", false)
		}
		pdf.Ln(5)

		pdf.SetFont("Arial", "B", 12)
		pdf.Cell(0, 10, "Identified Risks")
		pdf.Ln(8)
		pdf.SetFont("Arial", "", 11)
		for _, r := range analysis.Risks {
			pdf.MultiCell(0, 6, fmt.Sprintf("- %s", r), "", "L", false)
		}
	} else {
		pdf.MultiCell(0, 6, "No AI analysis has been performed for this case yet.", "", "L", false)
	}

	// --- Geo-Intelligence ---
	geoEntities := []*core.Entity{}
	for _, e := range entities {
		if e.Metadata != nil && (e.Metadata["lat"] != nil || e.Metadata["country"] != nil) {
			geoEntities = append(geoEntities, e)
		}
	}

	if len(geoEntities) > 0 {
		pdf.AddPage()
		pdf.SetFont("Arial", "B", 16)
		pdf.Cell(0, 10, "Geo-Intelligence")
		pdf.Ln(15)

		pdf.SetFont("Arial", "B", 10)
		pdf.SetFillColor(240, 240, 240)
		pdf.CellFormat(60, 8, "Target", "1", 0, "", true, 0, "")
		pdf.CellFormat(130, 8, "Location / ISP", "1", 1, "", true, 0, "")

		pdf.SetFont("Arial", "", 9)
		for _, e := range geoEntities {
			pdf.CellFormat(60, 8, e.Value, "1", 0, "", false, 0, "")
			loc := fmt.Sprintf("%v, %v (%v)", e.Metadata["city"], e.Metadata["country"], e.Metadata["isp"])
			pdf.CellFormat(130, 8, loc, "1", 1, "", false, 0, "")
		}
	}

	// --- Visual Evidence ---
	screenshots := []*core.Evidence{}
	for _, ev := range evidence {
		if ev.Collector == "screenshot" {
			screenshots = append(screenshots, ev)
		}
	}

	if len(screenshots) > 0 {
		for _, s := range screenshots {
			pdf.AddPage()
			pdf.SetFont("Arial", "B", 16)
			pdf.Cell(0, 10, fmt.Sprintf("Visual Evidence: %s", s.Metadata["target"]))
			pdf.Ln(15)

			// Add Image (Scale to fit page width)
			opt := gofpdf.ImageOptions{
				ImageType: "PNG",
				ReadDpi:   true,
			}
			pdf.ImageOptions(s.FilePath, 10, 35, 190, 0, false, opt, 0, "")
			
			pdf.SetY(260)
			pdf.SetFont("Arial", "I", 8)
			pdf.Cell(0, 10, fmt.Sprintf("Collected At: %s", s.CollectedAt.Format("2006-01-02 15:04:05")))
		}
	}

	// --- Entities ---
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(0, 10, "Discovered Entities")
	pdf.Ln(15)

	// Table Header
	pdf.SetFont("Arial", "B", 10)
	pdf.SetFillColor(240, 240, 240)
	pdf.CellFormat(30, 8, "Type", "1", 0, "", true, 0, "")
	pdf.CellFormat(110, 8, "Value", "1", 0, "", true, 0, "")
	pdf.CellFormat(50, 8, "Source", "1", 1, "", true, 0, "")

	// Table Rows
	pdf.SetFont("Arial", "", 9)
	for _, e := range entities {
		pdf.CellFormat(30, 8, e.Type, "1", 0, "", false, 0, "")
		
		// Truncate value if too long
		val := e.Value
		if len(val) > 60 {
			val = val[:57] + "..."
		}
		pdf.CellFormat(110, 8, val, "1", 0, "", false, 0, "")
		pdf.CellFormat(50, 8, e.Source, "1", 1, "", false, 0, "")
	}

	// --- Timeline ---
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(0, 10, "Investigation Timeline")
	pdf.Ln(15)

	pdf.SetFont("Arial", "", 10)
	for _, ev := range timeline {
		ts := ev.Timestamp.Format("2006-01-02 15:04:05")
		pdf.SetFont("Arial", "B", 10)
		pdf.Cell(0, 6, fmt.Sprintf("%s - %s", ts, ev.Type))
		pdf.Ln(6)
		pdf.SetFont("Arial", "I", 9)
		pdf.MultiCell(0, 5, fmt.Sprintf("%s (Source: %s)", ev.Description, ev.Source), "", "L", false)
		pdf.Ln(4)
	}

	outfile := fmt.Sprintf("report_%s.pdf", caseID)
	err = pdf.OutputFileAndClose(outfile)
	if err != nil {
		return "", err
	}

	return outfile, nil
}
