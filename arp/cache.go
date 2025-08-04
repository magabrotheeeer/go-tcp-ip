package arp

type ARPCache struct {
	cache map[string][]byte // кэш для ip и mac
}

// создает новый arpcache
func NewARPCache() *ARPCache {
	return &ARPCache{
		cache: make(map[string][]byte),
	}
}

// добавляет новый ip-mac
func (ac *ARPCache) Add(ip string, mac []byte) bool {
	if len(mac) != 6 {
		return false
	}
	ac.cache[ip] = mac

	return true
}

// возвращает мак для ip
func (ac *ARPCache) Get(ip string) ([]byte, bool) {
	res, ok := ac.cache[ip]

	return res, ok
}

// удаляет ip-mac
func (ac *ARPCache) Remove(ip string) {
	delete(ac.cache, ip)
}
