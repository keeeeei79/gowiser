package main

import (
	"fmt"
	"log"
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

type Document struct {
	ID    int64
	Title string
	Body  string
}

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
				docs := []*Document{
					{Title: "Test Title", Body: "Test Body"},
					{Title: "Test Title Two", Body: "Test Body Two"},
					{Title: "Test Title Three", Body: "Test Body Three Test"},
					{Title: "Test Title Four", Body: "Test Body Four"},
					{Title: "Test Title Five", Body: "Test Body Five Test"},
				}
				err := addDocs(docs)
				if err != nil {
					log.Fatalln("Fail add Docs: ", err)
				}
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
