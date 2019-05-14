package main

import (
	"fmt"
	"os"

	"github.com/oa-pass/pass-tools/lib/assign"
	"github.com/oa-pass/pass-tools/lib/log"
	"github.com/urfave/cli"
)

type grantPIOpts struct {
	doSubmissions bool
	verbose       int
	dryRun        bool
}

func assignActions() cli.Command {

	return cli.Command{
		Name:  "assign",
		Usage: "Assign ownership of a PASS resource to a user",
		Description: `
			Depending on the nature of the object, the commands herein assign 
			"ownership" of a PASS resource to another individual.  For example, 
			changing the submitter of a submission, or the PI of a grant.
		`,
		Subcommands: []cli.Command{
			grantPI(),
		},
	}
}

func grantPI() cli.Command {
	opts := grantPIOpts{}

	return cli.Command{
		Name:      "pi",
		Usage:     "Assign a new PI to a grant",
		ArgsUsage: "[user grant1 grant2 ...]",
		Description: `
			Assigns a new PI to a grant, optionally re-assigning all submissions
			submitted by the former PI as well when provided with the -s flag 
			(note:  this is dangerous, its only real use case is for massaging demo data).

			The first argument is a URI or localKey to a user, subsequent args are 
			grants.  For example

			pass-tools assign pi -s johnshopkins.edu:jhed:foo1 johnshopkins.edu:grant:12345
		`,
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:        "submissions, s",
				Usage:       "If set, will also assign submissions to new grant holder",
				EnvVar:      "GRANT_ASSIGN_SUBMISSIONS",
				Destination: &opts.doSubmissions,
			},
			cli.IntFlag{
				Name:        "verbosity, v",
				Usage:       "Set the level of log verbosity. accepts values -1, 0, 1, 2",
				EnvVar:      "VERBOSE",
				Destination: &opts.verbose,
			},
			cli.BoolFlag{
				Name:        "dry-run",
				Usage:       "Logs which actions it will perform, but does not perform them",
				Destination: &opts.dryRun,
			},
		},
		Action: func(c *cli.Context) error {
			return grantPIAction(opts, c.Args())
		},
	}
}

func grantPIAction(opts grantPIOpts, args []string) (err error) {

	LOG := log.New(opts.verbose)

	if len(args) < 2 {
		return fmt.Errorf("at least two arguments (user, grant) expected")
	}

	user, grants := args[0], args[1:]

	for _, grant := range grants {
		err = assign.Grant{
			ID:          grant,
			To:          user,
			Submissions: opts.doSubmissions,
			BaseURI:     fedoraBaseURI(),
			Fedora:      fedoraClient(LOG),
			Elastic:     elasticClient(LOG),
			DryRun:      opts.dryRun,
			Log:         LOG,
		}.Perform()

		if err != nil {
			LOG.Warnf("failed assigning grant %s to user %s: %s", grant, user, err)
		} else if !opts.dryRun {
			LOG.Printf("Assigned grant %s to user %s", grant, user)
		}
	}

	if err != nil {
		LOG.Warn.Printf("Finished, but errors encountered.  Check the output")
		os.Exit(1)
	}

	return nil
}
