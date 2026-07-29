package category

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/farrasmumtaz/RentVibe/internal/catalog"
)

type repositoryStub struct {
	findAllCalls int
	items        []catalog.Category
	total        int64
}

func (r *repositoryStub) Create(*catalog.Category) error { return nil }
func (r *repositoryStub) FindAll(string, int, int) ([]catalog.Category, int64, error) {
	r.findAllCalls++
	return r.items, r.total, nil
}
func (r *repositoryStub) FindByID(uint) (*catalog.Category, error) { return nil, nil }
func (r *repositoryStub) Update(*catalog.Category) error           { return nil }
func (r *repositoryStub) Patch(*catalog.Category) error            { return nil }
func (r *repositoryStub) Delete(uint) error                        { return nil }

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

func TestFindAllCachesRepositoryResult(t *testing.T) {
	repository := &repositoryStub{
		items: []catalog.Category{{Name: "Camera"}},
		total: 1,
	}
	cacheStore := &cacheStub{values: make(map[string][]byte)}
	service := NewService(repository, cacheStore)

	first, firstTotal, err := service.FindAll("", 1, 10)
	if err != nil {
		t.Fatalf("first FindAll returned error: %v", err)
	}
	second, secondTotal, err := service.FindAll("", 1, 10)
	if err != nil {
		t.Fatalf("second FindAll returned error: %v", err)
	}

	if repository.findAllCalls != 1 {
		t.Fatalf("repository called %d times, want 1", repository.findAllCalls)
	}
	if firstTotal != 1 || secondTotal != 1 || first[0].Name != second[0].Name {
		t.Fatal("cached response does not match repository response")
	}
}

func TestCreateInvalidatesCategoryCache(t *testing.T) {
	repository := &repositoryStub{}
	cacheStore := &cacheStub{values: map[string][]byte{
		"categories:all:page=1": []byte(`{"items":[],"total":0}`),
		"items:all:page=1":      []byte(`{"items":[],"total":0}`),
	}}
	service := NewService(repository, cacheStore)

	if err := service.Create(&catalog.Category{Name: "Audio"}); err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	if _, exists := cacheStore.values["categories:all:page=1"]; exists {
		t.Fatal("category cache was not invalidated")
	}
	if _, exists := cacheStore.values["items:all:page=1"]; !exists {
		t.Fatal("unrelated item cache was invalidated")
	}
}
