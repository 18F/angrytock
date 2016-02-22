// Package safeDict contains the structs and methods for creates a thread saft map
package safeDict

// keyValue struct stores the a key-value pair for use in a channel
type keyValue struct {
	key   string
	value string
}

// SafeDict struct is a data type that contains a map along with channels
// that allow users to preform CRUD operations. The main difference between
// the SafeDict and a map is that all operations are serialized and thus thread safe.
type SafeDict struct {
	storage        map[string]string
	readChannel    chan string
	updateChannel  chan *keyValue
	deleteChannel  chan string
	replaceChannel chan map[string]string
}

// newSafeDict initalizes a new SafeDict and opens all the channels.
func newSafeDict() *SafeDict {
	return &SafeDict{
		make(map[string]string),
		make(chan string),
		make(chan *keyValue),
		make(chan string),
		make(chan map[string]string),
	}
}

// InitSafeDict initalizes a new SafeDict and creates for loop that listens for
// new operations.
func InitSafeDict() *SafeDict {
	safeDict := newSafeDict()
	go func() {
		for {
			select {
			case key := <-safeDict.readChannel:
				safeDict.readChannel <- safeDict.storage[key]
			case keyValuePair := <-safeDict.updateChannel:
				safeDict.storage[keyValuePair.key] = keyValuePair.value
			case key := <-safeDict.deleteChannel:
				delete(safeDict.storage, key)
			case newDict := <-safeDict.replaceChannel:
				safeDict.storage = newDict
			}
		}
	}()
	return safeDict
}

// DestorySaftDict closes all the dictionary channels
func DestorySaftDict(dict *SafeDict) {
	close(dict.readChannel)
	close(dict.updateChannel)
	close(dict.deleteChannel)
	close(dict.replaceChannel)
}
