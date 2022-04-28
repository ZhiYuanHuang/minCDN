package cmd

import (
	"context"
	"errors"
	"fmt"
	"log"
	"runtime"
	"strings"
	"sync"
	"time"

	pathutil "path"
)

type RWLocker interface {
	GetLock(ctx context.Context, timeout *dynamicTimeout) (lkCtx LockContext, timedOutErr error)
	Unlock(cancel context.CancelFunc)
	GetRLock(ctx context.Context, timeout *dynamicTimeout) (lkCtx LockContext, timedOutErr error)
	RUnlock(cancel context.CancelFunc)
}

type LockContext struct {
	ctx    context.Context
	cancel context.CancelFunc
}

func (l LockContext) Context() context.Context {
	return l.ctx
}

func (l LockContext) Cancel() {
	if l.cancel != nil {
		l.cancel()
	}
}

type nsLock struct {
	ref int32
	*LRWMutex
}

type nsLockMap struct {
	lockMap      map[string]*nsLock
	lockMapMutex sync.Mutex
}

func newNSLock() *nsLockMap {
	var nsMutex nsLockMap
	nsMutex.lockMap = make(map[string]*nsLock)
	return &nsMutex
}

func (n *nsLockMap) lock(ctx context.Context, uri string, lockSource, opsID string, readLock bool, timeout time.Duration) (locked bool) {
	resource := uri

	n.lockMapMutex.Lock()
	nsLk, found := n.lockMap[resource]
	if !found {
		nsLk = &nsLock{
			LRWMutex: NewLRWMutex(),
		}
	}
	nsLk.ref++
	n.lockMap[resource] = nsLk
	n.lockMapMutex.Unlock()

	if readLock {
		locked = nsLk.GetRLock(ctx, opsID, lockSource, timeout)
	} else {
		locked = nsLk.GetLock(ctx, opsID, lockSource, timeout)
	}

	if !locked {
		n.lockMapMutex.Lock()
		n.lockMap[resource].ref--
		if n.lockMap[resource].ref == 0 {
			log.Fatal(errors.New("resource reference count was lower than 0"))
		}
		if n.lockMap[resource].ref == 0 {
			delete(n.lockMap, resource)
		}
		n.lockMapMutex.Unlock()
	}

	return
}

func (n *nsLockMap) unlock(uri string, readLock bool) {
	resource := uri

	n.lockMapMutex.Lock()
	defer n.lockMapMutex.Unlock()
	if _, found := n.lockMap[resource]; !found {
		return
	}
	if readLock {
		n.lockMap[resource].RUnlock()
	} else {
		n.lockMap[resource].Unlock()
	}
	n.lockMap[resource].ref--
	if n.lockMap[resource].ref < 0 {
		log.Fatal(errors.New("resource reference count was lower than 0"))
	}
	if n.lockMap[resource].ref == 0 {
		delete(n.lockMap, resource)
	}
}

type localLockInstance struct {
	ns    *nsLockMap
	uri   string
	opsID string
}

func (n *nsLockMap) NewNSLock(uri string) RWLocker {
	opsID := mustGetUUID()
	return &localLockInstance{n, uri, opsID}
}

func (li *localLockInstance) GetLock(ctx context.Context, timeout *dynamicTimeout) (_ LockContext, timedOutErr error) {
	lockSource := getSource(2)
	start := UTCNow()
	const readLock = false
	if !li.ns.lock(ctx, li.uri, lockSource, li.opsID, readLock, timeout.Timeout()) {
		timeout.LogFailure()
		switch err := ctx.Err(); err {
		case context.Canceled:
			return LockContext{}, err
		}
		return LockContext{}, OperationTimedOut{}
	}
	timeout.LogSuccess(UTCNow().Sub(start))
	return LockContext{ctx: ctx, cancel: func() {}}, nil
}

func (li *localLockInstance) Unlock(cancel context.CancelFunc) {
	if cancel != nil {
		cancel()
	}
	const readLock = false
	li.ns.unlock(li.uri, readLock)
}

func (li *localLockInstance) GetRLock(ctx context.Context, timeout *dynamicTimeout) (_ LockContext, timedOutErr error) {
	lockSource := getSource(2)
	start := UTCNow()
	const readLock = true
	if !li.ns.lock(ctx, li.uri, lockSource, li.opsID, readLock, timeout.Timeout()) {
		timeout.LogFailure()
		switch err := ctx.Err(); err {
		case context.Canceled:
			return LockContext{}, err
		}
		return LockContext{}, OperationTimedOut{}
	}
	timeout.LogSuccess(UTCNow().Sub(start))
	return LockContext{ctx: ctx, cancel: func() {}}, nil
}

func (li *localLockInstance) RUnlock(cancel context.CancelFunc) {
	if cancel != nil {
		cancel()
	}
	const readLock = true
	li.ns.unlock(li.uri, readLock)
}

func getSource(n int) string {
	var funcName string
	pc, filename, lineNum, ok := runtime.Caller(n)
	if ok {
		filename = pathutil.Base(filename)
		funcName = strings.TrimPrefix(runtime.FuncForPC(pc).Name(),
			"github.com/ZhiYuanHuang/minCDN/cmd.")
	} else {
		filename = "<unknown>"
		lineNum = 0
	}

	return fmt.Sprintf("[%s:%d:%s()]", filename, lineNum, funcName)
}
