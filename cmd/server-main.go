package cmd

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/urfave/cli"
)

func serverMain(ctx *cli.Context) {
	signal.Notify(globalOSSignalCh, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

	fmt.Println("serverMain")

	<-globalOSSignalCh
}

func handleSignals() {
	// Custom exit function
	exit := func(success bool) {
		// If global profiler is set stop before we exit.
		globalProfilerMu.Lock()
		defer globalProfilerMu.Unlock()
		for _, p := range globalProfiler {
			p.Stop()
		}

		if success {
			os.Exit(0)
		}

		os.Exit(1)
	}

	stopProcess := func() bool {
		var err, oerr error

		// send signal to various go-routines that they need to quit.
		cancelGlobalContext()

		if globalNotificationSys != nil {
			globalNotificationSys.RemoveAllRemoteTargets()
		}

		if httpServer := newHTTPServerFn(); httpServer != nil {
			err = httpServer.Shutdown()
			if !errors.Is(err, http.ErrServerClosed) {
				logger.LogIf(context.Background(), err)
			}
		}

		if objAPI := newObjectLayerFn(); objAPI != nil {
			oerr = objAPI.Shutdown(context.Background())
			logger.LogIf(context.Background(), oerr)
		}

		if srv := newConsoleServerFn(); srv != nil {
			logger.LogIf(context.Background(), srv.Shutdown())
		}

		return (err == nil && oerr == nil)
	}

	for {
		select {
		case <-globalHTTPServerErrorCh:
			exit(stopProcess())
		case osSignal := <-globalOSSignalCh:
			logger.Info("Exiting on signal: %s", strings.ToUpper(osSignal.String()))
			exit(stopProcess())
		case signal := <-globalServiceSignalCh:
			switch signal {
			case serviceRestart:
				logger.Info("Restarting on service signal")
				stop := stopProcess()
				rerr := restartProcess()
				logger.LogIf(context.Background(), rerr)
				exit(stop && rerr == nil)
			case serviceStop:
				logger.Info("Stopping on service signal")
				exit(stopProcess())
			}
		}
	}
}
