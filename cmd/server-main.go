package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/urfave/cli"
)

func serverMain(ctx *cli.Context) {
	signal.Notify(globalOSSignalCh, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

	go handleSignals()

	fmt.Println("serverMain")

	time.Sleep(time.Duration(10) * time.Second)

	<-globalOSSignalCh
}
