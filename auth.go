package bitrix24

import (
	"errors"
	"net/url"
)

type AuthData struct {
	AccessToken      string `json:"access_token"`  //token for access, [0-9a-z]{32}
	RefreshToken     string `json:"refresh_token"` //token for refresh token access, [0-9a-z]{32}
	MemberId         string `json:"member_id"`     //the unique Bitrix24 portal ID
	ApplicationScope string `json:"scope"`
}

//Url to request authorization from the user
func (c *Client) GetUrlClientAuth(params *url.Values) string {

	params.Set("response_type", "code")

	return c.GetUrlAuth("", params)
}

//Use with the received code after request by getUrlClientAuth
func (c *Client) GetFirstAccessToken(params *url.Values, update bool) (string, AuthData, error) {
	if params.Get("code") == "" {
		return "", AuthData{}, errors.New("code must be set")
	}

	params.Set("grant_type", "authorization_code")

	url := c.GetUrlOAuthToken("", params)

	_, resp, err := c.execute(url, nil)
	if err != nil {
		return url, AuthData{}, err
	}
	defer resp.Close()

	var authData = AuthData{}
	err = resp.ParseJson(&authData)
	if err != nil {
		return url, authData, err
	}

	if update {
		c.updateAccessParams(authData)
	}

	return url, authData, nil
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

func (c *Client) updateAccessParams(data AuthData) {
	c.memberId = data.MemberId
	c.accessToken = data.AccessToken
	c.refreshToken = data.RefreshToken
	c.applicationScope = data.ApplicationScope
}

func (c *Client) GetUrlOAuthToken(url string, params *url.Values) string {
	return c.GetUrl(c.domain+OAUTH_TOKEN+url, params)
}

func (c *Client) GetUrlOAuth(url string, params *url.Values) string {
	return c.GetUrl(OAUTH_SERVER+url, params)
}

func (c *Client) GetUrlAuth(url string, params *url.Values) string {
	return c.GetUrl(c.domain+AUTH_URL+url, params)
}

func (c *Client) GetUrl(url string, params *url.Values) string {
	urlParams := ""
	if params != nil {
		urlParams = params.Encode()
	}

	return PROTOCOL + url + "?" + urlParams
}
