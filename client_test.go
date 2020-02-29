package bitrix24

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/kelseyhightower/envconfig"
	"github.com/skratchdot/open-golang/open"
)

type Config struct {
	ServerPort        string `envconfig:"SERVER_PORT"`
	ApplicationDomain string `envconfig:"APPLICATION_DOMAIN"`
	ApplicationID     string `envconfig:"APPLICATION_ID"`
	ApplicationSecret string `envconfig:"APPLICATION_SECRET"`
}

var bx24 *Client

func init() {

	var c Config
	if err := envconfig.Process("myapp", &c); err != nil {
		panic(err)
	}

	serverStringChannel := make(chan string)
	startServer(c.ServerPort, serverStringChannel)

	settings := Settings{
		ApplicationDomain: c.ApplicationDomain,
		ApplicationSecret: c.ApplicationSecret,
		ApplicationId:     c.ApplicationID,
	}

	var err error
	bx24, err = NewClient(settings)
	if err != nil {
		panic(err)
	}

	open.Run(bx24.GetUrlForRequestCode())

	code := <-serverStringChannel
	err = bx24.Authorize(code)
	if err != nil {
		panic(err)
	}
}

func startServer(port string, channel chan<- string) *http.Server {
	srv := &http.Server{Addr: ":" + port}
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
