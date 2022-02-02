package storage

import "testing"

func TestStorage(t *testing.T) {
	expectedCreative := "the-creative"

	s := Make()
	s.CreateCampaign(expectedCreative)
	campaign := s.RetrieveCampaign()

	if campaign.Creative != expectedCreative {
		t.Fatalf("Expected creative=%s, got=%s", campaign.Creative, expectedCreative)
	}
}
