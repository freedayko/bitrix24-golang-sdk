package bitrix24

import (
	"github.com/antonholmquist/jason"
	"net/url"
)

//Url to request authorization from the user
func (b *Bitrix24) GetUrlClientAuth(params *url.Values) string {
	b.generateParams(params, "applicationId", "applicationScope")
	params.Set("response_type", "code")

	return b.GetUrlAuth("", params)
}

//Use with the received code after request by getUrlClientAuth
func (b *Bitrix24) GetFirstAccessToken(params *url.Values, update bool) (string, *jason.Object, []error) {
	if params.Get("code") == "" {
		panic("Get code, request token returned by the server (the token default lifetime is 30 sec).")
	}
	b.generateParams(params, "applicationId", "applicationScope", "")
	params.Set("grant_type", "authorization_code")

	urlRequest := b.GetUrlOAuthToken("", params)

	_, _, data, errs := b.execute(urlRequest, nil)

	if len(errs) > 0 {
		return urlRequest, nil, errs
	}

	if update {
		b.updateAccessParams(data)
	}

	return urlRequest, data, nil
}

//func (t *Bitrix24) GetNewAccessToken(update bool) (string, *jason.Object, []error) {
//
//	return "", &jason.Object{}, nil
//}

//func (t *Bitrix24) UpdateToken(update bool) (string, *jason.Object, []error) {
//	params := url.Values{
//		"": {""},
//		"": {""},
//		"": {""},
//		"": {""},
//	}
//
//	url := t.GetUrlOAuthToken("", &params)
//
//}

func (b *Bitrix24) updateAccessParams(data *jason.Object) {
	b.memberId, _ = data.GetString("member_id")
	b.accessToken, _ = data.GetString("access_token")
	b.refreshToken, _ = data.GetString("refresh_token")
	b.applicationScope, _ = data.GetString("scope")
}

func (b *Bitrix24) GetUrlOAuthToken(url string, params *url.Values) string {
	return b.GetUrl(b.domain+OAUTH_TOKEN+url, params)
}

func (b *Bitrix24) GetUrlOAuth(url string, params *url.Values) string {
	return b.GetUrl(OAUTH_SERVER+url, params)
}

func (b *Bitrix24) GetUrlAuth(url string, params *url.Values) string {
	return b.GetUrl(b.domain+AUTH_URL+url, params)
}

func (b *Bitrix24) GetUrl(url string, params *url.Values) string {
	urlParams := ""
	if params != nil {
		urlParams = params.Encode()
	}

	return PROTOCOL + url + "?" + urlParams
}

func (b *Bitrix24) generateParams(params *url.Values, listParams ...string) []error {
	errs := []error{}

	//tReflect := reflect.ValueOf(&b)
	//
	//if tReflect.Kind() == reflect.Ptr {
	//	tReflect = tReflect.Elem()
	//}
	//
	//for _, value := range listParams {
	//	if len(realationNames[value]) > 0 {
	//		methodName := strings.Title(value)
	//		params.Set(realationNames[value], tReflect.MethodByName(methodName).Call([]reflect.Value{})[0].String())
	//	} else {
	//		errs = append(errs, errors.New(value+" isn't exist"))
	//	}
	//}

	return errs
}


