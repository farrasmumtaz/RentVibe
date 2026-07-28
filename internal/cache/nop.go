package cache

import (
	"context"
	"time"
)

type NopStore struct{}

func NewNopStore() Store {
	return NopStore{}
}

func (NopStore) Get(context.Context, string, interface{}) (bool, error) {
	return false, nil
}

func (NopStore) Set(context.Context, string, interface{}, time.Duration) error {
	return nil
}

func (NopStore) DeleteByPrefix(context.Context, string) error {
	return nil
}
