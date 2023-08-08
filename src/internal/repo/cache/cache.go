package cache

import (
	"errors"
	"fmt"

	"gitlab.com/Hofsiedge/l0/internal/repo"
)

var ErrorNotFound = errors.New("not found in cache")

// *Cache[T, I] implements repo.Repo[T, I]
type Cache[T any, I comparable] struct {
	inner repo.Repo[T, I]
	cache map[I]T
	valid bool
}

func New[T any, I comparable](inner repo.Repo[T, I]) (*Cache[T, I], error) {
	cache := Cache[T, I]{
		inner: inner,
		cache: make(map[I]T),
	}
	values, err := cache.inner.GetAll()
	if err != nil {
		return nil, err
	}
	for _, v := range values {
		cache.cache[inner.GetKey(v)] = v
	}
	return &cache, nil
}

func (c *Cache[T, I]) Get(id I) (T, error) {
	obj, ok := c.cache[id]
	if !ok {
		return obj, fmt.Errorf("%w: missing key %v", ErrorNotFound, id)
	}
	return obj, nil
}

func (c *Cache[T, I]) List() ([]I, error) {
	ids := make([]I, len(c.cache))
	idx := 0
	for key := range c.cache {
		ids[idx] = key
		idx++
	}
	return ids, nil
}

func (c *Cache[T, I]) GetAll() ([]T, error) {
	values := make([]T, len(c.cache))
	idx := 0
	for _, v := range c.cache {
		values[idx] = v
		idx++
	}
	return values, nil
}

func (c *Cache[T, I]) Save(obj T) error {
	if err := c.inner.Save(obj); err != nil {
		return err
	}
	c.cache[c.inner.GetKey(obj)] = obj
	return nil
}

func (c *Cache[T, I]) GetKey(obj T) I {
	return c.inner.GetKey(obj)
}
