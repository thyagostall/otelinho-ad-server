package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"thyago.com/otelinho/beacon"
	"thyago.com/otelinho/campaign"
	"thyago.com/otelinho/openrtb"
	"thyago.com/otelinho/pacing"
	"thyago.com/otelinho/storage"
)

type beaconRequest struct {
	Event         string
	EncodedBeacon string
}

func main() {
	// storage.CreateCampaign(db, "t:m4WgTi-BIDEdAu04G3DEaw;637797729088765952", "2022-03-01T00:00:00Z", "2022-04-10T00:00:00Z")
	// storage.CreateCampaign(db, "t:pWE-0YwL2ycRagbqsCSBuQ;642229909710946305", "2022-01-01T00:00:00Z", "2022-02-10T00:00:00Z")

	beacons := make(chan beaconRequest, 1000)
	go processBeacons(beacons)

	db := createDB()
	defer db.Close()

	r := gin.Default()
	r.POST("/openrtb", func(c *gin.Context) {
		var bid openrtb.BidRequest

		if err := c.BindJSON(&bid); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		go storage.TickAdRequest(db)
		campaign, err := storage.RetrieveCampaign(db)
		if err != nil {
			fmt.Println(err)
		}

		if campaign != nil {
			c.JSON(http.StatusOK, createBidResponse(db, *campaign))
		} else {
			c.Status(http.StatusNoContent)
		}
	})
	r.GET("/event/:event-type/:event-metadata", func(c *gin.Context) {
		eventType := c.Param("event-type")
		eventMetadata := c.Param("event-metadata")

		beacons <- beaconRequest{Event: eventType, EncodedBeacon: eventMetadata}

		pixel := []byte("\x47\x49\x46\x38\x39\x61\x01\x00\x01\x00\x80\x00\x00\x00\x00\x00\x00\x00\x00\x21\xF9\x04\x01\x00\x00\x00\x00\x2C\x00\x00\x00\x00\x01\x00\x01\x00\x00\x02\x02\x44\x01\x00\x3B")
		c.Data(http.StatusOK, "image/gif", pixel)
	})
	r.GET("/velocity/:campaign-id", func(c *gin.Context) {
		campaignID, _ := strconv.Atoi(c.Param("campaign-id"))

		db := createDB()
		defer db.Close()

		campaign := storage.RetrieveCampaignByID(db, campaignID)
		res, err := pacing.AdjustVelocity(db, campaign)
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("%v", err))
		} else {
			c.JSON(http.StatusOK, res)
		}
	})
	r.Run("localhost:3001")
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

func createBidResponse(db *sql.DB, c campaign.Campaign) openrtb.BidResponse {
	impressionID := uuid.New().String()
	return openrtb.BidResponse{
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

func createAdMarkup(c campaign.Campaign, impressionID string) string {
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
