package legendasdivx

import (
	"time"

	"github.com/patrickmn/go-cache"
)

type SubtitleCache struct {
	subtitles *cache.Cache
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

func (sc *SubtitleCache) Get(imdbId, id string) (interface{}, bool) {
	key := imdbId + id

	return sc.subtitles.Get(key)
}

func (sc *SubtitleCache) Set(imdbId, id string, subtitle interface{}) {
	key := imdbId + id

	sc.subtitles.Set(key, subtitle, cache.DefaultExpiration)
}
