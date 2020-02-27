package bitrix24

import (
	"errors"
	"github.com/antonholmquist/jason"
	"github.com/parnurzeal/gorequest"
	"github.com/sirupsen/logrus"
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
type Bitrix24 struct {
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

	log Logger

	request gorequest.SuperAgent
}

func NewClient(settigns Settigns) (Bitrix24, error) {

	b := Bitrix24{}

	if err := settigns.checkAccessParams(); err != nil {
		return b, err
	}

	b.domain = settigns.Domain
	b.applicationSecret = settigns.ApplicationSecret
	b.applicationId = settigns.ApplicationId

	// Set logger by default
	if b.log == nil {
		b.log = logrus.New()
	}

	//b.timeout = 30
	//b.request = *gorequest.New()

	return b, nil
}

func (b *Bitrix24) SetLogger(logger Logger) {
	b.log = logger
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

func (b *Bitrix24) execute(url string,
	additionalParameters url.Values) (*http.Request, *http.Response, *jason.Object, []error) {
	request.Post(url).SendMap(additionalParameters).Timeout(30 * time.Second)
	request.TargetType = "form"

	resp, body, errs := request.End()

	req, _ := request.MakeRequest()

	if len(errs) > 0 {
		return req, resp, nil, errs
	}

	data, _ := jason.NewObjectFromBytes([]byte(body))

	//json.Unmarshal([]byte(body), &data)

	return req, resp, data, nil
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

	//if len(b.accessToken) == 0 {
	//	errs = append(errs, errors.New("AccessToken is empty"))
	//}
	//if len(b.refreshToken) == 0 {
	//	errs = append(errs, errors.New("RefreshToken is empty"))
	//}
	//if len(b.memberId) == 0 {
	//	errs = append(errs, errors.New("MemberId is empty"))
	//}
	//
	//if len(b.applicationScope) == 0 {
	//	errs = append(errs, errors.New("ApplicationScope is empty"))
	//}
	//
	//b.isAccessParams = len(errs) == 0

	return nil
}


///

type AuthData struct {
	AccessToken  string //token for access, [0-9a-z]{32}
	RefreshToken string //token for refresh token access, [0-9a-z]{32}
	MemberId     string //the unique Bitrix24 portal ID
	ApplicationScope string
}

func GetSettingsByJson(data *jason.Object) *AuthData {
	memberId, _ := data.GetString("member_id")
	accessToken, _ := data.GetString("access_token")
	refreshToken, _ := data.GetString("refresh_token")
	applicationScope, _ := data.GetString("scope")

	return &AuthData{
		MemberId:         memberId,
		AccessToken:      accessToken,
		RefreshToken:     refreshToken,
		ApplicationScope: applicationScope,
	}
}