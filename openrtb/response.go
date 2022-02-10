package openrtb

type BidResponse struct {
	ID      string    `json:"id"`
	SeatBid []SeatBid `json:"seatbid"`
}

type SeatBid struct {
	Bid  []BidItem `json:"bid"`
	Seat string    `json:"seat"`
}

type BidItem struct {
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
