package forecast

import (
	"database/sql"
	"errors"
)

func ComputeForecast(db *sql.DB) (float64, error) {
	query := `
	SELECT avg(count) FROM (
		SELECT
			to_char(timestamp, 'YYYY-MM-DD HH24:MI') AS group,
			count(*) AS count
		FROM ad_requests
		GROUP BY 1
		ORDER BY 1 DESC
		LIMIT 3 OFFSET 1
	) AS forecast;`

	stmt, err := db.Prepare(query)
	if err != nil {
		return 0, err
	}
	rows, err := stmt.Query()
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	if rows.Next() {
		var f float64
		rows.Scan(&f)

		return f, nil
	}

	return 0, errors.New("no result")
}
