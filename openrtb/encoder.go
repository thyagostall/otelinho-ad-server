package openrtb

import (
	"encoding/json"
	"strconv"

	"github.com/google/uuid"
	"thyago.com/otelinho/beacon"
	"thyago.com/otelinho/campaign"
)

func Encode(c *campaign.Campaign) *BidResponse {
	impressionID := uuid.New().String()
	return &BidResponse{
		ID: "1",
		SeatBid: []SeatBid{
			{
				Seat: "1",
				Bid: []BidItem{
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
	adm := AdMarkup{
		Native: Native{
			Assets: []Asset{
				{
					ID: 1,
					Data: map[string]interface{}{
						"type":  501,
						"value": c.Creative,
					},
					Required: 1,
				},
			},
			EventTrackers: []EventTracker{
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
