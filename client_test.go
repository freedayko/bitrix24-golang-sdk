package bitrix24

import (
	"encoding/json"
	"fmt"
	"github.com/skratchdot/open-golang/open"
	. "github.com/smartystreets/goconvey/convey"
	"log"
	"net/http"
	"net/url"
	"testing"
	"time"
)

const (
	URL_SERVER  = "http://localhost"
	PORT_SERVER = ":8081"

	DOMAIN             = "b24-7wfqbj.bitrix24.ua"
	APPLICATION_ID     = "local.5e5824645d9023.70610449"
	APPLICATION_SECRET = "SGbNashnp6dHVChwcsa3KJQMRHbfdZeiO100f5t8uMI0H57pug"

	ACCESS_TOKEN  = "90vucimcz2bn349ableoe9z0gchswes5"
	REFRESH_TOKEN = "4dd9pmb1x2iml0hd821qoz5bg39k8qek"
	MEMBER_ID     = "56de4be4aba5516795b585fbaf0798ea"

	APPLICATION_SCOPE = "crm, lists"
	REDIRECT_URL      = URL_SERVER + "/authorization/"
)

func TestBitrix24(t *testing.T) {
	serverStringChannel := make(chan url.Values)
	srv := startServer(serverStringChannel)

	title := ""

	Convey("Check Bitrix24", t, func() {

		settings := Settigns{
			Domain:            DOMAIN,
			ApplicationSecret: APPLICATION_SECRET,
			ApplicationId:     APPLICATION_ID,
		}

		bx24, err := NewClient(settings)
		Convey("Check initiation Bitrix24", func() {
			So(err, ShouldEqual, nil)
		})

		//Convey("Execute Bitrix24", func() {
		//	data := url.Values{
		//		"key1": {"value1"},
		//		"key2": {"value2"},
		//		"key3": {"value3"},
		//		"key4": {"value4"},
		//	}
		//
		//	_, resp, _ := bx24.execute(URL_SERVER+PORT_SERVER+"/simplePost/", data)
		//
		//	jsData, _ := json.Marshal(data)
		//
		//	So(string(jsData), ShouldEqual, result.String())
		//})

		Convey("Check auth Bitrix24", func() {

			params := url.Values{
				"client_id":     {bx24.applicationId},
				"state":         {time.Now().String()},
				"redirect_uri":  {bx24.redirectUri},
				"response_type": {"code"},
				"scope":         {bx24.applicationScope},
			}

			urlAuthClient := PROTOCOL + bx24.domain + AUTH_URL + "?" + params.Encode()

			Convey("getUrl Bitrix24", func() {
				So(urlAuthClient, ShouldEqual, bx24.GetUrlClientAuth(&params))
			})

			title = "Check authorization Bitrix24"

			if testing.Short() {
				SkipConvey(title, func() {})
			} else {
				clientAuthTest := func(update bool) {
					open.Run(urlAuthClient)

					params := <-serverStringChannel

					params.Set("grant_type", "authorization_code")
					params.Set("client_id", bx24.applicationId)
					params.Set("client_secret", bx24.applicationSecret)
					params.Set("scope", bx24.applicationScope)

					urlAuthToken := PROTOCOL + bx24.domain + OAUTH_TOKEN + "?" + params.Encode()

					urlAccessTokenCheck, authData, err := bx24.GetFirstAccessToken(&params, update)
					So(err, ShouldEqual, nil)

					So(urlAccessTokenCheck, ShouldEqual, urlAuthToken)

					if update {
						So(bx24.accessToken, ShouldEqual, authData.AccessToken)
						So(bx24.refreshToken, ShouldEqual, authData.RefreshToken)

						Print("\nAccessToken = " + authData.AccessToken + "\n" +
							"RefreshToken = " + authData.RefreshToken + "\n" +
							"MemberId = " + authData.MemberId + "\n")
					} else {
						So(bx24.accessToken, ShouldNotEqual, authData.AccessToken)
						So(bx24.refreshToken, ShouldNotEqual, authData.RefreshToken)
					}
				}
				Convey(title, func() {
					clientAuthTest(true)
					clientAuthTest(false)
				})
			}
		})
	})

	srv.Close()
}

func startServer(channel chan<- url.Values) *http.Server {
	srv := &http.Server{Addr: PORT_SERVER}
	http.HandleFunc("/simplePost/", func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")

		r.ParseMultipartForm(32 << 20)
		js, _ := json.Marshal(r.PostForm)
		fmt.Fprintf(w, "%s", string(js))
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		channel <- r.URL.Query()

		w.Header().Set("Content-Type", "application/json")

		fmt.Fprint(w, "ok")
	})

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			// cannot panic, because this probably is an intentional close
			log.Printf("Httpserver: ListenAndServe() error: %s", err)
		}
	}()

	return srv
}
