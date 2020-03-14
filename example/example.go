package main

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"

	"github.com/freedayko/bitrix24-golang-sdk"
)

var config struct {
	ApplicationDomain string `envconfig:"APPLICATION_DOMAIN"`
	ApplicationID     string `envconfig:"APPLICATION_ID"`
	ApplicationSecret string `envconfig:"APPLICATION_SECRET"`
}

func init() {
	if err := envconfig.Process("myapp", &config); err != nil {
		panic(err)
	}
}

func main() {

	bx24, err := GetAuthorizedClient()
	if err != nil {
		panic(err)
	}

	deal := bitrix24.CrmDeal{
		Title:        "123",
		TypeID:       bitrix24.CRM_TYPE_ID_GOODS,
		StageID:      bitrix24.CRM_STAGE_ID_NEW,
		//CompanyID:    1,
		//ContactID:    1,
		//Opened:       "123",
		//AssignedByID: 1,
		//Probability:  1,
		//CurrencyID:   "123",
		//Opportunity:  1,
		//CategoryID:   1,
		//BeginDate:    time.Now(),
		//CloseDate:    time.Now().Add(time.Hour * 24 * 30),
	}

	err = bx24.CrmDealAdd(deal)
	if err != nil {
		panic(err)
	}
}

func GetAuthorizedClient() (*bitrix24.Client, error) {
	settings := bitrix24.Settings{
		ApplicationDomain: config.ApplicationDomain,
		ApplicationSecret: config.ApplicationSecret,
		ApplicationId:     config.ApplicationID,
	}

	bx24, err := bitrix24.NewClient(settings)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Please folow the link for authorization: %s\n", bx24.GetUrlForRequestCode())
	fmt.Print("Enter code: ")

	var code string
	fmt.Scanln(&code)

	err = bx24.Authorize(code)
	if err != nil {
		return nil, err
	}

	return bx24, nil
}
