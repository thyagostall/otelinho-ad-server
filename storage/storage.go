package storage

import (
	"database/sql"
	"time"

	"thyago.com/otelinho/campaign"
)

func CreateCampaign(db *sql.DB, creative string, strStartDate string, strEndDate string, goal uint) {
	stmt, _ := db.Prepare("INSERT INTO campaigns (creative, start_date, end_date, goal) VALUES ($1, $2, $3, $4)")
	_, _ = stmt.Exec(creative, strStartDate, strEndDate, goal)
}

func CreateBudgetReversalControlRecord(db *sql.DB, impressionID string, maxBid float64) error {
	stmt, err := db.Prepare("INSERT INTO budget_reversal_control (impression_id, expires_at, bid, consumed) VALUES ($1, $2, $3, $4);")
	if err != nil {
		return err
	}

	expiration := time.Now().Add(30 * time.Duration(time.Second))
	_, err = stmt.Exec(impressionID, expiration, maxBid, false)
	if err != nil {
		return err
	}

	return nil
}

func RetrieveCampaign(db *sql.DB) *campaign.Campaign {
	query := `
		SELECT id, creative, start_date, end_date, goal, max_bid
		FROM campaigns
		JOIN pacing ON campaigns.id = pacing.campaign_id
		WHERE start_date <= $1 AND end_date >= $2
		AND floor(random() * (2^32-1))::bigint < velocity
		AND remaining_budget >= (0.01 / 1000)
		ORDER BY max_bid DESC
	`

	now := time.Now()
	stmt, _ := db.Prepare(query)
	rows, _ := stmt.Query(now, now)
	defer rows.Close()

	var id int
	var creative string
	var startDate time.Time
	var endDate time.Time
	var goal uint
	var maxBid float64

	var firstCampaign campaign.Campaign
	var secondCampaign campaign.Campaign
	if rows.Next() {
		rows.Scan(&id, &creative, &startDate, &endDate, &goal, &maxBid)
		firstCampaign = campaign.Campaign{ID: id, Creative: creative, StartDate: startDate, EndDate: endDate, Goal: goal, MaxBid: maxBid}

		if rows.Next() {
			rows.Scan(&id, &creative, &startDate, &endDate, &goal, &maxBid)
			secondCampaign = campaign.Campaign{ID: id, Creative: creative, StartDate: startDate, EndDate: endDate, Goal: goal, MaxBid: maxBid}
			firstCampaign.MaxBid = secondCampaign.MaxBid + 0.01
		} else {
			firstCampaign.MaxBid = 0.01
		}

		stmt, _ = db.Prepare("UPDATE campaigns SET remaining_budget = remaining_budget - $1 WHERE id = $2 AND remaining_budget - $1 >= 0")
		res, _ := stmt.Exec(firstCampaign.MaxBid/1000, firstCampaign.ID)
		rowsAffected, _ := res.RowsAffected()
		if rowsAffected > 0 {
			return &firstCampaign
		}
	}

	return nil
}

func RetrieveCampaignByID(db *sql.DB, campaignID int) *campaign.Campaign {
	stmt, _ := db.Prepare("SELECT id, creative, start_date, end_date, goal FROM campaigns WHERE id = $1")
	rows, _ := stmt.Query(campaignID)
	defer rows.Close()

	var id int
	var creative string
	var startDate time.Time
	var endDate time.Time
	var goal uint

	for rows.Next() {
		rows.Scan(&id, &creative, &startDate, &endDate, &goal)
		return &campaign.Campaign{ID: id, Creative: creative, StartDate: startDate, EndDate: endDate, Goal: goal}
	}

	return nil
}

func TickAdRequest(db *sql.DB) error {
	stmt, err := db.Prepare("INSERT INTO ad_requests DEFAULT VALUES;")
	if err != nil {
		return err
	}
	_, err = stmt.Exec()
	if err != nil {
		return err
	}

	return nil
}
