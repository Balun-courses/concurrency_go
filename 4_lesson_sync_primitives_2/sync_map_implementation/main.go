package main

import (
	"sync"
	"sync/atomic"
)

// expunged pointer that marks entries
// which have been deleted from the dirty, but present in read
var expunged = new(string)

type entry struct {
	pointer atomic.Pointer[string]
}

func newEntry(value string) *entry {
	e := &entry{}
	e.pointer.Store(&value)
	return e
}

func (e *entry) load() (string, bool) {
	pointer := e.pointer.Load()
	if pointer == nil || pointer == expunged {
		return "", false
	}

	return *pointer, true
}

func (e *entry) delete() {
	for {
		pointer := e.pointer.Load()
		if pointer == nil || pointer == expunged {
			return
		}

		if e.pointer.CompareAndSwap(pointer, nil) {
			return
		}
	}
}

func (e *entry) set(value *string) bool {
	for {
		pointer := e.pointer.Load()
		if pointer == expunged {
			return false
		}

		if e.pointer.CompareAndSwap(pointer, value) {
			return true
		}
	}
}

func (e *entry) setLocked(value *string) {
	e.pointer.Store(value)
}

func (e *entry) expungeLocked() bool {
	pointer := e.pointer.Load()
	for pointer == nil {
		if e.pointer.CompareAndSwap(nil, expunged) {
			return true
		}

		pointer = e.pointer.Load()
	}

	return pointer == expunged
}

func (e *entry) unexpungeLocked() bool {
	return e.pointer.CompareAndSwap(expunged, nil)
}

type readOnly struct {
	data map[string]*entry

	// amended true if the dirty contains some key not in read
	amended bool
}

type Map struct {
	read atomic.Pointer[readOnly]

	mutex sync.Mutex
	dirty map[string]*entry

	// misses counts the number of loads from dirty
	misses int
}

func (m *Map) Load(key string) (string, bool) {
	read := m.loadLocalReadOnly()
	e, found := read.data[key]
	if !found && read.amended {
		m.mutex.Lock()
		read = m.loadLocalReadOnly()
		e, found = read.data[key]
		if !found && read.amended {
			e, found = m.dirty[key]
			m.missLocked()
		}

		m.mutex.Unlock()
	}

	if !found {
		return "", false
	}

	return e.load()
}

func (m *Map) Range(action func(key, value string)) {
	if action == nil {
		return
	}

	read := m.loadLocalReadOnly()
	if read.amended {
		m.mutex.Lock()
		if read = m.loadLocalReadOnly(); read.amended {
			m.promoteLocked()
		}
		m.mutex.Unlock()
	}

	for k, e := range read.data {
		value, ok := e.load()
		if !ok {
			continue
		}

		action(k, value)
	}
}

func (m *Map) Delete(key string) {
	read := m.loadLocalReadOnly()
	e, found := read.data[key]
	if !found && read.amended {
		m.mutex.Lock()
		read = m.loadLocalReadOnly()
		e, found = read.data[key]
		if !found && read.amended {
			if e, found = m.dirty[key]; found {
				delete(m.dirty, key)
				m.missLocked()
			}
		}

		m.mutex.Unlock()
	}

	if found {
		e.delete()
	}
}

func (m *Map) Store(key, value string) {
	read := m.loadLocalReadOnly()
	if e, found := read.data[key]; found {
		if ok := e.set(&value); ok {
			return
		}
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	read = m.loadLocalReadOnly()
	if e, found := read.data[key]; found {
		if e.unexpungeLocked() {
			m.dirty[key] = e
		}

		e.setLocked(&value)
	} else if e, found = m.dirty[key]; found {
		e.setLocked(&value)
	} else {
		if !read.amended {
			m.copyLocked()
			m.read.Store(&readOnly{data: read.data, amended: true})
		}

		m.dirty[key] = newEntry(value)
	}
}

func (m *Map) loadLocalReadOnly() readOnly {
	if pointer := m.read.Load(); pointer != nil {
		return *pointer
	}

	return readOnly{}
}

func (m *Map) copyLocked() {
	if m.dirty != nil {
		return
	}

	read := m.loadLocalReadOnly()
	m.dirty = make(map[string]*entry, len(read.data))

	for k, e := range read.data {
		if !e.expungeLocked() {
			m.dirty[k] = e
		}
	}
}

func (m *Map) missLocked() {
	m.misses++
	if m.misses < len(m.dirty) {
		return
	}

	m.promoteLocked()
}

func (m *Map) promoteLocked() {
	m.read.Store(&readOnly{data: m.dirty})
	m.dirty = nil
	m.misses = 0
}
