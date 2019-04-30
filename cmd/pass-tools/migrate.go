package main

import (
	"github.com/urfave/cli"
)

func migrate() cli.Command {

	return cli.Command{
		Name:  "migrate",
		Usage: "Migrate PASS data from an old format/schema/context to a new one",
		Description: `
			Use one of the sub-commands to perform a specific migration
		`,
		Subcommands: []cli.Command{
			migrateBlob(),
		},
	}
}

type migrateOpts struct {
	fedoraBaseurl string
	elasticURL    string
}

func migrateBlob() cli.Command {
	opts := migrateOpts{}

	return cli.Command{
		Name:  "metadata",
		Usage: "Migrate submission metadata blobs to a new format",
		Description: `
			Finds submissions that contain submission metadata, and attempts
			to migrate them to the desired format (JSON schema). 
		`,
		Flags: []cli.Flag{
			flagFedoraBaseURL(&opts.fedoraBaseurl),
			flagElasticURL(&opts.elasticURL),
		},
		Action: func(c *cli.Context) error {
			return migrateBlobAction(opts, c.Args())
		},
	}
}

func migrateBlobAction(opts migrateOpts, args []string) error {
	return nil
}
