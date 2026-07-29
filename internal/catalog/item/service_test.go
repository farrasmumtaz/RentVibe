package item

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/farrasmumtaz/RentVibe/internal/catalog"
)

type repositoryStub struct {
	findByIDCalls int
	item          *catalog.Item
}

func (r *repositoryStub) Create(*catalog.Item) error { return nil }
func (r *repositoryStub) FindAll(string, int, int) ([]catalog.Item, int64, error) {
	return nil, 0, nil
}
func (r *repositoryStub) FindByID(uint) (*catalog.Item, error) {
	r.findByIDCalls++
	return r.item, nil
}
func (r *repositoryStub) Update(uint, *catalog.Item) error                    { return nil }
func (r *repositoryStub) Patch(uint, PatchItemRequest) (*catalog.Item, error) { return r.item, nil }
func (r *repositoryStub) Delete(uint) error                                   { return nil }

type cacheStub struct {
	values map[string][]byte
}

func (c *cacheStub) Get(_ context.Context, key string, destination interface{}) (bool, error) {
	value, exists := c.values[key]
	if !exists {
		return false, nil
	}
	return true, json.Unmarshal(value, destination)
}
func (c *cacheStub) Set(_ context.Context, key string, value interface{}, _ time.Duration) error {
	encoded, err := json.Marshal(value)
	if err == nil {
		c.values[key] = encoded
	}
	return err
}
func (c *cacheStub) DeleteByPrefix(_ context.Context, prefix string) error {
	for key := range c.values {
		if len(key) >= len(prefix) && key[:len(prefix)] == prefix {
			delete(c.values, key)
		}
	}
	return nil
}

func TestFindByIDUsesCacheAfterFirstQuery(t *testing.T) {
	repository := &repositoryStub{item: &catalog.Item{Name: "Sony A7"}}
	cacheStore := &cacheStub{values: make(map[string][]byte)}
	service := NewService(repository, cacheStore)

	first, err := service.FindByID(1)
	if err != nil {
		t.Fatalf("first FindByID returned error: %v", err)
	}
	second, err := service.FindByID(1)
	if err != nil {
		t.Fatalf("second FindByID returned error: %v", err)
	}

	if repository.findByIDCalls != 1 {
		t.Fatalf("repository called %d times, want 1", repository.findByIDCalls)
	}
	if first.Name != second.Name {
		t.Fatal("cached item does not match repository item")
	}
}

func TestDeleteInvalidatesItemAndCategoryCaches(t *testing.T) {
	repository := &repositoryStub{}
	cacheStore := &cacheStub{values: map[string][]byte{
		"items:1":      []byte(`{}`),
		"categories:1": []byte(`{}`),
	}}
	service := NewService(repository, cacheStore)

	if err := service.Delete(1); err != nil {
		t.Fatalf("Delete returned error: %v", err)
	}
	if len(cacheStore.values) != 0 {
		t.Fatalf("stale cache remains: %v", cacheStore.values)
	}
}
