package beacon

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"

	"thyago.com/otelinho/campaign"
)

type beacon struct {
	CampaignID   int     `json:"campaign_id"`
	ImpressionID string  `json:"impression_id"`
	BidValue     float64 `json:"bid_value"`
}

var secretKey = []byte("thesecretkey1234thesecretkey1234")
var host = "6ab9-177-220-174-231.ngrok.io"

func GenerateBeacon(campaign *campaign.Campaign, impressionID string, event string) string {
	beacon, _ := json.Marshal(beacon{CampaignID: campaign.ID, ImpressionID: impressionID, BidValue: campaign.MaxBid})
	encrypted := encrypt(secretKey, beacon)
	encoded := base64.URLEncoding.EncodeToString(encrypted)
	return fmt.Sprintf("https://%s/event/%s/%s", host, event, encoded)
}

func RecordBeaconReceived(db *sql.DB, metadata string, event string) error {
	decoded, err := base64.URLEncoding.DecodeString(metadata)
	if err != nil {
		return err
	}

	decrypted := decrypt(secretKey, decoded)
	var b beacon
	json.Unmarshal(decrypted, &b)

	stmt, err := db.Prepare("INSERT INTO beacons (campaign_id, event) VALUES ($1, $2)")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(b.CampaignID, event)
	if err != nil {
		return err
	}

	if event == "loss" {
		bidValue := b.BidValue / 1000.0

		stmt, err = db.Prepare("UPDATE campaigns SET remaining_budget = remaining_budget + $1 WHERE id = $2")
		if err != nil {
			return err
		}

		_, err = stmt.Exec(bidValue, b.CampaignID)
		if err != nil {
			return err
		}
	}

	return nil
}

func encrypt(key []byte, plaintext []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	return ciphertext
}

func decrypt(key []byte, ciphertext []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	if len(ciphertext) < aes.BlockSize {
		panic("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return ciphertext
}
