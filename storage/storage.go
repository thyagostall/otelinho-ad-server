package storage

import (
	"database/sql"
	"fmt"
	"time"

	"thyago.com/otelinho/campaign"
)

func CreateCampaign(db *sql.DB, creative string, strStartDate string, strEndDate string) {
	stmt, _ := db.Prepare("INSERT INTO campaigns (creative, start_date, end_date) VALUES ($1, $2, $3)")
	_, _ = stmt.Exec(creative, strStartDate, strEndDate)
}

func RetrieveCampaign(db *sql.DB) *campaign.Campaign {
	now := time.Now()
	stmt, err := db.Prepare("SELECT id, creative, start_date, end_date FROM campaigns WHERE start_date <= $1 AND end_date >= $2")
	fmt.Println(err)
	rows, err := stmt.Query(now, now)
	fmt.Println(err)
	defer rows.Close()

	var id int
	var creative string
	var startDate time.Time
	var endDate time.Time

	for rows.Next() {
		rows.Scan(&id, &creative, &startDate, &endDate)
		return &campaign.Campaign{ID: id, Creative: creative, StartDate: startDate, EndDate: endDate}
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
