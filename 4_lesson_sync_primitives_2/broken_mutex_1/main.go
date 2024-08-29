package main

const (
	unlocked = false
	locked   = true
)

type BrokenMutex struct {
	state bool
}

// Здесь есть data race и нет гарантии взаимного исключения (safety),
// так как несколько горутин могут попасть совместно в критическую секцию

func (m *BrokenMutex) Lock() {
	for m.state {
		// iteration by iteration...
	}

	m.state = locked
}

func (m *BrokenMutex) Unlock() {
	m.state = unlocked
}
