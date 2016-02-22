package safeDict

// Get returns the value given a specific key
func (dict *SafeDict) Get(key string) string {
	dict.readChannel <- key
	return <-dict.readChannel
}

// Update sets a key to a specific value
func (dict *SafeDict) Update(key string, value string) {
	dict.updateChannel <- &keyValue{key, value}
}

// Delete removes a key-value pair given a key
func (dict *SafeDict) Delete(key string) {
	dict.deleteChannel <- key
}

// Replace replaces the internal hashmap
func (dict *SafeDict) Replace(newDict map[string]string) {
	dict.replaceChannel <- newDict
}
