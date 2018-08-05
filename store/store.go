package store

import "github.com/ppincak/rysen/crypto"

type Store struct {
	Symbols *crypto.Symbols
	Limits  []*Limit
}

func NewStore() *Store {
	return &Store{}
}
