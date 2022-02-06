package pacing

import (
	"database/sql"
	"math/rand"

	"thyago.com/otelinho/campaign"
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
