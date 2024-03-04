package wapi

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rauldotgit/pricemage/utils"
)

type AuctionHouses struct {
	AuctionHouseID int
	Type           string
	LastModified   int64
}

type RealmItems struct {
	RealmID       int
	Name          string
	LocalizedName string
	Locale        string
	AuctionHouses []AuctionHouses
}

// {"key":{"href":"https://eu.api.blizzard.com/data/wow/realm/1408?namespace=dynamic-eu"},"name":"Norgannon","id":1408,"slug":"norgannon"}

type Realm struct {
	Name string `json:"name"`
	Id   int    `json:"id"`
	Slug string `json:"slug"`
}

type RealmResponse struct {
	Realms []Realm `json:"realms"`
}

func GetRealms(region string, namespace string, locale string, client *http.Client, authToken *string) ([]Realm, error) {
	url := fmt.Sprintf("https://%s.api.blizzard.com/data/wow/realm/index?namespace=%s&locale=%s", region, namespace, locale)

	res, err := utils.Fetch(url, "GET", nil, client, authToken)
	if err != nil {
		return []Realm{}, fmt.Errorf("error in GetRealms: %w", err)
	}

	var realms RealmResponse

	umErr := json.Unmarshal(res, &realms)
	if umErr != nil {
		return []Realm{}, fmt.Errorf("error in GetRealms: %w", umErr)
	}

	return realms.Realms, nil
}
