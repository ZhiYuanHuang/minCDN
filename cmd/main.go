package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/urfave/cli"
)

var GlobalFlags = []cli.Flag{
	cli.BoolFlag{
		Name:  "UseEtcd",
		Usage: "use etcd to register minCDN service",
	},
	cli.StringFlag{
		Name:   "EtcdAddress",
		Usage:  "etcd Address",
		EnvVar: "MINCDN_ETCD_ADDRESS",
	},
	cli.StringFlag{
		Name:   "address",
		Value:  ":" + GlobalMinCDNDefaultPort,
		Usage:  "bind to service address",
		EnvVar: "MINCDN_ADDRESS",
	},
}

func newApp(name string) *cli.App {

	app := cli.NewApp()
	app.Name = name
	app.Usage = "min static file CDN"
	app.Description = "Build min static CDN"
	app.Flags = GlobalFlags
	app.Action = func(c *cli.Context) error {
		useEtcd := c.Bool("UseEtcd")
		if useEtcd {
			fmt.Println("use etcd to register service")
		}
		fmt.Println("")
		return nil
	}

	return app
}

func Main(args []string) {
	appName := filepath.Base(args[0])

	if err := newApp(appName).Run(args); err != nil {
		os.Exit(1)
	}
}
