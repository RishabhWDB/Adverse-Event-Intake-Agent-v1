package router

func AssignLane(seriousness string) string {
	switch seriousness {
	case "Life-threatening":
		return "7_day_expedited"
	case "Serious":
		return "15_day_review"
	default:
		return "30_day_batch"
	}
}
