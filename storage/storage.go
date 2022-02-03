package storage

import (
	"database/sql"
	"time"
)

type Campaign struct {
	ID        int
	Creative  string
	StartDate time.Time
	EndDate   time.Time
}

func CreateCampaign(db *sql.DB, creative string, strStartDate string, strEndDate string) {
	stmt, _ := db.Prepare("INSERT INTO campaigns (creative, start_date, end_date) VALUES (?, ?, ?)")
	_, _ = stmt.Exec(creative, strStartDate, strEndDate)
}

func RetrieveCampaign(db *sql.DB) *Campaign {
	stmt, _ := db.Prepare("SELECT id, creative, start_date, end_date FROM campaigns WHERE start_date <= ? AND end_date >= ?")
	rows, _ := stmt.Query(time.Now(), time.Now())
	defer rows.Close()

	var id int
	var creative string
	var startDate time.Time
	var endDate time.Time

	for rows.Next() {
		rows.Scan(&id, &creative, &startDate, &endDate)
		return &Campaign{ID: id, Creative: creative, StartDate: startDate, EndDate: endDate}
	}

	return nil
}
