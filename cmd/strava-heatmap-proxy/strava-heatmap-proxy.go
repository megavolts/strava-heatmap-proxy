package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/patrickziegler/strava-heatmap-proxy/internal/strava"
)

type Param struct {
	Config *string
	Port   *string
}

func main() {
	param := &Param{
		Config: flag.String("config", "config.json", "Path to configuration file"),
		Port:   flag.String("port", "8080", "Local proxy port"),
	}
	flag.Parse()

	config, err := strava.ParseConfig(*param.Config)
	if err != nil {
		log.Fatalf("Failed to get configuration: %s", err)
	}

	client := strava.NewStravaClient()

	if err = client.Authenticate(config.Email, config.Password); err != nil {
		log.Fatalf("Failed to authenticate client: %s", err)
	}

	for k, v := range client.GetCloudFrontCookies() {
		fmt.Printf("%s\t%s\n", k, v)
	}

	log.Printf("Starting heatmap proxy on port %s ..", *param.Port)

	http.Handle("/", strava.NewStravaProxy(client))
	log.Fatal(http.ListenAndServe(":"+*param.Port, nil))
}
