package main

import (
	"log"
	"os"

	"github.com/urfave/cli"
)

var fatalf = log.Fatalf

func main() {
	app := cli.NewApp()
	app.Name = "pass-utils"
	app.Usage = "PASS utilities"
	app.EnableBashCompletion = true
	app.Commands = []cli.Command{
		assign(),
		migrate(),
	}
	err := app.Run(os.Args)
	if err != nil {
		fatalf("%s", err)
	}
}

func flagFedoraBaseURL(dest *string) cli.Flag {
	return cli.StringFlag{
		Name:        "fedora, pass.fedora.baseurl",
		Usage:       "Fedora baseURL",
		EnvVar:      "PASS_FEDORA_BASEURL",
		Value:       "http://localhost:8080/fcrepo/rest/",
		Destination: dest,
	}
}

func flagElasticURL(dest *string) cli.Flag {
	return cli.StringFlag{
		Name:        "es, pass.elasticsearch.url",
		Usage:       "Elasticsearch URL",
		EnvVar:      "PASS_ELASTICSEARCH_URL",
		Value:       "http://localhost:9200/pass",
		Destination: dest,
	}
}
