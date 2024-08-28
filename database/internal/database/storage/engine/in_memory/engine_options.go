package in_memory

type EngineOption func(*Engine)

func WithPartitions(partitionsNumber uint) EngineOption {
	return func(engine *Engine) {
		engine.partitions = make([]*HashTable, partitionsNumber)
		for i := 0; i < int(partitionsNumber); i++ {
			engine.partitions[i] = NewHashTable()
		}
	}
}
