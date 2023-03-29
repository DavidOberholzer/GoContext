package storage

import (
	"errors"
	"server"
	"sync"
)

var (
	ErrNotFound      = errors.New("object not found")
	ErrAlreadyExists = errors.New("object already exists")
)

type db struct {
	store map[string]server.Object
	mutex sync.RWMutex
}

func newDB() *db {
	return &db{
		store: make(map[string]server.Object, 0),
		mutex: sync.RWMutex{},
	}
}

func (d *db) Get(key string) (server.Object, error) {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	obj, ok := d.store[key]
	if !ok {
		return server.Object{}, ErrNotFound
	}
	return obj, nil
}

func (d *db) Set(key string, value server.Object) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	_, ok := d.store[key]
	if ok {
		return ErrAlreadyExists
	}

	d.store[key] = value
	return nil
}
