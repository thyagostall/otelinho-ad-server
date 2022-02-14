package targeting

import (
	"sort"

	"thyago.com/otelinho/campaign"
	"thyago.com/otelinho/openrtb"
)

type RankedCampaign struct {
	Campaign       *campaign.Campaign
	MatchingFactor uint
}

type ByMatchingFactor []*RankedCampaign

func (m ByMatchingFactor) Len() int           { return len(m) }
func (m ByMatchingFactor) Swap(i, j int)      { m[i], m[j] = m[j], m[i] }
func (m ByMatchingFactor) Less(i, j int) bool { return m[i].MatchingFactor > m[j].MatchingFactor }

func RankAndExclude(request *openrtb.BidRequest, campaigns []*campaign.Campaign) []*campaign.Campaign {
	var rankedCampaigns []*RankedCampaign
	for _, campaign := range campaigns {
		matchingFactor, exclude := calculateMatchingFactor(request, campaign)
		if exclude {
			continue
		}

		var r RankedCampaign
		r.Campaign = campaign
		r.MatchingFactor = matchingFactor

		rankedCampaigns = append(rankedCampaigns, &r)
	}

	sort.Sort(ByMatchingFactor(rankedCampaigns))

	var sortedCampaigns []*campaign.Campaign
	for _, rankedCampaign := range rankedCampaigns {
		sortedCampaigns = append(sortedCampaigns, rankedCampaign.Campaign)
	}

	return sortedCampaigns
}

func calculateMatchingFactor(request *openrtb.BidRequest, campaign *campaign.Campaign) (uint, bool) {
	var matchingFactor uint = 0
	for _, targetingRule := range campaign.Targeting {
		if !targetingRule.ShouldInclude(request) {
			return 0, true
		}

		matchingFactor++
	}
	return matchingFactor, false
}
