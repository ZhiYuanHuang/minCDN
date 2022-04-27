package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	xhttp "github.com/ZhiYuanHuang/minCDN/internal/http"

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

	spileIndex := strings.Index(addr, ":")
	if spileIndex != -1 && spileIndex < len(addr)-1 {
		port := addr[spileIndex+1:]
		globalMinCDNPort = port
	}

	handler, err := configureServerHandler()
	if err != nil {
		log.Fatal("unable to configure server handler")
	}

	addrs := make([]string, 0, 1)
	addrs = append(addrs, globalMinCDNAddr)

	httpServer := xhttp.NewServer(addrs).
		UseHandler(corsHandler(handler)).
		UseBaseContext(GlobalContext).
		UseCustomLogger(log.New(ioutil.Discard, "", 0))

	go func() {
		globalHTTPServerErrorCh <- httpServer.Start(GlobalContext)
	}()

	setHTTPServer(httpServer)
}
