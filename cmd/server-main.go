package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/urfave/cli"
)

func serverMain(ctx *cli.Context) {
	signal.Notify(globalOSSignalCh, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

	fmt.Println("serverMain")

	<-globalOSSignalCh
}
