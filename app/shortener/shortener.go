package shortener

var store *urlStore

func shortener() {
	store = newStore()

	store.cacheKeys()
}

func Put(url string) (string, error) {
	return store.Put(url)
}

func Get(key string) (string, error) {
	return store.Get(key)
}

func init() {
	go shortener()
}
