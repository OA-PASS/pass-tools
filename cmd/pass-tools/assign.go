package main

import (
	"fmt"
	"log"

	"github.com/oa-pass/pass-tools/lib/assign"
	"github.com/urfave/cli"
)

type grantPIOpts struct {
	doSubmissions bool
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
		},
		Action: func(c *cli.Context) error {
			return grantPIAction(opts, c.Args())
		},
	}
}

func grantPIAction(opts grantPIOpts, args []string) (err error) {

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
			Fedora:      fedoraClient(),
			Elastic:     elasticClient(),
		}.Perform()

		if err != nil {
			log.Printf("ERROR: failed assigning grant %s to user %s: %s", grant, user, err)
		} else {
			log.Printf("Assigned grant %s to user %s", grant, user)
		}
	}

	if err != nil {
		log.Fatalf("Finished, but errors encountered.  Check the output")
	}

	return nil
}
