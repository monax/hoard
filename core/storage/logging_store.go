package storage

import (
	"encoding/base64"

	"reflect"

	"github.com/go-kit/kit/log"
)

type loggingStore struct {
	store  Store
	logger log.Logger
}

// Decorates a Store with some simple logging of method/address pairs
func NewLoggingStore(store Store, logger log.Logger) *loggingStore {
	return &loggingStore{
		store: store,
		logger: log.With(logger, "module", "storage",
			"store", storeType(store)),
	}
}

var _ Store = (*loggingStore)(nil)

func (ls *loggingStore) Put(address, data []byte) error {
	ls.logger.Log("method", "Put", "address", formatAddress(address))
	return ls.store.Put(address, data)
}

func (ls *loggingStore) Get(address []byte) ([]byte, error) {
	ls.logger.Log("method", "Get", "address", formatAddress(address))
	return ls.store.Get(address)
}

func (ls *loggingStore) Stat(address []byte) (*StatInfo, error) {
	ls.logger.Log("method", "Stat", "address", formatAddress(address))
	return ls.store.Stat(address)
}

func (ls *loggingStore) Location(address []byte) string {
	ls.logger.Log("method", "Location", "address", formatAddress(address))
	return ls.store.Location(address)
}

func formatAddress(address []byte) string {
	return base64.StdEncoding.EncodeToString(address)
}

func storeType(store Store) string {
	storeType := reflect.TypeOf(store)
	var storeTypeName string
	if storeType.Kind() == reflect.Ptr {
		storeTypeName = storeType.Elem().String()
	} else {
		storeTypeName = storeType.String()
	}
	return storeTypeName
}
