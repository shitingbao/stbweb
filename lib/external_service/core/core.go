package core

import (
	"log"
	"sync"
	"time"
)

var (
	UserHub *Hub
)

type Hub struct {
	data map[string]time.Time
	lock *sync.RWMutex
}

func init() {
	UserHub = &Hub{
		data: make(map[string]time.Time),
		lock: new(sync.RWMutex),
	}
}

func (h *Hub) PutData(key string) {
	h.lock.Lock()
	defer h.lock.Unlock()
	h.data[key] = time.Now()
	log.Println("PutData:", key)
}

func (h *Hub) DeleteData(key string) {
	h.lock.Lock()
	defer h.lock.Unlock()
	delete(h.data, key)
	log.Println("delete:", key)
}

func (h *Hub) Length() int {
	h.lock.Lock()
	defer h.lock.Unlock()
	return len(h.data)
}

func (h *Hub) GetData(key string) time.Time {
	h.lock.RLock()
	defer h.lock.RUnlock()
	log.Println("GetData:", key)

	return h.data[key]
}
