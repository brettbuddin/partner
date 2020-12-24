package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/brettbuddin/partner/internal/command"
	"github.com/urfave/cli/v2"
)

func main() {
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	paths, err := command.DefaultPaths(pwd)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	app := cli.NewApp()
	app.Name = "partner"
	app.Usage = "Manage git coauthors"
	app.Commands = []*cli.Command{
		{
			Name:  "manifest",
			Usage: "Coauthor manifest operations",
			Subcommands: []*cli.Command{
				{
					Name:      "github-add",
					Aliases:   []string{"gh-add"},
					Usage:     "Add a coauthor by fetching their information from GitHub",
					ArgsUsage: "[username, ...]",
					Action: func(c *cli.Context) error {
						if c.Args().Len() == 0 {
							cli.ShowCommandHelp(c, c.Command.Name)
							return codeError{
								error: fmt.Errorf("At least one GitHub username is required"),
								code:  2,
							}
						}
						fetcher := &command.GitHubFetcher{
							BaseURL: "https://api.github.com",
							Client: &http.Client{
								Timeout: 10 * time.Second,
							},
						}
						return command.New(paths).ManifestFetchAdd(fetcher, c.Args().Slice()...)
					},
				},
				{
					Name:  "add",
					Usage: "Add a coauthor by manually entering their information",
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:     "id",
							Usage:    "(required) Identifier for referring to the coauthor",
							Required: true,
							Value:    "",
						},
						&cli.StringFlag{
							Name:     "email",
							Usage:    "(required) Email address",
							Required: true,
							Value:    "",
						},
						&cli.StringFlag{
							Name:     "name",
							Usage:    "(required) Full name",
							Required: true,
							Value:    "",
						},
					},
					Action: func(c *cli.Context) error {
						return command.New(paths).ManifestManualAdd(
							c.String("id"),
							c.String("name"),
							c.String("email"),
						)
					},
				},
				{
					Name:    "list",
					Aliases: []string{"ls"},
					Usage:   "List coauthors",
					Action: func(c *cli.Context) error {
						return command.New(paths).ManifestList(os.Stdout)
					},
				},
				{
					Name:      "remove",
					Aliases:   []string{"rm"},
					Usage:     "Remove coauthors",
					ArgsUsage: "[id, ...]",
					Action: func(c *cli.Context) error {
						if c.Args().Len() == 0 {
							cli.ShowCommandHelp(c, c.Command.Name)
							return codeError{
								error: fmt.Errorf("At least one ID is required"),
								code:  2,
							}
						}
						return command.New(paths).ManifestRemove(c.Args().Slice()...)
					},
				},
			},
		},
		{
			Name:  "status",
			Usage: "Show active coauthors",
			Action: func(c *cli.Context) error {
				return command.New(paths).TemplateStatus(os.Stdout)
			},
		},
		{
			Name:      "set",
			Aliases:   []string{"activate"},
			Usage:     "Set active coauthors",
			ArgsUsage: "[id, ...]",
			Action: func(c *cli.Context) error {
				if c.Args().Len() == 0 {
					cli.ShowCommandHelp(c, c.Command.Name)
					return codeError{
						error: fmt.Errorf("At least one ID is required"),
						code:  2,
					}
				}
				return command.New(paths).TemplateSet(c.Args().Slice()...)
			},
		},
		{
			Name:  "clear",
			Usage: "Clear active coauthors",
			Action: func(c *cli.Context) error {
				return command.New(paths).TemplateClear()
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

type codeError struct {
	error
	code int
}

func (e codeError) ExitCode() int {
	return e.code
}
