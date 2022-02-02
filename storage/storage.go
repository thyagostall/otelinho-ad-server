package storage

import (
	"fmt"
	"time"
)

type Storage struct {
	Campaigns []Campaign
}

type Campaign struct {
	Creative  string
	StartDate time.Time
	EndDate   time.Time
}

func (s *Storage) CreateCampaign(creative string, strStartDate string, strEndDate string) {
	startDate, _ := time.Parse(time.RFC3339, strStartDate)
	endDate, _ := time.Parse(time.RFC3339, strEndDate)
	c := Campaign{Creative: creative, StartDate: startDate, EndDate: endDate}
	s.Campaigns = append(s.Campaigns, c)
}

func (s *Storage) RetrieveCampaign() *Campaign {
	now := time.Now()
	candidates := []Campaign{}
	fmt.Println(now)
	for i := 0; i < len(s.Campaigns); i++ {
		fmt.Println(s.Campaigns[i])

		if !now.After(s.Campaigns[i].StartDate) {
			continue
		}

		if !now.Before(s.Campaigns[i].EndDate) {
			continue
		}

		candidates = append(candidates, s.Campaigns[i])
	}

	index := time.Now().Unix() % int64(len(candidates))
	return &candidates[index]
}

func Make() *Storage {
	return &Storage{
		Campaigns: []Campaign{},
	}
}
