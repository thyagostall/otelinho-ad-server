package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
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

func main() {
	r := gin.Default()
	r.POST("/openrtb", func(c *gin.Context) {
		var bid bidRequest

		if err := c.BindJSON(&bid); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, bidResponse{
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
							AdMarkup:     "{\"native\":{\"assets\":[{\"id\":1,\"data\":{\"type\":501,\"value\":\"t:pWE-0YwL2ycRagbqsCSBuQ;642229909710946305\"},\"required\":1}],\"eventtrackers\":[{\"method\":1,\"url\":\"https://tspmagic.tumblr.com/supply_imp2/611b2667-35d7-4732-90b3-dcbf2e374300/yIDej3IEPzQ7dCmK5__d_eFZg6fFkyhhWBDUgekmYUvIBzXw70buLUD9NXz2oeg1zIpMEl33F1GWVwhQQZBjh48Q2lYqxswkmyO4SU8XnkX-504SqWHZkAYC8zwZQ6V2hncRhqbGGUAepoE3S5GshI755LtH_iF4jIRipbXVe8OqmyhvWBI5KsPhO3zVX0e7LPaRCDhZUNgnf45t0sPzQlof6u6az9DMQPRoT103BtnO-9e3FoUyafJNRrZO44woSgRjUM9__IXP_F6iAUw8YSWxbNeaLTrPMKTrl9z1A2rN0E7Z31yzB19OaGn3r9Cpc9XLHxFue1U14thZoP6qjO6PbEBah9bdms3uC7HVqsLa-hxcG2-Ye_cmCKp3kO2b5eOB54VXzsFEgp7xaRZpuwKn1VcuxqBmLkWVg6BSfxJvH4XeyX6lHt-K1W6fCIAIFDqJzbZj0sdLhzvjY3S1oXz37Syoaq-f08hgzI0leOfDITD_scwjzRpXWG6a3qSOihPD0RZ3Kp3tFI49XvNXM1XguyaAdyVKufWIrnsp0v_OSeu5HOv_Z5Ybh2OLXxW-xYUstG9Qv42F7fOkcjrFjh28X_ljmJrIwrDUIcJbP6R_Wb90ArCQAst1KyVJlmHiKIAiNTGVEgcioCnlyGi6qMGmGemAPPcRr2c4wiKun8SAgVglNxS4uOICFKx2LeC8afBZ4GdFMGFwer0KfeI64FRiQHzqzpA69KM1uxMWoWE0Vi8ojoHyhFi1FR3iyW1LuEKHSO6BypNzB2EJ0me4zVatLn7Og0tYbu3cBVkrYonRSrlhYJZXR0N5qPxszlSi1tz552g0VfXqLFu29mgascDEgZduRUcPgXE8wwglJ1Qx_HXtDqlM7gjitnB5pNUyDyrv5yT5-GTjVYkUbwiy0zuuozaDPwWQY8UKRRwdXRrIuWUCy6rrebKqU7sK0Ax7rAzWqEHWyJ5KrCgB_G8Y0jR1tTCPnvsHwMXmFgOFWvob4sLE0cfIpcV17pSJ56FWoLLHRhFJpCoVGOkYAvDNg9kI7ps1n8OmEWP-wAGeeruYcsxArl3zz0TEF8MPEI3biWdPwxUTH7UUtkZr-AfykrU//\",\"event\":1},{\"method\":1,\"url\":\"https://tspmagic.tumblr.com/native_event/yIDej3IEPzQ7dCmK5__d_eFZg6fFkyhhWBDUgekmYUvIBzXw70buLUD9NXz2oeg1zIpMEl33F1GWVwhQQZBjh48Q2lYqxswkmyO4SU8XnkX-504SqWHZkAYC8zwZQ6V2hncRhqbGGUAepoE3S5GshI755LtH_iF4jIRipbXVe8OqmyhvWBI5KsPhO3zVX0e7LPaRCDhZUNgnf45t0sPzQlof6u6az9DMQPRoT103BtnO-9e3FoUyafJNRrZO44woSgRjUM9__IXP_F6iAUw8YSWxbNeaLTrPMKTrl9z1A2rN0E7Z31yzB19OaGn3r9Cpc9XLHxFue1U14thZoP6qjO6PbEBah9bdms3uC7HVqsLa-hxcG2-Ye_cmCKp3kO2b5eOB54VXzsFEgp7xaRZpuwKn1VcuxqBmLkWVg6BSfxJvH4XeyX6lHt-K1W6fCIAIFDqJzbZj0sdLhzvjY3S1oXz37Syoaq-f08hgzI0leOfDITD_scwjzRpXWG6a3qSOihPD0RZ3Kp3tFI49XvNXM1XguyaAdyVKufWIrnsp0v_OSeu5HOv_Z5Ybh2OLXxW-xYUstG9Qv42F7fOkcjrFjh28X_ljmJrIwrDUIcJbP6R_Wb90ArCQAst1KyVJlmHiKIAiNTGVEgcioCnlyGi6qMGmGemAPPcRr2c4wiKun8SAgVglNxS4uOICFKx2LeC8afBZ4GdFMGFwer0KfeI64FRiQHzqzpA69KM1uxMWoWE0Vi8ojoHyhFi1FR3iyW1LuEKHSO6BypNzB2EJ0me4zVatLn7Og0tYbu3cBVkrYonRSrlhYJZXR0N5qPxszlSi1tz552g0VfXqLFu29mgascDEgZduRUcPgXE8wwglJ1Qx_HXtDqlM7gjitnB5pNUyDyrv5yT5-GTjVYkUbwiy0zuuozaDPwWQY8UKRRwdXRrIuWUCy6rrebKqU7sK0Ax7rAzWqEHWyJ5KrCgB_G8Y0jR1tTCPnvsHwMXmFgOFWvob4sLE0cfIpcV17pSJ56FWoLLHRhFJpCoVGOkYAvDNg9kI7ps1n8OmEWP-wAGeeruYcsxArl3zz0TEF8MPEI3biWdPwxUTH7UUtkZr-AfykrU/${EVENT_TYPE}/?u=611b2667-35d7-4732-90b3-dcbf2e374300&sponsored=${SPONSORED}\",\"event\":600}],\"ver\":\"1.2\"}}",
							LossURL:      "https://dev.tumblr.iponweb.net/supply_win_notice/osCsGP4Cuv_jPYqwAlZLxlumReXwQb1p2u-uqcfEVeeiLXgbcxDhRoyppDF0m2MwMFrtULXUYvKxd8Bh2xKkJQeCoUjfgI4TC6snQx91Dl01JC9IGJRCTmzClvQfqFz4Dd7ymyaijjgQYp-R2V3YwFcmrItw4548i-lu8-X2RhpzqOkqYCoo5kE8RAU6pHWYsH7i3TrfSYgs7YAm2Czn4npc07xzcApD9RNfUnmgJnea_W6AZOmaHTrqzxdDSocMEgv9ce9smJ9SUN5lLd0xEL3Ti2hlg6UIyesIWzmVtdKxFyzGY_t9IIcTflRjacvgRRhTMCdO",
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
		})
	})
	r.Run("localhost:3001")
}
