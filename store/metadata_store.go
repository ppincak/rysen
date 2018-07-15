package store

import "github.com/ppincak/rysen/core"

type MetadataStore struct {
	Symbols *core.Symbols
}

func NewMetadataStore() *MetadataStore {
	return &MetadataStore{}
}
