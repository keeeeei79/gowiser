package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
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
				db, err := sqlx.Connect("pgx", "host=127.0.0.1 port=35432 dbname=gowiser user=gowiser password=gowiser sslmode=disable")
				defer db.Close()
				if err != nil {
					log.Fatalln("Fail to connect DB", err)
				}
				docs := []*Document{
					{Title: "Test Title", Body: "Test Body"},
					{Title: "Test Title Two", Body: "Test Hand"},
					{Title: "Test Title Three", Body: "Test Body Two"},
					{Title: "Test Title Four", Body: "Test Foot"},
					{Title: "Test Title Five", Body: "Test Body Two Three"},
				}
				err = addDocs(db, docs)
				if err != nil {
					log.Fatalln("Fail to add Docs: ", err)
				}
				return nil
			},
		},
		{
			Name:    "search",
			Aliases: []string{"s"},
			Usage:   "search",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "keyword",
					Usage: "keyword you want to search",
				},
			},
			Action: func(c *cli.Context) error {
				fmt.Println("Starting to search...")
				db, err := sqlx.Connect("pgx", "host=127.0.0.1 port=35432 dbname=gowiser user=gowiser password=gowiser sslmode=disable")
				defer db.Close()
				if err != nil {
					log.Fatalln("Fail to connect DB", err)
				}
				keyword := c.String("keyword")
				search(db, keyword)
				return nil
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}
