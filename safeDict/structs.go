package safeDict

type keyValue struct {
	key   string
	value string
}

type SafeDict struct {
	storage        map[string]string
	readChannel    chan string
	updateChannel  chan *keyValue
	deleteChannel  chan string
	replaceChannel chan map[string]string
}

func newSafeDict() *SafeDict {
	return &SafeDict{
		make(map[string]string),
		make(chan string),
		make(chan *keyValue),
		make(chan string),
		make(chan map[string]string),
	}
}

func InitSafeDict() *SafeDict {
	manager := newSafeDict()
	go func() {
		for {
			select {
			case key := <-manager.readChannel:
				manager.readChannel <- manager.storage[key]
			case keyValuePair := <-manager.updateChannel:
				manager.storage[keyValuePair.key] = keyValuePair.value
			case key := <-manager.deleteChannel:
				delete(manager.storage, key)
			case newDict := <-manager.replaceChannel:
				manager.storage = newDict
			}
		}
	}()
	return manager
}
