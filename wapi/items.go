package wapi

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/rauldotgit/pricemage/utils"
)

type GameItem struct {
	ID      int32       `json:"id"`
	Name    string      `json:"name"`
	Quality interface{} `json:"quality"`
	Media   struct {
		Key struct {
			Href string `json:"href"`
		} `json:"key"`
	} `json:"media"`
}

type ItemClass struct {
	Name string
}

type ItemClassIndex struct {
	ItemClasses []struct {
		Name string `json:"name"`
		ID   int    `json:"id"`
	} `json:"item_classes"`
}

// map of all items indexed by id
type LocalItemIndex struct {
	Items map[int32]GameItem
	Size  int32
}

func (L *LocalItemIndex) Initialize() {
	L.Items = make(map[int32]GameItem)
}

func (L *LocalItemIndex) Add(newItem GameItem) {
	_, ok := L.Items[newItem.ID]
	if !ok {
		log.Println("Local index added:", newItem.Name)
		L.Items[newItem.ID] = newItem
	}
}

func GetItemInfo(itemID int32, region string, locale string, client *http.Client, authToken *string) (GameItem, error) {
	url := fmt.Sprintf("https://%s.api.blizzard.com/data/wow/item/%d?namespace=%s&locale=%s", region, itemID, "static-"+region, locale)

	resJSON, err := utils.Fetch(url, "GET", nil, client, authToken)
	if err != nil {
		return GameItem{}, fmt.Errorf("error in GetItemInfo: %w", err)
	}

	var item GameItem

	umErr := json.Unmarshal(resJSON, &item)
	if umErr != nil {
		return GameItem{}, fmt.Errorf("error in GetItemInfo: %w", umErr)
	}

	return item, nil
}

func GetItemClassIndex(region string, locale string, client *http.Client, authToken *string) (ItemClassIndex, error) {
	url := fmt.Sprintf("https://%s.api.blizzard.com/data/wow/item-class/index?namespace=%s&locale=%s", region, "static-"+region, locale)

	resJSON, err := utils.Fetch(url, "GET", nil, client, authToken)
	if err != nil {
		return ItemClassIndex{}, fmt.Errorf("error in GetItemClassIndex: %w", err)
	}

	var itemClassIndex ItemClassIndex

	umErr := json.Unmarshal(resJSON, &itemClassIndex)
	if umErr != nil {
		return ItemClassIndex{}, fmt.Errorf("error in GetItemClassIndex: %w", umErr)
	}

	return itemClassIndex, nil
}

func GetItemClass(itemClassID int, region string, locale string, client *http.Client, authToken *string) (ItemClass, error) {
	url := fmt.Sprintf("https://%s.api.blizzard.com/data/wow/item-class/%d?namespace=%s&locale=%s", region, itemClassID, "static-"+region, locale)

	resJSON, err := utils.Fetch(url, "GET", nil, client, authToken)
	if err != nil {
		return ItemClass{}, fmt.Errorf("error in GetItemClass: %w", err)
	}

	log.Println(string(resJSON))
	// var itemClassIndex ItemClassIndex

	// umErr := json.Unmarshal(resJSON, &itemClassIndex)
	// if umErr != nil {
	// 	return ItemClass{}, fmt.Errorf("error in GetItemClassIndex: %w", umErr)
	// }

	return ItemClass{}, nil
}

// time consuming, fetches all available auctions, retrieves the item info and
func BuildLocalItemIndex(realmID int, region string, locale string, itemIndex *LocalItemIndex, client *http.Client, authToken *string) error {
	rawAuctions, err := GetRawAuctions(realmID, region, locale, client, authToken)
	if err != nil {
		return fmt.Errorf("error in BuildLocalItemIndex: %w", err)
	}

	localIndex := *itemIndex

	log.Println("Starting full local item index build.")
	for _, rawAuc := range rawAuctions {
		gameItem, err := GetItemInfo(rawAuc.Item.ID, region, locale, client, authToken)
		if err != nil {
			errMsg := fmt.Errorf("non-fatal error in BuildLocalItemIndex: %w", err)
			log.Println(errMsg)
			continue
		}
		localIndex.Add(gameItem)
	}
	log.Println("Full local item index build complete.")
	return nil
}


// less time consuming, fetches only gameItem info of items not yet in the index taken from new auctions
func UpdateLocalItemIndex(realmID int, region string, locale string, itemIndex *LocalItemIndex, client *http.Client, authToken *string) error {
	rawAuctions, err := GetRawAuctions(realmID, region, locale, client, authToken)
	if err != nil {
		return fmt.Errorf("error in BuildLocalItemIndex: %w", err)
	}

	localIndex := *itemIndex

	log.Println("Starting local item index refresh.")
	for _, rawAuc := range rawAuctions {

		_, ok := localIndex.Items[rawAuc.Item.ID]
		if !ok {
			gameItem, err := GetItemInfo(rawAuc.Item.ID, region, locale, client, authToken)
			if err != nil {
				errMsg := fmt.Errorf("non-fatal error in BuildLocalItemIndex: %w", err)
				log.Println(errMsg)
				continue
			}
			localIndex.Add(gameItem)
		}

	}
	log.Println("Local item index refresh complete.")
	return nil
}
