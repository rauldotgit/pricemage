package logic

import (
	"log"
	"net/http"

	"github.com/rauldotgit/pricemage/wapi"
)

func LocalExec() {
	client := &http.Client{}

	auth, err := wapi.GetAuth(client)
	if err != nil {
		log.Println(err.Error())
		return
	}

	var itemIndex wapi.LocalItemIndex = wapi.LocalItemIndex{}
	itemIndex.Initialize()

	localIndexErr := wapi.BuildLocalItemIndex(581, "eu", "en_US", &itemIndex, client, &auth.AccessToken)
	if localIndexErr != nil {
		log.Println(localIndexErr.Error())
		return
	}

	log.Println(itemIndex)

}
