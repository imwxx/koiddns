package dns

import (
	"sync"
)

// RecordCache 按 vendor、子域名、主域名、记录类型缓存 DNS 记录信息。
type RecordCache struct {
	sync.RWMutex
	cache map[string]*RecordCacheItem
}

// RecordCacheItem represents a single DNS record in the cache.
type RecordCacheItem struct {
	RecordId    string
	RecordValue string
}

var recordCache = &RecordCache{
	cache: make(map[string]*RecordCacheItem),
}

// GetRecord retrieves a DNS record from the cache.
func GetRecord(key string) (*RecordCacheItem, bool) {
	recordCache.RLock()
	defer recordCache.RUnlock()
	item, found := recordCache.cache[key]
	return item, found
}

// SetRecord adds or updates a DNS record in the cache.
func SetRecord(key, recordId, recordValue string) {
	recordCache.Lock()
	defer recordCache.Unlock()
	recordCache.cache[key] = &RecordCacheItem{
		RecordId:    recordId,
		RecordValue: recordValue,
	}
}

// ClearCache clears the DNS cache.
func ClearCache() {
	recordCache.Lock()
	defer recordCache.Unlock()
	recordCache.cache = make(map[string]*RecordCacheItem)
}

func generateCacheKey(vendor, subDomain, primaryDomain, recordType string) string {
	return vendor + ":" + subDomain + ":" + primaryDomain + ":" + recordType
}
