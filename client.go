package bitrix24

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/parnurzeal/gorequest"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const (
	PROTOCOL     = "https://"
	OAUTH_SERVER = "oauth.bitrix.info"
	OAUTH_TOKEN  = "/oauth/token/"
	AUTH_URL     = "/oauth/authorize/"
)

type Settigns struct {
	Domain            string // domain bitrix24 application
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

//var realationNames = map[string]string{
//	"domain":            "domain",
//	"applicationSecret": "client_secret",
//	"applicationId":     "client_id",
//
//	"accessToken":  "access_token",
//	"refreshToken": "refresh_token",
//	"memberId":     "member_id",
//
//	"applicationScope": "scope",
//	"redirectUri":      "redirect_uri",
//}

//Consist data for authorization
type Client struct {
	//isAccessParams bool //Specifies that all access settings are set

	domain            string // domain bitrix24 application
	applicationSecret string // secret code application [0-9A-z]{50,50} "client_secret"
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

func NewClient(settigns Settigns) (Client, error) {

	b := Client{}

	if err := settigns.checkAccessParams(); err != nil {
		return b, err
	}

	b.domain = settigns.Domain
	b.applicationSecret = settigns.ApplicationSecret
	b.applicationId = settigns.ApplicationId

	//b.timeout = 30
	//b.request = *gorequest.New()

	return b, nil
}

////Set settings attributes
//func (b *Bitrix24) SetAttributes(attributes types.ApplicationSettings) {
//	tReflect := reflect.ValueOf(&b)
//
//	if tReflect.Kind() == reflect.Ptr {
//		tReflect = tReflect.Elem()
//	}
//
//	settings := structs.Map(&attributes)
//
//	for key, value := range settings {
//		if value == nil || value == "" {
//			continue
//		}
//
//		attribute := tReflect.MethodByName("Set" + key)
//
//		if attribute.IsValid() {
//			attribute.Call([]reflect.Value{reflect.ValueOf(value)})
//		} else {
//			panic(key + " not exitst in " + tReflect.Type().Name())
//		}
//	}
//
//	b.CheckAccessParams()
//}

type Response struct {
	resp *http.Response
}

func (r Response) ParseJson(v interface{}) error {
	if r.resp == nil {
		return errors.New("response is empty")
	}

	body, err := ioutil.ReadAll(r.resp.Body)
	if err != nil {
		return err
	}

	fmt.Println(string(body))

	return json.Unmarshal(body, v)
}

func (r Response) Close() error {
	return r.resp.Body.Close()
}

func (c *Client) execute(url string, additionalParameters url.Values) (*http.Request, Response, error) {

	// TODO move to http.
	request.Post(url).SendMap(additionalParameters).Timeout(30 * time.Second)
	request.TargetType = "form"

	resp, _, errs := request.End()

	req, _ := request.MakeRequest()

	if len(errs) > 0 {
		return req, Response{resp: resp}, errs[0]
	}

	//data, _ := jason.NewObjectFromBytes([]byte(body))

	//json.Unmarshal([]byte(body), &data)

	return req, Response{resp: resp}, nil
}

func (s *Settigns) checkAccessParams() error {

	if len(s.Domain) == 0 {
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
