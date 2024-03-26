package legendasdivx

import (
	"time"

	"github.com/patrickmn/go-cache"
)

type SubtitleCache struct {
	subtitles *cache.Cache
}

type CachedData struct {
	Name        string
	Credentials string
}

const (
	defaultExpiration = 5 * time.Minute
	cleanupInterval   = 10 * time.Minute
)

func NewSubtitleCache() *SubtitleCache {
	c := cache.New(defaultExpiration, cleanupInterval)
	return &SubtitleCache{
		subtitles: c,
	}
}

func (sc *SubtitleCache) Get(username, imdbId, id string) (*CachedData, bool) {
	key := username + imdbId + id

	data, exists := sc.subtitles.Get(key)

	if !exists {
		return &CachedData{}, false
	}

	return data.(*CachedData), true
}

func (sc *SubtitleCache) Set(username, imdbId, id string, subtitle *CachedData) {
	key := username + imdbId + id

	sc.subtitles.Set(key, subtitle, cache.DefaultExpiration)
}
