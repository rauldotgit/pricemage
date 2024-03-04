package wapi

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/rauldotgit/pricemage/utils"
)

// {"id":510724466,"item":{"id":210088,"context":28,"bonus_lists":[9547,6652,1481,8766],"modifiers":[{"type":28,"value":2699}]},"buyout":70003200,"quantity":1,"time_left":"LONG"}

type RawAuctionItemInfo struct {
	ID         int32 `json:"id"`
	Context    int   `json:"context"`
	BonusLists []int `json:"bonus_lists"`
	Modifiers  []interface{}
}

// {"id":510724466,"item":{"id":210088,"context":28,"bonus_lists":[9547,6652,1481,8766],"modifiers":[{"type":28,"value":2699}]},"buyout":70003200,"quantity":1,"time_left":"LONG"}
type RawAuction struct {
	AuctionID int64              `json:"id"`
	Item      RawAuctionItemInfo `json:"item"`
	Buyout    int64              `json:"buyout"`
	Bid       int64              `json:"bid"`
	Quantity  int                `json:"quantity"`
	TimeLeft  string             `json:"time_left"`
}

type RawAuctionReponse struct {
	Auctions []RawAuction `json:"auctions"`
}

// married RawAuction item with item info from GetItemInfo
type PMAuction struct {
	AuctionID int64    `json:"id"`
	Item      GameItem `json:"item"`
	Buyout    int64    `json:"buyout"`
	Bid       int64    `json:"bid"`
	Quantity  int      `json:"quantity"`
	TimeLeft  string   `json:"time_left"`
}

func GetRawAuctions(realmID int, region string, locale string, client *http.Client, authToken *string) ([]RawAuction, error) {
	url := fmt.Sprintf("https://%s.api.blizzard.com/data/wow/connected-realm/%d/auctions?namespace=%s&locale=%s", region, realmID, "dynamic-"+region, locale)

	log.Println("Starting raw auction fetch.")
	res, err := utils.Fetch(url, "GET", nil, client, authToken)
	if err != nil {
		return []RawAuction{}, fmt.Errorf("error in GetRawAuctionData: %w", err)
	}
	log.Println("Raw auction fetch complete.")

	var rawAuctions RawAuctionReponse

	umErr := json.Unmarshal(res, &rawAuctions)
	if umErr != nil {
		return []RawAuction{}, fmt.Errorf("error in GetRawAuctionData: %w", umErr)
	}

	// {"id":510724466,"item":{"id":210088,"context":28,"bonus_lists":[9547,6652,1481,8766],"modifiers":[{"type":28,"value":2699}]},"buyout":70003200,"quantity":1,"time_left":"LONG"}

	return rawAuctions.Auctions, nil
}

func GetPMAuctions(realmID int, region string, locale string, client *http.Client, authToken *string) ([]PMAuction, error) {
	rawAuctions, err := GetRawAuctions(realmID, region, locale, client, authToken)
	if err != nil {
		return []PMAuction{}, fmt.Errorf("error in GetPMAuctions: %w", err)
	}

	pmAuctions := make([]PMAuction, len(rawAuctions))

	log.Println("Starting PM auction fetch.")
	for index, rawAuc := range rawAuctions {
		gameItemInfo, err := GetItemInfo(rawAuc.Item.ID, region, locale, client, authToken)
		if err != nil {
			errMsg := fmt.Errorf("non-fatal error in GetPMAuctions: %w", err)
			log.Println(errMsg)
			continue
		}

		newPMAuction := PMAuction{
			AuctionID: rawAuc.AuctionID,
			Item:      gameItemInfo,
			Buyout:    rawAuc.Buyout,
			Bid:       rawAuc.Bid,
			Quantity:  rawAuc.Quantity,
			TimeLeft:  rawAuc.TimeLeft,
		}

		pmAuctions[index] = newPMAuction

		searchKey := "Vibrant"
		if len(newPMAuction.Item.Name) >= len(searchKey) && newPMAuction.Item.Name[:len(searchKey)] == searchKey {
			log.Println(newPMAuction.Item.Name, newPMAuction.Item.Quality)
		}
	}
	log.Println("PM auction fetch complete.")
	return pmAuctions, nil
}
