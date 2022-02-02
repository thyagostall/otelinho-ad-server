package storage

import (
	"time"
)

type Storage struct {
	Campaigns []Campaign
}

type Campaign struct {
	Creative string
}

func (s *Storage) CreateCampaign(creative string) {
	c := Campaign{Creative: creative}
	s.Campaigns = append(s.Campaigns, c)
}

func (s *Storage) RetrieveCampaign() *Campaign {
	index := time.Now().Unix() % int64(len(s.Campaigns))
	return &s.Campaigns[index]
}

func Make() *Storage {
	return &Storage{
		Campaigns: []Campaign{},
	}
}
