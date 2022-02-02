package storage

import "testing"

func TestStorage(t *testing.T) {
	expectedCreative := "the-creative"

	s := Make()
	s.CreateCampaign(expectedCreative, "2022-01-01T00:00:00Z", "2099-02-10T00:00:00Z")
	campaign := s.RetrieveCampaign()

	if campaign.Creative != expectedCreative {
		t.Fatalf("Expected creative=%s, got=%s", campaign.Creative, expectedCreative)
	}
}
