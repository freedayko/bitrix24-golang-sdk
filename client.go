package bitrix24

import (
	"encoding/json"
	"errors"
	"github.com/parnurzeal/gorequest"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const (
	PROTOCOL             = "https://"
	OAUTH_TOKEN_PATH     = "/oauth/token/"
	OAUTH_AUTHORIZE_PATH = "/oauth/authorize/"
)

//Settings is a struct for init client
type Settings struct {
	ApplicationDomain string // domain bitrix24 application
	ApplicationSecret string // secret code application [0-9A-z]{50} "client_secret"
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
	resp *http.Response
}

func (r Response) BindJSON(v interface{}) error {
	if r.resp == nil {
		return errors.New("response is empty")
	}

	body, err := ioutil.ReadAll(r.resp.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, v)
}

func (r Response) Close() error {
	return r.resp.Body.Close()
}

func (c *Client) execute(url string, additionalParameters url.Values) (Response, error) {

	// TODO move to http.
	request.Post(url).SendMap(additionalParameters).Timeout(30 * time.Second)
	request.TargetType = "form"

	resp, _, errs := request.End()

	if len(errs) > 0 {
		return Response{resp: resp}, errs[0]
	}

	return Response{resp: resp}, nil
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
