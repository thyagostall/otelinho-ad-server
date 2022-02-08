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
	"thyago.com/otelinho/pacing"
	"thyago.com/otelinho/storage"
)

// bid request
type bidRequest struct {
	ID     string `json:"id"`
	At     int    `json:"at"`
	Device device `json:"device"`
	User   user   `json:"user"`
	Regs   regs   `json:"regs"`
	Site   site   `json:"site"`
	Imp    []imp  `json:"imp"`
}

type device struct {
}

type user struct {
}

type regs struct {
}

type imp struct {
}

type site struct {
}

// bid response
type bidResponse struct {
	ID      string    `json:"id"`
	SeatBid []seatBid `json:"seatbid"`
}

type seatBid struct {
	Bid  []bidItem `json:"bid"`
	Seat string    `json:"seat"`
}

type bidItem struct {
	DemandSource string   `json:"demand_source"`
	Price        float64  `json:"price"`
	CampaignID   string   `json:"cid"`
	ID           string   `json:"id"`
	AdMarkup     string   `json:"adm"`
	WinURL       string   `json:"nurl"`
	LossURL      string   `json:"lurl"`
	ADomain      []string `json:"adomain"`
	Cat          []string `json:"cat"`
	CrID         string   `json:"crid"`
	ImpressionID string   `json:"impid"`
	AdID         string   `json:"adid"`
	AdmMediaType string   `json:"adm_media_type"`
}

// admarkup
type adMarkup struct {
	Native native `json:"native"`
}

type native struct {
	Assets        []asset        `json:"assets"`
	EventTrackers []eventTracker `json:"eventtrackers"`
	Ver           string         `json:"ver"`
}

type asset struct {
	ID       int                    `json:"id"`
	Data     map[string]interface{} `json:"data"`
	Required int                    `json:"required"`
}

type eventTracker struct {
	Method int    `json:"method"`
	URL    string `json:"url"`
	Event  int    `json:"event"`
}

func main() {
	// storage.CreateCampaign(db, "t:m4WgTi-BIDEdAu04G3DEaw;637797729088765952", "2022-03-01T00:00:00Z", "2022-04-10T00:00:00Z")
	// storage.CreateCampaign(db, "t:pWE-0YwL2ycRagbqsCSBuQ;642229909710946305", "2022-01-01T00:00:00Z", "2022-02-10T00:00:00Z")

	r := gin.Default()
	r.POST("/openrtb", func(c *gin.Context) {
		var bid bidRequest

		if err := c.BindJSON(&bid); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		db := createDB()
		defer db.Close()

		storage.TickAdRequest(db)
		campaign := storage.RetrieveCampaign(db)
		if campaign != nil {
			c.JSON(http.StatusOK, createBidResponse(*campaign))
		} else {
			c.Status(http.StatusNoContent)
		}
	})
	r.GET("/event/:event-type/:event-metadata", func(c *gin.Context) {
		eventType := c.Param("event-type")
		eventMetadata := c.Param("event-metadata")

		db := createDB()
		defer db.Close()

		err := beacon.RecordBeaconReceived(db, eventMetadata, eventType)
		if err != nil {
			fmt.Println(err)
		}

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

func createDB() *sql.DB {
	db, _ := sql.Open("postgres", "host=localhost port=5432 user=otelinho password=devpassword dbname=otelinho sslmode=disable")
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(20)
	db.SetConnMaxIdleTime(1 * time.Second)
	db.SetConnMaxLifetime(30 * time.Second)
	return db
}

func createBidResponse(c campaign.Campaign) bidResponse {
	impressionID := uuid.New().String()
	return bidResponse{
		ID: "1",
		SeatBid: []seatBid{
			{
				Seat: "1",
				Bid: []bidItem{
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
	adm := adMarkup{
		Native: native{
			Assets: []asset{
				{
					ID: 1,
					Data: map[string]interface{}{
						"type":  501,
						"value": c.Creative,
					},
					Required: 1,
				},
			},
			EventTrackers: []eventTracker{
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
