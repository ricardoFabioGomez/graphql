package cache

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"regexp"
	"sync"

	"github.com/dgraph-io/ristretto"
)

type Cache interface {
	Get(query string) bool
	Set(query string)
}

type cache struct {
	mux    sync.Mutex
	client *ristretto.Cache
}

type itemCached struct{}

var expOfuscateUser = regexp.MustCompile("[0-9]+")

func NewCache() Cache {
	client, _ := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1000, // Num keys to track frequency of (10M).
		MaxCost:     100,  // Maximum cost of cache (1GB).
		BufferItems: 64,   // Number of keys per Get buffer.
	})
	return &cache{
		client: client,
	}
}

func getKey(query string) string {
	b, _ := json.Marshal(expOfuscateUser.ReplaceAllString(query, "0"))
	hash := md5.Sum(b)
	return hex.EncodeToString(hash[:])
}

func (c *cache) Get(query string) bool {
	key := getKey(query)
	_, ok := c.client.Get(key)
	return ok
}

func (c *cache) Set(query string) {
	c.mux.Lock()
	defer c.mux.Unlock()
	key := getKey(query)
	c.client.Set(key, itemCached{}, 1)
	c.client.Wait()
}
