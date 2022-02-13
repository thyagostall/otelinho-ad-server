package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"thyago.com/otelinho/beacon"
	"thyago.com/otelinho/campaign"
	"thyago.com/otelinho/index"
	"thyago.com/otelinho/openrtb"
	"thyago.com/otelinho/pacing"
	"thyago.com/otelinho/storage"
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
		campaign := storage.RetrieveCampaign(campaigns)
		if campaign != nil {
			response, _ := json.Marshal(createBidResponse(campaign))
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

func createBidResponse(c *campaign.Campaign) *openrtb.BidResponse {
	impressionID := uuid.New().String()
	return &openrtb.BidResponse{
		ID: "1",
		SeatBid: []openrtb.SeatBid{
			{
				Seat: "1",
				Bid: []openrtb.BidItem{
					{
						DemandSource: "direct",
						Price:        c.MaxBid,
						CampaignID:   strconv.Itoa(c.ID),
						ID:           "3c8e88f7-9be3-46c3-8c83-26a69fd68e6d",
						AdMarkup:     createAdMarkup(c, impressionID),
						WinURL:       beacon.GenerateBeacon(c, impressionID, "win"),
						LossURL:      beacon.GenerateBeacon(c, impressionID, "loss"),
						ADomain: []string{
							"",
						},
						Cat: []string{
							"IAB12-3",
						},
						CrID:         "1",
						ImpressionID: impressionID,
						AdID:         "1",
						AdmMediaType: "native",
					},
				},
			},
		},
	}
}

func createAdMarkup(c *campaign.Campaign, impressionID string) string {
	adm := openrtb.AdMarkup{
		Native: openrtb.Native{
			Assets: []openrtb.Asset{
				{
					ID: 1,
					Data: map[string]interface{}{
						"type":  501,
						"value": c.Creative,
					},
					Required: 1,
				},
			},
			EventTrackers: []openrtb.EventTracker{
				{
					Method: 1,
					URL:    beacon.GenerateBeacon(c, impressionID, "impression"),
					Event:  1,
				},
				{
					Method: 1,
					URL:    beacon.GenerateBeacon(c, impressionID, "${EVENT_TYPE}"),
					Event:  600,
				},
			},
			Ver: "1.2",
		},
	}
	marshalledAdm, _ := json.Marshal(adm)
	return string(marshalledAdm)
}
