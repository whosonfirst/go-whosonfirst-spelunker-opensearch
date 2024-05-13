package ristretto

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	_ "net/url"
	"sync/atomic"

	dg_ristretto "github.com/dgraph-io/ristretto"
	"github.com/whosonfirst/go-cache"
	"github.com/whosonfirst/go-ioutil"
)

type RistrettoCache struct {
	cache.Cache
	misses    int64
	hits      int64
	evictions int64
	client    *dg_ristretto.Cache
}

func init() {
	ctx := context.Background()
	cache.RegisterCache(ctx, "ristretto", NewRistrettoCache)
}

func NewRistrettoCache(ctx context.Context, uri string) (cache.Cache, error) {

	max_cost := int64(500000000)
	buffer_items := int64(64)

	cfg := &dg_ristretto.Config{
		NumCounters: 1e7, // number of keys to track frequency of (10M).
		MaxCost:     max_cost,
		BufferItems: buffer_items,
	}

	client, err := dg_ristretto.NewCache(cfg)

	if err != nil {
		return nil, fmt.Errorf("Failed to create new cache, %w", err)
	}

	c := &RistrettoCache{
		client:    client,
		hits:      int64(0),
		misses:    int64(0),
		evictions: int64(0),
	}

	return c, nil
}

func (c *RistrettoCache) Close(ctx context.Context) error {
	c.client.Close()
	return nil
}

func (c *RistrettoCache) Name() string {
	return "ristretto"
}

func (c *RistrettoCache) Get(ctx context.Context, key string) (io.ReadSeekCloser, error) {

	slog.Debug("GET", "key", key)

	v, exists := c.client.Get(key)

	if !exists {

		atomic.AddInt64(&c.misses, 1)
		return nil, new(cache.CacheMiss)
	}

	br := bytes.NewReader(v.([]byte))
	r, err := ioutil.NewReadSeekCloser(br)

	if err != nil {
		return nil, fmt.Errorf("Failed to create new readseekcloser, %w", err)
	}

	atomic.AddInt64(&c.hits, 1)
	return r, nil
}

func (c *RistrettoCache) Set(ctx context.Context, key string, r io.ReadSeekCloser) (io.ReadSeekCloser, error) {

	slog.Debug("SET", "key", key)

	body, err := io.ReadAll(r)

	if err != nil {
		return nil, fmt.Errorf("Failed to read body, %w", err)
	}

	ok := c.client.Set(key, body, 0) // int64(len(body)))

	if !ok {
		return nil, fmt.Errorf("Failed to set ristretto item")
	}

	_, err = r.Seek(0, 0)

	if err != nil {
		return nil, fmt.Errorf("Failed to rewind body, %w", err)
	}

	c.client.Wait()

	return r, nil
}

func (c *RistrettoCache) Unset(ctx context.Context, key string) error {

	slog.Debug("UNSET", "key", key)

	c.client.Del(key)

	atomic.AddInt64(&c.evictions, 1)
	return nil
}

func (c *RistrettoCache) Size() int64 {
	return 0
}

func (c *RistrettoCache) SizeWithContext(ctx context.Context) int64 {
	return 0
}

func (c *RistrettoCache) Hits() int64 {
	return atomic.LoadInt64(&c.hits)
}

func (c *RistrettoCache) Misses() int64 {
	return atomic.LoadInt64(&c.misses)
}

func (c *RistrettoCache) Evictions() int64 {
	return atomic.LoadInt64(&c.evictions)
}
