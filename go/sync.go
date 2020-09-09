package main

import "sync"

var (
	ItemSyncMap       sync.Map
	ItemSyncMapForBuy sync.Map

	maxConn = MaxConnection{}
)

type MaxConnection struct {
	CurrentConnCount int
	Mux              sync.Mutex
}

const (
	MAX_CONNECTION int = 250
)

func LockItem(itemID int64) bool {
	_, loaded := ItemSyncMap.LoadOrStore(itemID, 1)
	return !loaded
}

func UnlockItem(itemID int64) {
	ItemSyncMap.Delete(itemID)
}

func LockItemForBuy(itemID int64) bool {
	_, loaded := ItemSyncMapForBuy.LoadOrStore(itemID, 1)
	return !loaded
}

func UnlockItemForBuy(itemID int64) {
	ItemSyncMapForBuy.Delete(itemID)
}

func WaitConnection() bool {
	maxConn.Mux.Lock()
	if maxConn.CurrentConnCount < MAX_CONNECTION {
		maxConn.CurrentConnCount += 1
		maxConn.Mux.Unlock()
		return true
	}
	return false
}
