package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/hectane/hectane/db"
	"github.com/hectane/hectane/imap"
	"github.com/hectane/hectane/receiver"
	"github.com/hectane/hectane/server"
	"github.com/hectane/hectane/storage"
	"github.com/howeyc/gopass"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func globalInit(c *cli.Context) error {

	// Enable debug logging if parameter set
	if c.GlobalBool("debug") {
		logrus.SetLevel(logrus.DebugLevel)
	}

	// Initialize the database
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

	// Migrate the database
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
			Name:   "storage-directory",
			Value:  "data",
			EnvVar: "STORAGE_DIRECTORY",
			Usage:  "directory for storing email content",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:  "createadmin",
			Usage: "create an administrator user",
			Action: func(c *cli.Context) error {

				if err := globalInit(c); err != nil {
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
					Name:   "imap-addr",
					Value:  ":imap",
					EnvVar: "IMAP_ADDR",
					Usage:  "address for IMAP connections",
				},
				cli.StringFlag{
					Name:   "receiver-addr",
					Value:  ":smtp",
					EnvVar: "RECEIVER_ADDR",
					Usage:  "address for incoming SMTP connections",
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
			},
			Action: func(c *cli.Context) error {

				if err := globalInit(c); err != nil {
					return err
				}

				// Initialize the storage backend
				st := storage.New(&storage.Config{
					Directory: c.GlobalString("storage-directory"),
				})

				// Create the incoming mail receiver
				r, err := receiver.New(&receiver.Config{
					Addr:    c.String("receiver-addr"),
					Storage: st,
				})
				if err != nil {
					return err
				}
				defer r.Close()

				// Create the IMAP server
				i, err := imap.New(&imap.Config{
					Addr:    c.String("imap-addr"),
					Storage: st,
				})
				if err != nil {
					return err
				}
				defer i.Close()

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
