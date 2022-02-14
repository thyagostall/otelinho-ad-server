package storage

import (
	"database/sql"
	"regexp"
	"strconv"
	"time"

	"thyago.com/otelinho/campaign"
)

func ActiveCampaignsFromDatabase(db *sql.DB) ([]*campaign.Campaign, error) {
	query := `
		SELECT id, creative, start_date, end_date, goal, max_bid, remaining_budget, budget, velocity
		FROM campaigns
		JOIN pacing ON campaigns.id = pacing.campaign_id
		WHERE start_date <= $1 AND end_date >= $1
		ORDER BY max_bid DESC
	`

	now := time.Now().UTC()

	stmt, err := db.Prepare(query)
	if err != nil {
		return nil, err
	}
	rows, err := stmt.Query(now)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := []*campaign.Campaign{}
	for rows.Next() {
		var c campaign.Campaign
		rows.Scan(&c.ID, &c.Creative, &c.StartDate, &c.EndDate, &c.Goal, &c.MaxBid, &c.RemainingBudget, &c.Budget, &c.PacingFactor)

		targeting, err := fetchTargeting(db, &c)
		if err != nil {
			return nil, err
		}
		c.Targeting = targeting

		result = append(result, &c)
	}

	return result, nil
}

func fetchTargeting(db *sql.DB, c *campaign.Campaign) ([]campaign.TargetingRule, error) {
	query := "SELECT key, value FROM campaign_targeting WHERE campaign_id = $1"
	stmt, err := db.Prepare(query)
	if err != nil {
		return nil, err
	}
	rows, err := stmt.Query(c.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := []campaign.TargetingRule{}
	var key string
	var value string
	for rows.Next() {
		rows.Scan(&key, &value)

		if key == "age" {
			r, _ := regexp.Compile(`(==|!=|<|<=|>|>=)(\d+)`)
			elements := r.FindStringSubmatch(value)

			var operator campaign.TargetingOperator
			if elements[1] == "==" {
				operator = campaign.Equal
			} else if elements[1] == "!=" {
				operator = campaign.NotEqual
			} else if elements[1] == "<" {
				operator = campaign.LessThan
			} else if elements[1] == "<=" {
				operator = campaign.LessThanOrEqual
			} else if elements[1] == ">" {
				operator = campaign.GreaterThan
			} else if elements[1] == ">=" {
				operator = campaign.GreaterThanOrEqual
			}

			value, _ := strconv.ParseUint(elements[2], 10, 32)

			var targetingRule campaign.TargetingRule = campaign.AgeTargetingRule{Operator: operator, Value: uint(value)}
			result = append(result, targetingRule)
		}
	}

	return result, nil
}

func CreateCampaign(db *sql.DB, creative string, strStartDate string, strEndDate string, goal uint) {
	stmt, _ := db.Prepare("INSERT INTO campaigns (creative, start_date, end_date, goal) VALUES ($1, $2, $3, $4)")
	_, _ = stmt.Exec(creative, strStartDate, strEndDate, goal)
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
