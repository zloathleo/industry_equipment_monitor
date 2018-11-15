package dstruct

import "sync"

//线程安全实时值Map
type ValueMap struct {
	values map[string]float64
	*sync.RWMutex
}

func NewValueMap() *ValueMap {
	return &ValueMap{make(map[string]float64), new(sync.RWMutex)}
}

func (m *ValueMap) Set(pn string ,value float64) {
	m.Lock()
	m.values[pn] = value
	m.Unlock()
}

func (m *ValueMap) Clear() {
	m.values = make(map[string]float64)
}

func (m *ValueMap) SafeCopyAndClear() map[string]float64 {
	m.Lock()
	count := len(m.values)
	cacheCopy := make(map[string]float64, count)
	if count != 0 {
		for key, value := range m.values {
			cacheCopy[key] = value
		}
	}
	m.values = make(map[string]float64)
	m.Unlock()
	return cacheCopy
}
