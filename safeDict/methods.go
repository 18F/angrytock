package safeDict

func (dict *SafeDict) Get(key string) string {
	dict.readChannel <- key
	return <-dict.readChannel
}

func (dict *SafeDict) Update(key string, value string) {
	dict.updateChannel <- &keyValue{key, value}
}

func (dict *SafeDict) Delete(key string) {
	dict.deleteChannel <- key
}

func (dict *SafeDict) Replace(newDict map[string]string) {
	dict.replaceChannel <- newDict
}

func DestorySaftDict(dict *SafeDict) {
	close(dict.readChannel)
	close(dict.updateChannel)
	close(dict.deleteChannel)
	close(dict.replaceChannel)
}
