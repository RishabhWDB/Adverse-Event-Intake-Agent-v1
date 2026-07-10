package classifier

import "strings"

type ClassificationResult struct {
	Seriousness string
	CriteriaMet string
}

func Classify(eventDescription string, outcome string, rawNarrative string) ClassificationResult {
	combined := strings.ToLower(eventDescription + " " + outcome + " " + rawNarrative)

	if strings.Contains(combined, "fatal") || strings.Contains(combined, "death") || strings.Contains(combined, "died") {
		return ClassificationResult{
			Seriousness: "Life-threatening",
			CriteriaMet: "Death",
		}
	}
	if strings.Contains(combined, "life-threatening") || strings.Contains(combined, "life threatening") {
		return ClassificationResult{
			Seriousness: "Life-threatening",
			CriteriaMet: "Life-threatening event",
		}
	}
	if strings.Contains(combined, "hospitalis") || strings.Contains(combined, "hospitaliz") {
		return ClassificationResult{
			Seriousness: "Serious",
			CriteriaMet: "Hospitalisation",
		}
	}
	if strings.Contains(combined, "disab") {
		return ClassificationResult{
			Seriousness: "Serious",
			CriteriaMet: "Disability",
		}
	}
	if strings.Contains(combined, "congenital") || strings.Contains(combined, "birth defect") {
		return ClassificationResult{
			Seriousness: "Serious",
			CriteriaMet: "Congenital anomaly",
		}
	}

	return ClassificationResult{
		Seriousness: "Non-serious",
		CriteriaMet: "None",
	}
}
