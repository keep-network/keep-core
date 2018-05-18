package keep

import (
	"fmt"
	"os"
	"time"

	"github.com/urfave/cli"
)

// ClientApp is the Keep client application
func ClientApp(version, revision string) {
	app := cli.NewApp()
	app.Name = "keep-client"
	app.Version = fmt.Sprintf("%s (revision %s)", version, revision)
	app.Compiled = time.Now()
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "Keep Authors",
			Email: "noreply@example.com",
		},
	}
	app.Copyright = "(c) 2018 Thesis, Inc."
	app.HelpName = "keep-client"
	app.Usage = "The Keep Client Application"

	app.Commands = KeepCommands
	//app.Flags = KeepFlags  //TODO: Any runtime flags?
	app.Action = func(c *cli.Context) error {
		return nil
	}
	app.Run(os.Args)
}
