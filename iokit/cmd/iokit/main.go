package main

import (
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"sudachen.xyz/pkg/go-forge/fu"
	"sudachen.xyz/pkg/go-forge/iokit"
)

func main() {
	app := &cli.App{
		Name:     "iokit",
		HelpName: "iokit",
		Usage:    "go-iokit package utility",
		Commands: []*cli.Command{
			{
				Name:      "encrypt",
				Usage:     "encrypt credentials file",
				UsageText: "encrypt <password> <filename>",
				Action: func(c *cli.Context) error {
					passwd := c.Args().Get(0)
					file := c.Args().Get(1)
					data1, err := iokit.File(file).ReadAll()
					if err != nil {
						return err
					}
					data, err := fu.Encrypt(passwd, data1)
					if err != nil {
						return err
					}
					err = iokit.File(file + ".enc").WriteAll(data)
					if err != nil {
						return err
					}
					return nil
				},
			}},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
