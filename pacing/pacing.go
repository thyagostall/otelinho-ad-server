package pacing

import (
	"database/sql"
	"math/rand"
	"time"

	"thyago.com/otelinho/campaign"
	"thyago.com/otelinho/forecast"
)

func ShouldServe(db *sql.DB, campaign campaign.Campaign) bool {
	stmt, _ := db.Prepare("SELECT velocity FROM pacing WHERE campaign_id = $1")
	rows, _ := stmt.Query(campaign.ID)
	if rows.Next() {
		var velocity uint32
		rows.Scan(&velocity)
		randValue := rand.Uint32()

		return randValue < velocity
	}

	return false
}

func AdjustVelocity(db *sql.DB, campaign *campaign.Campaign) (map[string]interface{}, error) {
	queryImpressions := "SELECT count(*) FROM beacons WHERE campaign_id = $1 AND event = $2;"
	stmt, err := db.Prepare(queryImpressions)
	if err != nil {
		return nil, err
	}
	rows, err := stmt.Query(campaign.ID, "impression")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var currentCount uint
	if rows.Next() {
		rows.Scan(&currentCount)
	} else {
		currentCount = 0
	}

	inventory, err := forecast.ComputeForecast(db)
	if err != nil {
		return nil, err
	}

	remaining := time.Until(campaign.EndDate)
	ratio := float64(campaign.Goal-currentCount) / remaining.Minutes() / inventory
	newVelocity := uint32(ratio * float64(^uint32(0)))

	res := make(map[string]interface{})
	res["inventory"] = inventory
	res["remaining_minutes"] = remaining.Minutes()
	res["impressions"] = currentCount
	res["goal"] = campaign.Goal
	res["new_velocity"] = newVelocity

	stmt, err = db.Prepare("UPDATE pacing SET velocity = $1 WHERE campaign_id = $2")
	if err != nil {
		return nil, err
	}
	stmt.Exec(newVelocity, campaign.ID)

	return res, nil
}
