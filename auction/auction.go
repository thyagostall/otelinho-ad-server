package auction

import "thyago.com/otelinho/campaign"

func RunAuction(campaigns []*campaign.Campaign) *campaign.Campaign {
	if len(campaigns) < 1 {
		return nil
	}

	firstCampaign := campaigns[0]
	if secondCampaign := campaigns[1]; secondCampaign != nil {
		secondCampaign := campaigns[1]
		firstCampaign.MaxBid = secondCampaign.MaxBid + 0.01
	} else {
		firstCampaign.MaxBid = 0.25
	}

	return firstCampaign
}
