package utils

import (
	"sync"
)

// 泛型包装器，使用sync.Map实现
type SyncMap[K comparable, V any] struct {
	Mi sync.Map
}

// 初始化
func NewSyncMap[K comparable, V any]() *SyncMap[K, V] {
	return &SyncMap[K, V]{}
}

// Store 存储键值对
func (sm *SyncMap[K, V]) Store(key K, value V) {
	sm.Mi.Store(key, value)
}

// Load 加载键对应的值
func (sm *SyncMap[K, V]) Load(key K) (V, bool) {
	value, ok := sm.Mi.Load(key)
	if !ok {
		var zero V // 返回类型的零值
		return zero, false
	}
	return value.(V), true
}

// LoadOrStore 如果键不存在，则存储键值对；如果键已存在，则返回已存在的值
func (sm *SyncMap[K, V]) LoadOrStore(key K, value V) (actual V, loaded bool) {
	actualValue, loaded := sm.Mi.LoadOrStore(key, value)
	if loaded {
		return actualValue.(V), loaded
	}
	return value, loaded
}

// Swap 替换键对应的值
func (sm *SyncMap[K, V]) Swap(key K, value V) (old V, loaded bool) {
	oldValue, loaded := sm.Mi.Swap(key, value)
	if loaded {
		old = oldValue.(V)
	}
	return old, loaded
}

// CompareAndSwap 如果键对应的值等于old，则替换为new
func (sm *SyncMap[K, V]) CompareAndSwap(key K, old, new V) bool {
	return sm.Mi.CompareAndSwap(key, old, new)
}

// LoadAndDelete 加载并删除键对应的值
func (sm *SyncMap[K, V]) LoadAndDelete(key K) (V, bool) {
	value, ok := sm.Mi.LoadAndDelete(key)
	if !ok {
		var zero V // 返回类型的零值
		return zero, false
	}
	return value.(V), true
}

// Delete 删除键对应的值
func (sm *SyncMap[K, V]) Delete(key K) {
	sm.Mi.Delete(key)
}

// Range 遍历映射
func (sm *SyncMap[K, V]) Range(f func(K, V) bool) {
	sm.Mi.Range(func(key, value interface{}) bool {
		return f(key.(K), value.(V))
	})
}

// 清空
func (sm *SyncMap[K, V]) Clear() {
	sm.Mi.Range(func(key, value interface{}) bool {
		sm.Mi.Delete(key)
		return true
	})
}

// 大小
func (sm *SyncMap[K, V]) Len() int {
	var count int
	sm.Mi.Range(func(key, value interface{}) bool {
		count++
		return true
	})
	return count
}
