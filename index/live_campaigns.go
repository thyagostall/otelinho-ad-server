package index

import (
	"sync"

	"thyago.com/otelinho/campaign"
)

type cacheLiveCampaigns struct {
	LiveCampaigns []*campaign.Campaign
	lock          sync.RWMutex
}

var cache cacheLiveCampaigns = cacheLiveCampaigns{
	LiveCampaigns: []*campaign.Campaign{},
	lock:          sync.RWMutex{},
}

func RetrieveLiveCampaigns() []*campaign.Campaign {
	cache.lock.RLock()
	campaigns := cache.LiveCampaigns
	cache.lock.RUnlock()
	return campaigns
}

func SetLiveCampaigns(campaigns []*campaign.Campaign) {
	cache.lock.Lock()
	cache.LiveCampaigns = campaigns
	cache.lock.Unlock()
}
