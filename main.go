package main

import (
	"fmt"
	"os"

	cli "github.com/urfave/cli/v2"
)

var (
	Version = "0.0.1"
)

const (
	name  = "gowiser"
	usage = ""
)

func main() {
	app := cli.NewApp()
	app.Version = Version
	app.Name = name
	app.Usage = usage
	app.Commands = []*cli.Command{
		{
			Name:    "index",
			Aliases: []string{"i"},
			Usage:   "build index",
			Action: func(c *cli.Context) error {
				fmt.Println("Starting to build index...")
				return nil
			},
		},
		{
			Name:    "search",
			Aliases: []string{"s"},
			Usage:   "search",
			Action: func(c *cli.Context) error {
				fmt.Println("Starting to search...")
				return nil
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}
