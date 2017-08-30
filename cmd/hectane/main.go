package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/hectane/hectane/db"
	"github.com/hectane/hectane/queue"
	"github.com/hectane/hectane/server"
	"github.com/hectane/hectane/storage"
	"github.com/howeyc/gopass"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func initDB(c *cli.Context) error {
	if err := db.Connect(
		c.GlobalString("db-driver"),
		c.GlobalString("db-name"),
		c.GlobalString("db-user"),
		c.GlobalString("db-password"),
		c.GlobalString("db-host"),
		c.GlobalInt("db-port"),
	); err != nil {
		return err
	}
	if err := db.Migrate(); err != nil {
		return err
	}
	return nil
}

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
			Name:   "db-driver",
			Value:  "postgres",
			EnvVar: "DB_DRIVER",
			Usage:  "database driver",
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
			Name:   "queue-addr",
			Value:  ":smtp",
			EnvVar: "QUEUE_ADDR",
			Usage:  "address for incoming SMTP connections",
		},
		cli.StringFlag{
			Name:   "storage-directory",
			Value:  ".",
			EnvVar: "STORAGE_DIRECTORY",
			Usage:  "directory for storing email content",
		},
	}
	app.Before = func(c *cli.Context) error {
		if c.Bool("debug") {
			logrus.SetLevel(logrus.DebugLevel)
		}
		return nil
	}
	app.Commands = []cli.Command{
		{
			Name:  "createadmin",
			Usage: "create an administrator user",
			Action: func(c *cli.Context) error {

				// Initialize the database
				if err := initDB(c); err != nil {
					return err
				}

				// Prompt for the username
				var username string
				fmt.Print("Username? ")
				fmt.Scanln(&username)

				// Prompt for the password, hiding the input
				fmt.Print("Password? ")
				p, err := gopass.GetPasswd()
				if err != nil {
					return err
				}

				// Generate a new user with the data
				u := &db.User{
					Username: username,
					IsAdmin:  true,
				}
				if err := u.SetPassword(string(p)); err != nil {
					return err
				}

				// Store the user in the database
				if err := db.C.Create(u).Error; err != nil {
					return err
				}

				return nil
			},
		},
		{
			Name:  "run",
			Usage: "run the application",
			Flags: []cli.Flag{
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
			},
			Action: func(c *cli.Context) error {

				// Initialize the storage backend
				st := storage.New(&storage.Config{
					Directory: c.String("storage-directory"),
				})

				// Create the incoming mail queue
				q, err := queue.New(&queue.Config{
					Addr:    c.String("queue-addr"),
					Storage: st,
				})
				if err != nil {
					return err
				}
				defer q.Close()

				// Initialize the database
				if err := initDB(c); err != nil {
					return err
				}

				// Create the web server
				s, err := server.New(&server.Config{
					Addr:      c.String("web-addr"),
					SecretKey: c.String("secret-key"),
				})
				if err != nil {
					return err
				}
				defer s.Close()

				// Wait for SIGINT or SIGTERM
				sigChan := make(chan os.Signal)
				signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
				<-sigChan

				return nil
			},
		},
	}

	// Watch for an error
	log := logrus.WithField("context", "main")
	if err := app.Run(os.Args); err != nil {
		log.Error(err)
		os.Exit(1)
	}
}
