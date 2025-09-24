package dns

import (
	"sync"
)

// var (
// 	cache = struct {
// 		sync.RWMutex
// 		Records map[string]map[string]string
// 	}{
// 		Records: make(map[string]map[string]string),
// 	}
// )

// func getRecordIdFromCache(vendor, domain, recordType string) string {
// 	cache.RLock()
// 	defer cache.RUnlock()
// 	if records, ok := cache.Records[vendor]; ok {
// 		return records[domain+recordType]
// 	}
// 	return ""
// }

// func setRecordIdToCache(vendor, domain, recordType, recordId string) {
// 	cache.Lock()
// 	defer cache.Unlock()
// 	if cache.Records[vendor] == nil {
// 		cache.Records[vendor] = make(map[string]string)
// 	}
// 	cache.Records[vendor][domain+recordType] = recordId
// }

// RecordCache stores DNS records and their associated information, considering vendor, domain, and record type.
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
