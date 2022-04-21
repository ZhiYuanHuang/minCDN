package cmd

import "os"

const (
	GlobalMinCDNDefaultPort = "9006"
)

var (
	globalOSSignalCh = make(chan os.Signal, 1)
)
