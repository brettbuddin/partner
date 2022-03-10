package internal

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/brettbuddin/partner/internal/command"
	"github.com/urfave/cli/v2"
)

func Main() {
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	app := cli.NewApp()
	app.Name = "partner"
	app.Usage = "Manage git coauthors"
	app.Commands = []*cli.Command{
		cmdManifest(pwd),
		cmdStatus(pwd),
		cmdSet(pwd),
		cmdClear(pwd),
	}

	if err = app.Run(os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

type codeError struct {
	error
	code int
}

func (e codeError) Error() string {
	return e.error.Error()
}

func (e codeError) ExitCode() int {
	return e.code
}

func newCodeError(err error, code int) error {
	return codeError{error: err, code: code}
}

func cmdManifest(pwd string) *cli.Command {
	return &cli.Command{
		Name:  "manifest",
		Usage: "Coauthor manifest operations",
		Subcommands: []*cli.Command{
			cmdManifestGitHubAdd(pwd),
			cmdManifestGitLabAdd(pwd),
			cmdManifestAdd(pwd),
			cmdManifestList(pwd),
			cmdManifestRemove(pwd),
		},
	}
}

func cmdManifestGitHubAdd(pwd string) *cli.Command {
	return &cli.Command{
		Name:      "github-add",
		Aliases:   []string{"gh-add"},
		Usage:     "Add a coauthor from GitHub usernames",
		ArgsUsage: "[username, ...]",
		Action: func(c *cli.Context) error {
			if c.Args().Len() == 0 {
				cli.ShowCommandHelp(c, c.Command.Name)
				return newCodeError(fmt.Errorf("at least one GitHub username is required"), 2)
			}
			fetcher := &command.GitHubFetcher{
				BaseURL: "https://api.github.com",
				Client: &http.Client{
					Timeout: 10 * time.Second,
				},
			}
			paths, err := command.DefaultPaths(pwd)
			if err != nil {
				return newCodeError(err, 1)
			}
			err = command.New(paths).ManifestFetchAdd(fetcher, c.Args().Slice()...)
			if err != nil {
				return newCodeError(err, 1)
			}
			return nil
		},
	}
}
func cmdManifestGitLabAdd(pwd string) *cli.Command {
	return &cli.Command{
		Name:      "gitlab-add",
		Aliases:   []string{"gl-add"},
		Usage:     "Add a coauthor from GitLab usernames",
		ArgsUsage: "[username, ...]",
		Action: func(c *cli.Context) error {
			if c.Args().Len() == 0 {
				cli.ShowCommandHelp(c, c.Command.Name)
				return newCodeError(fmt.Errorf("at least one GitLab username is required"), 2)
			}
			fetcher := &command.GitLabFetcher{
				BaseURL: "https://gitlab.com",
				Client: &http.Client{
					Timeout: 10 * time.Second,
				},
			}
			paths, err := command.DefaultPaths(pwd)
			if err != nil {
				return newCodeError(err, 1)
			}
			err = command.New(paths).ManifestFetchAdd(fetcher, c.Args().Slice()...)
			if err != nil {
				return newCodeError(err, 1)
			}
			return nil
		},
	}
}
func cmdManifestAdd(pwd string) *cli.Command {
	return &cli.Command{
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
			paths, err := command.DefaultPaths(pwd)
			if err != nil {
				return newCodeError(err, 1)
			}
			err = command.New(paths).ManifestAdd(
				c.String("id"),
				c.String("name"),
				c.String("email"),
			)
			if err != nil {
				return newCodeError(err, 1)
			}
			return nil
		},
	}
}
func cmdManifestList(pwd string) *cli.Command {
	return &cli.Command{
		Name:    "list",
		Aliases: []string{"ls"},
		Usage:   "List coauthors",
		Action: func(c *cli.Context) error {
			paths, err := command.DefaultPaths(pwd)
			if err != nil {
				return newCodeError(err, 1)
			}
			if err := command.New(paths).ManifestList(os.Stdout); err != nil {
				return newCodeError(err, 1)
			}
			return nil
		},
	}
}
func cmdManifestRemove(pwd string) *cli.Command {
	return &cli.Command{
		Name:      "remove",
		Aliases:   []string{"rm"},
		Usage:     "Remove coauthors",
		ArgsUsage: "[id, ...]",
		Action: func(c *cli.Context) error {
			if c.Args().Len() == 0 {
				cli.ShowCommandHelp(c, c.Command.Name)
				return newCodeError(fmt.Errorf("at least one ID is required"), 2)
			}
			paths, err := command.DefaultPaths(pwd)
			if err != nil {
				return newCodeError(err, 1)
			}
			if err := command.New(paths).ManifestRemove(c.Args().Slice()...); err != nil {
				return newCodeError(err, 1)
			}
			return nil
		},
	}
}

func cmdStatus(pwd string) *cli.Command {
	return &cli.Command{
		Name:  "status",
		Usage: "Show active coauthors",
		Action: func(c *cli.Context) error {
			paths, err := command.DefaultPaths(pwd)
			if err != nil {
				return newCodeError(err, 1)
			}
			if err := command.New(paths).TemplateStatus(os.Stdout); err != nil {
				return newCodeError(err, 1)
			}
			return nil
		},
	}
}

func cmdSet(pwd string) *cli.Command {
	return &cli.Command{
		Name:      "set",
		Aliases:   []string{"activate"},
		Usage:     "Set active coauthors",
		ArgsUsage: "[id, ...]",
		Action: func(c *cli.Context) error {
			if c.Args().Len() == 0 {
				cli.ShowCommandHelp(c, c.Command.Name)
				return newCodeError(fmt.Errorf("at least one ID is required"), 2)
			}
			paths, err := command.DefaultPaths(pwd)
			if err != nil {
				return newCodeError(err, 1)
			}
			if err := command.New(paths).TemplateSet(c.Args().Slice()...); err != nil {
				return newCodeError(err, 1)
			}
			return nil
		},
	}
}

func cmdClear(pwd string) *cli.Command {
	return &cli.Command{
		Name:  "clear",
		Usage: "Clear active coauthors",
		Action: func(c *cli.Context) error {
			paths, err := command.DefaultPaths(pwd)
			if err != nil {
				return newCodeError(err, 1)
			}
			if err := command.New(paths).TemplateClear(); err != nil {
				return newCodeError(err, 1)
			}
			return nil
		},
	}
}
