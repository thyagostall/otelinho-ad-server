package openrtb

type AdMarkup struct {
	Native Native `json:"native"`
}

type Native struct {
	Assets        []Asset        `json:"assets"`
	EventTrackers []EventTracker `json:"eventtrackers"`
	Ver           string         `json:"ver"`
}

type Asset struct {
	ID       int                    `json:"id"`
	Data     map[string]interface{} `json:"data"`
	Required int                    `json:"required"`
}

type EventTracker struct {
	Method int    `json:"method"`
	URL    string `json:"url"`
	Event  int    `json:"event"`
}
