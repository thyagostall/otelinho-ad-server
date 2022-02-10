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

func RetrieveCampaign(db *sql.DB) (*campaign.Campaign, error) {
	query := `
		SELECT id, creative, start_date, end_date, goal, max_bid
		FROM campaigns
		JOIN pacing ON campaigns.id = pacing.campaign_id
		WHERE start_date <= $1 AND end_date >= $1
		AND floor(random() * (2^32-1))::bigint < velocity
		AND remaining_budget >= ($2 / 1000.00)
		ORDER BY max_bid DESC
	`

	now := time.Now().UTC()
	floorPrice := 0.25

	stmt, err := db.Prepare(query)
	if err != nil {
		return nil, err
	}
	rows, err := stmt.Query(now, floorPrice)
	if err != nil {
		return nil, err
	}
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
			firstCampaign.MaxBid = floorPrice
		}

		stmt, err = db.Prepare("UPDATE campaigns SET remaining_budget = remaining_budget - $1 WHERE id = $2 AND remaining_budget - $1 >= 0")
		if err != nil {
			return nil, err
		}
		res, err := stmt.Exec(firstCampaign.MaxBid/1000, firstCampaign.ID)
		if err != nil {
			return nil, err
		}
		rowsAffected, err := res.RowsAffected()
		if err != nil {
			return nil, err
		}
		if rowsAffected > 0 {
			return &firstCampaign, nil
		}
	}

	return nil, nil
}

func RetrieveCampaignByID(db *sql.DB, campaignID int) *campaign.Campaign {
	stmt, _ := db.Prepare("SELECT id, creative, start_date, end_date, goal, budget, remaining_budget, max_bid FROM campaigns WHERE id = $1")
	rows, _ := stmt.Query(campaignID)
	defer rows.Close()

	var id int
	var creative string
	var startDate time.Time
	var endDate time.Time
	var goal uint
	var remainingBudget float64
	var budget float64
	var maxBid float64

	for rows.Next() {
		rows.Scan(&id, &creative, &startDate, &endDate, &goal, &budget, &remainingBudget, &maxBid)
		return &campaign.Campaign{ID: id, Creative: creative, StartDate: startDate, EndDate: endDate, Goal: goal, RemainingBudget: remainingBudget, Budget: budget, MaxBid: maxBid}
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
