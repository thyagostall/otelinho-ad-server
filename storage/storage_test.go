package storage

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestStorage(t *testing.T) {
	db, _ := sql.Open("sqlite3", "./campaigns_test.db")
	defer db.Close()

	expectedCreative := "the-creative"

	CreateCampaign(db, expectedCreative, "2022-01-01T00:00:00Z", "2099-02-10T00:00:00Z")
	campaign := RetrieveCampaign(db)

	if campaign.Creative != expectedCreative {
		t.Fatalf("Expected creative=%s, got=%s", campaign.Creative, expectedCreative)
	}
}
