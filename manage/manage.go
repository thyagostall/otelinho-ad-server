package manage

import (
	"database/sql"
	"time"
)

type Campaign struct {
	ID              int       `json:"id"`
	Creative        string    `json:"creative"`
	StartDate       time.Time `json:"startDate"`
	EndDate         time.Time `json:"endDate"`
	Goal            uint      `json:"goal"`
	MaxBid          float64   `json:"maxBid"`
	RemainingBudget float64   `json:"remainingBudget"`
	Budget          float64   `json:"budget"`
}

func RetrieveCampaigns(db *sql.DB) ([]*Campaign, error) {
	query := `
		SELECT id, creative, start_date, end_date, goal, max_bid, remaining_budget, budget
		FROM campaigns
		ORDER BY id DESC
	`

	stmt, err := db.Prepare(query)
	if err != nil {
		return nil, err
	}
	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := []*Campaign{}
	for rows.Next() {
		var c Campaign
		rows.Scan(&c.ID, &c.Creative, &c.StartDate, &c.EndDate, &c.Goal, &c.MaxBid, &c.RemainingBudget, &c.Budget)

		result = append(result, &c)
	}

	return result, nil
}
