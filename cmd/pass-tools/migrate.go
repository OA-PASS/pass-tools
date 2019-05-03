package main

import (
	"github.com/oa-pass/pass-tools/lib/migrate"
	"github.com/urfave/cli"
)

func migrateActions() cli.Command {

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
	dryRun bool
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
			cli.BoolFlag{
				Name:        "dry-run",
				Usage:       "Retrieves and transforms metadata, but does not update submission records",
				Destination: &opts.dryRun,
			},
		},
		Action: func(c *cli.Context) error {
			return migrateBlobAction(opts, c.Args())
		},
	}
}

func migrateBlobAction(opts migrateOpts, args []string) error {
	return migrate.MetadataV0toV1{
		DryRun:  opts.dryRun,
		BaseURI: fedoraBaseURI(),
		Fedora:  fedoraClient(),
		Elastic: elasticClient(),
	}.Perform()
}
