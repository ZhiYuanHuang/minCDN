package cmd

import (
	"context"
)

type serviceSignal int

const (
	serviceRestart serviceSignal = iota
	serviceStop
)

var globalServiceSignalCh chan serviceSignal

var GlobalContext context.Context

var cancelGlobalContext context.CancelFunc

func initGlobalContext() {
	GlobalContext, cancelGlobalContext = context.WithCancel(context.Background())
	globalServiceSignalCh = make(chan serviceSignal)
}
