package main

import (
	"log"
	"net/http"
	"os"

	"github.com/oa-pass/pass-tools/lib/client"
	"github.com/urfave/cli"
)

var fatalf = log.Fatalf

var globalOpts struct {
	fedoraBaseurl string
	elasticURL    string
	username      string
	password      string
}

func main() {
	app := cli.NewApp()
	app.Name = "pass-utils"
	app.Usage = "PASS utilities"
	app.EnableBashCompletion = true
	app.Commands = []cli.Command{
		assignActions(),
		migrateActions(),
	}
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "fedora, pass.fedora.baseurl",
			Usage:       "Fedora baseURL",
			EnvVar:      "PASS_FEDORA_BASEURL",
			Value:       "http://localhost:8080/fcrepo/rest/",
			Destination: &globalOpts.fedoraBaseurl,
		},

		cli.StringFlag{
			Name:        "es, pass.elasticsearch.url",
			Usage:       "Elasticsearch URL",
			EnvVar:      "PASS_ELASTICSEARCH_URL",
			Value:       "http://localhost:9200/pass/_search",
			Destination: &globalOpts.elasticURL,
		},
		cli.StringFlag{
			Name:        "pass.fedora.user, username, u",
			Usage:       "Username for basic auth to Fedora",
			EnvVar:      "PASS_FEDORA_USER",
			Value:       "fedoraAdmin",
			Destination: &globalOpts.username,
		},
		cli.StringFlag{
			Name:        "pass.fedora.password, password, p",
			Usage:       "Password for basic auth to Fedora",
			EnvVar:      "PASS_FEDORA_PASSWORD",
			Value:       "moo",
			Destination: &globalOpts.password,
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		fatalf("%s", err.Error())
	}
}

func fedoraClient() *client.Simple {
	var credentials *client.Credentials
	if globalOpts.username != "" {
		credentials = &client.Credentials{
			Username: globalOpts.username,
			Password: globalOpts.password,
		}
	}

	return &client.Simple{
		Requester:   &http.Client{},
		BaseURI:     fedoraBaseURI(),
		Credentials: credentials,
	}
}

func fedoraBaseURI() client.BaseURI {
	return client.BaseURI(globalOpts.fedoraBaseurl)
}

func elasticClient() *client.Simple {

	return &client.Simple{
		Requester: &http.Client{},
		BaseURI:   client.BaseURI(globalOpts.elasticURL),
	}
}
