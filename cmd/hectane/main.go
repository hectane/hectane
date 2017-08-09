package main

import (
	"os"

	"github.com/hectane/hectane/db"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "hectane"
	app.Usage = "SMTP mail server"
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:   "debug",
			EnvVar: "DEBUG",
			Usage:  "enable debug logging",
		},
		cli.StringFlag{
			Name:   "db-host",
			Value:  "localhost",
			EnvVar: "DB_HOST",
			Usage:  "PostgreSQL database host",
		},
		cli.IntFlag{
			Name:   "db-port",
			Value:  5432,
			EnvVar: "DB_PORT",
			Usage:  "PostgreSQL database port",
		},
		cli.StringFlag{
			Name:   "db-name",
			Value:  "postgres",
			EnvVar: "DB_NAME",
			Usage:  "PostgreSQL database name",
		},
		cli.StringFlag{
			Name:   "db-user",
			Value:  "postgres",
			EnvVar: "DB_USER",
			Usage:  "PostgreSQL database user",
		},
		cli.StringFlag{
			Name:   "db-password",
			Value:  "postgres",
			EnvVar: "DB_PASSWORD",
			Usage:  "PostgreSQL database password",
		},
	}
	app.Action = func(c *cli.Context) {

		// Configure logging
		log := logrus.WithField("context", "main")
		if c.Bool("debug") {
			logrus.SetLevel(logrus.DebugLevel)
		}

		// Connect to the database
		if err := db.Connect(
			c.String("db-name"),
			c.String("db-user"),
			c.String("db-password"),
			c.String("db-host"),
			c.Int("db-port"),
		); err != nil {
			log.Error(err)
			return
		}

		// Perform all pending migrations
		log.Info("performing pending migrations...")
		if err := db.Migrate(); err != nil {
			log.Error(err)
			return
		}
	}
	app.Run(os.Args)
}
