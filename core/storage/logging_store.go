package storage

import (
	"encoding/base64"

	"fmt"

	"github.com/go-kit/kit/log"
	"github.com/monax/hoard/core/logging"
)

type loggingStore struct {
	store  Store
	logger log.Logger
}

// Decorates a Store with some simple logging of method/address pairs
func NewLoggingStore(store Store, logger log.Logger) *loggingStore {
	ls := &loggingStore{
		store:  store,
		logger: logging.TraceLogger(log.With(logger, "module", "storage")),
	}
	ls.logger = log.With(ls.logger, "store", ls.Name())
	return ls
}

var _ Store = (*loggingStore)(nil)

func (ls *loggingStore) Put(address, data []byte) error {
	logger := log.With(ls.logger, "method", "Put", "address", formatAddress(address))
	logger.Log()
	return logging.Err(logger, ls.store.Put(address, data))
}

func (ls *loggingStore) Get(address []byte) ([]byte, error) {
	logger := log.With(ls.logger, "method", "Get", "address", formatAddress(address))
	logger.Log()
	data, err := ls.store.Get(address)
	return data, logging.Err(logger, err)
}

func (ls *loggingStore) Stat(address []byte) (*StatInfo, error) {
	logger := log.With(ls.logger, "method", "Stat", "address", formatAddress(address))
	logger.Log()
	statInfo, err := ls.store.Stat(address)
	return statInfo, logging.Err(logger, err)
}

func (ls *loggingStore) Location(address []byte) string {
	ls.logger.Log("method", "Location", "address", formatAddress(address))
	return ls.store.Location(address)
}

func (ls *loggingStore) Name() string {
	return fmt.Sprintf("loggingStore(%s)", ls.store.Name())
}

func formatAddress(address []byte) string {
	return base64.StdEncoding.EncodeToString(address)
}
