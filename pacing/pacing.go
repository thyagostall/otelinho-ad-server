package pacing

import (
	"database/sql"
	"fmt"
	"time"

	"thyago.com/otelinho/campaign"
	"thyago.com/otelinho/forecast"
)

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

	remainingTime := time.Until(campaign.EndDate)
	// elapsedTime := time.Since(campaign.StartDate)

	var ratio float64
	if currentCount > 0 {
		projectedCount := (float64(currentCount) / (campaign.Budget - campaign.RemainingBudget)) * campaign.RemainingBudget
		ratio = (projectedCount / remainingTime.Minutes()) / inventory
	} else {
		ratio = ((campaign.Budget / campaign.MaxBid * 1000.0) / remainingTime.Minutes()) / inventory
		fmt.Printf("Impressions: %f\n", ratio)
	}
	newVelocity := uint32(ratio * float64(^uint32(0)))

	res := make(map[string]interface{})
	res["inventory"] = inventory
	res["remaining_minutes"] = remainingTime.Minutes()
	res["impressions"] = currentCount
	res["remaining_budget"] = campaign.RemainingBudget
	res["new_velocity"] = newVelocity

	stmt, err = db.Prepare("UPDATE pacing SET velocity = $1 WHERE campaign_id = $2")
	if err != nil {
		return nil, err
	}
	stmt.Exec(newVelocity, campaign.ID)

	return res, nil
}
