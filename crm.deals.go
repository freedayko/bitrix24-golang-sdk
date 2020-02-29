package bitrix24

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/parnurzeal/gorequest"
)

type CrmDeal struct {
	Title        string    `json:"TITLE"`
	TypeID       string    `json:"TYPE_ID"`
	StageID      string    `json:"STAGE_ID"`
	CompanyID    int64     `json:"COMPANY_ID"`
	ContactID    int64     `json:"CONTACT_ID"`
	Opened       string    `json:"OPENED"`
	AssignedByID int64     `json:"ASSIGNED_BY_ID"`
	Probability  int64     `json:"PROBABILITY"`
	CurrencyID   string    `json:"CURRENCY_ID"`
	Opportunity  int64     `json:"OPPORTUNITY"`
	CategoryID   int64     `json:"CATEGORY_ID"`
	BeginDate    time.Time `json:"BEGINDATE"`
	CloseDate    time.Time `json:"CLOSEDATE"`
}

func (c *Client) CrmDealAdd(deal CrmDeal) error {

	var body = struct {
		Fields CrmDeal
		Params map[string]interface{}
	}{
		Fields: deal,
		Params: map[string]interface{}{"REGISTER_SONET_EVENT": "Y"},
	}

	u := c.getUrlMethod("crm.deal.add")
	resp, err := c.execute(gorequest.TypeJSON, u, body)
	if err != nil {
		return err
	}
	defer resp.Close()

	all, err := ioutil.ReadAll(resp.resp.Body)

	return fmt.Errorf(string(all))
}
