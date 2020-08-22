package main

import (
	"context"
	"sync"
	"time"
)

var (
	itemMap          = sync.Map{}
	lockWaitDuration = 20 * time.Millisecond
)

func LockItem(itemID int64) bool {
	ctx, cancel := context.WithTimeout(context.Background(), lockWaitDuration)
	defer cancel()
	ch := make(chan bool)
	go func() {
		for {
			_, loaded := itemMap.LoadOrStore(itemID, true)
			if !loaded {
				break
			}
		}
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
	itemMap.Delete(itemID)
}
