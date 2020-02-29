package bitrix24

import (
	"errors"
	"net/url"
	"time"

	"github.com/parnurzeal/gorequest"
)

type AuthData struct {
	AccessToken      string `json:"access_token"`  //token for access, [0-9a-z]{32}
	RefreshToken     string `json:"refresh_token"` //token for refresh token access, [0-9a-z]{32}
	MemberId         string `json:"member_id"`     //the unique Bitrix24 portal ID
	ApplicationScope string `json:"scope"`
}

//GetUrlForRequestCode returns url for request authorization code
func (c *Client) GetUrlForRequestCode() string {

	params := url.Values{
		"client_id":     {c.applicationId},
		"state":         {time.Now().String()},
		"redirect_uri":  {c.redirectUri},
		"response_type": {"code"},
		"scope":         {c.applicationScope},
	}

	return c.getUrlAuth(params)
}

//Authorize use with the received code
func (c *Client) Authorize(code string) error {
	if code == "" {
		return errors.New("code must be set")
	}

	var params = url.Values{
		"code":          {code},
		"grant_type":    {"authorization_code"},
		"client_id":     {c.applicationId},
		"client_secret": {c.applicationSecret},
		"scope":         {c.applicationScope},
	}

	u := c.getUrlOAuthToken(params)

	resp, err := c.execute(gorequest.TypeForm, u, nil)
	if err != nil {
		return err
	}
	defer resp.Close()

	var authData = AuthData{}
	err = resp.BindJSON(&authData)
	if err != nil {
		return err
	}

	c.updateAccessParams(authData)

	return nil
}

func (c *Client) updateAccessParams(data AuthData) {
	c.memberId = data.MemberId
	c.accessToken = data.AccessToken
	c.refreshToken = data.RefreshToken
	c.applicationScope = data.ApplicationScope
}

func (c *Client) getUrlOAuthToken(params url.Values) string {
	return c.getUrl(OAUTH_TOKEN_PATH, params)
}

func (c *Client) getUrlAuth(params url.Values) string {
	return c.getUrl(OAUTH_AUTHORIZE_PATH, params)
}

func (c *Client) getUrl(path string, params url.Values) string {
	urlParams := ""
	if params != nil {
		urlParams = "?" + params.Encode()
	}

	return PROTOCOL + c.domain + path + urlParams
}

func (c *Client) getUrlMethod(method string) string {
	return c.getUrl("/rest/"+method+"?auth="+c.accessToken, nil)
}
