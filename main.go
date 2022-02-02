package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"thyago.com/hello-world-golang-gin/storage"
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
	Price        float32  `json:"price"`
	CampaignID   string   `json:"cid"`
	ID           string   `json:"id"`
	AdMarkup     string   `json:"adm"`
	LossURL      string   `json:"nurl"`
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
	s := storage.Make()
	s.CreateCampaign("t:m4WgTi-BIDEdAu04G3DEaw;637797729088765952", "2022-03-01T00:00:00Z", "2022-04-10T00:00:00Z")
	s.CreateCampaign("t:pWE-0YwL2ycRagbqsCSBuQ;642229909710946305", "2022-01-01T00:00:00Z", "2022-02-10T00:00:00Z")

	r := gin.Default()
	r.POST("/openrtb", func(c *gin.Context) {
		var bid bidRequest

		if err := c.BindJSON(&bid); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		res := s.RetrieveCampaign()
		c.JSON(http.StatusOK, createBidResponse(*res))
	})
	r.GET("/event/:event-id/:event-type", func(c *gin.Context) {
		eventID := c.Param("event-id")
		eventType := c.Param("event-type")

		fmt.Printf("event-id: %s, event-type: %s\n", eventID, eventType)

		pixel := []byte("\x47\x49\x46\x38\x39\x61\x01\x00\x01\x00\x80\x00\x00\x00\x00\x00\x00\x00\x00\x21\xF9\x04\x01\x00\x00\x00\x00\x2C\x00\x00\x00\x00\x01\x00\x01\x00\x00\x02\x02\x44\x01\x00\x3B")
		c.Data(http.StatusOK, "image/gif", pixel)
	})
	r.Run("localhost:3001")
}

func createBidResponse(c storage.Campaign) bidResponse {
	return bidResponse{
		ID: "1",
		SeatBid: []seatBid{
			{
				Seat: "1",
				Bid: []bidItem{
					{
						DemandSource: "direct",
						Price:        7,
						CampaignID:   "1",
						ID:           "3c8e88f7-9be3-46c3-8c83-26a69fd68e6d",
						AdMarkup:     createAdMarkup(c.Creative),
						LossURL:      "https://localhost:3333/lossurl",
						ADomain: []string{
							"",
						},
						Cat: []string{
							"IAB12-3",
						},
						CrID:         "1",
						ImpressionID: "25eed2e8-6520-47cb-a22c-15ef9b6af4c1",
						AdID:         "1",
						AdmMediaType: "native",
					},
				},
			},
		},
	}
}

func createAdMarkup(c string) string {
	adm := adMarkup{
		Native: native{
			Assets: []asset{
				{
					ID: 1,
					Data: map[string]interface{}{
						"type":  501,
						"value": c,
					},
					Required: 1,
				},
			},
			EventTrackers: []eventTracker{
				{
					Method: 1,
					URL:    "https://localhost:3333/event1",
					Event:  1,
				},
				{
					Method: 1,
					URL:    "https://localhost:3333/universalevent",
					Event:  600,
				},
			},
			Ver: "1.2",
		},
	}
	marshalledAdm, _ := json.Marshal(adm)
	return string(marshalledAdm)
}
