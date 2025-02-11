package service

import "sync"

type MerchCache struct {
	prices map[int]int
	sync.RWMutex
}

func NewMerchCache(prices map[int]int) *MerchCache {
	return &MerchCache{
		prices: prices,
	}
}

func (m *MerchCache) Get(key int) int {
	m.RLock()
	defer m.RUnlock()
	return m.prices[key]
}

func (m *MerchCache) Patch(key int, value int) {
	m.Lock()
	defer m.Unlock()
	m.prices[key] = value
}

func (m *MerchCache) Del(key int) {
	m.Lock()
	defer m.Unlock()
	delete(m.prices, key)
}
