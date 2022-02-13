package auction

import "thyago.com/otelinho/campaign"

func RunAuction(campaigns []*campaign.Campaign) *campaign.Campaign {
	if len(campaigns) < 1 {
		return nil
	} else if len(campaigns) == 1 {
		result := campaigns[0]
		result.MaxBid = 0.25
		return result
	}

	first := campaigns[0]
	second := campaigns[1]
	first.MaxBid = second.MaxBid + 0.01

	return first
}
