package memcache

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/url"
	"sync/atomic"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/whosonfirst/go-cache"
	"github.com/whosonfirst/go-ioutil"
)

type MemcacheCache struct {
	cache.Cache
	misses    int64
	hits      int64
	evictions int64
	client    *memcache.Client
}

func init() {
	ctx := context.Background()
	cache.RegisterCache(ctx, "memcache", NewMemcacheCache)
}

func NewMemcacheCache(ctx context.Context, uri string) (cache.Cache, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse URI, %w", err)
	}

	q := u.Query()

	hosts := q["host"]

	if len(hosts) == 0 {
		return nil, fmt.Errorf("Missing hosts, %w", err)
	}

	client := memcache.New(hosts...)

	err = client.Ping()

	if err != nil {
		return nil, fmt.Errorf("Failed to ping client, %w", err)
	}

	c := &MemcacheCache{
		client:    client,
		hits:      int64(0),
		misses:    int64(0),
		evictions: int64(0),
	}

	return c, nil
}

func (c *MemcacheCache) Close(ctx context.Context) error {
	return c.client.Close()
}

func (c *MemcacheCache) Name() string {
	return "memcache"
}

func (c *MemcacheCache) Get(ctx context.Context, key string) (io.ReadSeekCloser, error) {

	slog.Debug("GET", "key", key)

	item, err := c.client.Get(key)

	if err != nil {

		if err == memcache.ErrCacheMiss {
			atomic.AddInt64(&c.misses, 1)
			return nil, new(cache.CacheMiss)
		}

		return nil, fmt.Errorf("Failed to retrieve cache item, %w", err)
	}

	br := bytes.NewReader(item.Value)
	r, err := ioutil.NewReadSeekCloser(br)

	if err != nil {
		return nil, fmt.Errorf("Failed to create new readseekcloser, %w", err)
	}

	atomic.AddInt64(&c.hits, 1)
	return r, nil
}

func (c *MemcacheCache) Set(ctx context.Context, key string, r io.ReadSeekCloser) (io.ReadSeekCloser, error) {

	slog.Debug("SET", "key", key)

	body, err := io.ReadAll(r)

	if err != nil {
		return nil, fmt.Errorf("Failed to read body, %w", err)
	}

	err = c.client.Set(&memcache.Item{Key: key, Value: body})

	if err != nil {
		return nil, fmt.Errorf("Failed to set memcache item, %w", err)
	}

	_, err = r.Seek(0, 0)

	if err != nil {
		return nil, fmt.Errorf("Failed to rewind body, %w", err)
	}

	return r, nil
}

func (c *MemcacheCache) Unset(ctx context.Context, key string) error {

	slog.Debug("UNSET", "key", key)

	err := c.client.Delete(key)

	if err != nil {
		return fmt.Errorf("Failed to delete key, %w", err)
	}

	atomic.AddInt64(&c.evictions, 1)
	return nil
}

func (c *MemcacheCache) Size() int64 {
	return 0
}

func (c *MemcacheCache) SizeWithContext(ctx context.Context) int64 {
	return 0
}

func (c *MemcacheCache) Hits() int64 {
	return atomic.LoadInt64(&c.hits)
}

func (c *MemcacheCache) Misses() int64 {
	return atomic.LoadInt64(&c.misses)
}

func (c *MemcacheCache) Evictions() int64 {
	return atomic.LoadInt64(&c.evictions)
}
