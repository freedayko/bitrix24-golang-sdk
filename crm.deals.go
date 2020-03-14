package bitrix24

import (
	"time"

	"github.com/parnurzeal/gorequest"
)

const (
	CRM_TYPE_ID_GOODS = "GOODS"
	CRM_TYPE_ID_SALE = "SALE"
	CRM_TYPE_ID_SERVICE = "SERVICE"

	CRM_STAGE_ID_NEW = "NEW"
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
		Fields CrmDeal                `json:"fields"`
		Params map[string]interface{} `json:"params"`
	}{
		Fields: deal,
		Params: map[string]interface{}{"REGISTER_SONET_EVENT": "Y"},
	}

	u := c.getUrlMethod("crm.deal.add")
	_, err := c.execute(gorequest.TypeJSON, u, body)
	return err
}
