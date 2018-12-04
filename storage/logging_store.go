package storage

import (
	"encoding/base64"

	"fmt"

	"github.com/go-kit/kit/log"
	"github.com/monax/hoard/logging"
)

type loggingStore struct {
	store  NamedStore
	logger log.Logger
}

// Decorates a Store with some simple logging of method/address pairs
func NewLoggingStore(store NamedStore, logger log.Logger) *loggingStore {
	ls := &loggingStore{
		store:  store,
		logger: logging.TraceLogger(log.With(logger, "module", "storage")),
	}
	ls.logger = log.With(ls.logger, "store", ls.Name())
	return ls
}

var _ NamedStore = (*loggingStore)(nil)

func (ls *loggingStore) Put(address []byte, data []byte) ([]byte, error) {
	address, err := ls.store.Put(address, data)
	return address, logErrorOrSuccess(log.With(ls.logger, "method", "Put", "address",
		formatAddress(address)), err)
}

func (ls *loggingStore) Get(address []byte) ([]byte, error) {
	data, err := ls.store.Get(address)
	return data, logErrorOrSuccess(log.With(ls.logger, "method", "Get", "address",
		formatAddress(address)), err)
}

func (ls *loggingStore) Stat(address []byte) (*StatInfo, error) {
	statInfo, err := ls.store.Stat(address)
	return statInfo, logErrorOrSuccess(log.With(ls.logger, "method", "Stat", "address",
		formatAddress(address)), err)
}

func (ls *loggingStore) Location(address []byte) string {
	ls.logger.Log("method", "Location", "address", formatAddress(address))
	return ls.store.Location(address)
}

func (ls *loggingStore) Name() string {
	return fmt.Sprintf("loggingStore(%s)", ls.store.Name())
}

func logErrorOrSuccess(logger log.Logger, err error) error {
	err = logging.Err(logger, err)
	if err == nil {
		logging.Msg(logger, "Success")
	}
	return err
}

func formatAddress(address []byte) string {
	return base64.StdEncoding.EncodeToString(address)
}
