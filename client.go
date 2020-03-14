package bitrix24

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/parnurzeal/gorequest"
)

const (
	PROTOCOL             = "https://"
	OAUTH_TOKEN_PATH     = "/oauth/token/"
	OAUTH_AUTHORIZE_PATH = "/oauth/authorize/"
)

//Settings is a struct for init client
type Settings struct {
	ApplicationDomain string //domain bitrix24 application
	ApplicationSecret string //secret code application [0-9A-z]{50} "client_secret"
	ApplicationId     string //application identity, (app|local).[0-9a-z]{14,14}.[0-9]{8} "client_id"

	/*
		permissions connection
		calendar, crm, disk, department, entity, im, imbot, lists, log,
		mailservice, sonet_group, task, tasks_extended, telephony, call, user,
		imopenlines, placement
	*/
	//RedirectUri      string //url for redirect after authorization
	//
	//Timeout int //timeout before disconnect
}

var request = gorequest.New()

//Client is a main struct of bitrix client
type Client struct {
	domain            string //domain bitrix24 application
	applicationSecret string //secret code application [0-9A-z]{50,50} "client_secret"
	applicationId     string //application identity, (app|local).[0-9a-z]{14,14}.[0-9]{8,8} "client_id"

	accessToken  string //token for access, [0-9a-z]{32}
	refreshToken string //token for refresh token access
	memberId     string //the unique Bitrix24 portal ID

	applicationScope string //array permissions connection
	redirectUri      string //url for redirect after authorization

	//timeout before disconnect (trying to connect + receiving a response)
	//https://github.com/parnurzeal/gorequest/blob/develop/gorequest.go#L452
	//timeout int

	//request gorequest.SuperAgent
}

func NewClient(settings Settings) (*Client, error) {

	if err := settings.checkAccessParams(); err != nil {
		return nil, err
	}

	b := Client{}

	b.domain = settings.ApplicationDomain
	b.applicationSecret = settings.ApplicationSecret
	b.applicationId = settings.ApplicationId

	return &b, nil
}

type Response struct {
	StatusCode int
	Body       []byte
}

type ErrorResponse struct {
	Error            *int64 `json:"error"`
	ErrorDescription string `json:"error_description"`
}

func (r Response) BindJSON(v interface{}) error {
	return json.Unmarshal(r.Body, v)
}

func (c *Client) execute(targetType string, url string, body interface{}) (*Response, error) {

	// TODO move to http.

	if targetType == gorequest.TypeForm {
		request.Post(url).SendMap(body).Timeout(30 * time.Second)
	} else if targetType == gorequest.TypeJSON {
		request.Post(url).SendStruct(body).Timeout(30 * time.Second)
	} else {
		return nil, errors.New("unknown target type")
	}

	request.TargetType = targetType

	resp, _, errs := request.End()

	defer func() {
		if resp.Body != nil {
			resp.Body.Close()
		}
	}()

	if len(errs) > 0 {
		return nil, errs[0]
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var respError ErrorResponse
	json.Unmarshal(b, &respError)
	if respError.Error != nil {
		return nil, fmt.Errorf("API Error (%d): %s", *respError.Error, respError.ErrorDescription)
	}

	result := Response{
		StatusCode: resp.StatusCode,
		Body:       b,
	}

	return &result, nil
}

func (s *Settings) checkAccessParams() error {

	if len(s.ApplicationDomain) == 0 {
		return errors.New("domain must be set")
	}
	if len(s.ApplicationSecret) == 0 {
		return errors.New("ApplicationSecret must be set")
	}
	if len(s.ApplicationId) == 0 {
		return errors.New("ApplicationId must be set")
	}
	return nil
}
