package timecheck

import (
	"time"
)

func CheckAppPriority(created_at time.Time) string {
	now := time.Now().UTC()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	threeDaysAgo := today.AddDate(0, 0, -3)
	oneWeeksAgo := today.AddDate(0, 0, -7)

	createdDate := time.Date(created_at.Year(), created_at.Month(), created_at.Day(), 0, 0, 0, 0, time.UTC)
	if createdDate.Before(oneWeeksAgo) {
		return "urgent"
	} else if createdDate.Before(threeDaysAgo) {
		return "medium"
	}
	return "low"

}
