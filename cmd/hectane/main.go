package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/hectane/hectane/db"
	"github.com/hectane/hectane/server"
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
		cli.StringFlag{
			Name:   "secret-key",
			EnvVar: "SECRET_KEY",
			Usage:  "key used for auth cookies",
		},
		cli.StringFlag{
			Name:   "web-addr",
			Value:  ":8000",
			EnvVar: "WEB_ADDR",
			Usage:  "address for the web interface",
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

		// Create the web server
		s, err := server.New(&server.Config{
			Addr:      c.String("web-addr"),
			SecretKey: c.String("secret-key"),
		})
		if err != nil {
			log.Error(err)
			return
		}
		defer s.Close()

		// Wait for SIGINT or SIGTERM
		sigChan := make(chan os.Signal)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
	}
	app.Run(os.Args)
}
