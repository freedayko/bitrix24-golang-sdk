package bitrix24

import (
	"encoding/json"
	"fmt"
	"github.com/skratchdot/open-golang/open"
	. "github.com/smartystreets/goconvey/convey"
	"log"
	"net/http"
	"testing"
)

const (
	SERVER_PORT = ":8081"

	APPLICATION_DOMAIN = "b24-7wfqbj.bitrix24.ua"
	APPLICATION_ID     = "local.5e5824645d9023.70610449"
	APPLICATION_SECRET = "SGbNashnp6dHVChwcsa3KJQMRHbfdZeiO100f5t8uMI0H57pug"
)

func TestBitrix24(t *testing.T) {
	serverStringChannel := make(chan string)
	srv := startServer(serverStringChannel)

	title := ""

	Convey("Check Bitrix24", t, func() {

		settings := Settings{
			ApplicationDomain: APPLICATION_DOMAIN,
			ApplicationSecret: APPLICATION_SECRET,
			ApplicationId:     APPLICATION_ID,
		}

		bx24, err := NewClient(settings)
		Convey("Check initiation Bitrix24", func() {
			So(err, ShouldEqual, nil)
		})

		Convey("Check auth Bitrix24", func() {

			title = "Check authorization Bitrix24"

			if testing.Short() {
				SkipConvey(title, func() {})
			} else {
				clientAuthTest := func() {
					open.Run(bx24.GetUrlForRequestCode())

					code := <-serverStringChannel

					authData, err := bx24.Authorization(code)
					So(err, ShouldEqual, nil)

					So(bx24.accessToken, ShouldEqual, authData.AccessToken)
					So(bx24.refreshToken, ShouldEqual, authData.RefreshToken)

					Print("\nAccessToken = " + authData.AccessToken + "\n" +
						"RefreshToken = " + authData.RefreshToken + "\n" +
						"ApplicationScope = " + authData.ApplicationScope + "\n" +
						"MemberId = " + authData.MemberId + "\n")
				}
				Convey(title, func() {
					clientAuthTest()
				})
			}
		})
	})

	srv.Close()
}

func startServer(channel chan<- string) *http.Server {
	srv := &http.Server{Addr: SERVER_PORT}
	http.HandleFunc("/simplePost/", func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")

		r.ParseMultipartForm(32 << 20)
		js, _ := json.Marshal(r.PostForm)
		fmt.Fprintf(w, "%s", string(js))
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		channel <- r.URL.Query().Get("code")
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
