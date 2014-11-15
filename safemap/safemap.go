//
// The thread safe of map
//
package safemap

import (
	"sync"
)

type SafeMap struct {
	lock *sync.RWMutex
	mp   map[interface{}]interface{}
}

// Return new safemap
func New() *SafeMap {
	return &SafeMap{
		lock: new(sync.RWMutex),
		mp:   make(map[interface{}]interface{}),
	}
}

// Return the key's value
func (this *SafeMap) Get(key interface{}) (interface{}, bool) {
	this.lock.RLock()
	defer this.lock.RUnlock()
	value, ok := this.mp[key]
	return value, ok
}

// Set (key --> value)
func (this *SafeMap) Set(key, value interface{}) {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.mp[key] = value
}

// Return true if key is exist in the safemap
func (this *SafeMap) IsExist(key interface{}) bool {
	this.lock.RLock()
	defer this.lock.RUnlock()
	_, ok := this.mp[key]
	return ok
}

// Delete (key --> value)
func (this *SafeMap) Delete(key interface{}) {
	this.lock.Lock()
	defer this.lock.Unlock()
	delete(this.mp, key)
}

// Return the size of safemap
func (this *SafeMap) Size() int {
	this.lock.RLock()
	defer this.lock.RUnlock()
	return len(this.mp)
}

// Return all items of safemap
func (this *SafeMap) Items() map[interface{}]interface{} {
	this.lock.RLock()
	defer this.lock.RUnlock()
	m := make(map[interface{}]interface{})
	for k, v := range this.mp {
		m[k] = v
	}
	return m
}
