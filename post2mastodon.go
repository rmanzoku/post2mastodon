package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/mattn/go-mastodon"
	"github.com/urfave/cli"
)

type config struct {
	Profile map[string]profile
}

type profile struct {
	ClientID     string `toml:"client-id"`
	ClientSecret string `toml:"client-secret"`
	Server       string `toml:"server"`
	EMail        string `toml:"e-mail"`
	Password     string `toml:"password"`
}

func makeApp() *cli.App {
	app := cli.NewApp()
	app.Name = "post2mastodon"
	app.Usage = "Post to mastodon simply"
	app.Version = "0.1"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config",
			Usage: "config file",
			Value: os.Getenv("HOME") + "/.config/post2mastodon.toml",
		},
		cli.StringFlag{
			Name:  "profile",
			Usage: "profile name",
			Value: "default",
		},
	}
	app.Setup()
	return app
}

func run() int {

	app := makeApp()

	var config config

	app.Before = func(c *cli.Context) error {

		_, err := toml.DecodeFile(c.String("config"), &config)
		if err != nil {
			return err
		}

		_, ok := config.Profile[c.String("profile")]
		if ok != true {
			return fmt.Errorf("Profile not found")
		}

		return nil
	}

	app.Action = func(c *cli.Context) error {
		v, _ := config.Profile[c.String("profile")]

		var toot string
		fmt.Scan(&toot)

		m := mastodon.NewClient(&mastodon.Config{
			Server:       v.Server,
			ClientID:     v.ClientID,
			ClientSecret: v.ClientSecret,
		})
		err := m.Authenticate(context.Background(), v.EMail, v.Password)
		if err != nil {
			log.Fatal(err)
		}

		_, err = m.PostStatus(context.Background(), &mastodon.Toot{
			Status: toot,
		})

		if err != nil {
			log.Fatal(err)
		}

		return nil
	}
	app.Run(os.Args)
	return 0
}

func main() {
	os.Exit(run())
}
