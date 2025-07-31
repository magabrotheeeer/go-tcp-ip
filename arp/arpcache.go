package arp

type ARPCache struct {
	cache map[string][]byte
}

func NewARPCache() *ARPCache {
	return &ARPCache{
		cache: make(map[string][]byte),
	}
}

type ManagerCache interface {
	Add(ip string, mac []byte) bool
	Get(ip string) ([]byte, bool)
	Remove (ip string)
}

func (ac *ARPCache) Add(ip string, mac []byte) bool {
	if len(mac) != 6 {
		return false
	}
	ac.cache[ip] = mac

	return true
}


func (ac *ARPCache) Get(ip string) ([]byte, bool) {
	res, ok := ac.cache[ip]

	return res, ok
}

func (ac *ARPCache) Remove(ip string) {
	delete(ac.cache, ip)
}