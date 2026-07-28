package category

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/farrasmumtaz/RentVibe/internal/cache"
	"github.com/farrasmumtaz/RentVibe/internal/models"
)

const (
	categoryCachePrefix = "categories:"
	cacheTTL            = 5 * time.Minute
)

type categoryListCache struct {
	Items []models.Category `json:"items"`
	Total int64             `json:"total"`
}

type Service interface {
	Create(category *models.Category) error
	FindAll(search string, page int, limit int) ([]models.Category, int64, error)
	FindByID(id uint) (*models.Category, error)
	Update(id uint, req UpdateCategoryRequest) (*models.Category, error)
	Patch(id uint, req PatchCategoryRequest) (*models.Category, error)
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

func (s *service) Create(category *models.Category) error {
	if err := s.repository.Create(category); err != nil {
		return err
	}
	s.invalidateCache()
	return nil
}

func (s *service) FindAll(search string, page int, limit int) ([]models.Category, int64, error) {
	ctx := context.Background()
	key := fmt.Sprintf("%sall:search=%s:page=%d:limit=%d", categoryCachePrefix, url.QueryEscape(search), page, limit)
	var cached categoryListCache
	if hit, err := s.cache.Get(ctx, key, &cached); err == nil && hit {
		return cached.Items, cached.Total, nil
	}

	items, total, err := s.repository.FindAll(search, page, limit)
	if err == nil {
		_ = s.cache.Set(ctx, key, categoryListCache{Items: items, Total: total}, cacheTTL)
	}
	return items, total, err
}

func (s *service) FindByID(id uint) (*models.Category, error) {
	ctx := context.Background()
	key := fmt.Sprintf("%s%d", categoryCachePrefix, id)
	var cached models.Category
	if hit, err := s.cache.Get(ctx, key, &cached); err == nil && hit {
		return &cached, nil
	}

	result, err := s.repository.FindByID(id)
	if err == nil {
		_ = s.cache.Set(ctx, key, result, cacheTTL)
	}
	return result, err
}

func (s *service) Update(id uint, req UpdateCategoryRequest) (*models.Category, error) {

	category, err := s.repository.FindByID(id)
	if err != nil {
		return nil, err
	}

	category.Name = req.Name
	category.Description = req.Description

	err = s.repository.Update(category)
	if err != nil {
		return nil, err
	}

	s.invalidateCache()
	return category, nil
}

func (s *service) Patch(id uint, req PatchCategoryRequest) (*models.Category, error) {

	category, err := s.repository.FindByID(id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		category.Name = *req.Name
	}

	if req.Description != nil {
		category.Description = *req.Description
	}

	err = s.repository.Patch(category)
	if err != nil {
		return nil, err
	}

	s.invalidateCache()
	return category, nil
}

func (s *service) Delete(id uint) error {

	_, err := s.repository.FindByID(id)
	if err != nil {
		return err
	}

	if err := s.repository.Delete(id); err != nil {
		return err
	}
	s.invalidateCache()
	return nil
}

func (s *service) invalidateCache() {
	_ = s.cache.DeleteByPrefix(context.Background(), categoryCachePrefix)
}
