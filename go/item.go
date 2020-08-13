package main

import "sync"

var (
	itemMap sync.Map
)

func LockItem(itemID int64) {
	itemMap.Store(itemID, true)
}

func UnLockItem(itemID int64) {
	itemMap.Delete(itemID)
}

func CheckOrLockItem(itemID int64) bool {
	_, loaded := itemMap.LoadOrStore(itemID, true)
	return loaded
}
