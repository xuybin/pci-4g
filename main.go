package main

import (
	"github.com/mkideal/cli"
	"github.com/xuybin/pci-4g/server"
)
const PCI_VERSION = "v0.0.1"
type cliArgs struct {
	cli.Helper
	ListenAddr    string `cli:"*l,*listen" usage:"pci listen host and port" dft:"$PCI_LS"`
	Version  bool `cli:"!v" usage:"force flag, note the !"`

}

func main() {
	cli.Run(new(cliArgs), func(ctx *cli.Context) error {
		argv := ctx.Argv().(*cliArgs)
		if argv.Version {
			ctx.String("%s\n",PCI_VERSION)
		}
		return server.NewPciServer().InitDocs().Start(argv.ListenAddr)
	})
}
