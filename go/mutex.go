package main

import (
	"context"
	"sync"
	"time"
)

var (
	itemMutex        = map[int64]*sync.Mutex{}
	lockWaitDuration = 20 * time.Millisecond
)

func LockItem(itemID int64) bool {
	ctx, cancel := context.WithTimeout(context.Background(), lockWaitDuration)
	defer cancel()
	ch := make(chan bool)
	go func() {
		itemMutex[itemID].Lock()
		ch <- true
	}()
	select {
	case <-ch:
		return true
	case <-ctx.Done():
		return false
	}

}

func UnlockItem(itemID int64) {
	itemMutex[itemID].Unlock()
}
