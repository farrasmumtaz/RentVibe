package item

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/farrasmumtaz/RentVibe/internal/cache"
	"github.com/farrasmumtaz/RentVibe/internal/catalog"
)

const (
	itemCachePrefix = "items:"
	itemCacheTTL    = 5 * time.Minute
)

type itemListCache struct {
	Items []catalog.Item `json:"items"`
	Total int64          `json:"total"`
}

type Service interface {
	Create(item *catalog.Item) error
	FindAll(search string, page int, limit int) ([]catalog.Item, int64, error)
	FindByID(id uint) (*catalog.Item, error)
	Update(id uint, item *catalog.Item) error
	Patch(id uint, req PatchItemRequest) (*catalog.Item, error)
	Delete(id uint) error
}

type service struct {
	repository Repository
	cache      cache.Store
}

func NewService(repository Repository, cacheStore cache.Store) Service {
	return &service{
		repository: repository,
		cache:      cacheStore,
	}
}

func (s *service) Create(item *catalog.Item) error {
	if err := s.repository.Create(item); err != nil {
		return err
	}
	s.invalidateCache()
	return nil
}

func (s *service) FindAll(search string, page int, limit int) ([]catalog.Item, int64, error) {
	ctx := context.Background()
	key := fmt.Sprintf("%sall:search=%s:page=%d:limit=%d", itemCachePrefix, url.QueryEscape(search), page, limit)
	var cached itemListCache
	if hit, err := s.cache.Get(ctx, key, &cached); err == nil && hit {
		return cached.Items, cached.Total, nil
	}

	items, total, err := s.repository.FindAll(search, page, limit)
	if err == nil {
		_ = s.cache.Set(ctx, key, itemListCache{Items: items, Total: total}, itemCacheTTL)
	}
	return items, total, err
}

func (s *service) FindByID(id uint) (*catalog.Item, error) {
	ctx := context.Background()
	key := fmt.Sprintf("%s%d", itemCachePrefix, id)
	var cached catalog.Item
	if hit, err := s.cache.Get(ctx, key, &cached); err == nil && hit {
		return &cached, nil
	}

	result, err := s.repository.FindByID(id)
	if err == nil {
		_ = s.cache.Set(ctx, key, result, itemCacheTTL)
	}
	return result, err
}

func (s *service) Update(id uint, item *catalog.Item) error {
	if err := s.repository.Update(id, item); err != nil {
		return err
	}
	s.invalidateCache()
	return nil
}

func (s *service) Patch(id uint, req PatchItemRequest) (*catalog.Item, error) {
	result, err := s.repository.Patch(id, req)
	if err == nil {
		s.invalidateCache()
	}
	return result, err
}

func (s *service) Delete(id uint) error {
	if err := s.repository.Delete(id); err != nil {
		return err
	}
	s.invalidateCache()
	return nil
}

func (s *service) invalidateCache() {
	ctx := context.Background()
	_ = s.cache.DeleteByPrefix(ctx, itemCachePrefix)
	_ = s.cache.DeleteByPrefix(ctx, "categories:")
}
