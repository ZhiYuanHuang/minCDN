package cmd

import (
	"sync"
)

var globalObjLayerMutex sync.RWMutex

var globalObjectAPI ObjectLayer
