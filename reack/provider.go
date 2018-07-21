package reack

import (
	"context"
)

type Provider struct {
	store Store
}

type Store interface {
	HasState
}

type HasState interface {
	GetState() interface{}
}

func NewProvider(store Store) *Provider {
	return &Provider{store: store}
}

func (this Provider) Context() context.Context {
	return nil
}