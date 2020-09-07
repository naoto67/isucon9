package main

import "sync"

var (
	ItemSyncMap sync.Map
)

func LockItem(itemID int64) bool {
	_, loaded := ItemSyncMap.LoadOrStore(itemID, 1)
	return !loaded
}

func UnlockItem(itemID int64) {
	ItemSyncMap.Delete(itemID)
}
