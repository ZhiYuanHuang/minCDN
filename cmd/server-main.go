package cmd

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/urfave/cli"
)

func serverMain(ctx *cli.Context) {
	signal.Notify(globalOSSignalCh, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

	go handleSignals()

	serverHandlCmdArgs(ctx)

	fmt.Println("serverMain")

	time.Sleep(time.Duration(10) * time.Second)

	<-globalOSSignalCh
}

func serverHandlCmdArgs(cts *cli.Context) {
	globalUseETCD = cts.IsSet("UseEtcd")

	etcdAddress := cts.GlobalString("EtcdAddress")
	if globalUseETCD {
		if etcdAddress == "" {
			log.Fatalln("etcd address cann't be empty when use etcd")
		} else {
			globalETCDAddress = etcdAddress
		}
	}

	addr := cts.GlobalString("address")
	globalMinCDNAddr = addr
}
