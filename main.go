package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	_ "github.com/lib/pq"
	"thyago.com/otelinho/auction"
	"thyago.com/otelinho/beacon"
	"thyago.com/otelinho/index"
	"thyago.com/otelinho/openrtb"
	"thyago.com/otelinho/openrtbencoder"
	"thyago.com/otelinho/pacing"
	"thyago.com/otelinho/storage"
	"thyago.com/otelinho/targeting"
)

type beaconRequest struct {
	Event         string
	EncodedBeacon string
}

func timeit(endpoint string, handler func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		handler(w, r)
		fmt.Printf("elapsed: %dμs\t[%s]\n", time.Since(start).Microseconds(), endpoint)
	}
}

func main() {
	beacons := make(chan beaconRequest, 1000)
	go processBeacons(beacons)

	db := createDB()
	defer db.Close()

	http.HandleFunc("/openrtb", timeit("openrtb", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.NotFound(w, r)
			return
		}

		var bid openrtb.BidRequest
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&bid)

		if err != nil {
			return
		}

		campaigns := index.RetrieveLiveCampaigns()
		campaigns = targeting.RankAndExclude(&bid, campaigns)
		campaigns = pacing.PaceCampaigns(campaigns)
		campaign := auction.RunAuction(campaigns)
		if campaign != nil {
			response, _ := json.Marshal(openrtbencoder.Encode(campaign))
			w.WriteHeader(http.StatusOK)
			w.Write(response)
		} else {
			w.WriteHeader(http.StatusNoContent)
		}
	}))
	http.HandleFunc("/event/", timeit("event", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		items := strings.Split(path, "/")

		eventType := items[2]
		eventMetadata := items[3]

		beacons <- beaconRequest{Event: eventType, EncodedBeacon: eventMetadata}

		pixel := []byte("\x47\x49\x46\x38\x39\x61\x01\x00\x01\x00\x80\x00\x00\x00\x00\x00\x00\x00\x00\x21\xF9\x04\x01\x00\x00\x00\x00\x2C\x00\x00\x00\x00\x01\x00\x01\x00\x00\x02\x02\x44\x01\x00\x3B")
		w.WriteHeader(http.StatusOK)
		w.Header().Add("Content-Type", "image/gif")
		w.Write(pixel)
	}))
	http.HandleFunc("/velocity/", timeit("velocity", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		items := strings.Split(path, "/")

		campaignID, _ := strconv.Atoi(items[2])

		db := createDB()
		defer db.Close()

		campaign := storage.RetrieveCampaignByID(db, campaignID)
		res, err := pacing.AdjustVelocity(db, campaign)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "%v", err)
		} else {
			response, _ := json.Marshal(res)
			w.WriteHeader(http.StatusOK)
			w.Write(response)
		}
	}))
	http.HandleFunc("/warm-cache", timeit("warm-cache", func(w http.ResponseWriter, r *http.Request) {
		campaigns, err := storage.ActiveCampaignsFromDatabase(db)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "%v", err)
			return
		}

		index.SetLiveCampaigns(campaigns)
		w.WriteHeader(http.StatusOK)
	}))

	http.ListenAndServe("localhost:3000", nil)
}

func processBeacons(beacons chan beaconRequest) {
	db := createDB()
	defer db.Close()

	for b := range beacons {
		err := beacon.RecordBeaconReceived(db, b.EncodedBeacon, b.Event)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func createDB() *sql.DB {
	db, _ := sql.Open("postgres", "host=localhost port=5432 user=otelinho password=devpassword dbname=otelinho sslmode=disable")
	db.SetMaxIdleConns(30)
	db.SetMaxOpenConns(60)
	db.SetConnMaxIdleTime(10 * time.Minute)
	db.SetConnMaxLifetime(30 * time.Minute)
	return db
}
