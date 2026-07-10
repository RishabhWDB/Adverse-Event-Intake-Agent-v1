package classifier

import (
	"strings"

	"github.com/RishabhWDB/Adverse-Event-Intake-Agent-v1/internal/extractor"
)

type ClassificationResult struct {
	Seriousness string
	CriteriaMet string
	RoutingLane string
}

func Classify(c *extractor.CaseData, rawNarrative string) ClassificationResult {
	// Check both the extracted event description AND the original narrative
	combined := strings.ToLower(c.EventDescription + " " + c.Outcome + " " + rawNarrative)

	if strings.Contains(combined, "fatal") || strings.Contains(combined, "death") || strings.Contains(combined, "died") {
		return ClassificationResult{
			Seriousness: "Life-threatening",
			CriteriaMet: "Death",
			RoutingLane: "7_day_expedited",
		}
	}
	if strings.Contains(combined, "life-threatening") || strings.Contains(combined, "life threatening") {
		return ClassificationResult{
			Seriousness: "Life-threatening",
			CriteriaMet: "Life-threatening event",
			RoutingLane: "7_day_expedited",
		}
	}
	if strings.Contains(combined, "hospitalis") || strings.Contains(combined, "hospitaliz") {
		return ClassificationResult{
			Seriousness: "Serious",
			CriteriaMet: "Hospitalisation",
			RoutingLane: "15_day_review",
		}
	}
	if strings.Contains(combined, "disab") {
		return ClassificationResult{
			Seriousness: "Serious",
			CriteriaMet: "Disability",
			RoutingLane: "15_day_review",
		}
	}
	if strings.Contains(combined, "congenital") || strings.Contains(combined, "birth defect") {
		return ClassificationResult{
			Seriousness: "Serious",
			CriteriaMet: "Congenital anomaly",
			RoutingLane: "15_day_review",
		}
	}

	return ClassificationResult{
		Seriousness: "Non-serious",
		CriteriaMet: "None",
		RoutingLane: "30_day_batch",
	}
}
