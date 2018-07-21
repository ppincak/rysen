package store

import "github.com/ppincak/rysen/crypto"

type Store struct {
	Symbols *crypto.Symbols
}

func NewStore() *Store {
	return &Store{}
}
