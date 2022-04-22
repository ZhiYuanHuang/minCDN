package cmd

import (
	"errors"
	"log"
	"net/http"
	"os"
	"strings"
)

func handleSignals() {
	// Custom exit function
	exit := func(success bool) {
		// If global profiler is set stop before we exit.

		if success {
			os.Exit(0)
		}

		os.Exit(1)
	}

	stopProcess := func() bool {
		var err, oerr error

		// send signal to various go-routines that they need to quit.
		cancelGlobalContext()

		if httpServer := newHTTPServerFn(); httpServer != nil {
			err = httpServer.Shutdown()
			if !errors.Is(err, http.ErrServerClosed) {
				log.Println(err)
			}
		}

		return (err == nil && oerr == nil)
	}

	for {
		select {
		// case <-globalHTTPServerErrorCh:
		// 	exit(stopProcess())
		case osSignal := <-globalOSSignalCh:
			log.Println("Exiting on signal: %s", strings.ToUpper(osSignal.String()))
			exit(stopProcess())
		case signal := <-globalServiceSignalCh:
			switch signal {
			// case serviceRestart:
			// 	logger.Info("Restarting on service signal")
			// 	stop := stopProcess()
			// 	rerr := restartProcess()
			// 	logger.LogIf(context.Background(), rerr)
			// 	exit(stop && rerr == nil)
			case serviceStop:
				log.Println("Stopping on service signal")
				exit(stopProcess())
			}
		}
	}
}
