package main

import (
	"github.com/urfave/cli"
)

func assign() cli.Command {

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

type grantPIOpts struct {
	fedoraBaseurl string
	elasticURL    string
}

func grantPI() cli.Command {
	opts := grantPIOpts{}

	return cli.Command{
		Name:  "pi",
		Usage: "Assign a new PI to a grant",
		Description: `
			Assigns a new PI to a grant, optionally re-assigning all submissions
			submitted by the former PI as well (note:  this is dangerous, its only
		    real use case is for massaging demo data).
		`,
		Flags: []cli.Flag{
			flagFedoraBaseURL(&opts.fedoraBaseurl),
			flagElasticURL(&opts.elasticURL),
		},
		Action: func(c *cli.Context) error {
			return grantPIAction(opts, c.Args())
		},
	}
}

func grantPIAction(opts grantPIOpts, args []string) error {
	return nil
}
